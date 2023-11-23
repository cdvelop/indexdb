package indexdb

import (
	"syscall/js"
)

// action create,read, delete, update
func (d *indexDB) getStore(action, table_name string) (st *js.Value, err string) {

	if err := d.checkTableStatus(action, table_name); err != "" {
		return nil, err
	}

	var readwrite = "readonly"
	if action != "read" {
		readwrite = "readwrite"
	}

	// Obtiene una transacción de escritura
	transaction := d.db.Call("transaction", table_name, readwrite)

	// Obtiene el almacén de objetos
	store := transaction.Call("objectStore", table_name)

	if !store.Truthy() {
		return nil, "error no se logro abrir el almacén: " + table_name + " db para la acción " + action
	}

	return &store, ""
}
