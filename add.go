package indexdb

import (
	"syscall/js"

	"github.com/cdvelop/model"
)

func Add(objects []*model.Object) *indexDB {

	newDb := indexDB{
		db:      js.Value{},
		objects: objects,
	}

	newDb.initDataBase()

	return &newDb
}
