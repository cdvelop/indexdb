package indexdb

import (
	"syscall/js"

	"github.com/cdvelop/model"
)

func (d *indexDB) ReadStringDataInDB(r model.ReadParams) (out []map[string]string, err string) {
	const this = "ReadStringDataInDB "

	d.Log("info COMIENZO LECTURA")

	chanResult := make(chan model.ReadResults)

	d.readDataTwo(r, chanResult)

	data := <-chanResult

	d.Log("info FIN LECTURA")
	d.Log("dataString", data.ResultsString)
	// d.Log("erro", data.Error)

	return
}

func (d *indexDB) readDataTwo(r model.ReadParams, chanResult chan model.ReadResults) {

	var result = model.ReadResults{
		ResultsString: []map[string]string{},
		ResultsAny:    []map[string]any{},
	}

	d.err = d.readPrepareCursor(r)
	if d.err != "" {
		// result.Error = err
		chanResult <- result
		return
	}

	then := make(chan bool)

	go func() {
		d.cursor.Call("addEventListener", "success", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
			cursor := p[0].Get("target").Get("result")
			if cursor.Truthy() {

				data := cursor.Get("value")

				// for _, where := range r.WHERE {
				// if strings.Contains(data.Get(where).String(), r.SEARCH_ARGUMENT) == 0 {
				// 	cursor.Call("continue")
				// 	return nil
				// }
				// }

				data_out_any := make(map[string]interface{})
				data_out_string := map[string]string{}

				// log("READ DATA:", data)

				keys := js.Global().Get("Object").Call("keys", data)

				for i := 0; i < keys.Length(); i++ {
					key := keys.Index(i).String()
					value_js := data.Get(key)

					if key == "blob" {

						if r.RETURN_ANY {
							data_out_any["url"] = CreateBlobURL(value_js)
							data_out_any[key] = value_js
						} else {
							data_out_string["url"] = CreateBlobURL(value_js)
						}

					} else {

						if r.RETURN_ANY {
							data_out_any[key] = value_js
						} else {
							data_out_string[key] = value_js.String()
						}

					}

				}

				if r.RETURN_ANY {
					// result.ResultsAny = append(result.ResultsAny, data_out_any)
				} else {
					// result.ResultsStringing = append(resuResultsStringString, data_out_string)
				}

				cursor.Call("continue")

			} else {
				d.Log("info Fin de los datos.")
				then <- true
			}
			return nil
		}))
	}()

	<-then

	chanResult <- result

}
