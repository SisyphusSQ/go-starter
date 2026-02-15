# go-starter

Go Web 服务脚手架，基于 Echo + fx，内置分层结构（controller/service/repository）、MySQL 与 MongoDB User CRUD 示例、Redis 分布式锁 cron 示例。

详细步骤见 [Guide.md](Guide.md)。

## 当前能力

- 提供 `http` / `version` 两个命令，默认执行 `http`
- 提供 MySQL 用户接口：`/mysql/users`
- 提供 Mongo 用户接口：`/mongo/users`（Mongo 不可用时自动降级，不注册该路由）
- 提供基础健康检查与监控：`/`、`/health`、`/metrics`
- 提供分页统一校验：`page > 0`、`pageSize > 0`、`pageSize <= 100`
- 提供 cron 示例任务（Redis 锁 + 计数落 Redis）

## 技术栈

| 类别 | 组件 | 说明 |
|---|---|---|
| Web | Echo | HTTP 框架 |
| CLI | Cobra | 命令行入口 |
| 配置 | Viper | 读取 `config.yml` |
| 依赖注入 | fx | 组件装配 |
| MySQL | GORM | 关系库访问 |
| MongoDB | qmgo | 文档库访问 |
| Redis | go-redis + redislock | 缓存与分布式锁 |
| 定时任务 | robfig/cron | 周期任务 |
| 日志 | zap | 统一日志输出 |
| 指标 | echoprometheus | `/metrics` 暴露 |

## 快速开始

```shell
git clone git@github.com:pinkhello/go-starter.git
cd go-starter
go mod tidy
make build
./bin/go_starter http -c config/config.yml
```

默认监听地址来自 `config/config.yml` 的 `server.address`（示例 `:8080`）。

## HTTP 路由

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/` | 健康检查 |
| GET | `/health` | 健康检查 |
| GET | `/metrics` | Prometheus 指标 |
| GET | `/mysql/users/:id` | MySQL 用户详情 |
| GET | `/mysql/users?page=&pageSize=` | MySQL 用户分页 |
| POST | `/mysql/users` | MySQL 创建用户 |
| PUT | `/mysql/users/:id` | MySQL 更新用户 |
| DELETE | `/mysql/users/:id` | MySQL 删除用户 |
| GET | `/mongo/users/:id` | Mongo 用户详情（降级时不可用） |
| GET | `/mongo/users?page=&pageSize=` | Mongo 用户分页（降级时不可用） |
| POST | `/mongo/users` | Mongo 创建用户（降级时不可用） |
| PUT | `/mongo/users/:id` | Mongo 更新用户（降级时不可用） |
| DELETE | `/mongo/users/:id` | Mongo 删除用户（降级时不可用） |

## 认证说明

认证类型由 `key.type` 控制：

- `basic`：使用 BasicAuth（全局生效）
- `key`：使用 `access_key` / `secret_key` 头（`/` 与 `/health` 放行）

## 项目结构

```text
go-starter/
├── app/
│   ├── main.go
│   └── cmd/
│       ├── root.go
│       ├── http.go
│       └── version.go
├── config/
│   ├── config.go
│   ├── config.yml
│   └── config_docker.yml
├── internal/
│   ├── controller/
│   ├── service/
│   │   ├── common_srv/
│   │   └── example_srv/
│   ├── repository/
│   │   ├── mysql/
│   │   └── mongo/
│   ├── models/
│   │   ├── do/
│   │   └── vo/
│   ├── cron/
│   ├── http/
│   └── lib/
├── utils/
├── vars/
├── Makefile
└── Guide.md
```

## 验证命令

```shell
go build ./...
go test ./internal/service/example_srv ./internal/controller ./internal/repository/mongo/example_repo ./internal/repository/mysql/example_repo ./internal/models/vo
```

注：仓库当前存在若干历史集成测试依赖本地 DB 环境，直接 `go test ./...` 可能因环境缺失失败。详见 [Guide.md](Guide.md)。
