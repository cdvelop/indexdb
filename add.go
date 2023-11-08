package indexdb

import (
	"syscall/js"

	"github.com/cdvelop/model"
	"github.com/cdvelop/timeclient"
	"github.com/cdvelop/unixid"
)

// run = RunBootData()
func Add(u model.UserAuthNumber, h *model.Handlers) error {

	newDb := indexDB{
		db_name: "localdb",
		db:      js.Value{},
		http:    h,
		objects: nil, //add in CreateTablesInDB func
		run:     nil,
		UnixID:  nil,
		Logger:  h.Logger,
	}

	h.DataBaseAdapter = &newDb

	uid, err := unixid.NewHandler(&timeclient.TimeCLient{}, newDb, u)
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
