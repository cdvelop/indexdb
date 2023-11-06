package indexdb

import "github.com/cdvelop/model"

func (d indexDB) UpdateObjectsInDB(table_name string, data ...map[string]string) error {

	if err := d.checkTableStatus("update", table_name); err != nil {
		return err
	}

	return model.Error("ERROR UpdateObjectsInDB NO IMPLEMENTADO EN indexDB")
}
