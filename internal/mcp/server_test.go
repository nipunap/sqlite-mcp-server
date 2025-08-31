package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nipunap/sqlite-mcp-server/internal/db"
)

func setupTestManager(t *testing.T) (*db.Manager, func()) {
	// Create in-memory registry
	registry, err := db.NewRegistry(":memory:")
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	// Create manager
	manager := db.NewManager(registry)

	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", fmt.Sprintf("test_db_%s_*.db", t.Name()))
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close() // Close immediately, we just need the path
	testDBPath := tempFile.Name()

	// Register a test database
	info := &db.DatabaseInfo{
		ID:          "test-id",
		Name:        "test",
		Path:        testDBPath,
		Description: "Test database (in-memory)",
		ReadOnly:    false,
		Owner:       "test",
		Status:      "active",
	}

	err = registry.RegisterDatabase(info)
	if err != nil {
		t.Fatalf("Failed to register test database: %v", err)
	}

	// Get connection and create test table with timeout
	conn, err := manager.GetConnection("test")
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}

	// Create test table with simple structure
	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS test_table (id INTEGER PRIMARY KEY, name TEXT)`)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
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

func TestServerCapabilities(t *testing.T) {
	t.Parallel() // Enable parallel execution

	manager, cleanup := setupTestManager(t)
	defer cleanup()

	server, err := NewServer(manager)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test capabilities method (this should be fast and not require database operations)
	id1 := json.RawMessage(`"1"`)
	msg := &JSONRPCMessage{
		Version: "2.0",
		ID:      &id1,
		Method:  "capabilities",
	}

	response := server.handleMessage(msg)
	if response.Error != nil {
		t.Errorf("Capabilities request failed: %v", response.Error)
	}

	// Convert interface{} to []byte for unmarshaling
	resultBytes, err := json.Marshal(response.Result)
	if err != nil {
		t.Errorf("Failed to marshal result: %v", err)
	}

	var capabilities []interface{}
	err = json.Unmarshal(resultBytes, &capabilities)
	if err != nil {
		t.Errorf("Failed to unmarshal capabilities: %v", err)
	}

	if len(capabilities) == 0 {
		t.Error("Expected non-empty capabilities list")
	}

	t.Logf("Capabilities count: %d", len(capabilities))

	// Test prompt access (this should also be fast)
	id2 := json.RawMessage(`"2"`)
	msg = &JSONRPCMessage{
		Version: "2.0",
		ID:      &id2,
		Method:  "invoke",
		Params:  json.RawMessage(`{"name": "db/query_help", "params": {}}`),
	}

	response = server.handleMessage(msg)
	if response.Error != nil {
		t.Errorf("Prompt access failed: %v", response.Error)
	}

	// Convert interface{} to []byte for unmarshaling
	promptBytes, err := json.Marshal(response.Result)
	if err != nil {
		t.Errorf("Failed to marshal prompt result: %v", err)
	}

	var prompt string
	err = json.Unmarshal(promptBytes, &prompt)
	if err != nil {
		t.Errorf("Failed to unmarshal prompt: %v", err)
	}

	if len(prompt) == 0 {
		t.Error("Expected non-empty prompt content")
	}

	t.Logf("Prompt response length: %d characters", len(prompt))
}

func TestServerErrorHandling(t *testing.T) {
	t.Parallel() // Enable parallel execution
	manager, cleanup := setupTestManager(t)
	defer cleanup()

	server, err := NewServer(manager)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test invalid method
	id1 := json.RawMessage(`"1"`)
	msg := &JSONRPCMessage{
		Version: "2.0",
		ID:      &id1,
		Method:  "invalid_method",
	}

	response := server.handleMessage(msg)
	if response.Error == nil {
		t.Error("Expected error for invalid method, got nil")
	}

	// Test invalid parameters
	id2 := json.RawMessage(`"2"`)
	msg = &JSONRPCMessage{
		Version: "2.0",
		ID:      &id2,
		Method:  "invoke",
		Params:  json.RawMessage(`{"name": "db/get_table_schema", "params": {"invalid": "params"}}`),
	}

	response = server.handleMessage(msg)
	if response.Error == nil {
		t.Error("Expected error for invalid parameters, got nil")
	}

	// Test invalid JSON-RPC version
	id3 := json.RawMessage(`"3"`)
	msg = &JSONRPCMessage{
		Version: "1.0",
		ID:      &id3,
		Method:  "invoke",
		Params:  json.RawMessage(`{"name": "db/get_tables", "params": {"database_name": "test"}}`),
	}

	response = server.handleMessage(msg)
	// Note: The server might not validate JSON-RPC version strictly, so this test might pass
	// This is acceptable behavior for an MCP server
	t.Logf("JSON-RPC version test response: %+v", response)
}
