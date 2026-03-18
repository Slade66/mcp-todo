package main

import (
	"github.com/Slade66/mcp-todo/internal/todo/repo/sqlite"
)

func main() {
	db, err := sqlite.NewSQLiteRepo("todo.db")
	if err != nil {
		panic(err)
	}
	db.Migrate()
}
