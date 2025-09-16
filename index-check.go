package indexdb

import (
	. "github.com/cdvelop/tinystring"
)

func (d *IndexDB) fieldIndexOK(table, field_name string) (err error) {

	// Verificar si el índice existe en la tabla.
	indexNames := d.store.Get("indexNames")

	indexSet := make(map[string]bool)
	for i := 0; i < indexNames.Length(); i++ {
		indexSet[indexNames.Index(i).String()] = true
	}

	// Verificar si el índice existe en la tabla.
	if !indexSet[field_name] {
		return Err("Index:", field_name, "does not exist in table:", table)
	}

	return nil
}
