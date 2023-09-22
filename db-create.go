package indexdb

import (
	"syscall/js"

	"github.com/cdvelop/model"
)

func (i *indexDB) initDataBase() {
	log("*** Conexión DB: localdb Establecida, Engine: indexDB ***")

	// Accede a la base de datos
	i.db = js.Global().Get("indexedDB").Call("open", "localdb")

	// Agrega eventos
	i.db.Call("addEventListener", "error", js.FuncOf(i.showDbError))
	i.db.Call("addEventListener", "success", js.FuncOf(i.openExistingDB))
	i.db.Call("addEventListener", "upgradeneeded", js.FuncOf(i.createNewDB))

}

func (d indexDB) CreateTablesInDB(...*model.Object) error {

	return nil
}

func (d indexDB) createNewDB(this js.Value, p []js.Value) interface{} {
	log("createNewDB", p[0])

	// Accede al resultado del evento
	res := p[0].Get("target").Get("result")

	for _, o := range d.objects {

		// Crea la tabla
		newTable := res.Call("createObjectStore", o.Name, map[string]interface{}{"keyPath": o.PrimaryKeyName()})

		// Crear un índices para búsqueda según campo
		for _, f := range o.Fields {

			if !f.NotRequiredInDB {
				newTable.Call("createIndex", f.Name, f.Name, map[string]interface{}{"unique": f.Unique})
			}
		}

	}

	return nil
}

func (indexDB) showDbError(this js.Value, p []js.Value) interface{} {
	log("indexDB Error", p[0])
	return nil
}

func (indexDB) openExistingDB(this js.Value, p []js.Value) interface{} {
	log("openExistingDB", p[0])
	return nil
}
