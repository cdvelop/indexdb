package indexdb

import (
	"syscall/js"

	"github.com/cdvelop/model"
)

func (i *indexDB) CreateTablesInDB(objects []*model.Object, action model.Subsequently) error {

	i.run = action
	i.objects = objects
	// Accede a la base de datos
	db := js.Global().Get("indexedDB").Call("open", i.db_name)

	// Agrega eventos
	db.Call("addEventListener", "error", js.FuncOf(i.showDbError))
	db.Call("addEventListener", "success", js.FuncOf(i.openExistingDB))
	db.Call("addEventListener", "upgradeneeded", js.FuncOf(i.upgradeneeded))

	return nil
}

func (d *indexDB) upgradeneeded(this js.Value, p []js.Value) interface{} {

	err := d.open(&p[0], "NEW CREATE")
	if err != nil {
		log(err)
		return nil
	}

	for _, o := range d.objects {

		// if !d.TableExist(o.Table) {

		log(o.Table, "keyPath:", o.PrimaryKeyName())
		// Crea la tabla
		newTable := d.db.Call("createObjectStore", o.Table, map[string]interface{}{"keyPath": o.PrimaryKeyName()})

		// Crear un índices para búsqueda según campo
		for _, f := range o.Fields {

			if !f.NotRequiredInDB {
				newTable.Call("createIndex", f.Name, f.Name, map[string]interface{}{"unique": f.Unique})
			}
		}
		// } else {
		// log("TABLA:", o.Table, "YA EXISTE EN LA DB!!!")
		// }

	}

	return nil
}

func (indexDB) showDbError(this js.Value, p []js.Value) interface{} {
	log("indexDB Error", p[0])
	return nil
}

func (d *indexDB) openExistingDB(this js.Value, p []js.Value) interface{} {
	err := d.open(&p[0], "OPEN")
	if err != nil {
		log(err)
		return nil
	}

	d.run.ActionExecutedLater()

	return nil
}

func (d *indexDB) open(p *js.Value, message string) error {

	d.db = p.Get("target").Get("result")

	if !d.db.Truthy() {
		return model.Error("error no se logro establecer conexión", d.db_name, "indexdb")
	}

	log("***", message, "IndexDB Connection:", d.db_name, "OK ***")

	// DB : localdb Established, Engine: indexedDB
	return nil
}
