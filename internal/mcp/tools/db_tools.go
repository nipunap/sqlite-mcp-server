package tools

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/nipunap/sqlite-mcp-server/internal/db"
)

// DBTools provides database-related MCP tools
type DBTools struct {
	manager *db.Manager
}

// NewDBTools creates a new DBTools instance
func NewDBTools(manager *db.Manager) *DBTools {
	return &DBTools{manager: manager}
}

// RegisterDatabase registers a new SQLite database
func (t *DBTools) RegisterDatabase(params json.RawMessage) (interface{}, error) {
	var req struct {
		Name        string `json:"name"`
		Path        string `json:"path"`
		Description string `json:"description,omitempty"`
		ReadOnly    bool   `json:"readonly,omitempty"`
		Owner       string `json:"owner"`
	}
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, fmt.Errorf("invalid_params: %w", err)
	}

	// Create database info
	info := &db.DatabaseInfo{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Path:        req.Path,
		Description: req.Description,
		ReadOnly:    req.ReadOnly,
		Owner:       req.Owner,
		Status:      "active",
	}

	// Register database
	if err := t.manager.Registry.RegisterDatabase(info); err != nil {
		return nil, fmt.Errorf("registration_error: %w", err)
	}

	return map[string]interface{}{
		"id":      info.ID,
		"name":    info.Name,
		"status":  "registered",
		"message": fmt.Sprintf("Database '%s' registered successfully", info.Name),
	}, nil
}

// ListDatabases lists all registered databases
func (t *DBTools) ListDatabases(params json.RawMessage) (interface{}, error) {
	databases, err := t.manager.Registry.ListDatabases()
	if err != nil {
		return nil, fmt.Errorf("registry_error: %w", err)
	}

	return map[string]interface{}{
		"databases": databases,
		"count":     len(databases),
	}, nil
}

// GetTableSchema returns the schema for a specific table
func (t *DBTools) GetTableSchema(params json.RawMessage) (interface{}, error) {
	var req struct {
		DatabaseName string `json:"database_name"`
		TableName    string `json:"table_name"`
	}
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, fmt.Errorf("invalid_params: %w", err)
	}

	// Get database connection
	database, err := t.manager.GetConnection(req.DatabaseName)
	if err != nil {
		return nil, fmt.Errorf("database_connection_error: %w", err)
	}

	// Query table schema
	rows, err := database.Query(`
		SELECT sql
		FROM sqlite_master
		WHERE type='table' AND name=?
	`, req.TableName)
	if err != nil {
		return nil, fmt.Errorf("db_error: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("table_not_found: table does not exist")
	}

	var schema string
	if err := rows.Scan(&schema); err != nil {
		return nil, fmt.Errorf("db_error: %w", err)
	}

	// Get column information
	columns, err := t.getTableColumns(database, req.TableName)
	if err != nil {
		return nil, err
	}

	// Get index information
	indexes, err := t.getTableIndexes(database, req.TableName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"table_name": req.TableName,
		"schema":     schema,
		"columns":    columns,
		"indexes":    indexes,
	}, nil
}

// InsertRecord inserts a new record into a table
func (t *DBTools) InsertRecord(params json.RawMessage) (interface{}, error) {
	var req struct {
		DatabaseName string                 `json:"database_name"`
		TableName    string                 `json:"table_name"`
		Data         map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, fmt.Errorf("invalid_params: %w", err)
	}

	// Get database connection
	database, err := t.manager.GetConnection(req.DatabaseName)
	if err != nil {
		return nil, fmt.Errorf("database_connection_error: %w", err)
	}

	// Build insert query
	columns := make([]string, 0, len(req.Data))
	values := make([]interface{}, 0, len(req.Data))
	placeholders := make([]string, 0, len(req.Data))

	for col, val := range req.Data {
		columns = append(columns, col)
		values = append(values, val)
		placeholders = append(placeholders, "?")
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		req.TableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	result, err := database.Exec(query, values...)
	if err != nil {
		return nil, fmt.Errorf("db_error: %w", err)
	}

	id, _ := result.LastInsertId()
	rows, _ := result.RowsAffected()

	return map[string]interface{}{
		"id":            id,
		"rows_affected": rows,
	}, nil
}

// ExecuteQuery executes a read-only SQL query
func (t *DBTools) ExecuteQuery(params json.RawMessage) (interface{}, error) {
	var req struct {
		DatabaseName string        `json:"database_name"`
		Query        string        `json:"query"`
		Args         []interface{} `json:"args,omitempty"`
	}
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, fmt.Errorf("invalid_params: %w", err)
	}

	// Verify query is read-only
	if !isReadOnlyQuery(req.Query) {
		return nil, fmt.Errorf("invalid_query: only SELECT queries are allowed")
	}

	// Get database connection
	database, err := t.manager.GetConnection(req.DatabaseName)
	if err != nil {
		return nil, fmt.Errorf("database_connection_error: %w", err)
	}

	// Execute query
	rows, err := database.Query(req.Query, req.Args...)
	if err != nil {
		return nil, fmt.Errorf("db_error: %w", err)
	}
	defer rows.Close()

	// Get columns
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("db_error: %w", err)
	}

	// Prepare result
	var result []map[string]interface{}
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))

	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("db_error: %w", err)
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = values[i]
		}
		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db_error: %w", err)
	}

	return map[string]interface{}{
		"columns": columns,
		"rows":    result,
	}, nil
}

// Helper functions

func (t *DBTools) getTableColumns(database *sql.DB, tableName string) ([]map[string]interface{}, error) {
	rows, err := database.Query(fmt.Sprintf("PRAGMA table_info('%s')", tableName))
	if err != nil {
		return nil, fmt.Errorf("db_error: %w", err)
	}
	defer rows.Close()

	var columns []map[string]interface{}
	for rows.Next() {
		var cid int
		var name, typ string
		var notnull, pk int
		var dfltValue interface{}

		if err := rows.Scan(&cid, &name, &typ, &notnull, &dfltValue, &pk); err != nil {
			return nil, fmt.Errorf("db_error: %w", err)
		}

		columns = append(columns, map[string]interface{}{
			"name":        name,
			"type":        typ,
			"nullable":    notnull == 0,
			"default":     dfltValue,
			"primary_key": pk == 1,
		})
	}

	return columns, rows.Err()
}

func (t *DBTools) getTableIndexes(database *sql.DB, tableName string) ([]map[string]interface{}, error) {
	rows, err := database.Query(`
		SELECT name, sql
		FROM sqlite_master
		WHERE type='index' AND tbl_name=?
	`, tableName)
	if err != nil {
		return nil, fmt.Errorf("db_error: %w", err)
	}
	defer rows.Close()

	var indexes []map[string]interface{}
	for rows.Next() {
		var name, sql string
		if err := rows.Scan(&name, &sql); err != nil {
			return nil, fmt.Errorf("db_error: %w", err)
		}

		indexes = append(indexes, map[string]interface{}{
			"name": name,
			"sql":  sql,
		})
	}

	return indexes, rows.Err()
}

func isReadOnlyQuery(query string) bool {
	query = strings.TrimSpace(strings.ToUpper(query))
	return strings.HasPrefix(query, "SELECT") || strings.HasPrefix(query, "EXPLAIN")
}
