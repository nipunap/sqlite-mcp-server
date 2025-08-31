/*
Package db provides database management functionality for the SQLite MCP server.

The package handles database registration, connection pooling, and query execution
for multiple SQLite databases. It includes the following main components:

  - Manager: Handles database connections and pooling
  - Registry: Manages database registration and metadata
  - Batch: Provides batch operation support

Example usage:

	// Create a new registry
	registry, err := db.NewRegistry("registry.db")
	if err != nil {
	    log.Fatal(err)
	}

	// Create database manager
	manager := db.NewManager(registry)

	// Register a database
	info := &db.DatabaseInfo{
	    Name: "mydb",
	    Path: "/path/to/db.sqlite",
	    Status: "active",
	}
	err = registry.RegisterDatabase(info)

	// Execute a query
	rows, err := manager.ExecuteQuery("mydb", "SELECT * FROM users WHERE id = ?", 1)

The package also supports batch operations and bulk inserts for efficient data manipulation.
*/
package db
