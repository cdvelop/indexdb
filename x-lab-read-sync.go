package indexdb

import (
	"syscall/js"
)

func (d *IndexDB) ReadStringDataInDB(r *ReadParams) (out []interface{}, err error) {
	const this = "ReadStringDataInDB"

	d.logger("info COMIENZO LECTURA")

	chanResult := make(chan ReadResults)

	go d.readDataTwo(r, chanResult)

	data := <-chanResult

	d.logger("info FIN LECTURA")

	if data.Error != nil {
		return nil, data.Error
	}

	return data.Results, nil
}

func (d *IndexDB) readDataTwo(r *ReadParams, chanResult chan ReadResults) {

	var result = ReadResults{
		Results: []interface{}{},
	}

	if d.err = d.readPrepareCursor(r); d.err != nil {
		result.Error = d.err
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
					result.Results = append(result.Results, data_out_any)
				} else {
					result.Results = append(result.Results, data_out_string)
				}

				cursor.Call("continue")

			} else {
				d.logger("info Fin de los datos.")
				then <- true
			}
			return nil
		}))
	}()

	<-then

	chanResult <- result

}
