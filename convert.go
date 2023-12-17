package indexdb

func (d *indexDB) prepareDataIN(item interface{}) {

	d.data_in_any = make([]map[string]interface{}, 0)
	d.data_in_str = make([]map[string]string, 0)

	switch item := item.(type) {
	case map[string]string:
		d.data_in_str = append(d.data_in_str, item)

	case []map[string]string:
		d.data_in_str = item

	case []map[string]interface{}:
		d.data_in_any = item

	case map[string]interface{}:
		d.data_in_any = append(d.data_in_any, item)
	}

	d.convertStringDataIN()

}

func (d *indexDB) convertStringDataIN() {

	for _, data := range d.data_in_str {
		newData := make(map[string]interface{})
		for k, v := range data {
			newData[k] = v
		}
		d.data_in_any = append(d.data_in_any, newData)
	}

}

func (d *indexDB) DataConvertToAnyOLD(item interface{}) (all_data []map[string]interface{}) {

	convert := func(data map[string]string) map[string]interface{} {
		newData := make(map[string]interface{})
		for k, v := range data {
			newData[k] = v
		}
		return newData
	}

	switch item := item.(type) {
	case map[string]string:
		all_data = append(all_data, convert(item))

	case []map[string]string:
		for _, data := range item {
			all_data = append(all_data, convert(data))
		}

	case []map[string]interface{}:
		return item

	case map[string]interface{}:
		all_data = append(all_data, item)

	}

	return
}
