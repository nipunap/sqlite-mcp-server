package testutil

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// CreateTempDB creates a temporary SQLite database for testing
func CreateTempDB(t *testing.T) (*sql.DB, string) {
	t.Helper()

	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
		os.Remove(dbPath)
	})

	return db, dbPath
}

// ExecuteSQL executes SQL statements on the test database
func ExecuteSQL(t *testing.T, db *sql.DB, sql string) {
	t.Helper()

	_, err := db.Exec(sql)
	if err != nil {
		t.Fatalf("Failed to execute SQL: %v", err)
	}
}
