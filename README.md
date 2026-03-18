
### 项目介绍

本项目是一个简单的 Todo 应用，用于实践所学的 MCP 知识。

项目的最终目标是通过对话，让 AI 借助 MCP 帮我管理 Todo（新增、删除、查询）。

### 目录结构

按业务聚合，而不是按类型分层。相关的代码放在一起。

### 编译依赖

`go get github.com/mattn/go-sqlite3` 库编译时依赖 GCC，需要先装 `gcc`，并 `go env -w CGO_ENABLED=1`。
