package indexdb

import (
	"syscall/js"

	. "github.com/cdvelop/tinystring"
	"github.com/cdvelop/unixid"
)

// run = RunBootData()
func Add(h *MainHandler) (err error) {

	newDb := indexDB{
		db_name:               "localdb",
		db:                    js.Value{},
		http:                  h,
		ObjectsHandlerAdapter: h,
		BackupHandlerAdapter:  h,
		response:              nil,
		UnixID:                nil,
		Logger:                h,
	}

	h.DataBaseAdapter = &newDb

	uid, err := unixid.NewUnixID(h)
	if err != nil {
		return Err("error creating new unixid", err)
	}

	newDb.UnixID = uid

	return nil
}

func (indexDB) RunOnClientDB() bool { //true base de datos corre en el browser
	return true
}

func (indexDB) Lock()   {}
func (indexDB) Unlock() {}
