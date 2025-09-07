# Integrating SQLite MCP Server with Claude

This guide shows how to integrate your `sqlite-mcp-server` with Claude Desktop and Claude API.

## 🖥️ Claude Desktop Integration (Easiest)

### 1. Install Claude Desktop
Download Claude Desktop from [claude.ai](https://claude.ai/download)

### 2. Configure MCP Server
Create or edit Claude's MCP configuration file:

**macOS Location:**
```
~/Library/Application Support/Claude/claude_desktop_config.json
```

**Windows Location:**
```
%APPDATA%/Claude/claude_desktop_config.json
```

### 3. Add Your SQLite MCP Server
```json
{
  "mcpServers": {
    "sqlite-conversations": {
      "command": "/Users/nipuna.perera/Code/personal/sqlite-mcp-server/sqlite-mcp-server",
      "args": ["--registry", "/Users/nipuna.perera/Code/personal/ai-prompt/backend/mcp_registry.db"],
      "env": {}
    }
  }
}
```

### 4. Restart Claude Desktop
After saving the config, restart Claude Desktop. You should see MCP server status in the interface.

### 5. Test Integration
In Claude Desktop, try queries like:
- "What databases are available?"
- "Show me the conversation tables"
- "How many conversations are stored?"

## 🔧 Claude API Integration (Advanced)

### 1. Create MCP Bridge Script
```python
#!/usr/bin/env python3
"""
Claude API MCP Bridge
Connects Claude API to your SQLite MCP Server
"""

import asyncio
import json
from anthropic import Anthropic
from mcp_client import MCPClient

class ClaudeMCPBridge:
    def __init__(self, anthropic_api_key: str, mcp_server_path: str, registry_path: str):
        self.anthropic = Anthropic(api_key=anthropic_api_key)
        self.mcp_client = MCPClient(mcp_server_path, registry_path)

    async def start(self):
        """Initialize MCP connection"""
        await self.mcp_client.start_server()

    async def query_with_context(self, user_message: str) -> str:
        """Query Claude with MCP context"""

        # Check if query needs database access
        db_keywords = ["database", "conversation", "table", "sql", "query"]
        needs_db = any(keyword in user_message.lower() for keyword in db_keywords)

        context = ""
        if needs_db:
            try:
                # Get database context
                databases = await self.mcp_client.list_databases()
                context = f"Available databases: {databases}\n"

                # If asking about conversations specifically
                if "conversation" in user_message.lower():
                    tables = await self.mcp_client.get_tables("conversations.db")
                    context += f"Conversation tables: {tables}\n"

            except Exception as e:
                context = f"Database access error: {e}\n"

        # Build prompt with context
        full_prompt = f"""You have access to SQLite databases through an MCP server.

Database Context:
{context}

User Query: {user_message}

Please provide a helpful response. If you need to query the database, let me know what specific query you'd like to run."""

        # Query Claude
        response = self.anthropic.messages.create(
            model="claude-3-sonnet-20240229",
            max_tokens=1000,
            messages=[{"role": "user", "content": full_prompt}]
        )

        return response.content[0].text

# Usage example
async def main():
    bridge = ClaudeMCPBridge(
        anthropic_api_key="your-api-key-here",
        mcp_server_path="/Users/nipuna.perera/Code/personal/sqlite-mcp-server/sqlite-mcp-server",
        registry_path="/Users/nipuna.perera/Code/personal/ai-prompt/backend/mcp_registry.db"
    )

    await bridge.start()

    # Interactive loop
    while True:
        user_input = input("You: ")
        if user_input.lower() in ['quit', 'exit']:
            break

        response = await bridge.query_with_context(user_input)
        print(f"Claude: {response}")

if __name__ == "__main__":
    asyncio.run(main())
```

## 🌐 Web Interface Integration

### 1. Create Claude MCP Web Interface
```html
<!DOCTYPE html>
<html>
<head>
    <title>Claude + SQLite MCP</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .chat-container { max-width: 800px; margin: 0 auto; }
        .message { margin: 10px 0; padding: 10px; border-radius: 5px; }
        .user { background-color: #e3f2fd; text-align: right; }
        .claude { background-color: #f5f5f5; }
        .mcp-status { background-color: #e8f5e8; padding: 10px; margin-bottom: 20px; }
        input[type="text"] { width: 70%; padding: 10px; }
        button { padding: 10px 20px; }
    </style>
</head>
<body>
    <div class="chat-container">
        <h1>Claude + SQLite MCP Integration</h1>

        <div class="mcp-status" id="mcpStatus">
            MCP Server Status: <span id="statusText">Checking...</span>
        </div>

        <div id="chatMessages"></div>

        <div style="margin-top: 20px;">
            <input type="text" id="messageInput" placeholder="Ask Claude about your databases...">
            <button onclick="sendMessage()">Send</button>
        </div>

        <div style="margin-top: 20px;">
            <h3>Quick Database Queries:</h3>
            <button onclick="quickQuery('What databases are available?')">List Databases</button>
            <button onclick="quickQuery('Show conversation statistics')">Conversation Stats</button>
            <button onclick="quickQuery('What tables exist in the conversations database?')">Show Tables</button>
        </div>
    </div>

    <script>
        async function checkMCPStatus() {
            try {
                const response = await fetch('/api/mcp/status');
                const status = await response.json();
                document.getElementById('statusText').textContent =
                    status.sqlite?.running ? 'Connected ✅' : 'Disconnected ❌';
            } catch (error) {
                document.getElementById('statusText').textContent = 'Error ❌';
            }
        }

        async function sendMessage() {
            const input = document.getElementById('messageInput');
            const message = input.value.trim();
            if (!message) return;

            addMessage(message, 'user');
            input.value = '';

            try {
                const response = await fetch('/api/claude/query', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ message: message })
                });

                const result = await response.json();
                addMessage(result.response, 'claude');
            } catch (error) {
                addMessage('Error: ' + error.message, 'claude');
            }
        }

        function quickQuery(query) {
            document.getElementById('messageInput').value = query;
            sendMessage();
        }

        function addMessage(text, sender) {
            const container = document.getElementById('chatMessages');
            const div = document.createElement('div');
            div.className = `message ${sender}`;
            div.textContent = text;
            container.appendChild(div);
            container.scrollTop = container.scrollHeight;
        }

        // Check MCP status on load
        checkMCPStatus();
        setInterval(checkMCPStatus, 30000); // Check every 30 seconds

        // Enter key support
        document.getElementById('messageInput').addEventListener('keypress', function(e) {
            if (e.key === 'Enter') sendMessage();
        });
    </script>
</body>
</html>
```

## 🚀 Quick Setup Script

### 1. Create Setup Script
```bash
#!/bin/bash
# setup_claude_mcp.sh

echo "🔌 Setting up Claude MCP Integration"

# Check if Claude Desktop is installed
if [ -d "/Applications/Claude.app" ] || [ -d "$HOME/Applications/Claude.app" ]; then
    echo "✅ Claude Desktop found"

    # Create config directory
    CONFIG_DIR="$HOME/Library/Application Support/Claude"
    mkdir -p "$CONFIG_DIR"

    # Create MCP configuration
    cat > "$CONFIG_DIR/claude_desktop_config.json" << EOF
{
  "mcpServers": {
    "sqlite-conversations": {
      "command": "$(pwd)/sqlite-mcp-server/sqlite-mcp-server",
      "args": ["--registry", "$(pwd)/ai-prompt/backend/mcp_registry.db"],
      "env": {}
    }
  }
}
EOF

    echo "✅ Claude Desktop MCP configuration created"
    echo "📝 Config location: $CONFIG_DIR/claude_desktop_config.json"
    echo "🔄 Please restart Claude Desktop to apply changes"

else
    echo "❌ Claude Desktop not found"
    echo "📥 Download from: https://claude.ai/download"
fi

# Create API bridge script
cat > "claude_mcp_bridge.py" << 'EOF'
#!/usr/bin/env python3
import asyncio
import sys
import os
sys.path.append('./ai-prompt/backend')

from mcp_client import MCPClient

async def test_mcp_connection():
    """Test MCP server connection"""
    try:
        client = MCPClient(
            mcp_server_path="./sqlite-mcp-server/sqlite-mcp-server",
            registry_path="./ai-prompt/backend/mcp_registry.db"
        )

        await client.start_server()
        print("✅ MCP Server connected successfully")

        # Test basic operations
        databases = await client.list_databases()
        print(f"📊 Available databases: {databases}")

        await client.cleanup()

    except Exception as e:
        print(f"❌ MCP Connection failed: {e}")

if __name__ == "__main__":
    asyncio.run(test_mcp_connection())
EOF

chmod +x claude_mcp_bridge.py

echo "✅ Setup complete!"
echo ""
echo "🧪 Test MCP connection:"
echo "   python3 claude_mcp_bridge.py"
echo ""
echo "🔧 Next steps:"
echo "   1. Restart Claude Desktop"
echo "   2. Test database queries in Claude"
echo "   3. Try: 'What conversations do I have?'"
```

## 📋 Testing Your Integration

### 1. Test Queries for Claude
Once integrated, try these queries in Claude:

```
"What databases are available through MCP?"
"Show me the structure of the conversations table"
"How many conversations are stored in the database?"
"What are the most recent 5 conversations?"
"Analyze the conversation patterns in my database"
```

### 2. Verify MCP Connection
```bash
# Test MCP server directly
./sqlite-mcp-server/sqlite-mcp-server --registry ./ai-prompt/backend/mcp_registry.db

# Test through your existing system
./bin/start_system.sh
curl http://localhost:8000/api/mcp/status
```

## 🎯 Benefits of Claude + MCP Integration

1. **Natural Language Database Queries**: Ask Claude to analyze your data in plain English
2. **Intelligent Insights**: Claude can provide sophisticated analysis of your conversation patterns
3. **Automated Reporting**: Generate summaries and reports from your database
4. **Interactive Exploration**: Explore your data conversationally
5. **Cross-Database Analysis**: Compare data across multiple registered databases

## 🔧 Troubleshooting

### Common Issues:
1. **MCP Server Not Found**: Check executable path in config
2. **Permission Denied**: Ensure MCP server binary is executable
3. **Database Access**: Verify registry database path is correct
4. **Claude Desktop Not Updating**: Restart the application after config changes

### Debug Commands:
```bash
# Check MCP server status
ps aux | grep sqlite-mcp-server

# Test MCP server manually
echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | ./sqlite-mcp-server/sqlite-mcp-server

# Verify database access
sqlite3 ./ai-prompt/backend/mcp_registry.db ".tables"
```

---

Your SQLite MCP server is now ready to work with Claude! 🎉
