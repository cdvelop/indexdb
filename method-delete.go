package indexdb

import (
	. "github.com/cdvelop/tinystring"
)

func (d *indexDB) DeleteObjectsInDB(table_name string, on_server_too bool, all_data ...map[string]string) (err error) {

	const e = "DeleteObjectsInDB"

	if on_server_too {
		d.BackupOneObjectType("delete", table_name, all_data)
	}

	if d.err = d.prepareStore("delete", table_name); d.err != nil {
		return Errf("%s %v", e, d.err)
	}

	for _, data := range all_data {
		if id, exist := data[PREFIX_ID_NAME+table_name]; exist {
			// elimina cada elemento en el almacén de objetos
			d.result = d.store.Call("delete", id)

			if d.result.IsNull() {
				return Errf("%s error when deleting in table: %s", e, table_name)
			}

		} else {
			return Errf("%s id not found in table: %s", e, table_name)
		}

	}

	return nil
}
