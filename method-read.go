package indexdb

import (
	"syscall/js"

	"github.com/cdvelop/model"
	"github.com/cdvelop/strings"
)

func (d *indexDB) ReadAsyncDataDB(p model.ReadParams, callback func(r model.ReadResult)) {

	var result = model.ReadResult{
		DataString: []map[string]string{},
		DataAny:    []map[string]any{},
		Error:      "",
	}

	cursor, err := d.readPrepareCursor(p)
	if err != "" {
		result.Error = err
		callback(result)
		return
	}

	cursor.Call("addEventListener", "success", js.FuncOf(func(this js.Value, v []js.Value) interface{} {
		item := v[0].Get("target").Get("result")
		if item.Truthy() {

			data := item.Get("value")

			for _, where := range p.WHERE {
				if strings.Contains(data.Get(where).String(), p.SEARCH_ARGUMENT) == 0 {
					item.Call("continue")
					return nil
				}
			}

			data_out_any := make(map[string]interface{})
			data_out_string := map[string]string{}

			// log("READ DATA:", data)

			keys := js.Global().Get("Object").Call("keys", data)

			for i := 0; i < keys.Length(); i++ {
				key := keys.Index(i).String()
				value_js := data.Get(key)
				// d.Log("key", key, "value", value_js)
				if key == "blob" {
					// url := CreateBlobURL(value_js)
					// d.Log("BLOB FOUND:", value_js)
					// d.Log("URL:", url)

					if p.RETURN_ANY {
						data_out_any["url"] = CreateBlobURL(value_js)
						data_out_any[key] = value_js
					} else {
						data_out_string["url"] = CreateBlobURL(value_js)
					}

				} else {

					if p.RETURN_ANY {
						data_out_any[key] = value_js
					} else {
						data_out_string[key] = value_js.String()
					}
				}
			}

			if p.RETURN_ANY {
				result.DataAny = append(result.DataAny, data_out_any)
			} else {
				result.DataString = append(result.DataString, data_out_string)
			}

			item.Call("continue")
		} else {
			// d.Log("Fin de los datos.")
			callback(result)
		}
		return nil
	}))

}
