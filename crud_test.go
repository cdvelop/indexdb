package indexdb

import (
	"fmt"
	"testing"
	"time"
)

// mockLogger is a mock implementation of the Logger interface for testing.
type mockLogger struct{}

func (m *mockLogger) Log(v ...interface{}) {
	fmt.Println(v...)
}

// mockObjectsHandler is a mock implementation of the ObjectsHandlerAdapter interface for testing.
type mockObjectsHandler struct {
	objects []*Object
}

func (m *mockObjectsHandler) GetAllObjects(all bool) []*Object {
	return m.objects
}

// mockBackupHandler is a mock implementation of the BackupHandlerAdapter interface for testing.
type mockBackupHandler struct{}

func (m *mockBackupHandler) BackupOneObjectType(action string, table_name string, items any) {}

func TestCreate(t *testing.T) {
	// Channel to wait for the test to complete
	done := make(chan bool)

	// Mock implementations
	logger := &mockLogger{}
	objectsHandler := &mockObjectsHandler{
		objects: []*Object{
			{
				Table: "user",
				Fields: []Field{
					{Name: "id_user", Unique: true},
					{Name: "name"},
				},
			},
		},
	}
	backupHandler := &mockBackupHandler{}

	// Initialize indexDB
	db := &indexDB{
		db_name:               "test_db_create",
		Logger:                logger,
		ObjectsHandlerAdapter: objectsHandler,
		BackupHandlerAdapter:  backupHandler,
	}

	// Create tables
	db.CreateTablesInDB(objectsHandler.objects, func(err string) {
		if err != "" {
			t.Errorf("Failed to create tables: %s", err)
			done <- true
			return
		}

		// Data to create
		userData := map[string]string{
			"id_user": "1",
			"name":    "John Doe",
		}

		// Create object
		errCreate := db.CreateObjectsInDB("user", false, []map[string]string{userData})
		if errCreate != nil {
			t.Errorf("Failed to create object: %s", errCreate)
			done <- true
			return
		}

		// Read back the data to verify
		readParams := &ReadParams{
			FROM_TABLE: "user",
			ID:         "1",
		}

		db.ReadAsyncDataDB(readParams, func(r *ReadResults, errRead error) {
			if errRead != nil {
				t.Errorf("Failed to read data: %s", errRead)
				done <- true
				return
			}

			if len(r.ResultsString) != 1 {
				t.Errorf("Expected 1 result, got %d", len(r.ResultsString))
				done <- true
				return
			}

			if r.ResultsString[0]["name"] != "John Doe" {
				t.Errorf("Expected name 'John Doe', got '%s'", r.ResultsString[0]["name"])
			}

			done <- true
		})
	})

	// Wait for the test to complete or timeout
	select {
	case <-done:
		// Test finished
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out")
	}
}

func TestDelete(t *testing.T) {
	// Channel to wait for the test to complete
	done := make(chan bool)

	// Mock implementations
	logger := &mockLogger{}
	objectsHandler := &mockObjectsHandler{
		objects: []*Object{
			{
				Table: "user",
				Fields: []Field{
					{Name: "id_user", Unique: true},
					{Name: "name"},
				},
			},
		},
	}
	backupHandler := &mockBackupHandler{}

	// Initialize indexDB
	db := &indexDB{
		db_name:               "test_db_delete",
		Logger:                logger,
		ObjectsHandlerAdapter: objectsHandler,
		BackupHandlerAdapter:  backupHandler,
	}

	// Create tables
	db.CreateTablesInDB(objectsHandler.objects, func(err string) {
		if err != "" {
			t.Errorf("Failed to create tables: %s", err)
			done <- true
			return
		}

		// Data to create
		userData := map[string]string{
			"id_user": "6",
			"name":    "David",
		}

		// Create object
		errCreate := db.CreateObjectsInDB("user", false, []map[string]string{userData})
		if errCreate != nil {
			t.Errorf("Failed to create object: %s", errCreate)
			done <- true
			return
		}

		// Data to delete
		deleteUserData := map[string]string{
			"id_user": "6",
		}

		// Delete object
		errDelete := db.DeleteObjectsInDB("user", false, deleteUserData)
		if errDelete != nil {
			t.Errorf("Failed to delete object: %s", errDelete)
			done <- true
			return
		}

		// Read back the data to verify it's gone
		readParams := &ReadParams{
			FROM_TABLE: "user",
			ID:         "6",
		}

		db.ReadAsyncDataDB(readParams, func(r *ReadResults, errRead error) {
			if errRead != nil {
				t.Errorf("Failed to read data: %s", errRead)
				done <- true
				return
			}

			if len(r.ResultsString) != 0 {
				t.Errorf("Expected 0 results, got %d", len(r.ResultsString))
			}

			done <- true
		})
	})

	// Wait for the test to complete or timeout
	select {
	case <-done:
		// Test finished
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out")
	}
}

