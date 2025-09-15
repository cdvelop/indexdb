package indexdb

import (
	"syscall/js"

	. "github.com/cdvelop/tinystring"
)


func (d *indexDB) open(p *js.Value, message string) (err error) {

	d.db = p.Get("target").Get("result")

	if !d.db.Truthy() {
		return Err("open indexdb error no se logro establecer conexión " + d.db_name)
	}

	// d.Log("***", message, "IndexDB Connection:", d.db_name, " ***")
	// DB : localdb Established, Engine: indexedDB
	return nil
}
