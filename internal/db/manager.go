package db

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"
)

type Manager struct {
	Registry    *Registry // Exported for API handlers
	connections map[string]*sql.DB
	mu          sync.RWMutex
}

func NewManager(registry *Registry) *Manager {
	return &Manager{
		Registry:    registry,
		connections: make(map[string]*sql.DB),
	}
}

func (m *Manager) GetConnection(name string) (*sql.DB, error) {
	m.mu.RLock()
	db, exists := m.connections[name]
	m.mu.RUnlock()

	if exists {
		return db, nil
	}

	return m.openConnection(name)
}

func (m *Manager) openConnection(name string) (*sql.DB, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check after acquiring lock
	if db, exists := m.connections[name]; exists {
		return db, nil
	}

	info, err := m.Registry.GetDatabase(name)
	if err != nil {
		return nil, err
	}

	if !filepath.IsAbs(info.Path) {
		return nil, errors.New("database path must be absolute")
	}

	db, err := sql.Open("sqlite3", info.Path)
	if err != nil {
		return nil, err
	}

	// Configure connection
	db.SetMaxOpenConns(1) // SQLite supports only one writer
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	m.connections[name] = db

	// Update last accessed time
	if err := m.Registry.UpdateLastAccessed(info.ID); err != nil {
		// Log error but don't fail the connection
		fmt.Printf("Error updating last accessed time: %v\n", err)
	}

	return db, nil
}

func (m *Manager) CloseConnection(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if db, exists := m.connections[name]; exists {
		delete(m.connections, name)
		return db.Close()
	}
	return nil
}

func (m *Manager) CloseAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var lastErr error
	for name, db := range m.connections {
		if err := db.Close(); err != nil {
			lastErr = err
		}
		delete(m.connections, name)
	}
	return lastErr
}

func (m *Manager) ExecuteQuery(name string, query string, args ...interface{}) (*sql.Rows, error) {
	db, err := m.GetConnection(name)
	if err != nil {
		return nil, err
	}

	return db.Query(query, args...)
}

func (m *Manager) ExecuteUpdate(name string, query string, args ...interface{}) (sql.Result, error) {
	db, err := m.GetConnection(name)
	if err != nil {
		return nil, err
	}

	return db.Exec(query, args...)
}
