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

	result func(err string)

	*unixid.UnixID

	model.Logger

	backups       []backup
	backupRespond func(err string)

	remaining int

	//READ NEW
	cursor     js.Value
	readParams model.ReadParams
}

type backup struct {
	object   *model.Object
	data     []map[string]interface{}
	finished bool
	err      string
}
