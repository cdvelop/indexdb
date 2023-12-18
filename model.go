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

	model.ObjectHandler

	result func(err string)

	*unixid.UnixID

	model.Logger

	backups       []backup
	backupRespond func(err string)

	remaining int

	//DATA IN TO CREATE, UPDATE
	data_in_any []map[string]interface{}
	data_in_str []map[string]string

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
