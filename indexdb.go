package indexdb

import (
	"sync"

	. "github.com/cdvelop/tinystring"
	"syscall/js"
)

type idGenerator interface {
	GetNewID() string
}

type structName interface {
	StructName() string
}

type IndexDB struct {
	db_name string
	db      js.Value

	tables []any

	logger func(...any)

	idGenerator

	//DATA IN TO CREATE, UPDATE
	data []map[string]any

	transaction js.Value
	store       js.Value
	cursor      js.Value
	result      js.Value
	err         error

	initDone      chan struct{}
	initOnce      sync.Once
	initCompleted bool
}

// New creates a new IndexDB instance with the given database name, ID generator, and logger.
//
//	type idGenerator interface {
//		GetNewID() string
//	}
func New(dbName string, idg idGenerator, logger func(...any)) *IndexDB {

	idb := IndexDB{
		db_name:     dbName,
		db:          js.Value{},
		idGenerator: idg,
		logger:      logger,
		initDone:    make(chan struct{}, 1),
	}

	return &idb
}

// InitDB require structs with interface structName { StructName() string } for table names
// eg: type User struct { NameField string }  =>  User.StructName() string { return "user" }
// eg: type Product struct { NameField string }  =>  Product.StructName() string { return "product" }
// then call: indexdb.InitDB(User{}, Product{})
// This is because Reflect support is not complete, it is the only way this library can be compatible with TinyGo.
func (d *IndexDB) InitDB(structTables ...any) {

	d.tables = structTables

	// Open connection to IndexedDB
	db := js.Global().Get("indexedDB").Call("open", d.db_name)

	// add events
	db.Call("addEventListener", "error", js.FuncOf(d.showDbError))
	db.Call("addEventListener", "success", js.FuncOf(d.openExistingDB))
	db.Call("addEventListener", "upgradeneeded", js.FuncOf(d.upgradeneeded))

	<-d.initDone // wait until init is done
}

func (d *IndexDB) open(p *js.Value, message string) (err error) {

	d.db = p.Get("target").Get("result")

	if !d.db.Truthy() {
		return Err("error open", d.db_name, message)
	}

	// d.Logger("***", message, "IndexDB Connection:", d.db_name, " ***")
	// DB : localdb Established, Engine: indexedDB
	return nil
}

func (d *IndexDB) upgradeneeded(this js.Value, p []js.Value) any {

	err := d.open(&p[0], "upgradeneeded")
	if err != nil {
		d.logger(err)
		return nil
	}

	for i, table := range d.tables {

		t, ok := table.(structName)
		if !ok {
			d.logger("error table", i, "does not implement structName interface (Name() string)")
			continue
		}

		// d.logger("**** CREANDO TABLA: ", t.StructName())

		err := d.createTable(t.StructName(), table)
		if err != nil {
			d.logger(err)
			continue
		}

	}

	// Wait for the version change transaction to complete
	transaction := p[0].Get("target").Get("transaction")
	transaction.Call("addEventListener", "complete", js.FuncOf(func(this js.Value, p []js.Value) any {
		d.initOnce.Do(func() { d.initCompleted = true; close(d.initDone) })
		return nil
	}))
	transaction.Call("addEventListener", "error", js.FuncOf(func(this js.Value, p []js.Value) any {
		d.logger("version change transaction error")
		d.initOnce.Do(func() { d.initCompleted = true; close(d.initDone) })
		return nil
	}))
	transaction.Call("addEventListener", "abort", js.FuncOf(func(this js.Value, p []js.Value) any {
		d.logger("version change transaction aborted")
		d.initOnce.Do(func() { d.initCompleted = true; close(d.initDone) })
		return nil
	}))

	return nil
}

func (d *IndexDB) showDbError(this js.Value, p []js.Value) any {
	d.logger("indexDB Error", p[0])
	return nil
}

func (d *IndexDB) openExistingDB(this js.Value, p []js.Value) any {
	err := d.open(&p[0], "OPEN")
	if err != nil {
		d.logger("open existing db error:", err)
		return nil
	}

	if !d.initCompleted {
		d.logger("open existing db success")
	}

	d.initOnce.Do(func() { d.initCompleted = true; close(d.initDone) })

	return nil
}
