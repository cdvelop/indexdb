package indexdb

import (
	"syscall/js"

	"github.com/cdvelop/model"
	"github.com/cdvelop/timeclient"
	"github.com/cdvelop/unixid"
)

// run = RunBootData()
func Add(u model.UserAuthNumber, l model.Logger) (*indexDB, error) {

	newDb := indexDB{
		db_name: "localdb",
		db:      js.Value{},
		objects: nil,
		run:     nil,
		UnixID:  nil,
		Logger:  l,
	}

	uid, err := unixid.NewHandler(&timeclient.TimeCLient{}, newDb, u)
	if err != nil {
		return nil, err
	}

	newDb.UnixID = uid

	return &newDb, nil
}

func (indexDB) RunOnClientDB() bool { //true base de datos corre en el browser
	return true
}

func (indexDB) Lock()   {}
func (indexDB) Unlock() {}
