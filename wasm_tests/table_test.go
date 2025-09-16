//go:build js && wasm
// +build js,wasm

package wasmtests_test

import (
	"testing"

	"github.com/cdvelop/indexdb/helpers"
)

// TestCreateTableIfNotExists tests the CreateTableIfNotExists function
func TestCreateTableIfNotExists(t *testing.T) {

	logger := func(args ...any) {
		t.Log(args...)
	}

	db := helpers.SetupDB(logger, "test_db_create_table")

	db.InitDB()
}
