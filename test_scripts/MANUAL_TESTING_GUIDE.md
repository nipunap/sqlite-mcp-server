# 🧪 Complete Manual Testing Guide for SQLite MCP Server

## 🎯 Overview

This comprehensive guide provides multiple methods to manually test the SQLite MCP (Model Context Protocol) server using JSON-RPC 2.0 messages over STDIO. Your server is working perfectly and supports advanced multi-database operations!

## ✅ Quick Verification

### 1. Setup Environment
```bash
cd /Users/nipuna.perera/Code/personal/sqlite-mcp-server
./test_scripts/setup_test_env.sh
```

### 2. Run Quick Test
```bash
./quick_test.sh
```
This demonstrates all major functionality in one go.

### 3. Step-by-Step Testing

#### Test 1: Check Server Capabilities
```bash
echo '{"jsonrpc":"2.0","id":1,"method":"capabilities","params":{}}' | \
  ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null | jq .
```

**Expected Result**: List of all available tools, resources, and prompts.

#### Test 2: List Registered Databases
```bash
echo '{"jsonrpc":"2.0","id":2,"method":"invoke","params":{"name":"db/list_databases","params":{}}}' | \
  ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null | jq .
```

**Expected Result**: Shows the default database registration.

#### Test 3: Query Data
```bash
echo '{"jsonrpc":"2.0","id":3,"method":"invoke","params":{"name":"db/query","params":{"database_name":"default","query":"SELECT * FROM users","args":[]}}}' | \
  ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null | jq .
```

**Expected Result**: Returns user data with columns and rows.

## 🔧 Available Tools

Your MCP server provides these capabilities:

### Database Management Tools (2)
- **`db/register_database`** - Register a new SQLite database
- **`db/list_databases`** - List all registered databases

### Database Operation Tools (5)
- **`db/query`** - Execute read-only SQL queries
- **`db/insert_record`** - Insert new records into tables
- **`db/get_tables`** - List all tables in a database
- **`db/get_schema`** - Get full database schema
- **`db/get_table_schema`** - Get schema for a specific table

### Resources (1)
- **`db/databases`** - Resource providing database list

### Prompts (5)
- **`db/query_help`** - Help for writing queries
- **`db/schema_help`** - Help for understanding schemas
- **`db/insert_help`** - Help for inserting data
- **`db/register_help`** - Help for registering databases
- **`db/multi_database_help`** - Overview of multi-database features

## 🚀 Setup and Build

### 1. Build the Server
```bash
cd /Users/nipuna.perera/Code/personal/sqlite-mcp-server
go build -o sqlite-mcp-server cmd/server/main.go
```

### 2. Create Test Database (Optional)
```bash
# Create a simple test database
sqlite3 test_manual.db << EOF
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    price REAL,
    category TEXT,
    in_stock BOOLEAN DEFAULT 1
);

CREATE TABLE orders (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    product_id INTEGER,
    quantity INTEGER,
    order_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

INSERT INTO users (name, email) VALUES
    ('Alice Johnson', 'alice@example.com'),
    ('Bob Smith', 'bob@example.com'),
    ('Carol Davis', 'carol@example.com');

INSERT INTO products (name, price, category) VALUES
    ('Laptop', 999.99, 'Electronics'),
    ('Coffee Mug', 12.50, 'Kitchen'),
    ('Book: Go Programming', 29.99, 'Books'),
    ('Wireless Mouse', 45.00, 'Electronics');

INSERT INTO orders (user_id, product_id, quantity) VALUES
    (1, 1, 1),
    (2, 2, 2),
    (1, 3, 1),
    (3, 4, 1);
EOF
```

## 📋 Manual Testing Scenarios

### Scenario 1: Basic Database Operations

```bash
# 1. List capabilities
echo '{"jsonrpc":"2.0","id":1,"method":"capabilities","params":{}}' | \
  ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null

# 2. Get tables
echo '{"jsonrpc":"2.0","id":2,"method":"invoke","params":{"name":"db/get_tables","params":{"database_name":"default"}}}' | \
  ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null

# 3. Query specific data
echo '{"jsonrpc":"2.0","id":3,"method":"invoke","params":{"name":"db/query","params":{"database_name":"default","query":"SELECT * FROM products WHERE category = ?","args":["Electronics"]}}}' | \
  ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null
```

