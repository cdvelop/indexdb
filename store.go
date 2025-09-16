package indexdb

import (
	. "github.com/cdvelop/tinystring"
)

// action create,read, delete, update
func (d *IndexDB) prepareStore(action, table_name string) (err error) {

	var readwrite = "readonly"
	if action != "read" {
		readwrite = "readwrite"
	}

	// Obtiene una transacción de escritura
	d.transaction = d.db.Call("transaction", table_name, readwrite)

	// Obtiene el almacén de objetos
	d.store = d.transaction.Call("objectStore", table_name)

	if !d.store.Truthy() {
		return Err("error could not open store:", table_name, "for action", action)
	}

	return nil
}

// prepareStoreWithTableCheck prepares the store and creates table if needed
func (d *IndexDB) prepareStoreWithTableCheck(action, table_name string, sampleStruct interface{}) (err error) {
	// Create table if it doesn't exist
	if err = d.CreateTableIfNotExists(table_name, sampleStruct); err != nil {
		return err
	}

	return d.prepareStore(action, table_name)
}
