package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/nipunap/sqlite-mcp-server/internal/testutil"
)

func TestBatchOperations(t *testing.T) {
	db, dbPath := testutil.CreateTempDB(t)
	defer db.Close()

	// Create test table
	testutil.ExecuteSQL(t, db, `
		CREATE TABLE test (
			id INTEGER PRIMARY KEY,
			name TEXT,
			value INTEGER
		)
	`)

	registry, err := NewRegistry(":memory:")
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}
	defer registry.Close()

	// Register test database
	err = registry.RegisterDatabase(&DatabaseInfo{
		ID:     "test-db",
		Name:   "test",
		Path:   dbPath,
		Status: "active",
	})
	if err != nil {
		t.Fatalf("Failed to register database: %v", err)
	}

	manager := NewManager(registry)
	defer manager.CloseAll()

	t.Run("ExecuteBatch", func(t *testing.T) {
		operations := []BatchOperation{
			{
				Database: "test",
				Query:    "INSERT INTO test (name, value) VALUES (?, ?)",
				Args:     []interface{}{"test1", 1},
			},
			{
				Database: "test",
				Query:    "INSERT INTO test (name, value) VALUES (?, ?)",
				Args:     []interface{}{"test2", 2},
			},
			{
				Database: "test",
				Query:    "SELECT * FROM test WHERE value > ?",
				Args:     []interface{}{1},
			},
		}

		results := manager.ExecuteBatch(context.Background(), operations)

		if len(results) != len(operations) {
			t.Errorf("Expected %d results, got %d", len(operations), len(results))
		}

		// Check individual results
		for i, result := range results {
			if !result.Success {
				t.Errorf("Operation %d failed: %s", i, result.Error)
			}
		}

		// Verify data was inserted
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM test").Scan(&count)
		if err != nil {
			t.Fatalf("Failed to count records: %v", err)
		}
		if count != 2 {
			t.Errorf("Expected 2 records, got %d", count)
		}
	})

	t.Run("BulkInsert", func(t *testing.T) {
		operation := BulkInsertOperation{
			Database: "test",
			Table:    "test",
			Columns:  []string{"name", "value"},
			Values: [][]interface{}{
				{"bulk1", 10},
				{"bulk2", 20},
				{"bulk3", 30},
			},
		}

		rowsAffected, err := manager.BulkInsert(context.Background(), operation)
		if err != nil {
			t.Fatalf("BulkInsert failed: %v", err)
		}

		if rowsAffected != 3 {
			t.Errorf("Expected 3 rows affected, got %d", rowsAffected)
		}

		// Verify bulk insert
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM test WHERE value >= 10").Scan(&count)
		if err != nil {
			t.Fatalf("Failed to count bulk inserted records: %v", err)
		}
		if count != 3 {
			t.Errorf("Expected 3 bulk inserted records, got %d", count)
		}
	})

	t.Run("ConcurrentBatchOperations", func(t *testing.T) {
		// Create multiple batch operations
		batchCount := 5
		operationsPerBatch := 3
		allResults := make([][]BatchResult, batchCount)

		// Execute batches concurrently
		done := make(chan bool)
		for i := 0; i < batchCount; i++ {
			go func(batchIndex int) {
				operations := make([]BatchOperation, operationsPerBatch)
				for j := 0; j < operationsPerBatch; j++ {
					operations[j] = BatchOperation{
						Database: "test",
						Query:    "INSERT INTO test (name, value) VALUES (?, ?)",
						Args:     []interface{}{fmt.Sprintf("concurrent-%d-%d", batchIndex, j), batchIndex*100 + j},
					}
				}
				allResults[batchIndex] = manager.ExecuteBatch(context.Background(), operations)
				done <- true
			}(i)
		}

		// Wait for all batches to complete
		for i := 0; i < batchCount; i++ {
			<-done
		}

		// Verify all operations succeeded
		for i, results := range allResults {
			if len(results) != operationsPerBatch {
				t.Errorf("Batch %d: Expected %d results, got %d", i, operationsPerBatch, len(results))
			}
			for j, result := range results {
				if !result.Success {
					t.Errorf("Batch %d, Operation %d failed: %s", i, j, result.Error)
				}
			}
		}

		// Verify total number of records
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM test WHERE name LIKE 'concurrent-%'").Scan(&count)
		if err != nil {
			t.Fatalf("Failed to count concurrent records: %v", err)
		}
		expectedCount := batchCount * operationsPerBatch
		if count != expectedCount {
			t.Errorf("Expected %d concurrent records, got %d", expectedCount, count)
		}
	})
}
