package indexdb

func (d *indexDB) DeleteObjectsInDB(table_name string, on_server_too bool, all_data ...map[string]string) (err string) {

	const e = "DeleteObjectsInDB "

	if on_server_too {
		d.BackupOneObjectType("delete", table_name, all_data)
	}

	d.err = d.prepareStore("delete", table_name)
	if d.err != "" {
		return e + d.err
	}

	for _, data := range all_data {
		if id, exist := data[PREFIX_ID_NAME+table_name]; exist {
			// elimina cada elemento en el almacén de objetos
			d.result = d.store.Call("delete", id)

			if d.result.IsNull() {
				return e + "al eliminar en la tabla: " + table_name
			}

		} else {
			return e + "id no encontrado tabla: " + table_name
		}

	}

	return
}
