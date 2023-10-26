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
		d.Log(err)
		return nil
	}

	for _, o := range d.objects {

		// log("CREANDO TABLA: ", o.Table)

		// if !d.TableExist(o.Table) {

		// log(o.Table, "keyPath:", o.PrimaryKeyName())
		// Crea la tabla
		pk_name := o.PrimaryKeyName()

		newTable := d.db.Call("createObjectStore", o.Table, map[string]interface{}{"keyPath": pk_name})

		// Crear un índices para búsqueda campos principales
		principal_fields, err := o.GetFieldsByNames(o.NamePrincipalFields...)

		if err == nil {
			for _, f := range principal_fields {

				// if !f.NotRequiredInDB && f.Name != pk_name {
				if !f.NotRequiredInDB {
					newTable.Call("createIndex", f.Name, f.Name, map[string]interface{}{"unique": f.Unique})
				}
			}
		}

		// log("TABLA: ", o.Table, "CREADA.....")

		// } else {
		// log("TABLA:", o.Table, "YA EXISTE EN LA DB!!!")
		// }

	}

	return nil
}

func (d indexDB) showDbError(this js.Value, p []js.Value) interface{} {
	d.Log("indexDB Error", p[0])
	return nil
}

func (d *indexDB) openExistingDB(this js.Value, p []js.Value) interface{} {
	err := d.open(&p[0], "OPEN")
	if err != nil {
		d.Log(err)
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

	d.Log("***", message, "IndexDB Connection:", d.db_name, "OK ***")

	// DB : localdb Established, Engine: indexedDB
	return nil
}
