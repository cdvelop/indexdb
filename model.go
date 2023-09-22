package indexdb

import (
	"syscall/js"

	"github.com/cdvelop/model"
)

type indexDB struct {
	db      js.Value
	objects []*model.Object
}
