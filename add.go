package indexdb

import (
	"syscall/js"

	"github.com/cdvelop/model"
	"github.com/cdvelop/unixid"
)

// run = RunBootData()
func Add(h *model.MainHandler) (err string) {

	newDb := indexDB{
		db_name:               "localdb",
		db:                    js.Value{},
		http:                  h,
		ObjectsHandlerAdapter: h,
		result:                nil,
		UnixID:                nil,
		Logger:                h.Logger,
	}

	h.DataBaseAdapter = &newDb

	uid, err := unixid.NewHandler(h.TimeAdapter, newDb, h.SessionFrontendAdapter)
	if err != "" {
		return err
	}

	newDb.UnixID = uid

	return ""
}

func (indexDB) RunOnClientDB() bool { //true base de datos corre en el browser
	return true
}

func (indexDB) Lock()   {}
func (indexDB) Unlock() {}
