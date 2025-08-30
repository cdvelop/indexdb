package indexdb

import (
	. "github.com/cdvelop/tinystring"
)

func (d *indexDB) ClearAllTableDataInDB(tables ...string) (err error) {
	const e = "ClearAllTableDataInDB error"
	for _, table_name := range tables {

		d.err = d.prepareStore("clear", table_name)
		if d.err != nil {
			return Errf("%s %v", e, d.err)
		}

		// elimina cada elemento en el almacén de objetos
		result := d.store.Call("clear")

		if result.IsNull() {
			return Errf("%s error clearing data from table: %s", e, table_name)
		}

	}

	return nil
}
