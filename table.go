package indexdb

import (
	"github.com/cdvelop/tinyreflect"
	. "github.com/cdvelop/tinystring"
)

// CreateTableIfNotExists creates a table for the given struct type if it doesn't exist
func (d *indexDB) CreateTableIfNotExists(tableName string, structType interface{}) error {
	// Check if table already exists
	if d.TableExist(tableName) {
		return nil
	}

	// Create the table
	return d.createTable(structType)
}

// createTable creates a table for the given struct type
func (d *indexDB) createTable(structType interface{}) error {
	st := tinyreflect.TypeOf(structType)

	if st.Kind() == K.Struct {
		structTypeInfo := st.StructType()

		table_name := st.Name()

		if len(structTypeInfo.Fields) != 0 {
			pk_name := "id_" + table_name

			// Create object store
			newTable := d.db.Call("createObjectStore", table_name, map[string]interface{}{"keyPath": pk_name})

			// Create indexes for fields
			for _, f := range structTypeInfo.Fields {
				fieldName := f.Name.String()
				tag := f.Tag().Get("db")

				// Skip ID field (it's the keyPath)
				if IsPrimaryKey(fieldName, table_name) {
					continue
				}

				// Create index if field has db tag or is marked as unique
				if tag != "" {
					unique := tag == "unique"
					newTable.Call("createIndex", fieldName, fieldName, map[string]interface{}{"unique": unique})
				}
			}
		}
	}

	return nil
}

// TableExist checks if a table exists in the database
func (d indexDB) TableExist(table_name string) bool {
	// Get the list of object store names from the database
	objectStoreNames := d.db.Get("objectStoreNames")
	length := objectStoreNames.Length()

	// Iterate through the table names and check if the table already exists
	for i := 0; i < length; i++ {
		name := objectStoreNames.Index(i).String()
		if name == table_name {
			return true
		}
	}

	return false
}
