#!/bin/bash

# Debug MCP Server - Send custom JSON-RPC messages

set -e

echo "🔧 MCP Server Debug Tool"
echo "========================"

show_help() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -m, --message FILE    Send JSON message from file"
    echo "  -i, --interactive     Interactive mode - type JSON messages"
    echo "  -t, --template TYPE   Use predefined message template"
    echo "  -l, --list-templates  List available templates"
    echo "  -h, --help           Show this help"
    echo ""
    echo "Templates:"
    echo "  init         - Initialize server"
    echo "  tools        - List tools"
    echo "  resources    - List resources"
    echo "  list-db      - List databases"
    echo "  query        - Sample query"
    echo ""
    echo "Examples:"
    echo "  $0 -t init                    # Send init message"
    echo "  $0 -m my_message.json        # Send custom message"
    echo "  $0 -i                        # Interactive mode"
}

send_message() {
    local message="$1"
    local temp_file=$(mktemp)

    echo "$message" > "$temp_file"
    echo "📤 Sending message:"
    echo "$message" | jq . 2>/dev/null || echo "$message"
    echo ""
    echo "📥 Response:"

    ./sqlite-mcp-server --registry registry.db --db test_manual.db < "$temp_file" 2>/dev/null || echo "⚠️ Server error"

    rm -f "$temp_file"
}

get_template() {
    case "$1" in
        "init")
            echo '{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {},
    "clientInfo": {
      "name": "debug-client",
      "version": "1.0.0"
    }
  }
}'
            ;;
        "tools")
            echo '{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/list",
  "params": {}
}'
            ;;
        "resources")
            echo '{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "resources/list",
  "params": {}
}'
            ;;
        "list-db")
            echo '{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "invoke",
  "params": {
    "name": "db/list_databases",
    "params": {}
  }
}'
            ;;
        "query")
            echo '{
  "jsonrpc": "2.0",
  "id": 5,
  "method": "invoke",
  "params": {
    "name": "db/query",
    "params": {
      "database_name": "default",
      "query": "SELECT * FROM users LIMIT 3",
      "args": []
    }
  }
}'
            ;;
        *)
            echo "Unknown template: $1"
            echo "Use -l to list available templates"
            exit 1
            ;;
    esac
}

interactive_mode() {
    echo "🎛️  Interactive Mode"
    echo "Type JSON-RPC messages (Ctrl+D when done, 'quit' to exit)"
    echo "Example: {\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/list\",\"params\":{}}"
    echo ""

    while true; do
        echo -n "📝 Enter JSON message (or 'quit'): "
        read -r input

        if [[ "$input" == "quit" ]]; then
            break
        fi

        if [[ -n "$input" ]]; then
            echo ""
            send_message "$input"
            echo ""
        fi
    done
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -m|--message)
            if [[ -f "$2" ]]; then
                message=$(cat "$2")
                send_message "$message"
            else
                echo "❌ File not found: $2"
                exit 1
            fi
            shift 2
            ;;
        -i|--interactive)
            interactive_mode
            shift
            ;;
        -t|--template)
            message=$(get_template "$2")
            send_message "$message"
            shift 2
            ;;
        -l|--list-templates)
            echo "Available templates:"
            echo "  init      - Initialize server"
            echo "  tools     - List tools"
            echo "  resources - List resources"
            echo "  list-db   - List databases"
            echo "  query     - Sample query"
            exit 0
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# If no arguments, show help
if [[ $# -eq 0 ]]; then
    show_help
fi
