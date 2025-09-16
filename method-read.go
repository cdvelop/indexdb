package indexdb

import (
	"syscall/js"

	. "github.com/cdvelop/tinystring"
)

type ReadParams struct {
	FROM_TABLE string
	SORT_DESC  bool
	ID         string
	ORDER_BY   string
	WHERE      []interface{}
	RETURN_ANY bool
}

type ReadResults struct {
	Results []interface{}
	Error   error
}

func (d *IndexDB) Read(p *ReadParams, callback func(r *ReadResults, err error)) {

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
				// d.Logger("key", key, "value", value_js)
				if key == "blob" {
					// url := CreateBlobURL(value_js)
					// d.Logger("BLOB FOUND:", value_js)
					// d.Logger("URL:", url)

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
			// d.Logger("Fin de los datos.")
			callback(result, nil)
		}
		return nil
	}))

}

func (d *IndexDB) readPrepareCursor(r *ReadParams) (err error) {
	const e = "readPrepareCursor error"

	sort_order := "next"
	if r.SORT_DESC {
		sort_order = "prev"
	}

	// Obtener el almacén
	if d.err = d.prepareStore("read", r.FROM_TABLE); d.err != nil {
		return Errf("%s %v", e, d.err)
	}

	switch {

	case r.ID != "":

		field_name := "id_" + r.FROM_TABLE

		if d.err = d.fieldIndexOK(r.FROM_TABLE, field_name); d.err != nil {
			return Errf("%s %v", e, d.err)
		}

		rangeObj := js.Global().Get("IDBKeyRange").Call("only", r.ID)
		index := d.store.Call("index", field_name)
		d.cursor = index.Call("openCursor", rangeObj)

	case r.ORDER_BY != "":

		if d.err = d.fieldIndexOK(r.FROM_TABLE, r.ORDER_BY); d.err != nil {
			return Errf("%s %v", e, d.err)
		}

		index := d.store.Call("index", r.ORDER_BY)
		// El valor nil como clave inicial significa que el cursor comenzará desde el primer registro en orden descendente y luego avanzará hacia registros posteriores en ese orden
		d.cursor = index.Call("openCursor", nil, sort_order)
	default:
		// normal
		d.cursor = d.store.Call("openCursor")
	}

	return nil
}
