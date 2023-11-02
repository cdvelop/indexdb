package indexdb

import (
	"github.com/cdvelop/model"
)

// items support: []map[string]string or map[string]string
func (d *indexDB) CreateObjectsInDB(table_name string, backup_required bool, items any) error {

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

	for _, data := range dataConvert(items) {

		pk_field := model.PREFIX_ID_NAME + table_name

		id, id_exist := data[pk_field]
		if !id_exist {

			if !backup_required { // si no requiere backup es un objeto sin id del servidor retornamos error
				return model.Error("error data proveniente del servidor sin id en tabla:", table_name, data)
			}

			//agregar id al objeto si este no existe
			id = d.GetNewID() //id nuevo

			// d.Log("NUEVO ID GENERADO:", id)

			// date, _ := unixid.UnixNanoToStringDate(id.(string))
			// d.Log("SU FECHA ES:", date)

			data[pk_field] = id
		}

		if backup_required { // necesita respaldo en servidor
			data["backup"] = false //estado backup = no respaldado
		}

		// Inserta cada elemento en el almacén de objetos
		store.Call("add", data)

		// retornamos url temporal para acceder al archivo
		if blob, exist := data["blob"]; exist {
			data["url"] = CreateBlobURL(blob)
		}
	}

	return nil
}
