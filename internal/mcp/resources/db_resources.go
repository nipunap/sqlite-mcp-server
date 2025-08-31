package resources

import (
	"encoding/json"
	"fmt"

	"github.com/nipunap/sqlite-mcp-server/internal/db"
)

// DBResources provides database-related MCP resources
type DBResources struct {
	manager *db.Manager
}

// NewDBResources creates a new DBResources instance
func NewDBResources(manager *db.Manager) *DBResources {
	return &DBResources{manager: manager}
}

// GetDatabases returns a list of all registered databases
func (r *DBResources) GetDatabases() (interface{}, error) {
	databases, err := r.manager.Registry.ListDatabases()
	if err != nil {
		return nil, fmt.Errorf("registry_error: %w", err)
	}

	return map[string]interface{}{
		"databases": databases,
	}, nil
}

// GetTables returns a list of all tables for a specific database
func (r *DBResources) GetTables(params json.RawMessage) (interface{}, error) {
	var req struct {
		DatabaseName string `json:"database_name"`
	}
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, fmt.Errorf("invalid_params: %w", err)
	}

	database, err := r.manager.GetConnection(req.DatabaseName)
	if err != nil {
		return nil, fmt.Errorf("database_connection_error: %w", err)
	}

	rows, err := database.Query(`
		SELECT
			name,
			type,
			sql
		FROM sqlite_master
		WHERE type='table'
		ORDER BY name
	`)
	if err != nil {
		return nil, fmt.Errorf("db_error: %w", err)
	}
	defer rows.Close()

	var tables []map[string]interface{}
	for rows.Next() {
		var name, typ, sql string
		if err := rows.Scan(&name, &typ, &sql); err != nil {
			return nil, fmt.Errorf("db_error: %w", err)
		}

		tables = append(tables, map[string]interface{}{
			"name": name,
			"type": typ,
			"sql":  sql,
		})
	}

	return map[string]interface{}{
		"database": req.DatabaseName,
		"tables":   tables,
	}, nil
}

// GetSchema returns the full database schema for a specific database
func (r *DBResources) GetSchema(params json.RawMessage) (interface{}, error) {
	var req struct {
		DatabaseName string `json:"database_name"`
	}
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, fmt.Errorf("invalid_params: %w", err)
	}

	database, err := r.manager.GetConnection(req.DatabaseName)
	if err != nil {
		return nil, fmt.Errorf("database_connection_error: %w", err)
	}

	// Query all tables and their schemas
	rows, err := database.Query(`
		SELECT
			m.name as table_name,
			m.sql as table_sql,
			i.name as index_name,
			i.sql as index_sql
		FROM sqlite_master m
		LEFT JOIN sqlite_master i ON m.name = i.tbl_name AND i.type = 'index'
		WHERE m.type = 'table'
		ORDER BY m.name, i.name
	`)
	if err != nil {
		return nil, fmt.Errorf("db_error: %w", err)
	}
	defer rows.Close()

	schema := make(map[string]map[string]interface{})
	for rows.Next() {
		var tableName, tableSql string
		var indexName, indexSql interface{}

		if err := rows.Scan(&tableName, &tableSql, &indexName, &indexSql); err != nil {
			return nil, fmt.Errorf("db_error: %w", err)
		}

		if _, exists := schema[tableName]; !exists {
			schema[tableName] = map[string]interface{}{
				"sql":     tableSql,
				"indexes": make(map[string]string),
			}
		}

		if indexName != nil && indexSql != nil {
			schema[tableName]["indexes"].(map[string]string)[indexName.(string)] = indexSql.(string)
		}
	}

	return map[string]interface{}{
		"database": req.DatabaseName,
		"schema":   schema,
	}, nil
}
