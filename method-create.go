package indexdb

import (
	"github.com/cdvelop/model"
)

func (d *indexDB) CreateObjectsInDB(table_name string, backup_required bool, data ...map[string]string) error {

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

		// chequear si tiene llave primaria el objeto
		if _, id_exist := items[model.PREFIX_ID_NAME+table_name]; id_exist {

			new := make(map[string]interface{})

			for k, v := range items {
				new[k] = v
			}

			if backup_required { // necesita respaldo en servidor
				new["backup"] = false //estado backup = no respaldado
			}

			// Inserta cada elemento en el almacén de objetos
			store.Call("add", new)
		} else {
			d.Log("error data sin id en tabla:", table_name, items)
		}

	}

	return nil
}
