package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

const ddl = `
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS projects (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT    NOT NULL UNIQUE,
    description TEXT    NOT NULL DEFAULT '',
    status      TEXT    NOT NULL DEFAULT 'Active' CHECK (status IN ('Active', 'Archived')),
    created_at  TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at  TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS api_keys (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    key_value    TEXT    NOT NULL UNIQUE,
    name         TEXT    NOT NULL,
    is_enabled   INTEGER NOT NULL DEFAULT 1,
    last_used_at TEXT,
    project_id   INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    created_at   TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_api_keys_project_id ON api_keys(project_id);
`

func Init(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, err
	}

	if _, err := db.Exec(ddl); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
