# 使用指南

本指南面向当前版本 `go-starter`，已移除 Swagger 相关能力，保留 MySQL/Mongo 示例、Redis+cron 示例与 fx 分层注入结构。

## 1. 获取代码

```shell
git clone git@github.com:pinkhello/go-starter.git
cd go-starter
```

## 2. 环境准备

| 组件 | 用途 | 是否必需 |
|---|---|---|
| MySQL | MySQL 用户示例（`/mysql/users`） | 是 |
| Redis | cron 分布式锁与示例计数 | 是 |
| MongoDB | Mongo 用户示例（`/mongo/users`） | 建议 |

说明：
- Mongo 连接失败时，服务会降级启动（仅不注册 `/mongo/users` 路由）。
- MySQL/Redis 不可用会影响对应模块启动。

## 3. 配置项说明

主配置文件：`config/config.yml`

关键字段：

| 字段 | 说明 |
|---|---|
| `server.address` | 监听地址，例如 `:8080` |
| `contextTimeout` | Service 层超时秒数 |
| `database.*` | MySQL 连接参数 |
| `redis.*` | Redis 连接参数 |
| `mongodb.*` | Mongo 连接参数 |
| `key.type` | 认证类型：`basic` 或 `key` |
| `cron.on` | 是否启用 cron |

## 4. 启动与构建

```shell
# 整理依赖
go mod tidy

# 构建 linux 可执行文件（bin/go_starter）
make build

# 启动服务（指定配置文件）
./bin/go_starter http -c config/config.yml
```

版本命令：

```shell
./bin/go_starter version
```

## 5. 路由与鉴权

### 5.1 基础路由

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/` | 健康检查 |
| GET | `/health` | 健康检查 |
| GET | `/metrics` | Prometheus 指标 |

### 5.2 MySQL 用户接口

| 方法 | 路径 |
|---|---|
| GET | `/mysql/users/:id` |
| GET | `/mysql/users?page=&pageSize=` |
| POST | `/mysql/users` |
| PUT | `/mysql/users/:id` |
| DELETE | `/mysql/users/:id` |

### 5.3 Mongo 用户接口（可降级）

| 方法 | 路径 |
|---|---|
| GET | `/mongo/users/:id` |
| GET | `/mongo/users?page=&pageSize=` |
| POST | `/mongo/users` |
| PUT | `/mongo/users/:id` |
| DELETE | `/mongo/users/:id` |

当 Mongo 初始化失败时，上述路由不会注册。

### 5.4 鉴权行为

| `key.type` | 行为 |
|---|---|
| `basic` | 全局 BasicAuth |
| `key` | 校验请求头 `access_key`、`secret_key`，仅 `/` 与 `/health` 放行 |

## 6. 分页约束（必须遵守）

统一通过 `internal/models/vo/validate.go` 校验：

- `page > 0`
- `pageSize > 0`
- `pageSize <= 100`

MySQL 与 Mongo 列表接口均复用该规则。

## 7. 开发扩展流程

推荐按以下顺序扩展业务：

1. `internal/models/do/...` 定义模型
2. `internal/repository/...` 增加数据访问
3. `internal/service/...` 编排业务逻辑
4. `internal/controller/...` 暴露 HTTP 接口
5. 在对应 `module.go` 中注册 provider/invoke

现有示例参考：

- MySQL：`internal/models/do/mysql/example_do/user.go` → `internal/repository/mysql/example_repo/user.go` → `internal/service/example_srv/user.go` → `internal/controller/user_handler.go`
- Mongo：`internal/models/do/mongo/example_do/user.go` → `internal/repository/mongo/example_repo/user.go` → `internal/service/example_srv/user_mongo.go` → `internal/controller/user_mongo_handler.go`

## 8. 常用验证命令

```shell
# 构建验证
go build ./...

# 变更相关定向测试（推荐）
go test ./internal/service/example_srv ./internal/controller ./internal/repository/mongo/example_repo ./internal/repository/mysql/example_repo ./internal/models/vo
```

说明：仓库内存在依赖本地 MySQL 的历史测试，直接执行 `go test ./...` 可能因环境缺失失败。

## 9. Docker 运行

```shell
docker build -t go-starter:latest .
docker run -p 8080:8080 -v $(pwd)/config:/app/config go-starter:latest
```

容器默认读取 `/app/config/config.yml`（由 `config_docker.yml` 复制而来，可通过挂载覆盖）。
