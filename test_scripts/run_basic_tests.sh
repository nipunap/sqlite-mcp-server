#!/bin/bash

# Basic MCP Server Tests
# This script sends JSON-RPC messages to test core functionality

set -e

echo "🧪 Running Basic MCP Server Tests"
echo "================================="

# Function to send JSON-RPC message
test_rpc_call() {
    local test_name="$1"
    local message="$2"

    echo ""
    echo "📋 Testing: $test_name"
    echo "-------------------"

    # Create temporary files for communication
    local input_file=$(mktemp)
    local output_file=$(mktemp)

    # Write message to input file
    echo "$message" > "$input_file"

    # Start server and send message
    timeout 5s ./sqlite-mcp-server --registry registry.db --db test_manual.db < "$input_file" > "$output_file" 2>&1 || true

    # Display result
    if [ -s "$output_file" ]; then
        echo "Response:"
        cat "$output_file" | jq . 2>/dev/null || cat "$output_file"
    else
        echo "No response received"
    fi

    # Cleanup
    rm -f "$input_file" "$output_file"
}

# Test 1: Initialize
test_rpc_call "Server Initialize" '{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {},
    "clientInfo": {
      "name": "test-client",
      "version": "1.0.0"
    }
  }
}'

# Test 2: List Tools
test_rpc_call "List Available Tools" '{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/list",
  "params": {}
}'

# Test 3: List Resources
test_rpc_call "List Available Resources" '{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "resources/list",
  "params": {}
}'

# Test 4: List Databases
test_rpc_call "List Registered Databases" '{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "tools/call",
  "params": {
    "name": "db/list_databases",
    "arguments": {}
  }
}'

# Test 5: Get Tables
test_rpc_call "Get Tables from Default Database" '{
  "jsonrpc": "2.0",
  "id": 5,
  "method": "tools/call",
  "params": {
    "name": "db/get_tables",
    "arguments": {
      "database_name": "default"
    }
  }
}'

# Test 6: Get Schema
test_rpc_call "Get Database Schema" '{
  "jsonrpc": "2.0",
  "id": 6,
  "method": "tools/call",
  "params": {
    "name": "db/get_schema",
    "arguments": {
      "database_name": "default"
    }
  }
}'

# Test 7: Query Users
test_rpc_call "Query Users Table" '{
  "jsonrpc": "2.0",
  "id": 7,
  "method": "tools/call",
  "params": {
    "name": "db/query",
    "arguments": {
      "database_name": "default",
      "query": "SELECT * FROM users",
      "args": []
    }
  }
}'

# Test 8: Query with Parameters
test_rpc_call "Query with Parameters" '{
  "jsonrpc": "2.0",
  "id": 8,
  "method": "tools/call",
  "params": {
    "name": "db/query",
    "arguments": {
      "database_name": "default",
      "query": "SELECT * FROM users WHERE name LIKE ?",
      "args": ["%Alice%"]
    }
  }
}'

# Test 9: Insert Record
test_rpc_call "Insert New User" '{
  "jsonrpc": "2.0",
  "id": 9,
  "method": "tools/call",
  "params": {
    "name": "db/insert_record",
    "arguments": {
      "database_name": "default",
      "table_name": "users",
      "data": {
        "name": "Test User",
        "email": "test@example.com"
      }
    }
  }
}'

# Test 10: Register New Database
test_rpc_call "Register Inventory Database" '{
  "jsonrpc": "2.0",
  "id": 10,
  "method": "tools/call",
  "params": {
    "name": "db/register_database",
    "arguments": {
      "name": "inventory",
      "path": "'$(pwd)'/inventory.db",
      "description": "Inventory management database",
      "readonly": false,
      "owner": "test_user"
    }
  }
}'

# Test 11: Query Newly Registered Database
test_rpc_call "Query Inventory Database" '{
  "jsonrpc": "2.0",
  "id": 11,
  "method": "tools/call",
  "params": {
    "name": "db/query",
    "arguments": {
      "database_name": "inventory",
      "query": "SELECT w.name as warehouse, s.product_name, s.quantity FROM stock s JOIN warehouses w ON s.warehouse_id = w.id",
      "args": []
    }
  }
}'

# Test 12: Error Handling - Invalid Query
test_rpc_call "Error Handling - Invalid SQL" '{
  "jsonrpc": "2.0",
  "id": 12,
  "method": "tools/call",
  "params": {
    "name": "db/query",
    "arguments": {
      "database_name": "default",
      "query": "SELECT * FROM nonexistent_table",
      "args": []
    }
  }
}'

echo ""
echo "🎉 Basic tests completed!"
echo ""
echo "To run more interactive tests:"
echo "  python3 ./test_scripts/run_interactive_test.py"
