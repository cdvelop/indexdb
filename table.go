package indexdb

import (
	. "github.com/cdvelop/tinystring"
)

func (d *indexDB) createTable(o *Object) {

	if !o.NoAddObjectInDB {
		// d.Log("**** CREANDO TABLA: ", o.Table, "INDEX DB")
		if len(o.Fields) != 0 {

			pk_name := o.PrimaryKeyName()

			newTable := d.db.Call("createObjectStore", o.Table, map[string]interface{}{"keyPath": pk_name})

			for _, f := range o.Fields {

				if !f.NotRequiredInDB {
					// Crear un índices para búsqueda campos principales
					newTable.Call("createIndex", f.Name, f.Name, map[string]interface{}{"unique": f.Unique})
				}
			}
		}
	}
}

func (d indexDB) checkTableStatus(operation, table_name string) (err error) {

	if !d.db.Truthy() {
		return Err(operation, "variable db not defined in index db", table_name)
	}

	if !d.TableExist(table_name) {
		return Err("error action:", operation, ". table", table_name, "does not exist in indexdb")
	}

	return nil
}

func (d indexDB) TableExist(table_name string) bool {

	// Obtiene la lista de nombres de almacenes de objetos en la base de datos
	objectStoreNames := d.db.Get("objectStoreNames")
	length := objectStoreNames.Length()

	// Itera a través de los nombres de las tablas y verifica si la tabla ya existe
	for i := 0; i < length; i++ {
		name := objectStoreNames.Index(i).String()
		if name == table_name {
			return true
		}
	}

	return false
}
