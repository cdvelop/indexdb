package indexdb

import "github.com/cdvelop/model"

func (d *indexDB) GetObjectByTableName(table_name string) (o *model.Object, err string) {
	// d.Log("total objetos:", len(d.objects))
	for _, o := range d.GetAllObjects() {
		// d.Log("BUSCANDO OBJETO:", o.ObjectName)
		if o.Table == table_name {
			return o, ""
		}
	}

	return nil, "error objeto: " + table_name + ",no encontrado"
}

func (d *indexDB) getObjectsDB() []*model.Object {

	if len(d.objects_db) == 0 {
		for _, o := range d.GetAllObjects() {
			if !o.NoAddObjectInDB {
				if len(o.Fields) != 0 {
					d.objects_db = append(d.objects_db, o)
				}
			}
		}
	}

	return d.objects_db
}
