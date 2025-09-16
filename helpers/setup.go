package helpers

import (
	"fmt"

	"github.com/cdvelop/indexdb"
)

// idGenerator implements the idGenerator interface for testing
type idGenerator struct {
	counter int
}

func (t *idGenerator) GetNewID() string {
	t.counter++
	return fmt.Sprintf("%d", t.counter) // Simple ID generation for tests
}

// SetupDB creates a new IndexDB instance for testing
func SetupDB(logger func(...any), dbName ...string) *indexdb.IndexDB {
	testDbName := "local_test_db"
	if len(dbName) > 0 {
		testDbName = dbName[0]
	}

	// Create a test ID generator
	idGen := &idGenerator{}

	return indexdb.New(testDbName, idGen, logger)
}

// User represents a sample struct for testing table creation
type User struct {
	ID    string
	Name  string
	Email string
}

func (u User) StructName() string {
	return "user"
}

// TestProduct represents another sample struct for testing
type Product struct {
	IDProduct string
	Name      string
	Price     float64
}
