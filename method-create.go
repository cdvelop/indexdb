package indexdb

import (
	"github.com/cdvelop/model"
)

// items support: []map[string]string or map[string]string
func (d *indexDB) CreateObjectsInDB(table_name string, backup_required bool, items any) error {

	store, err := d.getStore("create", table_name)
	if err != nil {
		return err
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
			data["backup"] = "false" //estado backup = no respaldado
		}

		// Inserta cada elemento en el almacén de objetos
		result := store.Call("add", data)
		if result.IsNull() {
			return model.Error("error al crear elemento en la db tabla:", table_name, "id:", id)
		}

		// retornamos url temporal para acceder al archivo
		if blob, exist := data["blob"]; exist {
			data["url"] = CreateBlobURL(blob)
		}
	}

	return nil
}