func TestUpdate(t *testing.T) {
	// Channel to wait for the test to complete
	done := make(chan bool)

	// Mock implementations
	logger := &mockLogger{}
	objectsHandler := &mockObjectsHandler{
		objects: []*Object{
			{
				Table: "user",
				Fields: []Field{
					{Name: "id_user", Unique: true},
					{Name: "name"},
				},
			},
		},
	}
	backupHandler := &mockBackupHandler{}

	// Initialize indexDB
	db := &indexDB{
		db_name:               "test_db_update",
		Logger:                logger,
		ObjectsHandlerAdapter: objectsHandler,
		BackupHandlerAdapter:  backupHandler,
	}

	// Create tables
	db.CreateTablesInDB(objectsHandler.objects, func(err string) {
		if err != "" {
			t.Errorf("Failed to create tables: %s", err)
			done <- true
			return
		}

		// Data to create
		userData := map[string]string{
			"id_user": "5",
			"name":    "Charlie",
		}

		// Create object
		errCreate := db.CreateObjectsInDB("user", false, []map[string]string{userData})
		if errCreate != nil {
			t.Errorf("Failed to create object: %s", errCreate)
			done <- true
			return
		}

		// Data to update
		updatedUserData := map[string]string{
			"id_user": "5",
			"name":    "Charlie Brown",
		}

		// Update object
		errUpdate := db.UpdateObjectsInDB("user", false, updatedUserData)
		if errUpdate != nil {
			t.Errorf("Failed to update object: %s", errUpdate)
			done <- true
			return
		}

		// Read back the data to verify
		readParams := &ReadParams{
			FROM_TABLE: "user",
			ID:         "5",
		}

		db.ReadAsyncDataDB(readParams, func(r *ReadResults, errRead error) {
			if errRead != nil {
				t.Errorf("Failed to read data: %s", errRead)
				done <- true
				return
			}

			if len(r.ResultsString) != 1 {
				t.Errorf("Expected 1 result, got %d", len(r.ResultsString))
				done <- true
				return
			}

			if r.ResultsString[0]["name"] != "Charlie Brown" {
				t.Errorf("Expected name 'Charlie Brown', got '%s'", r.ResultsString[0]["name"])
			}

			done <- true
		})
	})

	// Wait for the test to complete or timeout
	select {
	case <-done:
		// Test finished
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out")
	}
}

func TestReadAll(t *testing.T) {
	// Channel to wait for the test to complete
	done := make(chan bool)

	// Mock implementations
	logger := &mockLogger{}
	objectsHandler := &mockObjectsHandler{
		objects: []*Object{
			{
				Table: "user",
				Fields: []Field{
					{Name: "id_user", Unique: true},
					{Name: "name"},
				},
			},
		},
	}
	backupHandler := &mockBackupHandler{}

	// Initialize indexDB
	db := &indexDB{
		db_name:               "test_db_read",
		Logger:                logger,
		ObjectsHandlerAdapter: objectsHandler,
		BackupHandlerAdapter:  backupHandler,
	}

	// Create tables
	db.CreateTablesInDB(objectsHandler.objects, func(err string) {
		if err != "" {
			t.Errorf("Failed to create tables: %s", err)
			done <- true
			return
		}

		// Data to create
		users := []map[string]string{
			{"id_user": "3", "name": "Alice"},
			{"id_user": "4", "name": "Bob"},
		}

		// Create objects
		errCreate := db.CreateObjectsInDB("user", false, users)
		if errCreate != nil {
			t.Errorf("Failed to create objects: %s", errCreate)
			done <- true
			return
		}

		// Read back all data to verify
		readParams := &ReadParams{
			FROM_TABLE: "user",
		}

		db.ReadAsyncDataDB(readParams, func(r *ReadResults, errRead error) {
			if errRead != nil {
				t.Errorf("Failed to read data: %s", errRead)
				done <- true
				return
			}

			if len(r.ResultsString) != 2 {
				t.Errorf("Expected 2 results, got %d", len(r.ResultsString))
			}

			done <- true
		})
	})

	// Wait for the test to complete or timeout
	select {
	case <-done:
		// Test finished
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out")
	}
}

func TestRead(t *testing.T) {
	// Channel to wait for the test to complete
	done := make(chan bool)

	// Mock implementations
	logger := &mockLogger{}
	objectsHandler := &mockObjectsHandler{
		objects: []*Object{
			{
				Table: "user",
				Fields: []Field{
					{Name: "id_user", Unique: true},
					{Name: "name"},
				},
			},
		},
	}
	backupHandler := &mockBackupHandler{}

	// Initialize indexDB
	db := &indexDB{
		db_name:               "test_db_read_all",
		Logger:                logger,
		ObjectsHandlerAdapter: objectsHandler,
		BackupHandlerAdapter:  backupHandler,
	}

	// Create tables
	db.CreateTablesInDB(objectsHandler.objects, func(err string) {
		if err != "" {
			t.Errorf("Failed to create tables: %s", err)
			done <- true
			return
		}

		// Data to create
		userData := map[string]string{
			"id_user": "2",
			"name":    "Jane Doe",
		}

		// Create object
		errCreate := db.CreateObjectsInDB("user", false, []map[string]string{userData})
		if errCreate != nil {
			t.Errorf("Failed to create object: %s", errCreate)
			done <- true
			return
		}

		// Read back the data to verify
		readParams := &ReadParams{
			FROM_TABLE: "user",
			ID:         "2",
		}

		db.ReadAsyncDataDB(readParams, func(r *ReadResults, errRead error) {
			if errRead != nil {
				t.Errorf("Failed to read data: %s", errRead)
				done <- true
				return
			}

			if len(r.ResultsString) != 1 {
				t.Errorf("Expected 1 result, got %d", len(r.ResultsString))
				done <- true
				return
			}

			if r.ResultsString[0]["name"] != "Jane Doe" {
				t.Errorf("Expected name 'Jane Doe', got '%s'", r.ResultsString[0]["name"])
			}

			done <- true
		})
	})

	// Wait for the test to complete or timeout
	select {
	case <-done:
		// Test finished
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out")
	}
}
