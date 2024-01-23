package indexdb

import (
	"syscall/js"

	"github.com/cdvelop/model"
)

func (d *indexDB) ReadStringDataInDBold(r *model.ReadParams) (out []map[string]string, err string) {
	const e = "ReadStringDataInDB error "

	d.Log("info COMIENZO LECTURA")

	// d.readParams = r

	// Abre un cursor para iterar sobre los objetos en el almacén
	d.err = d.readPrepareCursor(r)
	if d.err != "" {
		return nil, e + d.err
	}

	// Define las funciones resolve y reject
	resolve := js.FuncOf(func(e js.Value, args []js.Value) interface{} {
		// manejar resolve
		d.Log("resolve nada:", args[0])
		return nil
	})
	reject := js.FuncOf(func(e js.Value, args []js.Value) interface{} {
		// manejar reject
		d.Log("reject:", args[0])
		return nil
	})

	// Llama a la función pasando resolve y reject
	promise := d.onSuccess(resolve, reject).Invoke()

	// Usa la promesa
	promise.Call("then", js.FuncOf(func(e js.Value, args []js.Value) interface{} {
		// manejar resultado

		d.Log("info then datos", args[0])

		return nil

	}))

	d.Log("info FIN LECTURA")

	return
}

func (d *indexDB) onSuccess(resolve, reject js.Func) js.Func {
	return js.FuncOf(func(e js.Value, args []js.Value) interface{} {

		// Manejador de esta Promesa.
		promiseFunc := js.FuncOf(func(promisee js.Value, promiseArgs []js.Value) interface{} {
			resolve := promiseArgs[0]
			reject := promiseArgs[1]

			// Ejecuta esto de forma asincrónica.
			go func() {
				var itemsOut []interface{}
				// Maneja los resultados del cursor
				d.cursor.Set("onsuccess", js.FuncOf(func(e js.Value, p []js.Value) interface{} {
					item := p[0].Get("target").Get("result")

					if item.Truthy() {
						item := item.Get("value")

						// for _, where := range d.readParams.WHERE {
						// 	if strings.Contains(item.Get(where).String(), d.readParams.SEARCH_ARGUMENT) == 0 {
						// 		item.Call("continue")
						// 	}
						// }

						itemsOut = append(itemsOut, item)

						// Mueve el cursor al siguiente objeto en el almacén
						item.Call("continue")
					} else {
						// El cursor ha llegado al final de los objetos en el almacén
						d.Log("Fin de los datos")
						// Aquí podrías resolver la promesa con los datos acumulados
						resolve.Invoke(itemsOut)
					}

					return nil
				}))

				// Maneja errores
				d.cursor.Set("onerror", js.FuncOf(func(e js.Value, p []js.Value) interface{} {
					// Aquí podrías rechazar la promesa con el error
					reject.Invoke("Error en el cursor")
					return nil
				}))

			}()

			// El controlador no devuelve nada directamente.
			return nil
		})

		// Crea la Promesa y regresa.
		return js.Global().Get("Promise").New(promiseFunc)
	})
}
