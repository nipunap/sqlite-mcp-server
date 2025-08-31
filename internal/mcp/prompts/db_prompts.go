package prompts

// DBPrompts provides database-related MCP prompts
var DBPrompts = map[string]string{
	"db/multi_database_help": `
This MCP server supports multiple SQLite databases. Each database must be registered before use.

Available Tools:
1. db/register_database - Register a new SQLite database
2. db/list_databases - List all registered databases
3. db/query - Execute SELECT queries on a specific database
4. db/get_table_schema - Get table schema from a specific database
5. db/insert_record - Insert records into a specific database

Available Resources:
1. db/databases - List all registered databases
2. db/tables - List tables in a specific database
3. db/schema - Get full schema of a specific database

All database operations require a "database_name" parameter to specify which database to use.
`,

	"db/register_help": `
To register a new SQLite database, use the db/register_database tool.

Example:
{
  "name": "my_app_db",
  "path": "/path/to/database.sqlite",
  "description": "Main application database",
  "readonly": false,
  "owner": "app_user"
}

Guidelines:
1. name: Unique identifier for the database
2. path: Absolute path to the SQLite file
3. description: Optional description
4. readonly: Set to true for read-only access
5. owner: Database owner identifier
`,

	"db/query_help": `
To query a database, use the db/query tool. This tool accepts SQL SELECT queries.

Example:
{
  "database_name": "my_app_db",
  "query": "SELECT * FROM users WHERE age > ?",
  "args": [18]
}

Guidelines:
1. database_name: Name of the registered database
2. Only SELECT queries are allowed
3. Use parameterized queries with ? placeholders
4. Provide args array for parameter values
5. Results include column names and row data
`,

	"db/schema_help": `
To understand database schemas:

1. Use db/databases resource to list all registered databases
2. Use db/tables resource to list tables in a specific database
3. Use db/schema resource to get full database schema
4. Use db/get_table_schema tool for specific table details

Table Schema Example:
{
  "database_name": "my_app_db",
  "table_name": "users"
}

Returns detailed schema information including columns, types, and indexes.
`,

	"db/insert_help": `
To insert records, use the db/insert_record tool.

Example:
{
  "database_name": "my_app_db",
  "table_name": "users",
  "data": {
    "name": "John Doe",
    "email": "john@example.com",
    "age": 25
  }
}

Guidelines:
1. database_name: Name of the registered database
2. table_name: Target table name
3. data: Object with column names as keys
4. Values must match column types
5. Returns inserted ID and rows affected
`,
}
