package biz

import (
	"context"
	"errors"
	"time"
)

// 🧠 把 id 变成“领域类型”，而不是到处裸用 string/int。
// 这样将来你要从 UUID 换成 int 或雪花 ID，只需要集中改 TodoID 及其生成/解析逻辑，业务代码基本不动。
// 这就是“类型隔离变化点”，能把改动范围压到最小。
type TodoID int

type Todo struct {
	ID          TodoID
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// 不要直接用 biz.Todo 作为 Create 入参，而是新建一个输入对象，只包含创建时需要的字段。
// 原因是 biz.Todo 表示“完整实体”，但创建时其实不该让调用方传这些字段：ID、CreatedAt、UpdatedAt
// 职责就很清楚：
// CreateTodoInput：调用方应该提供什么
// biz.Todo：系统最终保存下来的完整结果
type CreateTodoInput struct {
	Title       string
	Description string
	Completed   bool
}

// 如果要支持部分更新，字段不能是普通值类型，否则你分不清“没传”还是“传了零值”。
type UpdateTodoInput struct {
	ID          TodoID
	Title       *string
	Description *string
	Completed   *bool
}

// 使业务不依赖具体数据库实现
type TodoRepo interface {
	// 创建
	Create(ctx context.Context, newTodo *CreateTodoInput) (*Todo, error)

	// 按 ID 获取单个
	// 🧠 GetByID 比 Get/Find 的语义更明确：它直接说明“按主键 ID 查单条”，调用方不用猜条件是什么。Get 太泛，Find 常被理解为“按条件查（可能多条）”。
	GetByID(ctx context.Context, id TodoID) (*Todo, error)

	// 列表查询
	List(ctx context.Context, filter *Filter) ([]*Todo, error)

	// 更新（全量/部分由 input 定义）
	Update(ctx context.Context, newTodo *UpdateTodoInput) (*Todo, error)

	// 删除
	Delete(ctx context.Context, id TodoID) error
}

type Filter struct {
	// bool 的零值是 false，结构体创建之后各字段是零值，无法表达“各种状态都要”的场景
	Completed *bool
	Limit     int
	Offset    int
}

var ErrTodoNotFound = errors.New("todo not found")
var ErrInvalidInput = errors.New("invalid input")
