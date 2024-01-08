package indexdb

func (d *indexDB) prepareDataIN(item interface{}, convert_data_to_any bool) {

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

	if convert_data_to_any {
		d.convertStringDataToAny()
	}
}

func (d *indexDB) convertStringDataToAny() {

	for _, data := range d.data_in_str {
		newData := make(map[string]interface{})
		for k, v := range data {
			newData[k] = v
		}
		d.data_in_any = append(d.data_in_any, newData)
	}

}
