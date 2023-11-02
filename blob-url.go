package indexdb

import "syscall/js"

// CreateBlobURL crea una URL Blob a partir de un blob.
func CreateBlobURL(blob any) string {

	jsBlob := js.ValueOf(blob)

	return js.Global().Get("URL").Call("createObjectURL", jsBlob).String()
}
