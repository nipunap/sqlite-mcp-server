package db

import (
	"context"
	"database/sql"
	"strings"
	"sync"
)

type BatchOperation struct {
	Database string        `json:"database"`
	Query    string        `json:"query"`
	Args     []interface{} `json:"args"`
}

type BatchResult struct {
	Database string      `json:"database"`
	Success  bool        `json:"success"`
	Results  interface{} `json:"results,omitempty"`
	Error    string      `json:"error,omitempty"`
}

func (m *Manager) ExecuteBatch(ctx context.Context, operations []BatchOperation) []BatchResult {
	results := make([]BatchResult, len(operations))
	var wg sync.WaitGroup

	// Create a buffered channel to limit concurrent operations
	semaphore := make(chan struct{}, 5) // Max 5 concurrent operations

	for i, op := range operations {
		wg.Add(1)
		go func(index int, operation BatchOperation) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := BatchResult{
				Database: operation.Database,
				Success:  false,
			}

			// Get database connection
			db, err := m.GetConnection(operation.Database)
			if err != nil {
				result.Error = err.Error()
				results[index] = result
				return
			}

			// Start transaction if needed
			tx, err := db.BeginTx(ctx, &sql.TxOptions{
				ReadOnly: false,
			})
			if err != nil {
				result.Error = err.Error()
				results[index] = result
				return
			}
			defer tx.Rollback()

			// Execute query
			rows, err := tx.QueryContext(ctx, operation.Query, operation.Args...)
			if err != nil {
				result.Error = err.Error()
				results[index] = result
				return
			}
			defer rows.Close()

			// Get column names
			columns, err := rows.Columns()
			if err != nil {
				result.Error = err.Error()
				results[index] = result
				return
			}

			// Prepare result set
			var resultSet []map[string]interface{}
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))

			for i := range columns {
				valuePtrs[i] = &values[i]
			}

			for rows.Next() {
				if err := rows.Scan(valuePtrs...); err != nil {
					result.Error = err.Error()
					results[index] = result
					return
				}

				row := make(map[string]interface{})
				for i, col := range columns {
					row[col] = values[i]
				}
				resultSet = append(resultSet, row)
			}

			if err := rows.Err(); err != nil {
				result.Error = err.Error()
				results[index] = result
				return
			}

			// Commit transaction
			if err := tx.Commit(); err != nil {
				result.Error = err.Error()
				results[index] = result
				return
			}

			result.Success = true
			result.Results = resultSet
			results[index] = result
		}(i, op)
	}

	wg.Wait()
	return results
}

type BulkInsertOperation struct {
	Database   string          `json:"database"`
	Table      string          `json:"table"`
	Columns    []string        `json:"columns"`
	Values     [][]interface{} `json:"values"`
	OnConflict string          `json:"on_conflict,omitempty"`
}

func (m *Manager) BulkInsert(ctx context.Context, operation BulkInsertOperation) (int64, error) {
	db, err := m.GetConnection(operation.Database)
	if err != nil {
		return 0, err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Build the query
	query := buildBulkInsertQuery(operation)

	// Flatten values for execution
	flatValues := make([]interface{}, 0, len(operation.Values)*len(operation.Columns))
	for _, row := range operation.Values {
		flatValues = append(flatValues, row...)
	}

	// Execute the bulk insert
	result, err := tx.ExecContext(ctx, query, flatValues...)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func buildBulkInsertQuery(operation BulkInsertOperation) string {
	// Build column list
	columns := "(" + strings.Join(operation.Columns, ", ") + ")"

	// Build placeholders for each row
	placeholders := make([]string, len(operation.Values))
	for i := range operation.Values {
		rowPlaceholders := make([]string, len(operation.Columns))
		for j := range operation.Columns {
			rowPlaceholders[j] = "?"
		}
		placeholders[i] = "(" + strings.Join(rowPlaceholders, ", ") + ")"
	}

	query := "INSERT INTO " + operation.Table + " " + columns + " VALUES " + strings.Join(placeholders, ", ")

	if operation.OnConflict != "" {
		query += " " + operation.OnConflict
	}

	return query
}
