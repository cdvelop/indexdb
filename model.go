package indexdb

import (
	"syscall/js"

	"github.com/cdvelop/model"
	"github.com/cdvelop/unixid"
)

type indexDB struct {
	db_name string
	db      js.Value
	objects []*model.Object

	run model.Subsequently

	*unixid.UnixID

	model.Logger
}
