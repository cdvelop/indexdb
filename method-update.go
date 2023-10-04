package indexdb

func (d indexDB) UpdateObjectsInDB(table_name string, data ...map[string]string) ([]map[string]string, error) {

	if err := d.checkTableStatus("update", table_name); err != nil {
		return nil, err
	}

	return nil, nil
}
