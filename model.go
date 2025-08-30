package indexdb

import (
	"syscall/js"

	"github.com/cdvelop/unixid"
)

const PREFIX_ID_NAME = "id_"

type ReadParams struct {
	FROM_TABLE string
	SORT_DESC  bool
	ID         string
	ORDER_BY   string
	WHERE      []map[string]string
	RETURN_ANY bool
}

type ReadResults struct {
	ResultsString []map[string]string
	ResultsAny    []map[string]any
}

type MainHandler struct {
	TimeAdapter
	SessionFrontendAdapter
	DataBaseAdapter
	FetchAdapter
	ObjectsHandlerAdapter
	BackupHandlerAdapter
	Logger
}

type Object struct {
	Table           string
	Fields          []Field
	NoAddObjectInDB bool
}

type Field struct {
	Name            string
	NotRequiredInDB bool
	Unique          bool
}

func (o *Object) PrimaryKeyName() string {
	return PREFIX_ID_NAME + o.Table
}

type FetchAdapter interface{}

type ObjectsHandlerAdapter interface {
	GetAllObjects(all bool) []*Object
}

type BackupHandlerAdapter interface {
	BackupOneObjectType(action string, table_name string, items any)
}

type Logger interface {
	Log(v ...interface{})
}

type DataBaseAdapter interface {
	RunOnClientDB() bool
	Lock()
	Unlock()
	CreateObjectsInDB(table_name string, on_server_too bool, items any) (err string)
	ReadAsyncDataDB(p *ReadParams, callback func(r *ReadResults, err string))
	ReadSyncDataDB(p *ReadParams, data ...map[string]string) (result []map[string]string, err string)
	DeleteObjectsInDB(table_name string, on_server_too bool, all_data ...map[string]string) (err string)
	UpdateObjectsInDB(table_name string, on_server_too bool, all_data ...map[string]string) (err string)
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

	http FetchAdapter
	ObjectsHandlerAdapter
	BackupHandlerAdapter
	Logger

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
}
