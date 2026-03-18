package sqlite

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteRepo struct {
	db *sql.DB
}

func NewSQLiteRepo(dbFileName string) (*SQLiteRepo, error) {
	db, err := sql.Open("sqlite3", dbFileName)

	if err != nil {
		return nil, err
	}

	return &SQLiteRepo{
		db: db,
	}, nil
}

func (r *SQLiteRepo) Migrate() error {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS todos (
	id INTEGER PRIMARY KEY,
	title TEXT NOT NULL,
	description TEXT,
	completed INTEGER NOT NULL DEFAULT 0,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)
	`
	_, err := r.db.ExecContext(context.Background(), sqlStmt)

	if err != nil {
		return err
	}

	return nil
}
