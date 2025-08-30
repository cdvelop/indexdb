package indexdb

import (
	. "github.com/cdvelop/tinystring"
)

func (d *indexDB) UpdateObjectsInDB(table_name string, on_server_too bool, all_data ...map[string]string) (err error) {

	const e = "UpdateObjectsInDB"

	if on_server_too {
		d.BackupOneObjectType("update", table_name, all_data)
	}
	// Obtener el almacén
	if d.err = d.prepareStore("update", table_name); d.err != nil {
		return Errf("%s %v", e, d.err)
	}

	d.prepareDataIN(all_data, true)

	// Iterar sobre los datos a actualizar
	for _, obj := range d.data_in_any {

		// Obtener el ID del objeto
		id, ok := obj[PREFIX_ID_NAME+table_name].(string)
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