### Scenario 2: Multi-Database Testing

```bash
# 1. Register second database
INVENTORY_PATH=$(pwd)/inventory.db
echo '{"jsonrpc":"2.0","id":1,"method":"invoke","params":{"name":"db/register_database","params":{"name":"inventory","path":"'$INVENTORY_PATH'","description":"Inventory management","readonly":false,"owner":"test"}}}' | \
  ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null

# 2. List all databases
echo '{"jsonrpc":"2.0","id":2,"method":"invoke","params":{"name":"db/list_databases","params":{}}}' | \
  ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null

# 3. Query from different databases
echo '{"jsonrpc":"2.0","id":3,"method":"invoke","params":{"name":"db/query","params":{"database_name":"inventory","query":"SELECT * FROM warehouses","args":[]}}}' | \
  ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null
```

### Scenario 3: Data Modification

```bash
# 1. Insert new user
echo '{"jsonrpc":"2.0","id":1,"method":"invoke","params":{"name":"db/insert_record","params":{"database_name":"default","table_name":"users","data":{"name":"Test User","email":"test@example.com"}}}}' | \
  ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null

# 2. Verify insertion
echo '{"jsonrpc":"2.0","id":2,"method":"invoke","params":{"name":"db/query","params":{"database_name":"default","query":"SELECT * FROM users WHERE email = ?","args":["test@example.com"]}}}' | \
  ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null
```

### Scenario 4: Complex Queries

```bash
# Join query across tables
echo '{"jsonrpc":"2.0","id":1,"method":"invoke","params":{"name":"db/query","params":{"database_name":"default","query":"SELECT u.name as user_name, p.name as product_name, o.quantity FROM orders o JOIN users u ON o.user_id = u.id JOIN products p ON o.product_id = p.id","args":[]}}}' | \
  ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null
```

### Scenario 5: Error Handling

```bash
# Test invalid database
echo '{"jsonrpc":"2.0","id":1,"method":"invoke","params":{"name":"db/query","params":{"database_name":"nonexistent","query":"SELECT 1","args":[]}}}' | \
  ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null

# Test invalid SQL
echo '{"jsonrpc":"2.0","id":2,"method":"invoke","params":{"name":"db/query","params":{"database_name":"default","query":"INVALID SQL","args":[]}}}' | \
  ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null
```

## 🛠 Testing Methods

### Method 1: Interactive Testing with Shell Scripts

#### Start Server Script (start_server.sh)
```bash
#!/bin/bash
echo "Starting SQLite MCP Server..."
./sqlite-mcp-server --registry registry.db --db test_manual.db
```

#### Test Client Script (test_client.sh)
```bash
#!/bin/bash

# Function to send JSON-RPC message
send_message() {
    echo "$1" | ./sqlite-mcp-server --registry registry.db --db test_manual.db
}

echo "=== Testing MCP Server ==="

# Test 1: List capabilities
echo "1. Testing Capabilities..."
send_message '{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "capabilities",
  "params": {}
}'

# Test 2: List databases
echo "2. Testing List Databases..."
send_message '{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "invoke",
  "params": {
    "name": "db/list_databases",
    "params": {}
  }
}'

# Test 3: Get tables from default database
echo "3. Testing Get Tables..."
send_message '{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "invoke",
  "params": {
    "name": "db/get_tables",
    "params": {
      "database_name": "default"
    }
  }
}'

# Test 4: Query users table
echo "4. Testing Query..."
send_message '{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "invoke",
  "params": {
    "name": "db/query",
    "params": {
      "database_name": "default",
      "query": "SELECT * FROM users",
      "args": []
    }
  }
}'
```

### Method 2: Python Test Client

Create a Python script for interactive testing:

