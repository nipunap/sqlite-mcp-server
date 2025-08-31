package resources

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nipunap/sqlite-mcp-server/internal/db"
)

func setupTestDB(t *testing.T) (*db.Manager, func()) {
	// Create in-memory registry
	registry, err := db.NewRegistry(":memory:")
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	// Create manager
	manager := db.NewManager(registry)

	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", fmt.Sprintf("test_db_resources_%s_*.db", t.Name()))
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close() // Close immediately, we just need the path
	testDBPath := tempFile.Name()

	testDB, err := sql.Open("sqlite3", testDBPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create test tables
	_, err = testDB.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE
		);
		CREATE INDEX idx_users_email ON users(email);

		CREATE TABLE posts (
			id INTEGER PRIMARY KEY,
			user_id INTEGER,
			title TEXT NOT NULL,
			content TEXT,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
		CREATE INDEX idx_posts_user_id ON posts(user_id);
	`)
	if err != nil {
		t.Fatalf("Failed to create test tables: %v", err)
	}
	testDB.Close()

	// Register the test database with the registry
	info := &db.DatabaseInfo{
		ID:          "test-resources-id",
		Name:        "test",
		Path:        testDBPath,
		Description: "Test database for resources",
		ReadOnly:    false,
		Owner:       "test",
		Status:      "active",
	}

	err = registry.RegisterDatabase(info)
	if err != nil {
		t.Fatalf("Failed to register test database: %v", err)
	}

	cleanup := func() {
		if err := manager.CloseAll(); err != nil {
			t.Logf("Error closing manager: %v", err)
		}
		if err := registry.Close(); err != nil {
			t.Logf("Error closing registry: %v", err)
		}
		os.Remove(testDBPath) // Clean up temp file
	}

	return manager, cleanup
}

func TestGetTables(t *testing.T) {
	t.Parallel()

	manager, cleanup := setupTestDB(t)
	defer cleanup()

	resources := NewDBResources(manager)

	params := []byte(`{"database_name": "test"}`)
	result, err := resources.GetTables(params)
	if err != nil {
		t.Errorf("GetTables failed: %v", err)
	}

	response := result.(map[string]interface{})
	tables := response["tables"].([]map[string]interface{})

	if len(tables) != 2 {
		t.Errorf("Expected 2 tables, got %d", len(tables))
	}

	// Verify table names
	tableNames := make(map[string]bool)
	for _, table := range tables {
		name := table["name"].(string)
		tableNames[name] = true

		if table["type"].(string) != "table" {
			t.Errorf("Expected type 'table' for %s, got %s", name, table["type"])
		}

		if table["sql"].(string) == "" {
			t.Errorf("Expected non-empty SQL for table %s", name)
		}
	}

	if !tableNames["users"] || !tableNames["posts"] {
		t.Error("Missing expected tables 'users' and/or 'posts'")
	}
}

func TestGetSchema(t *testing.T) {
	t.Parallel()

	manager, cleanup := setupTestDB(t)
	defer cleanup()

	resources := NewDBResources(manager)

	params := []byte(`{"database_name": "test"}`)
	result, err := resources.GetSchema(params)
	if err != nil {
		t.Errorf("GetSchema failed: %v", err)
	}

	response := result.(map[string]interface{})
	schema := response["schema"].(map[string]map[string]interface{})

	// Check users table
	if users, ok := schema["users"]; ok {
		if users["sql"].(string) == "" {
			t.Error("Expected non-empty SQL for users table")
		}

		indexes := users["indexes"].(map[string]string)
		if _, ok := indexes["idx_users_email"]; !ok {
			t.Error("Missing expected index 'idx_users_email' for users table")
		}
	} else {
		t.Error("Missing users table in schema")
	}

	// Check posts table
	if posts, ok := schema["posts"]; ok {
		if posts["sql"].(string) == "" {
			t.Error("Expected non-empty SQL for posts table")
		}

		indexes := posts["indexes"].(map[string]string)
		if _, ok := indexes["idx_posts_user_id"]; !ok {
			t.Error("Missing expected index 'idx_posts_user_id' for posts table")
		}
	} else {
		t.Error("Missing posts table in schema")
	}
}
