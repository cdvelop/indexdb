package indexdb

import (
	"syscall/js"
)

// run = RunBootData()
func Add() *indexDB {

	newDb := indexDB{
		db_name: "localdb",
		db:      js.Value{},
		objects: nil,
		run:     nil,
	}
	return &newDb
}

func (indexDB) RunOnClientDB() bool { //true base de datos corre en el browser
	return true
}
