package indexdb

import "github.com/cdvelop/model"

func (d *indexDB) GetObjectByTableName(table_name string) (*model.Object, error) {
	// d.Log("total objetos:", len(d.objects))
	for _, o := range d.objects {
		// d.Log("BUSCANDO OBJETO:", o.Name)
		if o.Table == table_name {
			return o, nil
		}
	}

	return nil, model.Error("error objeto:", table_name, ",no encontrado")
}