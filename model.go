package indexdb

import (
	"syscall/js"

	"github.com/cdvelop/model"
	"github.com/cdvelop/unixid"
)

type indexDB struct {
	db_name string
	db      js.Value

	http model.FetchAdapter

	model.ObjectsHandler

	result func(error)

	*unixid.UnixID

	model.Logger

	backups       []backup
	backupRespond func(error)

	remaining int
}

type backup struct {
	object   *model.Object
	data     []map[string]interface{}
	finished bool
	err      error
}
