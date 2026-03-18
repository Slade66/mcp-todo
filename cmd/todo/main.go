package main

import (
	"context"
	"fmt"

	"github.com/Slade66/mcp-todo/internal/todo/biz"
	"github.com/Slade66/mcp-todo/internal/todo/data/sqlite"
)

func main() {
	db, err := sqlite.NewSQLiteRepo("todo.db")
	if err != nil {
		panic(err)
	}
	db.Migrate()

	ctx := context.Background()

	db.Create(ctx, &biz.CreateTodoInput{
		Title:       "吃噜噜肉",
		Description: "我都好久没吃噜噜肉了，今天一定要把噜噜吃了！",
		Completed:   false,
	})

	db.Create(ctx, &biz.CreateTodoInput{
		Title:       "回家拿礼物",
		Description: "噜噜今天心情不好，得拿点小礼物收买她！",
		Completed:   false,
	})

	rows, _ := db.List(ctx, nil)
	for _, row := range rows {
		fmt.Println(row)
	}

	desc := "明天是噜噜的生日，今晚过去陪她？"

	db.Update(ctx, &biz.UpdateTodoInput{
		ID:          biz.TodoID(2),
		Description: &desc,
	})

	todo, _ := db.GetByID(ctx, biz.TodoID(2))
	fmt.Println(todo)

	db.Delete(ctx, biz.TodoID(2))

}
