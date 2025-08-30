package indexdb

func (d *indexDB) UpdateObjectsInDB(table_name string, on_server_too bool, all_data ...map[string]string) (err string) {

	const e = "UpdateObjectsInDB "

	if on_server_too {
		d.BackupOneObjectType("update", table_name, all_data)
	}
	// Obtener el almacén
	d.err = d.prepareStore("update", table_name)
	if d.err != "" {
		return e + d.err
	}

	d.prepareDataIN(all_data, true)

	// Iterar sobre los datos a actualizar
	for _, obj := range d.data_in_any {

		// Obtener el ID del objeto
		id, ok := obj[PREFIX_ID_NAME+table_name].(string)
		if !ok || id == "" {
			return e + "objeto invalido sin ID para actualizar "
		}

		// Guardar los cambios
		d.result = d.store.Call("put", obj)
		if d.result.IsNull() {
			return e + "al actualizar objeto " + id + " en la db"
		}

	}

	return ""
}
