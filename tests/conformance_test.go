//go:build wasm

package tests_test

import (
	"fmt"
	"testing"

	"webtyp.com/model"
	"webtyp.com/storage"
	dbconf "webtyp.com/storage/conformance"
)

func TestIndexDB_DBConformance(t *testing.T) {
	var n int
	dbconf.Run(t, dbconf.Factory{
		Name: "indexdb",
		New: func(t *testing.T, models ...model.Model) storage.Conn {
			n++
			dbName := fmt.Sprintf("conformance_db_%d", n) // fresh IndexedDB per clause
			structs := make([]any, len(models))
			for i, m := range models {
				structs[i] = m // declared as object stores up front
			}
			return SetupDB(func(...any) {}, dbName, structs...) // from tests/setup_test.go
		},
	})
}
