package db

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Registry struct {
	db *sql.DB
}

type DatabaseInfo struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Path         string     `json:"path"`
	Description  string     `json:"description"`
	ReadOnly     bool       `json:"readonly"`
	CreatedAt    time.Time  `json:"created_at"`
	LastAccessed *time.Time `json:"last_accessed,omitempty"`
	Owner        string     `json:"owner"`
	Status       string     `json:"status"`
}

const createRegistryTableSQL = `
CREATE TABLE IF NOT EXISTS registered_databases (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    path TEXT NOT NULL,
    description TEXT,
    readonly BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_accessed TIMESTAMP,
    owner TEXT NOT NULL,
    status TEXT CHECK(status IN ('active', 'inactive', 'error')) DEFAULT 'active'
);

CREATE TABLE IF NOT EXISTS database_metadata (
    database_id TEXT REFERENCES registered_databases(id),
    key TEXT NOT NULL,
    value TEXT,
    PRIMARY KEY (database_id, key)
);`

func NewRegistry(path string) (*Registry, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(createRegistryTableSQL); err != nil {
		db.Close()
		return nil, err
	}

	return &Registry{db: db}, nil
}

func (r *Registry) Close() error {
	return r.db.Close()
}

func (r *Registry) RegisterDatabase(info *DatabaseInfo) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO registered_databases (id, name, path, description, readonly, owner, status)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		info.ID,
		info.Name,
		info.Path,
		info.Description,
		info.ReadOnly,
		info.Owner,
		info.Status,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Registry) GetDatabase(name string) (*DatabaseInfo, error) {
	var info DatabaseInfo
	err := r.db.QueryRow(`
		SELECT id, name, path, description, readonly, created_at, last_accessed, owner, status
		FROM registered_databases
		WHERE name = ?
	`, name).Scan(
		&info.ID,
		&info.Name,
		&info.Path,
		&info.Description,
		&info.ReadOnly,
		&info.CreatedAt,
		&info.LastAccessed,
		&info.Owner,
		&info.Status,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("database not found")
	}
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (r *Registry) UpdateLastAccessed(id string) error {
	_, err := r.db.Exec(`
		UPDATE registered_databases
		SET last_accessed = CURRENT_TIMESTAMP
		WHERE id = ?
	`, id)
	return err
}

func (r *Registry) ListDatabases() ([]DatabaseInfo, error) {
	rows, err := r.db.Query(`
		SELECT id, name, path, description, readonly, created_at, last_accessed, owner, status
		FROM registered_databases
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []DatabaseInfo
	for rows.Next() {
		var info DatabaseInfo
		err := rows.Scan(
			&info.ID,
			&info.Name,
			&info.Path,
			&info.Description,
			&info.ReadOnly,
			&info.CreatedAt,
			&info.LastAccessed,
			&info.Owner,
			&info.Status,
		)
		if err != nil {
			return nil, err
		}
		databases = append(databases, info)
	}
	return databases, rows.Err()
}
