package indexdb

import (
	"github.com/cdvelop/tinyreflect"
	. "github.com/cdvelop/tinystring"
)

func (d *indexDB) Delete(table_name string, all_data ...interface{}) (err error) {

	const e = "Delete"

	if d.err = d.prepareStore("delete", table_name); d.err != nil {
		return Errf("%s %v", e, d.err)
	}

	for _, item := range all_data {

		v := tinyreflect.ValueOf(item)

		st := v.Type()

		if st.Kind() == K.Struct {

			structType := st.StructType()

			found := false

			for j, f := range structType.Fields {

				// Check if this is the ID field by name
				if IsPrimaryKey(f.Name.String(), table_name) {

					fieldValue, _ := v.Field(j)

					id, _ := fieldValue.Interface()

					d.result = d.store.Call("delete", id)

					if d.result.IsNull() {

						return Errf("%s error when deleting in table: %s", e, table_name)

					}

					found = true

					break

				}

			}

			if !found {

				return Errf("%s id not found in table: %s", e, table_name)

			}

		}

	}

	return nil
}
