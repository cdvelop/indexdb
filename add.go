package indexdb

import (
	"syscall/js"

	"github.com/cdvelop/model"
	"github.com/cdvelop/unixid"
)

// run = RunBootData()
func Add(h *model.Handlers) error {

	newDb := indexDB{
		db_name:        "localdb",
		db:             js.Value{},
		http:           h,
		ObjectsHandler: h,
		result:         nil,
		UnixID:         nil,
		Logger:         h.Logger,
	}

	h.DataBaseAdapter = &newDb

	uid, err := unixid.NewHandler(h.TimeAdapter, newDb, h.AuthAdapter)
	if err != nil {
		return err
	}

	newDb.UnixID = uid

	return nil
}

func (indexDB) RunOnClientDB() bool { //true base de datos corre en el browser
	return true
}

func (indexDB) Lock()   {}
func (indexDB) Unlock() {}
