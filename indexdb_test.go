package indexdb_test

import (
	"testing"

	"github.com/cdvelop/wasmtest"
)

// TestTableCreation runs WebAssembly tests for table creation functionality
func TestIndexDB(t *testing.T) {
	// Run WebAssembly tests in the browser
	if err := wasmtest.RunTests(); err != nil {
		t.Fatal(err)
	}
}
