#!/bin/bash

# Quick MCP Server Test Script

echo "🧪 Quick MCP Server Test"
echo "========================"

echo ""
echo "1. Testing Server Capabilities:"
echo '{"jsonrpc":"2.0","id":1,"method":"capabilities","params":{}}' | ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null | jq .

echo ""
echo "2. Testing List Databases:"
echo '{"jsonrpc":"2.0","id":2,"method":"invoke","params":{"name":"db/list_databases","params":{}}}' | ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null | jq .

echo ""
echo "3. Testing Query Users:"
echo '{"jsonrpc":"2.0","id":3,"method":"invoke","params":{"name":"db/query","params":{"database_name":"default","query":"SELECT * FROM users","args":[]}}}' | ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null | jq .

echo ""
echo "4. Testing Get Tables:"
echo '{"jsonrpc":"2.0","id":4,"method":"invoke","params":{"name":"db/get_tables","params":{"database_name":"default"}}}' | ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null | jq .

echo ""
echo "5. Testing Insert Record:"
echo '{"jsonrpc":"2.0","id":5,"method":"invoke","params":{"name":"db/insert_record","params":{"database_name":"default","table_name":"users","data":{"name":"Test User","email":"test@mcp.com"}}}}' | ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null | jq .

echo ""
echo "6. Testing Register New Database:"
INVENTORY_PATH=$(pwd)/inventory.db
echo '{"jsonrpc":"2.0","id":6,"method":"invoke","params":{"name":"db/register_database","params":{"name":"inventory","path":"'$INVENTORY_PATH'","description":"Inventory database","readonly":false,"owner":"test"}}}' | ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null | jq .

echo ""
echo "7. Testing Query Inventory Database:"
echo '{"jsonrpc":"2.0","id":7,"method":"invoke","params":{"name":"db/query","params":{"database_name":"inventory","query":"SELECT * FROM warehouses","args":[]}}}' | ./sqlite-mcp-server --registry registry.db --db test_manual.db 2>/dev/null | jq .

echo ""
echo "🎉 Quick test completed!"