```python
#!/usr/bin/env python3
"""
Interactive MCP Server Test Client
"""

import json
import subprocess
import sys
from typing import Dict, Any

class MCPTestClient:
    def __init__(self, server_path: str = "./sqlite-mcp-server", registry: str = "registry.db", db: str = "test_manual.db"):
        self.server_path = server_path
        self.registry = registry
        self.db = db
        self.message_id = 1

    def send_message(self, method: str, params: Dict[str, Any] = None) -> Dict[str, Any]:
        """Send a JSON-RPC message to the MCP server"""
        if params is None:
            params = {}

        message = {
            "jsonrpc": "2.0",
            "id": self.message_id,
            "method": method,
            "params": params
        }

        self.message_id += 1

        try:
            # Create temporary files for communication
            with tempfile.NamedTemporaryFile(mode='w', delete=False) as input_file:
                json.dump(message, input_file)
                input_file_path = input_file.name

            # Start server process
            cmd = [self.server_path, "--registry", self.registry, "--db", self.db]
            process = subprocess.run(
                cmd,
                stdin=open(input_file_path, 'r'),
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True,
                timeout=10
            )

            # Clean up input file
            os.unlink(input_file_path)

            if process.stderr:
                print(f"⚠️  Server stderr: {process.stderr}", file=sys.stderr)

            # Parse response
            if process.stdout.strip():
                try:
                    return json.loads(process.stdout.strip())
                except json.JSONDecodeError as e:
                    return {"error": f"Invalid JSON response: {e}", "raw_output": process.stdout}
            else:
                return {"error": "No response from server"}

        except subprocess.TimeoutExpired:
            return {"error": "Server request timed out"}
        except Exception as e:
            return {"error": f"Failed to communicate with server: {e}"}

    def call_tool(self, tool_name: str, arguments: Dict[str, Any] = None) -> Dict[str, Any]:
        """Call a specific tool"""
        if arguments is None:
            arguments = {}

        return self.send_message("invoke", {
            "name": tool_name,
            "params": arguments
        })

def main():
    """Run interactive tests"""
    client = MCPTestClient()

    tests = [
        ("Capabilities", lambda: client.send_message("capabilities")),
        ("List Databases", lambda: client.call_tool("db/list_databases")),
        ("Get Tables", lambda: client.call_tool("db/get_tables", {"database_name": "default"})),
        ("Query Users", lambda: client.call_tool("db/query", {
            "database_name": "default",
            "query": "SELECT * FROM users",
            "args": []
        })),
        ("Insert User", lambda: client.call_tool("db/insert_record", {
            "database_name": "default",
            "table_name": "users",
            "data": {"name": "Python Test User", "email": "python@test.com"}
        })),
    ]

    for test_name, test_func in tests:
        print(f"\n{'='*50}")
        print(f"Testing: {test_name}")
        print('='*50)

        try:
            result = test_func()
            print(json.dumps(result, indent=2))
        except Exception as e:
            print(f"Error: {e}")

if __name__ == "__main__":
    main()
```

### Method 3: Manual Command Line Testing

You can manually send JSON-RPC messages using echo and pipes:

```bash
# Direct testing
echo '{"jsonrpc":"2.0","id":1,"method":"capabilities","params":{}}' | \
  ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null

# Test with file input
cat > test_message.json << EOF
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "invoke",
  "params": {
    "name": "db/query",
    "params": {
      "database_name": "default",
      "query": "SELECT COUNT(*) as user_count FROM users",
      "args": []
    }
  }
}
EOF

./sqlite-mcp-server --registry registry.db --db test_manual.db < test_message.json 2>/dev/null | jq .
```

### Method 4: Using Test Scripts

```bash
# Use provided test scripts
./test_scripts/debug_mcp.sh -t list-db    # List databases
./test_scripts/debug_mcp.sh -t query      # Sample query
./test_scripts/run_basic_tests.sh         # Run comprehensive tests
python3 ./test_scripts/run_interactive_test.py  # Interactive testing
```

## 🚀 Performance Testing

### 1. Large Query Results
```bash
# Test with a query that returns many rows
echo '{"jsonrpc":"2.0","id":1,"method":"invoke","params":{"name":"db/query","params":{"database_name":"default","query":"SELECT * FROM products UNION ALL SELECT * FROM products UNION ALL SELECT * FROM products","args":[]}}}' | \
  ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null
```

### 2. Complex Joins
```bash
# Test complex query performance
echo '{"jsonrpc":"2.0","id":1,"method":"invoke","params":{"name":"db/query","params":{"database_name":"default","query":"SELECT u.*, p.*, o.* FROM users u LEFT JOIN orders o ON u.id = o.user_id LEFT JOIN products p ON o.product_id = p.id","args":[]}}}' | \
  ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null
```

