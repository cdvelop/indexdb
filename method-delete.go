package indexdb

import "github.com/cdvelop/model"

func (d indexDB) DeleteObjectsInDB(table_name string, all_data ...map[string]string) (err string) {

	store, err := d.getStore("delete", table_name)
	if err != "" {
		return err
	}

	for _, data := range all_data {
		if id, exist := data[model.PREFIX_ID_NAME+table_name]; exist {
			// elimina cada elemento en el almacén de objetos
			result := store.Call("delete", id)

			if result.IsNull() {
				return "error al eliminar en la tabla: " + table_name
			}

		} else {
			return "error en datos enviados a eliminar, id no encontrado tabla: " + table_name
		}

	}

	return ""
}
