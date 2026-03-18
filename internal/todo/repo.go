package todo

import (
	"context"
)

// 🧠 把 id 变成“领域类型”，而不是到处裸用 string/int。
// 这样将来你要从 UUID 换成 int 或雪花 ID，只需要集中改 TodoID 及其生成/解析逻辑，业务代码基本不动。
// 这就是“类型隔离变化点”，能把改动范围压到最小。
type TodoID string

// 使业务不依赖具体数据库实现
type TodoRepo interface {
	// 创建
	Create(ctx context.Context, newTodo Todo) (Todo, error)

	// 按 ID 获取单个
	// 🧠 GetByID 比 Get/Find 的语义更明确：它直接说明“按主键 ID 查单条”，调用方不用猜条件是什么。Get 太泛，Find 常被理解为“按条件查（可能多条）”。
	GetByID(ctx context.Context, id TodoID) (Todo, error)

	// 列表查询
	List(ctx context.Context, filter Filter) ([]*Todo, error)

	// 更新（全量/部分由 input 定义）
	Update(ctx context.Context, newTodo Todo) (Todo, error)

	// 删除
	Delete(ctx context.Context, id TodoID) error
}

type Filter struct {
	// bool 的零值是 false，结构体创建之后各字段是零值，无法表达“各种状态都要”的场景
	Completed *bool
	Limit     int
	Offset    int
}
