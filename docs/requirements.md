# MCP TODO APP 需求清单（MVP）

## 1. 目标

构建一个基于 HTTP 的 Todo 应用程序，支持任务的创建、删除、列表查询、详情查询，并补充更新与完成状态变更能力。  

## 2. 数据模型

### 2.1 Todo

| 字段 | 类型 | 必填 | 默认值 | 说明 | 约束 |
| --- | --- | --- | --- | --- | --- |
| `id` | string | 是 | 服务端生成 | 任务唯一标识 | 使用 UUID |
| `title` | string | 是 | 无 | 任务标题 | 长度 1~120，不能为空 |
| `description` | string | 否 | `""` | 任务描述 | 最大长度 2000 |
| `completed` | boolean | 否 | `false` | 是否完成 | 仅允许 `true/false` |
| `created_at` | string | 是 | 服务端生成 | 创建时间 | 使用 RFC3339 |
| `updated_at` | string | 是 | 服务端生成 | 更新时间 | 使用 RFC3339；`PATCH` 更新任意字段后必须刷新 |

## 3. HTTP API

### 3.1 创建任务

- 方法与路径：`POST /todos`

- 请求体：
  ```json
  {
    "title": "Buy milk",
    "description": "2 bottles"
  }
  ```

- 成功响应：`201 Created`
  ```json
  {
    "id": "uuid",
    "title": "Buy milk",
    "description": "2 bottles",
    "completed": false,
    "created_at": "2026-03-17T10:00:00Z",
    "updated_at": "2026-03-17T10:00:00Z"
  }
  ```

### 3.2 列出任务

- 方法与路径：`GET /todos`

- 查询参数：
  - `completed`：可选，`true|false`
  - `limit`：可选，默认 20，范围 1~100
  - `offset`：可选，默认 0，最小 0

- 成功响应：`200 OK`
  ```json
  {
    "items": [
      {
        "id": "uuid",
        "title": "Buy milk",
        "description": "2 bottles",
        "completed": false,
        "created_at": "2026-03-17T10:00:00Z",
        "updated_at": "2026-03-17T10:00:00Z"
      }
    ],
    "total": 1,
    "limit": 20,
    "offset": 0
  }
  ```

### 3.3 获取任务详情

- 方法与路径：`GET /todos/{id}`

- 成功响应：`200 OK`（返回单个 Todo）

- 不存在：`404 Not Found`

### 3.4 更新任务（部分更新）

- 方法与路径：`PATCH /todos/{id}`

- 请求体（字段可选）：
  ```json
  {
    "title": "Buy milk and eggs",
    "description": "2 bottles + 12 eggs",
    "completed": true
  }
  ```

- 成功响应：`200 OK`（返回更新后的 Todo）

- 不存在：`404 Not Found`

### 3.5 删除任务

- 方法与路径：`DELETE /todos/{id}`

- 成功响应：`204 No Content`

- 不存在：`404 Not Found`

## 4. 错误响应规范

所有错误统一返回 JSON：

  ```json
  {
    "code": "INVALID_ARGUMENT",
    "message": "title is required",
    "details": {}
  }
  ```

### 4.1 错误码

- `400 Bad Request`：参数校验失败、请求体格式错误

- `404 Not Found`：资源不存在

- `500 Internal Server Error`：服务内部异常

## 5. 数据存储

- 使用 SQLite 持久化，数据库文件放在项目本地（例如 `data/todo.db`）

- 服务启动时自动建表（若不存在）

- 写操作使用事务，保证原子性

## 6. 非功能要求

- 记录基础访问日志（方法、路径、状态码、请求数据）
