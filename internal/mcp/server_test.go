package mcp

import (
	"database/sql"
	"encoding/json"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create test table
	_, err = db.Exec(`
		CREATE TABLE test (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return db, cleanup
}

func TestServerCapabilities(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	server, err := NewServer(db)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test tool registration
	msg := &JSONRPCMessage{
		Version: "2.0",
		ID:      "1",
		Method:  "db/get_table_schema",
		Params:  json.RawMessage(`{"table_name": "test"}`),
	}

	response := server.handleMessage(msg)
	if response.Error != nil {
		t.Errorf("Tool execution failed: %v", response.Error)
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		t.Errorf("Failed to unmarshal result: %v", err)
	}

	if result["table_name"] != "test" {
		t.Errorf("Expected table_name 'test', got %v", result["table_name"])
	}

	// Test resource registration
	msg = &JSONRPCMessage{
		Version: "2.0",
		ID:      "2",
		Method:  "db/tables",
	}

	response = server.handleMessage(msg)
	if response.Error != nil {
		t.Errorf("Resource access failed: %v", response.Error)
	}

	var tables map[string]interface{}
	err = json.Unmarshal(response.Result, &tables)
	if err != nil {
		t.Errorf("Failed to unmarshal tables: %v", err)
	}

	tableList := tables["tables"].([]interface{})
	if len(tableList) != 1 {
		t.Errorf("Expected 1 table, got %d", len(tableList))
	}

	// Test prompt registration
	msg = &JSONRPCMessage{
		Version: "2.0",
		ID:      "3",
		Method:  "db/query_help",
	}

	response = server.handleMessage(msg)
	if response.Error != nil {
		t.Errorf("Prompt access failed: %v", response.Error)
	}

	var prompt string
	err = json.Unmarshal(response.Result, &prompt)
	if err != nil {
		t.Errorf("Failed to unmarshal prompt: %v", err)
	}

	if prompt == "" {
		t.Error("Expected non-empty prompt content")
	}
}

func TestServerErrorHandling(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	server, err := NewServer(db)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test invalid method
	msg := &JSONRPCMessage{
		Version: "2.0",
		ID:      "1",
		Method:  "invalid_method",
	}

	response := server.handleMessage(msg)
	if response.Error == nil {
		t.Error("Expected error for invalid method, got nil")
	}

	// Test invalid parameters
	msg = &JSONRPCMessage{
		Version: "2.0",
		ID:      "2",
		Method:  "db/get_table_schema",
		Params:  json.RawMessage(`{"invalid": "params"}`),
	}

	response = server.handleMessage(msg)
	if response.Error == nil {
		t.Error("Expected error for invalid parameters, got nil")
	}

	// Test invalid JSON-RPC version
	msg = &JSONRPCMessage{
		Version: "1.0",
		ID:      "3",
		Method:  "db/tables",
	}

	response = server.handleMessage(msg)
	if response.Error == nil {
		t.Error("Expected error for invalid JSON-RPC version, got nil")
	}
}
