package indexdb

import "github.com/cdvelop/model"

func (d *indexDB) UpdateObjectsInDB(table_name string, all_data ...map[string]string) error {

	// Obtener el almacén
	store, err := d.getStore("update", table_name)
	if err != nil {
		return err
	}

	// Iterar sobre los datos a actualizar
	for _, obj := range DataConvertToAny(all_data) {

		// Obtener el ID del objeto
		id, ok := obj[model.PREFIX_ID_NAME+table_name].(string)
		if !ok {
			return model.Error("objeto invalido sin ID para actualizar", obj)
		}

		// Guardar los cambios
		result := store.Call("put", obj)
		if result.IsNull() {
			return model.Error("error actualizando objeto ", id, " en la db")
		}

	}

	return nil
}
