package indexdb

func (d *indexDB) fieldIndexOK(table, field_name string) (err string) {

	// Verificar si el índice existe en la tabla.
	indexNames := d.store.Get("indexNames")

	indexSet := make(map[string]bool)
	for i := 0; i < indexNames.Length(); i++ {
		indexSet[indexNames.Index(i).String()] = true
	}

	// Verificar si el índice existe en la tabla.
	if !indexSet[field_name] {
		return "El índice: " + field_name + "no existe en la tabla:" + table
	}

	return ""
}
