package indexdb

import (
	"syscall/js"

	"github.com/cdvelop/unixid"
)

// run = RunBootData()
func Add(h *MainHandler) (err string) {

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
