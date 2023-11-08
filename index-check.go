package indexdb

import (
	"syscall/js"

	"github.com/cdvelop/model"
)

func fieldIndexOK(table, field_name string, store *js.Value) error {

	// Verificar si el índice existe en la tabla.
	indexNames := store.Get("indexNames")

	indexSet := make(map[string]bool)
	for i := 0; i < indexNames.Length(); i++ {
		indexSet[indexNames.Index(i).String()] = true
	}

	// Verificar si el índice existe en la tabla.
	if !indexSet[field_name] {
		return model.Error("El índice:", field_name, "no existe en la tabla:", table)
	}

	return nil
}
