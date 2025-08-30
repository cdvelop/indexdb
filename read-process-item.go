package indexdb

import (
	"syscall/js"
)

func (d indexDB) processReadItem(p *ReadParams, data js.Value, r *ReadResults) {

	out_any := make(map[string]interface{})
	out_string := map[string]string{}

	// log("READ DATA:", data)

	keys := js.Global().Get("Object").Call("keys", data)

	for i := 0; i < keys.Length(); i++ {
		key := keys.Index(i).String()
		value_js := data.Get(key)
		if key == "blob" {

			if p.RETURN_ANY {
				out_any["url"] = CreateBlobURL(value_js)
				out_any[key] = value_js
			} else {
				out_string["url"] = CreateBlobURL(value_js)
			}

		} else {

			if p.RETURN_ANY {
				out_any[key] = value_js
			} else { //STRING RETURN
				out_string[key] = value_js.String()
			}
		}
	}

	if p.RETURN_ANY {
		r.ResultsAny = append(r.ResultsAny, out_any)
	} else { //STRING RETURN
		r.ResultsString = append(r.ResultsString, out_string)
	}
}
