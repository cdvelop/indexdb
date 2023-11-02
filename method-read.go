package indexdb

import (
	"strings"
	"syscall/js"

	"github.com/cdvelop/model"
)

func (d indexDB) ReadDataAsyncInDB(from_tables string, params []map[string]string, callback func([]map[string]string, error)) {

	if err := d.checkTableStatus("read", from_tables); err != nil {
		callback(nil, err)
		return
	}

	var (
		results       = []map[string]string{}
		order_by      string
		cursorRequest js.Value
		sort_order    = "next"
		where         string
		args          string
	)

	transaction := d.db.Call("transaction", from_tables, "readonly")
	store := transaction.Call("objectStore", from_tables)

	for _, items := range params {
		for key, value := range items {

			switch {
			case key == "ORDER_BY":
				order_by = value

			case key == "SORT" && value == "DESC":
				sort_order = "prev"

			case key == "WHERE":
				where = value

			case key == "ARGS":
				args = value
			}
		}
	}

	switch {

	// case where != "" && args != "":
	// if err := fieldIndexOK(from_tables, where, store); err != nil {
	// 	callback(nil, err)
	// 	return
	// }
	// index := store.Call("index", where)

	// rangeObj := js.Global().Get("IDBKeyRange").Call("only", args)

	// cursorRequest = index.Call("openCursor", rangeObj)

	case order_by != "":

		if err := fieldIndexOK(from_tables, order_by, store); err != nil {
			callback(nil, err)
			return
		}

		index := store.Call("index", order_by)
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

			if where != "" && args != "" {
				if !strings.Contains(data.Get(where).String(), args) {
					cursor.Call("continue")
					return nil
				}
			}

			dataMap := make(map[string]string)

			// log("READ DATA:", data)

			keys := js.Global().Get("Object").Call("keys", data)

			for i := 0; i < keys.Length(); i++ {
				key := keys.Index(i).String()
				value_js := data.Get(key)

				// log("key", key, "value", value)

				if key == "blob" {
					// url := CreateBlobURL(value_js)
					// d.Log("BLOB FOUND:", value_js)
					// d.Log("URL:", url)
					dataMap["url"] = CreateBlobURL(value_js)
				} else {
					value := value_js.String()

					dataMap[key] = value
				}

			}

			results = append(results, dataMap)

			cursor.Call("continue")
		} else {
			callback(results, nil) // log("Fin de los datos.")
		}
		return nil
	}))
}

func fieldIndexOK(table, field_name string, store js.Value) error {

	// Verificar si el índice existe en la tabla.
	indexNames := store.Get("indexNames")

	indexSet := make(map[string]bool)
	for i := 0; i < indexNames.Length(); i++ {
		indexSet[indexNames.Index(i).String()] = true
	}

	// Verificar si el índice existe en la tabla.
	if !indexSet[field_name] {
		return model.Error("El índice:", field_name, "no existe en la tabla:", table)
	}

	return nil
}

func (d *indexDB) ReadObjectsInDB(from_tables string, data ...map[string]string) ([]map[string]string, error) {
	return nil, model.Error("error ReadObjectsInDB no implementado en indexDB")
}
