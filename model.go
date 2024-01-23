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
	model.ObjectsHandlerAdapter
	model.BackupHandlerAdapter
	model.Logger

	response func(err string)

	*unixid.UnixID

	//DATA IN TO CREATE, UPDATE
	data_in_any []map[string]interface{}
	data_in_str []map[string]string

	transaction js.Value
	store       js.Value
	cursor      js.Value
	result      js.Value
	err         string

	// readParams model.ReadParams
}
