package indexdb

import "syscall/js"

// variables globales con algunas de las variables globales más comunes para evitar asignaciones adicionales
var (
	promiseConstructor = js.Global().Get("Promise")
	errorJS            = js.Global().Get("Error")
)

func Await(cb func() (js.Value, error)) js.Value {

	promiseFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		//Esta sería la función con la que quiero tratar
		go func() {
			result, err := cb()
			if err != nil {
				reject.Invoke(PromiseError(err))
				return
			}

			resolve.Invoke(result)
		}()

		return nil
	})
	defer promiseFunc.Release()

	return promiseConstructor.New(promiseFunc)
}

// https://go-review.googlesource.com/c/go/+/402455/3/src/syscall/js/promise.go
// Función prometedora que satisface los requisitos de MakePromise.
type PromiseAbleFunc = func(js.Value, []js.Value) (interface{}, error)

// MakePromise hace una promesa de una función que toma una serie de argumentos.
func PromiseOf(fn PromiseAbleFunc) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Manejador de esta Promesa.
		handler := js.FuncOf(func(promiseThis js.Value, promiseArgs []js.Value) interface{} {
			resolve := promiseArgs[0]
			reject := promiseArgs[1]

			// Ejecuta esto de forma asincrónica.
			go func() {
				// Recuperar rechazando.
				defer func() {
					if err := recover(); err != nil {
						reject.Invoke(PromiseError(err))
					}
				}()

				// Ejecute la función PromiseAble con this y args originales.
				res, err := fn(this, args)
				if err != nil {
					reject.Invoke(PromiseError(err))
					return
				}
				resolve.Invoke(res)
			}()

			// El controlador no devuelve nada directamente.
			return nil
		})

		// Crea la Promesa y regresa.
		return promiseConstructor.New(handler)
	})
}

// PromiseError makes sure to return some error that Invoke will understand.
func PromiseError(e interface{}) (err js.Value) {
	switch x := e.(type) {
	case string:
		err = errorJS.New(x)
	case error:
		err = errorJS.New(x.Error())
	default:
		err = errorJS.New("unknown panic")
	}
	return err
}
