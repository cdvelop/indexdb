package indexdb

import (
	"syscall/js"
)

func (d *indexDB) readPrepareCursor(r *ReadParams) (err string) {
	const e = "readPrepareCursor error "

	if d.err = d.checkTableStatus("read", r.FROM_TABLE); d.err != "" {
		return e + d.err
	}

	sort_order := "next"
	if r.SORT_DESC {
		sort_order = "prev"
	}

	// Obtener el almacén
	d.err = d.prepareStore("read", r.FROM_TABLE)
	if d.err != "" {
		return e + d.err
	}

	switch {

	case r.ID != "":

		field_name := PREFIX_ID_NAME + r.FROM_TABLE

		if d.err = d.fieldIndexOK(r.FROM_TABLE, field_name); d.err != "" {
			return e + d.err
		}

		rangeObj := js.Global().Get("IDBKeyRange").Call("only", r.ID)
		index := d.store.Call("index", field_name)
		d.cursor = index.Call("openCursor", rangeObj)

	case r.ORDER_BY != "":

		if d.err = d.fieldIndexOK(r.FROM_TABLE, r.ORDER_BY); d.err != "" {
			return e + d.err
		}

		index := d.store.Call("index", r.ORDER_BY)
		// El valor nil como clave inicial significa que el cursor comenzará desde el primer registro en orden descendente y luego avanzará hacia registros posteriores en ese orden
		d.cursor = index.Call("openCursor", nil, sort_order)
	default:
		// normal
		d.cursor = d.store.Call("openCursor")
	}

	return
}
