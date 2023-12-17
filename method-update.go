package indexdb

import "github.com/cdvelop/model"

func (d *indexDB) UpdateObjectsInDB(table_name string, all_data ...map[string]string) (err string) {

	// Obtener el almacén
	store, err := d.getStore("update", table_name)
	if err != "" {
		return err
	}

	d.prepareDataIN(all_data)

	// Iterar sobre los datos a actualizar
	for _, obj := range d.data_in_any {

		// Obtener el ID del objeto
		id, ok := obj[model.PREFIX_ID_NAME+table_name].(string)
		if !ok {
			return "objeto invalido sin ID para actualizar "
		}

		// Guardar los cambios
		result := store.Call("put", obj)
		if result.IsNull() {
			return "error actualizando objeto " + id + " en la db"
		}

	}

	return ""
}
