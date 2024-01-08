package indexdb

// action create,read, delete, update
func (d *indexDB) prepareStore(action, table_name string) (err string) {

	if d.err = d.checkTableStatus(action, table_name); d.err != "" {
		return d.err
	}

	var readwrite = "readonly"
	if action != "read" {
		readwrite = "readwrite"
	}

	// Obtiene una transacción de escritura
	d.transaction = d.db.Call("transaction", table_name, readwrite)

	// Obtiene el almacén de objetos
	d.store = d.transaction.Call("objectStore", table_name)

	if !d.store.Truthy() {
		return "error no se logro abrir el almacén: " + table_name + " db para la acción " + action
	}

	return ""
}
