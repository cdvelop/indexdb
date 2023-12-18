package indexdb

import (
	"syscall/js"

	"github.com/cdvelop/model"
)

func (d *indexDB) CreateTablesInDB(objects []*model.Object, result func(err string)) {

	d.result = result

	// Accede a la base de datos
	db := js.Global().Get("indexedDB").Call("open", d.db_name)

	// Agrega eventos
	db.Call("addEventListener", "error", js.FuncOf(d.showDbError))
	db.Call("addEventListener", "success", js.FuncOf(d.openExistingDB))
	db.Call("addEventListener", "upgradeneeded", js.FuncOf(d.upgradeneeded))

}

func (d *indexDB) upgradeneeded(this js.Value, p []js.Value) interface{} {

	err := d.open(&p[0], "NEW CREATE")
	if err != "" {
		d.result(err)
		return nil
	}

	for _, o := range d.MainHandlerGetAllObjects() {
		if !o.NoAddObjectInDB {
			// d.Log("**** CREANDO TABLA: ", o.Table, "INDEX DB")
			if len(o.Fields) != 0 {

				pk_name := o.PrimaryKeyName()

				newTable := d.db.Call("createObjectStore", o.Table, map[string]interface{}{"keyPath": pk_name})

				for _, f := range o.Fields {

					if !f.NotRequiredInDB {
						// Crear un índices para búsqueda campos principales
						newTable.Call("createIndex", f.Name, f.Name, map[string]interface{}{"unique": f.Unique})
					}
				}
			}
		}
	}

	return nil
}

func (d indexDB) showDbError(this js.Value, p []js.Value) interface{} {
	d.Log("indexDB Error", p[0])
	return nil
}

func (d *indexDB) openExistingDB(this js.Value, p []js.Value) interface{} {
	err := d.open(&p[0], "OPEN")
	if err != "" {
		d.result(err)
		return nil
	}

	d.result("")

	return nil
}

func (d *indexDB) open(p *js.Value, message string) (err string) {

	d.db = p.Get("target").Get("result")

	if !d.db.Truthy() {
		return "open indexdb error no se logro establecer conexión " + d.db_name
	}

	// d.Log("***", message, "IndexDB Connection:", d.db_name, " ***")
	// DB : localdb Established, Engine: indexedDB
	return ""
}
