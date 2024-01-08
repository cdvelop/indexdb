package indexdb

func (d *indexDB) ClearAllTableDataInDB(tables ...string) (err string) {
	const e = "ClearAllTableDataInDB error "
	for _, table_name := range tables {

		d.err = d.prepareStore("clear", table_name)
		if d.err != "" {
			return e + d.err
		}

		// elimina cada elemento en el almacén de objetos
		result := d.store.Call("clear")

		if result.IsNull() {
			return e + "al borrar la data de la tabla: " + table_name
		}

	}

	return ""
}