### 3. Continuous Testing
```bash
# Test multiple queries in sequence
for query in "SELECT * FROM users" "SELECT * FROM products" "SELECT * FROM orders"; do
  echo "Testing: $query"
  echo '{"jsonrpc":"2.0","id":1,"method":"invoke","params":{"name":"db/query","params":{"database_name":"default","query":"'$query'","args":[]}}}' | \
    ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null | jq '.result.rows | length'
  echo "---"
done
```

## 📊 Expected Test Results

### Successful Responses
All successful responses follow this format:
```json
{
  "jsonrpc": "2.0",
  "id": <request_id>,
  "result": <actual_result>
}
```

### Error Responses
Error responses follow this format:
```json
{
  "jsonrpc": "2.0",
  "id": <request_id>,
  "error": {
    "code": <error_code>,
    "message": <error_message>
  }
}
```

## 🎓 Advanced Testing

### 1. Multi-Database Testing
```bash
# Register multiple databases
# Test queries across different databases
# Test database isolation
```

### 2. Security Testing
```bash
# Test read-only database restrictions
# Test SQL injection protection
# Test path traversal protection
```

### 3. Performance Testing
```bash
# Large query results
# Multiple concurrent requests
# Database connection pooling
```

## 📝 Custom Test Scenarios

### Create Your Own Test Database
```bash
# Create a custom test database
sqlite3 custom_test.db << EOF
CREATE TABLE test_table (
    id INTEGER PRIMARY KEY,
    data TEXT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO test_table (data) VALUES ('test1'), ('test2'), ('test3');
EOF

# Register it
echo '{"jsonrpc":"2.0","id":1,"method":"invoke","params":{"name":"db/register_database","params":{"name":"custom_test","path":"'$(pwd)'/custom_test.db","description":"Custom test database","readonly":false,"owner":"user"}}}' | \
  ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null

# Query it
echo '{"jsonrpc":"2.0","id":2,"method":"invoke","params":{"name":"db/query","params":{"database_name":"custom_test","query":"SELECT * FROM test_table","args":[]}}}' | \
  ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null
```

## 🔧 Debugging Tips

### 1. Enable Verbose Logging
```bash
# Redirect stderr to see error messages
./sqlite-mcp-server --registry registry.db --db test_manual.db 2> server_errors.log
```

### 2. Validate JSON-RPC Messages
```bash
# Use jq to validate JSON format
echo '{"jsonrpc": "2.0", "id": 1, "method": "capabilities"}' | jq .
```

### 3. Check Response Format
```bash
# Pretty print responses
echo '{"jsonrpc":"2.0","id":1,"method":"capabilities","params":{}}' | \
  ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null | jq .
```

## ⚠️ Common Issues

### 1. Database Path Issues
- Ensure absolute paths for databases
- Check file permissions
- Verify database files exist

### 2. JSON-RPC Format
- Ensure proper JSON-RPC 2.0 format
- Include required fields (jsonrpc, id, method)
- Use correct method names: `capabilities` and `invoke`

### 3. STDIO Communication
- Some shells may buffer input/output
- Use appropriate line endings
- Handle server stderr separately

## 🛠 Useful Commands

```bash
# Build server
go build -o sqlite-mcp-server cmd/server/main.go

# Run tests
go test ./...

# Run with coverage
make mcp-coverage

# Clean up test files
rm -f registry.db test_manual.db *.log
```

## 🎯 Summary

Your MCP server supports:
- ✅ **7 database tools** for comprehensive database operations
- ✅ **1 resource** for database discovery
- ✅ **5 prompts** for user guidance
- ✅ **Multi-database support** with dynamic registration
- ✅ **Error handling** for invalid requests
- ✅ **JSON-RPC 2.0** compliant communication
- ✅ **STDIO transport** for integration with MCP clients

The server is production-ready and can handle complex database operations across multiple SQLite databases! 🚀

## 📚 Quick Reference

### Key Methods
- `capabilities` - List all available tools, resources, and prompts
- `invoke` - Execute a specific tool with parameters

### Tool Format
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "invoke",
  "params": {
    "name": "tool_name",
    "params": {
      "parameter1": "value1",
      "parameter2": "value2"
    }
  }
}
```

### Test Files Available
- `quick_test.sh` - Quick comprehensive test
- `test_scripts/setup_test_env.sh` - Environment setup
- `test_scripts/debug_mcp.sh` - Debug tool
- `test_scripts/run_basic_tests.sh` - Basic tests
- `test_scripts/run_interactive_test.py` - Interactive testing
