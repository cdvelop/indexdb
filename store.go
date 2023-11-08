package indexdb

import (
	"syscall/js"

	"github.com/cdvelop/model"
)

// action create,read, delete, update
func (d *indexDB) getStore(action, table_name string) (*js.Value, error) {

	if err := d.checkTableStatus(action, table_name); err != nil {
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
		return nil, model.Error("error no se logro abrir el almacén:", table_name, "db para la acción", action)
	}

	return &store, nil
}
