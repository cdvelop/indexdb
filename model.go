package indexdb

import (
	"syscall/js"

	"github.com/cdvelop/model"
)

type indexDB struct {
	db_name string
	db      js.Value
	objects []*model.Object

	run model.Subsequently
}
