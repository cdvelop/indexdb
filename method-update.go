package indexdb

import (
	"github.com/cdvelop/tinyreflect"
	. "github.com/cdvelop/tinystring"
)

func (d *indexDB) Update(table_name string, all_data ...interface{}) (err error) {

	const e = "Update"
	// Obtener el almacén
	if len(all_data) == 0 {
		return nil
	}

	// Create table if it doesn't exist using the first item as template
	if d.err = d.prepareStoreWithTableCheck("update", table_name, all_data[0]); d.err != nil {
		return Errf("%s %v", e, d.err)
	}

	items := all_data

	d.data_in_any = make([]map[string]interface{}, len(items))

	d.data_in_str = nil

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

					pk_field := "id_" + table_name

					m[pk_field] = val

				} else {

					m[fieldName] = val

				}

			}

			d.data_in_any[i] = m

		}

	}

	// Iterar sobre los datos a actualizar
	for _, obj := range d.data_in_any {

		// Obtener el ID del objeto
		id, ok := obj["id_"+table_name].(string)
		if !ok || id == "" {
			return Errf("%s invalid object without ID to update", e)
		}

		// Guardar los cambios
		d.result = d.store.Call("put", obj)
		if d.result.IsNull() {
			return Errf("%s when updating object %s in the db", e, id)
		}

	}

	return nil
}
