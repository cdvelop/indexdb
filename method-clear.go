package indexdb

func (d *indexDB) ClearAllTableDataInDB(tables ...string) (err string) {
	const this = "ClearAllTableDataInDB error "
	for _, table_name := range tables {

		store, err := d.getStore("clear", table_name)
		if err != "" {
			return this + err
		}

		// elimina cada elemento en el almacén de objetos
		result := store.Call("clear")

		if result.IsNull() {
			return this + "al borrar la data de la tabla: " + table_name
		}

	}

	return ""
}
