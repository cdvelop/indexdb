package indexdb

import (
	"syscall/js"

	"github.com/cdvelop/unixid"
)

type ReadParams struct {
	FROM_TABLE string
	SORT_DESC  bool
	ID         string
	ORDER_BY   string
	WHERE      []interface{}
	RETURN_ANY bool
}

type ReadResults struct {
	Results []interface{}
	Error   error
}

type MainHandler struct {
	TimeAdapter
	SessionFrontendAdapter
	DataBaseAdapter
	Logger
}

type Logger interface {
	Log(v ...interface{})
}

type DataBaseAdapter interface {
	Create(table_name string, items []interface{}) (err error)
	Read(p *ReadParams, callback func(r *ReadResults, err error))
	ReadSync(p *ReadParams, data ...interface{}) (result []interface{}, err error)
	Delete(table_name string, all_data ...interface{}) (err error)
	Update(table_name string, all_data ...interface{}) (err error)
}

type TimeAdapter interface {
	Now() int64
}

type SessionFrontendAdapter interface {
	UserSessionNumber() string
}

type indexDB struct {
	db_name string
	db      js.Value

	Logger

	*unixid.UnixID

	//DATA IN TO CREATE, UPDATE
	data_in_any []map[string]interface{}
	data_in_str []map[string]string

	transaction js.Value
	store       js.Value
	cursor      js.Value
	result      js.Value
	err         error
}
