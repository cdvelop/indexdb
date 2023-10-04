package indexdb

func (d indexDB) DeleteObjectsInDB(table_name string, data ...map[string]string) ([]map[string]string, error) {

	if err := d.checkTableStatus("delete", table_name); err != nil {
		return nil, err
	}

	return nil, nil
}
