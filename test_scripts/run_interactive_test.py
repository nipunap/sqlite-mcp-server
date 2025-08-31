#!/usr/bin/env python3
"""
Interactive MCP Server Test Client
"""

import json
import subprocess
import sys
import os
import tempfile
from typing import Dict, Any, Optional

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

        return self.send_message("tools/call", {
            "name": tool_name,
            "arguments": arguments
        })

def print_response(response: Dict[str, Any], test_name: str):
    """Pretty print a response"""
    print(f"\n📋 {test_name}")
    print("=" * (len(test_name) + 4))

    if "error" in response:
        print(f"❌ Error: {response['error']}")
        if "raw_output" in response:
            print(f"Raw output: {response['raw_output']}")
    else:
        print("✅ Success:")
        print(json.dumps(response, indent=2))

def main():
    """Run interactive tests"""
    print("🧪 Interactive MCP Server Testing")
    print("=================================")

    # Check if server exists
    if not os.path.exists("./sqlite-mcp-server"):
        print("❌ Server not found. Run setup script first:")
        print("   ./test_scripts/setup_test_env.sh")
        sys.exit(1)

    client = MCPTestClient()

    # Define test scenarios
    tests = [
        # Basic protocol tests
        ("Initialize Server", lambda: client.send_message("initialize", {
            "protocolVersion": "2024-11-05",
            "capabilities": {},
            "clientInfo": {"name": "interactive-test-client", "version": "1.0.0"}
        })),

        ("List Tools", lambda: client.send_message("tools/list")),
        ("List Resources", lambda: client.send_message("resources/list")),

        # Database management tests
        ("List Databases", lambda: client.call_tool("db/list_databases")),
        ("Get Tables", lambda: client.call_tool("db/get_tables", {"database_name": "default"})),
        ("Get Schema", lambda: client.call_tool("db/get_schema", {"database_name": "default"})),

        # Query tests
        ("Query All Users", lambda: client.call_tool("db/query", {
            "database_name": "default",
            "query": "SELECT * FROM users",
            "args": []
        })),

        ("Query Products by Category", lambda: client.call_tool("db/query", {
            "database_name": "default",
            "query": "SELECT * FROM products WHERE category = ?",
            "args": ["Electronics"]
        })),

        ("Complex Join Query", lambda: client.call_tool("db/query", {
            "database_name": "default",
            "query": """
                SELECT u.name as user_name, p.name as product_name, o.quantity, o.order_date
                FROM orders o
                JOIN users u ON o.user_id = u.id
                JOIN products p ON o.product_id = p.id
                ORDER BY o.order_date DESC
            """,
            "args": []
        })),

        # Insert tests
        ("Insert New User", lambda: client.call_tool("db/insert_record", {
            "database_name": "default",
            "table_name": "users",
            "data": {"name": "Interactive Test User", "email": "interactive@test.com"}
        })),

        ("Insert New Product", lambda: client.call_tool("db/insert_record", {
            "database_name": "default",
            "table_name": "products",
            "data": {"name": "Test Product", "price": 19.99, "category": "Test"}
        })),

        # Multi-database tests
        ("Register Inventory Database", lambda: client.call_tool("db/register_database", {
            "name": "inventory_test",
            "path": os.path.abspath("inventory.db"),
            "description": "Test inventory database",
            "readonly": False,
            "owner": "interactive_test"
        })),

        ("Query Inventory Database", lambda: client.call_tool("db/query", {
            "database_name": "inventory_test",
            "query": "SELECT * FROM warehouses",
            "args": []
        })),

        # Error handling tests
        ("Invalid Database Name", lambda: client.call_tool("db/query", {
            "database_name": "nonexistent",
            "query": "SELECT 1",
            "args": []
        })),

        ("Invalid SQL Query", lambda: client.call_tool("db/query", {
            "database_name": "default",
            "query": "INVALID SQL",
            "args": []
        })),
    ]

    # Run tests
    for test_name, test_func in tests:
        try:
            result = test_func()
            print_response(result, test_name)
        except Exception as e:
            print_response({"error": str(e)}, test_name)

        # Ask user if they want to continue (except for last test)
        if test_name != tests[-1][0]:
            response = input("\n⏸️  Press Enter to continue, 'q' to quit, 's' to skip remaining: ").strip().lower()
            if response == 'q':
                break
            elif response == 's':
                # Show remaining test names
                remaining = [name for name, _ in tests[tests.index((test_name, test_func)) + 1:]]
                print(f"Skipping remaining tests: {', '.join(remaining)}")
                break

    print("\n🎉 Interactive testing completed!")
    print("\nFor custom testing, you can modify the test scenarios in this script.")

if __name__ == "__main__":
    main()
