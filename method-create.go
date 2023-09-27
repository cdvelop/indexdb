package indexdb

import (
	"fmt"

	"github.com/cdvelop/model"
)

func (d *indexDB) CreateObjectsInDB(table_name string, data ...map[string]string) error {

	fmt.Println("EJECUTANDO CreateObjectsInDB", table_name)

	if !d.db.Truthy() {
		return model.Error("ERROR DB ESTA NULA..", table_name)
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
