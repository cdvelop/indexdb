package indexdb

import (
	"strings"
	"syscall/js"

	"github.com/cdvelop/model"
)

func (d *indexDB) ReadStringDataAsyncInDB(r model.ReadDBParams, callback func(result []map[string]string, err string)) {

	d.readData(r, true, func(data []map[string]string, err string) {
		callback(data, err)
	}, nil)

}

func (d *indexDB) ReadAnyDataAsyncInDB(r model.ReadDBParams, callback func(result []map[string]interface{}, err string)) {

	d.readData(r, false, nil, func(data []map[string]interface{}, err string) {
		callback(data, err)
	})
}

func (d *indexDB) readData(r model.ReadDBParams, string_return bool, outString func(result []map[string]string, err string), outAny func(result []map[string]interface{}, err string)) {

	if err := d.checkTableStatus("read", r.FROM_TABLE); err != "" {
		if string_return {
			outString(nil, err)
		} else {
			outAny(nil, err)
		}
		return
	}

	var (
		results_string = []map[string]string{}
		results_any    = []map[string]interface{}{}
		cursorRequest  js.Value
		sort_order     = "next"
	)

	if r.SORT_DESC {
		sort_order = "prev"
	}

	// transaction := d.db.Call("transaction", r.FROM_TABLE, "readonly")

	// store := transaction.Call("objectStore", r.FROM_TABLE)

	// Obtener el almacén
	store, err := d.getStore("read", r.FROM_TABLE)
	if err != "" {
		if string_return {
			outString(nil, err)
		} else {
			outAny(nil, err)
		}
		return
	}

	switch {

	case r.ID != "":

		field_name := model.PREFIX_ID_NAME + r.FROM_TABLE

		if err := fieldIndexOK(r.FROM_TABLE, field_name, store); err != "" {
			if string_return {
				outString(nil, err)
			} else {
				outAny(nil, err)
			}
			return
		}

		rangeObj := js.Global().Get("IDBKeyRange").Call("only", r.ID)
		index := store.Call("index", field_name)
		cursorRequest = index.Call("openCursor", rangeObj)

	case r.ORDER_BY != "":

		if err := fieldIndexOK(r.FROM_TABLE, r.ORDER_BY, store); err != "" {
			if string_return {
				outString(nil, err)
			} else {
				outAny(nil, err)
			}
			return
		}

		index := store.Call("index", r.ORDER_BY)
		// El valor nil como clave inicial significa que el cursor comenzará desde el primer registro en orden descendente y luego avanzará hacia registros posteriores en ese orden
		cursorRequest = index.Call("openCursor", nil, sort_order)
	default:
		// normal
		cursorRequest = store.Call("openCursor")

	}

	cursorRequest.Call("addEventListener", "success", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		cursor := p[0].Get("target").Get("result")
		if cursor.Truthy() {

			data := cursor.Get("value")

			for _, where := range r.WHERE {
				if !strings.Contains(data.Get(where).String(), r.SEARCH_ARGUMENT) {
					cursor.Call("continue")
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
					data_out_string["url"] = CreateBlobURL(value_js)

					if !string_return {
						data_out_any["url"] = CreateBlobURL(value_js)
						data_out_any[key] = value_js
					}

				} else {

					if string_return {
						data_out_string[key] = value_js.String()
					} else {
						data_out_any[key] = value_js
					}

				}

			}

			if string_return {
				results_string = append(results_string, data_out_string)

			} else {
				results_any = append(results_any, data_out_any)
			}

			cursor.Call("continue")
		} else {

			if string_return {
				outString(results_string, "") // log("Fin de los datos.")
			} else {
				outAny(results_any, "") // log("Fin de los datos.")
			}

		}
		return nil
	}))

}

func (d *indexDB) ReadObjectsInDB(FROM_TABLE string, data ...map[string]string) (result []map[string]string, err string) {
	return nil, "error ReadObjectsInDB no implementado en indexDB"
}
