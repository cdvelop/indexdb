package indexdb

import "syscall/js"

func (i *indexDB) createFormData(formDataMap map[string]interface{}) js.Value {
	// Crear un nuevo objeto FormData
	formData := js.Global().Get("FormData").New()

	// Iterar a través del mapa y agregar cada campo al FormData
	for key, value := range formDataMap {
		switch v := value.(type) {
		case string:
			formData.Call("append", key, v)
		case []string:
			for _, item := range v {
				formData.Call("append", key, item)
			}
		default:
			i.Log("Campo no admitido en el FormData:", key)
		}
	}

	return formData
}
