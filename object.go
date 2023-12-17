package indexdb

import "github.com/cdvelop/model"

func (d *indexDB) GetObjectByTableName(table_name string) (o *model.Object, err string) {
	// d.Log("total objetos:", len(d.objects))
	for _, o := range d.GetAllObjectsFromMainHandler() {
		// d.Log("BUSCANDO OBJETO:", o.ObjectName)
		if o.Table == table_name {
			return o, ""
		}
	}

	return nil, "error objeto: " + table_name + ",no encontrado"
}
