# SQLite MCP (Model Context Protocol) Server

[![CI](https://github.com/nipunap/sqlite-mcp-server/actions/workflows/ci.yml/badge.svg)](https://github.com/nipunap/sqlite-mcp-server/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/nipunap/sqlite-mcp-server/branch/main/graph/badge.svg)](https://codecov.io/gh/nipunap/sqlite-mcp-server)
[![Go Report Card](https://goreportcard.com/badge/github.com/nipunap/sqlite-mcp-server)](https://goreportcard.com/report/github.com/nipunap/sqlite-mcp-server)

A server-side implementation of the Model Context Protocol (MCP) for SQLite databases, enabling AI applications to interact with **multiple SQLite databases** through a standardized protocol. Each database must be registered before use, allowing dynamic database management and multi-database operations.

## Project Structure

```
sqlite-mcp-server/
├── cmd/
│   └── server/           # Main application entry point
├── internal/
│   ├── mcp/             # MCP implementation
│   │   ├── tools/       # Tool implementations
│   │   ├── resources/   # Resource implementations
│   │   └── prompts/     # Prompt templates
│   └── db/              # Database management
│       └── migrations/  # Database migrations
```

## Features

### Database Management Tools
- `db/register_database`: Register a new SQLite database for use
- `db/list_databases`: List all registered databases

### Database Operation Tools
- `db/get_table_schema`: Get schema for a specific table in a database
- `db/insert_record`: Insert a new record into a table
- `db/query`: Execute a read-only SQL query on a specific database
- `db/get_tables`: List all tables in a specific database
- `db/get_schema`: Get full schema of a specific database

### Resources
- `db/databases`: List of all registered databases

### Prompts
- `db/multi_database_help`: Overview of multi-database capabilities
- `db/register_help`: Help for registering databases
- `db/query_help`: Help text for constructing queries
- `db/schema_help`: Help text for understanding schemas
- `db/insert_help`: Help text for inserting records

## Prerequisites

- Go 1.21 or later
- SQLite 3

## Installation

```bash
go install github.com/nipunap/sqlite-mcp-server@latest
```

## Usage

### Starting the Server

Run as a local MCP server:

```bash
# Start with registry only (no default database)
sqlite-mcp-server --registry registry.db

# Start with registry and register a default database
sqlite-mcp-server --registry registry.db --db path/to/default.sqlite
```

The server communicates via STDIO using JSON-RPC 2.0 messages.

### Multi-Database Workflow

1. **Register a database**:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "invoke",
  "params": {
    "name": "db/register_database",
    "params": {
      "name": "users_db",
      "path": "/absolute/path/to/users.sqlite",
      "description": "User management database",
      "readonly": false,
      "owner": "app_user"
    }
  }
}
```

2. **List registered databases**:
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "invoke",
  "params": {
    "name": "db/list_databases",
    "params": {}
  }
}
```

3. **Query a specific database**:
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "invoke",
  "params": {
    "name": "db/query",
    "params": {
      "database_name": "users_db",
      "query": "SELECT * FROM users WHERE id = ?",
      "args": [1]
    }
  }
}
```

4. **Get tables from a specific database**:
```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "invoke",
  "params": {
    "name": "db/get_tables",
    "params": {
      "database_name": "users_db"
    }
  }
}
```

5. **Insert into a specific database**:
```json
{
  "jsonrpc": "2.0",
  "id": 5,
  "method": "invoke",
  "params": {
    "name": "db/insert_record",
    "params": {
      "database_name": "users_db",
      "table_name": "users",
      "data": {
        "name": "John Doe",
        "email": "john@example.com"
      }
    }
  }
}
```

## Development

1. Clone the repository:
```bash
git clone https://github.com/nipunap/sqlite-mcp-server.git
cd sqlite-mcp-server
```

2. Install dependencies:
```bash
go mod download
```

3. Run tests:
```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.