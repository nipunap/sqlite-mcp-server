#!/bin/bash
# setup_claude_mcp.sh - Setup Claude Desktop MCP Integration

set -e

echo "🔌 Setting up Claude MCP Integration"
echo "=================================="

# Get absolute paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
MCP_SERVER_PATH="$(cd "$PROJECT_ROOT/../sqlite-mcp-server" && pwd)/sqlite-mcp-server"
REGISTRY_PATH="$PROJECT_ROOT/backend/mcp_registry.db"

echo "📁 Project root: $PROJECT_ROOT"
echo "🔧 MCP server: $MCP_SERVER_PATH"
echo "📊 Registry: $REGISTRY_PATH"

# Check if MCP server exists
if [ ! -f "$MCP_SERVER_PATH" ]; then
    echo "❌ SQLite MCP server not found at: $MCP_SERVER_PATH"
    echo "🔧 Please build it first:"
    echo "   cd ../sqlite-mcp-server && go build -o sqlite-mcp-server ./cmd/server/"
    exit 1
fi

echo "✅ SQLite MCP server found"

# Check if registry exists
if [ ! -f "$REGISTRY_PATH" ]; then
    echo "❌ MCP registry not found at: $REGISTRY_PATH"
    echo "🔧 Please run the AI Prompt system first to create the registry:"
    echo "   ./bin/start_system.sh"
    exit 1
fi

echo "✅ MCP registry found"

# Check if Claude Desktop is installed
CLAUDE_INSTALLED=false
CONFIG_DIR=""

if [ -d "/Applications/Claude.app" ]; then
    CLAUDE_INSTALLED=true
    CONFIG_DIR="$HOME/Library/Application Support/Claude"
elif [ -d "$HOME/Applications/Claude.app" ]; then
    CLAUDE_INSTALLED=true
    CONFIG_DIR="$HOME/Library/Application Support/Claude"
fi

if [ "$CLAUDE_INSTALLED" = true ]; then
    echo "✅ Claude Desktop found"

    # Create config directory
    mkdir -p "$CONFIG_DIR"

    # Backup existing config if it exists
    if [ -f "$CONFIG_DIR/claude_desktop_config.json" ]; then
        echo "📋 Backing up existing Claude config..."
        cp "$CONFIG_DIR/claude_desktop_config.json" "$CONFIG_DIR/claude_desktop_config.json.backup.$(date +%s)"
    fi

    # Create MCP configuration
    cat > "$CONFIG_DIR/claude_desktop_config.json" << EOF
{
  "mcpServers": {
    "sqlite-conversations": {
      "command": "$MCP_SERVER_PATH",
      "args": ["--registry", "$REGISTRY_PATH"],
      "env": {}
    }
  }
}
EOF

    echo "✅ Claude Desktop MCP configuration created"
    echo "📝 Config location: $CONFIG_DIR/claude_desktop_config.json"
    echo ""
    echo "🔄 IMPORTANT: Please restart Claude Desktop to apply changes"

else
    echo "❌ Claude Desktop not found"
    echo "📥 Download from: https://claude.ai/download"
    echo ""
    echo "🔧 Manual setup instructions:"
    echo "   1. Install Claude Desktop"
    echo "   2. Create config file at: ~/Library/Application Support/Claude/claude_desktop_config.json"
    echo "   3. Add this configuration:"
    echo ""
    cat << EOF
{
  "mcpServers": {
    "sqlite-conversations": {
      "command": "$MCP_SERVER_PATH",
      "args": ["--registry", "$REGISTRY_PATH"],
      "env": {}
    }
  }
}
EOF
fi

# Create test script
echo ""
echo "🧪 Creating MCP connection test script..."

cat > "$PROJECT_ROOT/test_claude_mcp.py" << 'EOF'
#!/usr/bin/env python3
"""
Test Claude MCP Integration
Tests the connection between Claude and SQLite MCP server
"""

import asyncio
import sys
import os
import json

# Add backend to path
sys.path.append(os.path.join(os.path.dirname(__file__), 'backend'))

try:
    from mcp_client import MCPClient
except ImportError as e:
    print(f"❌ Cannot import MCP client: {e}")
    print("Make sure you're running from the ai-prompt directory")
    sys.exit(1)

