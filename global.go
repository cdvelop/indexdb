package indexdb

import "syscall/js"

func log(message ...any) {
	js.Global().Get("console").Call("log", message...)
}
