package indexdb

import (
	"syscall/js"

	. "github.com/cdvelop/tinystring"
)

func (d *indexDB) Read(p *ReadParams, callback func(r *ReadResults, err error)) {

	var result = &ReadResults{
		Results: []interface{}{},
	}

	d.err = d.readPrepareCursor(p)
	if d.err != nil {
		callback(nil, d.err)
		return
	}

	d.cursor.Call("addEventListener", "success", js.FuncOf(func(this js.Value, v []js.Value) interface{} {
		item := v[0].Get("target").Get("result")
		if item.Truthy() {

			data := item.Get("value")

			for _, wheres := range p.WHERE {
				if w, ok := wheres.(map[string]string); ok {
					for key, search := range w {
						if !Contains(data.Get(key).String(), search) {
							item.Call("continue")
							return nil
						}
					}
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
				result.Results = append(result.Results, data_out_any)
			} else {
				result.Results = append(result.Results, data_out_string)
			}

			item.Call("continue")
		} else {
			// d.Log("Fin de los datos.")
			callback(result, nil)
		}
		return nil
	}))

}

func (d *indexDB) ReadSync(p *ReadParams, data ...interface{}) (result []interface{}, err error) {
	return nil, Err("error ReadSyncDataDB not implemented in indexDB")
}