async def test_mcp_connection():
    """Test MCP server connection for Claude integration"""
    print("🧪 Testing MCP Connection for Claude Integration")
    print("=" * 50)

    try:
        # Get paths
        script_dir = os.path.dirname(os.path.abspath(__file__))
        mcp_server_path = os.path.join(script_dir, "../sqlite-mcp-server/sqlite-mcp-server")
        registry_path = os.path.join(script_dir, "backend/mcp_registry.db")

        print(f"🔧 MCP Server: {mcp_server_path}")
        print(f"📊 Registry: {registry_path}")

        # Check files exist
        if not os.path.exists(mcp_server_path):
            print(f"❌ MCP server not found: {mcp_server_path}")
            return False

        if not os.path.exists(registry_path):
            print(f"❌ Registry not found: {registry_path}")
            return False

        # Test connection
        client = MCPClient(
            mcp_server_path=mcp_server_path,
            registry_path=registry_path
        )

        print("🔌 Starting MCP server...")
        await client.start_server()
        print("✅ MCP Server connected successfully")

        # Test basic operations that Claude will use
        print("\n📋 Testing MCP operations...")

        # List databases
        try:
            databases = await client.list_databases()
            print(f"✅ List databases: {len(databases)} found")
            for db in databases[:3]:  # Show first 3
                print(f"   - {db}")
        except Exception as e:
            print(f"❌ List databases failed: {e}")

        # Test conversation database if it exists
        try:
            tables = await client.get_tables("conversations.db")
            print(f"✅ Get tables: {len(tables)} tables in conversations.db")
            for table in tables:
                print(f"   - {table}")
        except Exception as e:
            print(f"⚠️  Conversations database not accessible: {e}")

        # Test a simple query
        try:
            result = await client.query_database("conversations.db", "SELECT COUNT(*) as count FROM conversations")
            if result and 'rows' in result and result['rows']:
                count = result['rows'][0].get('count', 0)
                print(f"✅ Query test: {count} conversations found")
            else:
                print("⚠️  Query returned no data")
        except Exception as e:
            print(f"⚠️  Query test failed: {e}")

        await client.cleanup()

        print("\n🎉 MCP Integration Test Complete!")
        print("\n📋 Next Steps:")
        print("1. Restart Claude Desktop if you haven't already")
        print("2. In Claude, try asking: 'What databases are available?'")
        print("3. Try: 'How many conversations are in my database?'")
        print("4. Try: 'Show me the structure of my conversation data'")

        return True

    except Exception as e:
        print(f"❌ MCP Connection failed: {e}")
        import traceback
        traceback.print_exc()
        return False

if __name__ == "__main__":
    success = asyncio.run(test_mcp_connection())
    sys.exit(0 if success else 1)
EOF

chmod +x "$PROJECT_ROOT/test_claude_mcp.py"

echo "✅ Test script created: $PROJECT_ROOT/test_claude_mcp.py"

# Create Claude query examples
cat > "$PROJECT_ROOT/claude_query_examples.md" << 'EOF'
# Claude + MCP Query Examples

Once you've set up Claude Desktop with MCP integration, try these example queries:

## 🗄️ Database Exploration
- "What databases are available through MCP?"
- "Show me all the tables in my conversation database"
- "What's the structure of the conversations table?"

## 📊 Data Analysis
- "How many conversations are stored in my database?"
- "What are the most recent 5 conversations?"
- "Show me conversation statistics by date"
- "Which sessions have the most messages?"

## 🔍 Content Queries
- "Find conversations that mention 'python' or 'programming'"
- "What topics do I discuss most frequently?"
- "Show me conversations from the last week"
- "Which conversations received positive feedback?"

## 📈 Advanced Analysis
- "Analyze my conversation patterns over time"
- "What's the average length of my conversations?"
- "Show me the distribution of conversation topics"
- "Generate a summary of my most active discussion themes"

## 🛠️ Technical Queries
- "What's the schema of my session_summaries table?"
- "Show me any conversations with feedback data"
- "List all unique session IDs from the last month"
- "What's the size and structure of my conversation database?"

## 💡 Tips for Better Results
1. Be specific about what data you want to see
2. Ask for summaries rather than raw data dumps
3. Request analysis and insights, not just data retrieval
4. Use natural language - Claude understands context well
5. Ask follow-up questions to dive deeper into interesting findings
EOF

echo "✅ Query examples created: $PROJECT_ROOT/claude_query_examples.md"

echo ""
echo "🎯 Setup Summary:"
echo "=================="
echo "✅ MCP server verified"
echo "✅ Registry database verified"
if [ "$CLAUDE_INSTALLED" = true ]; then
    echo "✅ Claude Desktop configuration created"
else
    echo "⚠️  Claude Desktop not found - manual setup required"
fi
echo "✅ Test script created"
echo "✅ Query examples created"

echo ""
echo "🧪 Test your setup:"
echo "   python3 $PROJECT_ROOT/test_claude_mcp.py"
echo ""
echo "📖 View query examples:"
echo "   cat $PROJECT_ROOT/claude_query_examples.md"

if [ "$CLAUDE_INSTALLED" = true ]; then
    echo ""
    echo "🔄 RESTART Claude Desktop now to activate MCP integration!"
fi
