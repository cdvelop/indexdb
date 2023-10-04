package indexdb

import (
	"github.com/cdvelop/model"
)

func (d *indexDB) CreateObjectsInDB(table_name string, data ...map[string]string) error {

	if err := d.checkTableStatus("create", table_name); err != nil {
		return err
	}

	// Obtiene una transacción de escritura usando d.res
	transaction := d.db.Call("transaction", table_name, "readwrite")

	// Obtiene el almacén de objetos
	store := transaction.Call("objectStore", table_name)

	if !store.Truthy() {
		return model.Error("error no se logro abrir el almacén:", table_name, "en indexdb")
	}

	for _, items := range data {
		new := make(map[string]interface{})

		for k, v := range items {
			new[k] = v
		}

		// Inserta cada elemento en el almacén de objetos
		store.Call("add", new)

	}

	return nil
}
