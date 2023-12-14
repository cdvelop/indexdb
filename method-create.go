package indexdb

import (
	"github.com/cdvelop/model"
)

// items support: []map[string]string or map[string]string
func (d *indexDB) CreateObjectsInDB(table_name string, backup_required bool, items any) (err string) {

	store, err := d.getStore("create", table_name)
	if err != "" {
		return err
	}

	for _, data := range DataConvertToAny(items) {

		pk_field := model.PREFIX_ID_NAME + table_name

		id, id_exist := data[pk_field]

		// d.Log("DATA table_name DB", table_name, "id_exist:", id_exist)

		if !id_exist || id.(string) == "" {

			if !backup_required { // si no requiere backup es un objeto sin id del servidor retornamos error
				err := "error data proveniente del servidor sin id en tabla: " + table_name
				d.Log(err, data)
				return err
			}

			//agregar id al objeto si este no existe
			id, err = d.GetNewID() //id nuevo
			if err != "" {
				return err
			}

			d.Log("NUEVO ID GENERADO:", id)

			// date, _ := unixid.UnixNanoToStringDate(id.(string))
			// d.Log("SU FECHA ES:", date)

			data[pk_field] = id
		}

		if backup_required { // necesita respaldo en servidor
			data["backup"] = "create" //estado backup = no respaldado
		}

		// Inserta cada elemento en el almacén de objetos
		result := store.Call("add", data)
		if result.IsNull() {

			return "error al crear elemento en la db tabla: " + table_name + " id: " + id.(string)
		}

		// retornamos url temporal para acceder al archivo
		if blob, exist := data["blob"]; exist {
			data["url"] = CreateBlobURL(blob)
		}
	}

	return ""
}
