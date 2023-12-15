package indexdb

import (
	"syscall/js"

	"github.com/cdvelop/model"
)

func (d *indexDB) readPrepareCursor(r model.ReadParams) (cursor js.Value, err string) {
	const this = "readPrepareCursor error "

	if err = d.checkTableStatus("read", r.FROM_TABLE); err != "" {
		return js.Value{}, this + err
	}

	sort_order := "next"
	if r.SORT_DESC {
		sort_order = "prev"
	}

	// Obtener el almacén
	store, err := d.getStore("read", r.FROM_TABLE)
	if err != "" {
		err = this + err
		return
	}

	switch {

	case r.ID != "":

		field_name := model.PREFIX_ID_NAME + r.FROM_TABLE

		if err := fieldIndexOK(r.FROM_TABLE, field_name, store); err != "" {
			return js.Value{}, this + err
		}

		rangeObj := js.Global().Get("IDBKeyRange").Call("only", r.ID)
		index := store.Call("index", field_name)
		cursor = index.Call("openCursor", rangeObj)

	case r.ORDER_BY != "":

		if err := fieldIndexOK(r.FROM_TABLE, r.ORDER_BY, store); err != "" {
			return js.Value{}, this + err
		}

		index := store.Call("index", r.ORDER_BY)
		// El valor nil como clave inicial significa que el cursor comenzará desde el primer registro en orden descendente y luego avanzará hacia registros posteriores en ese orden
		cursor = index.Call("openCursor", nil, sort_order)
	default:
		// normal
		cursor = store.Call("openCursor")
	}

	return
}
