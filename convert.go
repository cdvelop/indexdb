package indexdb

func dataConvert(item interface{}) (all_data []map[string]interface{}) {

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
