package indexdb

import (
	"github.com/tinywasm/tinyreflect"
	. "github.com/tinywasm/fmt"
)

func (d *IndexDB) Update(table_name string, items ...any) (err error) {

	const e = "Update"
	// Obtener el almacén
	if len(items) == 0 {
		return Err("no data to update table", table_name)
	}

	// Create table if it doesn't exist using the first item as template
	if d.err = d.prepareStoreWithTableCheck("update", table_name, items[0]); d.err != nil {
		return Errf("%s %v", e, d.err)
	}

	d.data = make([]map[string]any, len(items))

	// Find primary key field
	pk_field := ""
	if len(items) > 0 {
		v := tinyreflect.ValueOf(items[0])
		isPtr := v.Kind() == K.Pointer
		if isPtr {
			elem, _ := v.Elem()
			v = elem
		}
		st := v.Type()
		if st.Kind() == K.Struct {
			structType := st.StructType()
			for _, f := range structType.Fields {
				fieldName := f.Name.String()
				_, isPK := IDorPrimaryKey(table_name, fieldName)
				if isPK {
					if pk_field != "" {
						return Errf("%s multiple primary keys found", e)
					}
					pk_field = fieldName
				}
			}
		}
	}
	if pk_field == "" {
		return Errf("%s no primary key found", e)
	}

	for i, item := range items {

		v := tinyreflect.ValueOf(item)
		isPtr := v.Kind() == K.Pointer
		if isPtr {
			elem, _ := v.Elem()
			v = elem
		}

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

				// Check if this is the primary key field
				_, isPK := IDorPrimaryKey(table_name, f.Name.String())
				if isPK {

					m[pk_field] = val

				} else {

					m[fieldName] = val

				}

			}

			d.data[i] = m

		}

	}

	// Iterar sobre los datos a actualizar
	for _, obj := range d.data {

		// Obtener el ID del objeto
		id, ok := obj[pk_field].(string)
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
