package indexdb

import (
	"syscall/js"

	"github.com/cdvelop/model"
	"github.com/cdvelop/strings"
)

// items support: []map[string]string or map[string]string
func (d *indexDB) CreateObjectsInDB(table_name string, backup_required bool, items any) (err string) {

	const e = "indexdb create "

	store, err := d.getStore("create", table_name)
	if err != "" {
		return e + err
	}

	d.prepareDataIN(items)

	// fmt.Println("DATA IN INDEX DB:", d.data_in_any)

	for i, data := range d.data_in_any {

		pk_field := model.PREFIX_ID_NAME + table_name

		id, id_exist := data[pk_field]

		// d.Log("DATA table_name DB", table_name, "id_exist:", id_exist)

		if !id_exist || id.(string) == "" {

			if !backup_required { // si no requiere backup es un objeto sin id del servidor retornamos error
				err := e + "error data proveniente del servidor sin id en tabla: " + table_name
				d.Log(err, data)
				return err
			}

			//agregar id al objeto si este no existe
			id, err = d.GetNewID() //id nuevo
			if err != "" {
				return e + err
			}
			// d.Log("NUEVO ID GENERADO:", id)
			if strings.Contains(id.(string), ".") == 0 {
				return e + "id generado no contiene numero de usuario"
			}

			// date, _ := unixid.UnixNanoToStringDate(id.(string))
			// d.Log("SU FECHA ES:", date)

			data[pk_field] = id
		}

		if backup_required { // necesita respaldo en servidor
			data["create"] = "backup" //estado backup = no respaldado
		}

		// d.Log("insertando en indexdb:", data)

		// Inserta cada elemento en el almacén de objetos
		result := store.Call("add", data)
		if result.IsNull() {

			return "error al crear elemento en la db tabla: " + table_name + " id: " + id.(string)
		}
		// d.Log("resultado:", result)

		// result.Call("addEventListener", "success", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		// 	d.Log("Elemento creado con éxito:", data)
		// 	return nil
		// }))

		// Manejar la respuesta de manera asincrónica
		result.Call("addEventListener", "error", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
			// Log más detalles sobre el error
			errorObject := p[0].Get("target").Get("error")
			errorMessage := errorObject.Get("message").String()
			d.Log("Error al crear elemento en la db tabla:", table_name, "id:", id.(string), errorMessage)
			return nil
		}))

		// retornamos url temporal para acceder al archivo
		if blob, exist := data["blob"]; exist {
			data["url"] = CreateBlobURL(blob)
		}

		// si todo esta ok retornamos el id
		if len(d.data_in_str) != 0 { // se envió la data en string

			d.data_in_str[i][pk_field] = id.(string)

		} else { // se envió la data de tipo any

			d.data_in_any[i][pk_field] = id.(string)
		}

	}

	return ""
}
