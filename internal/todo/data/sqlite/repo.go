package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Slade66/mcp-todo/internal/todo/biz"
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

func (r *SQLiteRepo) Create(ctx context.Context, newTodo *biz.CreateTodoInput) (*biz.Todo, error) {
	createStmt := `INSERT INTO todos (title, description, completed, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`

	// 建议你自己在 Go 里生成 CreatedAt 和 UpdatedAt，然后写入数据库。
	// 原因：如果用数据库生成，当你想把完整对象返回给上层，通常还得再查一遍才能拿到值。
	now := time.Now().UTC()

	res, err := r.db.ExecContext(ctx, createStmt, newTodo.Title, newTodo.Description, newTodo.Completed, now, now)
	if err != nil {
		return nil, err
	}

	// LastInsertId() 返回这次插入行的自增主键值
	// 注意：如果以后你把 id 改成 TEXT（比如 UUID、雪花字符串），LastInsertId() 就没有意义了。
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &biz.Todo{
		ID:          biz.TodoID(id),
		Title:       newTodo.Title,
		Description: newTodo.Description,
		Completed:   newTodo.Completed,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (r *SQLiteRepo) GetByID(ctx context.Context, id biz.TodoID) (*biz.Todo, error) {

	queryStmt := `SELECT id, title, description, completed, created_at, updated_at FROM todos WHERE id = ?`
	row := r.db.QueryRowContext(ctx, queryStmt, id)

	var todo biz.Todo
	if err := row.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			// 如果 return nil, nil。这会让上层很难区分“没查到”和“查到了一个空对象”。
			// 定义一个专门的错误（比如 biz.ErrTodoNotFound），让上层通过 errors.Is(err, biz.ErrTodoNotFound) 来判断。
			return nil, biz.ErrTodoNotFound
		}

		return nil, err
	}

	return &todo, nil
}

func (r *SQLiteRepo) List(ctx context.Context, filter *biz.Filter) ([]*biz.Todo, error) {

	queryStmt := `SELECT id, title, description, completed, created_at, updated_at FROM todos`

	limit := 50
	offset := 0
	var args []any

	if filter != nil {

		if filter.Completed != nil {
			queryStmt += ` WHERE completed = ?`
			args = append(args, *filter.Completed)
		}

		// 不区分“没传”和“传 0”，统一设置默认值，如果传了就设置为用户的值。
		if filter.Limit > 0 {
			limit = filter.Limit
		}

		if filter.Offset > 0 {
			offset = filter.Offset
		}
	}

	// 固定排序
	queryStmt += ` ORDER BY created_at DESC`

	// 加上 LIMIT/OFFSET
	queryStmt += ` LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, queryStmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*biz.Todo
	for rows.Next() {
		var todo biz.Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *SQLiteRepo) Update(ctx context.Context, newTodo *biz.UpdateTodoInput) (*biz.Todo, error) {

	if newTodo == nil {
		return nil, fmt.Errorf("newTodo 为 nil：%w", biz.ErrInvalidInput)
	}

	// COALESCE(title, 'default title')：如果 title 不是 NULL，就用 title
	// 否则用 'default title'。这样就能实现“只更新传了的字段，没传的字段保持原值不变”的效果。
	updateStmt := `UPDATE todos SET title = COALESCE(?, title), description = COALESCE(?, description), completed = COALESCE(?, completed), updated_at = ? WHERE id = ?`

	now := time.Now().UTC()

	res, err := r.db.ExecContext(ctx, updateStmt, newTodo.Title, newTodo.Description, newTodo.Completed, now, newTodo.ID)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, biz.ErrTodoNotFound
	}

	// 为什么最后还要 GetByID：
	// 你接口约定要返回 *Todo
	// UPDATE 本身不会把更新后的整行数据直接给你
	// 你还需要拿到数据库里的最终状态，比如：
	// 1. 数据库层可能做过的默认值/转换
	// 2. 还有 CreatedAt 这些字段并不通过参数传入
	// 这不是多余，而是在用一次查询换取“返回值一定准确”。
	return r.GetByID(ctx, newTodo.ID)
}

func (r *SQLiteRepo) Delete(ctx context.Context, id biz.TodoID) error {
	deleteStmt := "DELETE FROM todos WHERE id = ?"

	res, err := r.db.ExecContext(ctx, deleteStmt, id)
	if err != nil {
		return err
	}
	
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return biz.ErrTodoNotFound
	}

	return nil
}
