package indexdb

import (
	"syscall/js"

	"github.com/cdvelop/tinyreflect"
	. "github.com/cdvelop/tinystring"
)

var blob_exist bool
var blob_file interface{}

// items support: []interface{} (struct instances)
func (d *indexDB) Create(table_name string, items []interface{}) (err error) {

	blob_exist = false
	blob_file = nil

	const e = "indexdb create"

	if len(items) == 0 {
		return nil
	}

	// Create table if it doesn't exist using the first item as template
	if d.err = d.prepareStoreWithTableCheck("create", table_name, items[0]); d.err != nil {
		return Errf("%s %v", e, d.err)
	}

	d.data_in_any = make([]map[string]interface{}, len(items))

	d.data_in_str = nil

	pk_field := "id_" + table_name

	for i, item := range items {

		v := tinyreflect.ValueOf(item)

		st := v.Type()

		if st.Kind() == K.Struct {

			m := make(map[string]interface{})

			structType := st.StructType()

			for j, f := range structType.Fields {

				fieldName := f.Name.String()

				tag := f.Tag().Get("db")

				// Use tag value as field name if present, otherwise use field name
				if tag != "" {
					fieldName = tag
				}

				fieldValue, _ := v.Field(j)

				val, _ := fieldValue.Interface()

				// Check if this is the ID field by name
				if IsPrimaryKey(f.Name.String(), table_name) {

					m[pk_field] = val

				} else {

					m[fieldName] = val

				}

			}

			d.data_in_any[i] = m

		}

	}

	// CHECK ID
	for i, data := range d.data_in_any {

		id, id_exist := data[pk_field]

		// d.Log("DATA table_name DB", table_name, "id_exist:", id_exist)

		if !id_exist || id.(string) == "" {
			//agregar id al objeto si este no existe
			id = d.GetNewID() //id nuevo
			// d.Log("NUEVO ID GENERADO:", id)
			if !Contains(id.(string), ".") {
				return Errf("%s generated id does not contain user number", e)
			}

			data[pk_field] = id
		}

		// si todo esta ok retornamos el id
		if len(d.data_in_str) != 0 { // se envió la data en string

			d.data_in_str[i][pk_field] = id.(string)

		} else { // se envió la data de tipo any

			d.data_in_any[i][pk_field] = id.(string)
		}
	}

	// fmt.Println("DATA IN INDEX DB:", d.data_in_any)

	for _, data := range d.data_in_any {

		// Inserta cada elemento en el almacén de objetos
		result := d.store.Call("add", data)
		if result.IsNull() {
			return Err("error creating element in db table:", table_name, "id:", data[pk_field].(string))
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
			d.Log("Error al crear elemento en la db tabla:", table_name, "id:", data[pk_field].(string), errorMessage)
			return nil
		}))

		// creamos url temporal para acceder al archivo en el dom
		if blob_file, blob_exist = data["blob"]; blob_exist {
			data["url"] = CreateBlobURL(blob_file)
		}

	}

	return nil
}
