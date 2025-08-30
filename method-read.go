package indexdb

import (
	"syscall/js"

	"github.com/cdvelop/model"
	"github.com/cdvelop/tinystring"
)

func (d *indexDB) ReadAsyncDataDB(p *model.ReadParams, callback func(r *model.ReadResults, err string)) {

	var result = &model.ReadResults{
		ResultsString: []map[string]string{},
		ResultsAny:    []map[string]any{},
	}

	d.err = d.readPrepareCursor(p)
	if d.err != "" {
		callback(nil, d.err)
		return
	}

	d.cursor.Call("addEventListener", "success", js.FuncOf(func(this js.Value, v []js.Value) interface{} {
		item := v[0].Get("target").Get("result")
		if item.Truthy() {

			data := item.Get("value")

			for _, wheres := range p.WHERE {
				for key, search := range wheres {
					if !tinystring.Contains(data.Get(key).String(), search) {
						item.Call("continue")
						return nil
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
				result.ResultsAny = append(result.ResultsAny, data_out_any)
			} else {
				result.ResultsString = append(result.ResultsString, data_out_string)
			}

			item.Call("continue")
		} else {
			// d.Log("Fin de los datos.")
			callback(result, "")
		}
		return nil
	}))

}

func (d *indexDB) ReadSyncDataDB(p *model.ReadParams, data ...map[string]string) (result []map[string]string, err string) {
	return nil, "error ReadSyncDataDB no implementado en indexDB"
}
