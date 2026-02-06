# 使用指南

go-starter 是一个基于 Go 的轻量级开发脚手架，采用分层架构与 fx 依赖注入，便于快速搭建 HTTP API 与定时任务。当前依赖包括：**MySQL**、**Redis**、**MongoDB**（可选）、**ClickHouse**（可选）等。

---

## 1. 获取代码

```shell
git clone git@github.com:pinkhello/go-starter.git
cd go-starter
```

---

## 2. 环境依赖

根据 `config/config.yml` 按需准备：

| 组件        | 说明           | 必选   |
|-------------|----------------|--------|
| MySQL       | 主业务库，GORM | 是     |
| Redis       | 缓存/分布式锁  | 是     |
| MongoDB     | 文档库，qmgo   | 按需   |
| ClickHouse  | 分析库         | 按需   |

本地开发时可用 Docker 自行启动 MySQL、Redis 等，或使用现有实例。配置路径：`config/config.yml`（复制并修改 `config_docker.yml` 用于容器内运行）。

---

## 3. 配置说明

- 主配置：`config/config.yml`
- Docker 运行：复制 `config/config_docker.yml` 为运行时的 `config.yml`（Dockerfile 中已使用）
- 认证方式：`key.type` 支持 `basic`（user/password）或 `key`（accessKey/secretKey），需在配置中填写对应字段

---

## 4. 本地运行

```shell
# 安装依赖
go mod tidy

# 编译（生成 bin/go_starter）
make build

# 启动 HTTP 服务（指定配置文件）
./bin/go_starter http -c config/config.yml
```

默认 HTTP 端口见配置中 `server.address`（如 `:8080`）。Swagger 文档地址：`http://{IP}:{PORT}/swagger/index.html`。

---

## 5. 项目结构

```
go-starter/
├── app/
│   ├── main.go              # 程序入口
│   └── cmd/
│       ├── root.go          # Cobra 根命令
│       ├── http.go          # HTTP 启动与 fx 注入
│       └── version.go       # 版本信息
├── config/
│   ├── config.go            # 配置结构体与加载
│   ├── config.yml           # 本地/默认配置
│   └── config_docker.yml    # Docker 环境配置
├── docs/                    # Swag 生成的 Swagger 文档
├── internal/
│   ├── controller/          # HTTP 控制器（REST 映射）
│   ├── cron/                # 定时任务（robfig/cron）
│   ├── http/                # Echo Server、中间件（日志、JWT、CORS 等）
│   ├── lib/                 # 基础设施
│   │   ├── gorm/            # GORM MySQL 适配
│   │   ├── mongodb/         # MongoDB（qmgo）
│   │   ├── redis/           # Redis 客户端
│   │   └── log/             # 日志（含 Lark）
│   ├── models/
│   │   ├── do/              # 数据对象
│   │   │   ├── mysql/       # MySQL 表模型
│   │   │   └── mongo/       # MongoDB 文档模型
│   │   └── vo/              # 请求/响应与校验
│   ├── repository/          # 数据访问层
│   │   ├── mysql/           # MySQL repo（按业务分子目录）
│   │   └── mongo/           # MongoDB repo
│   └── service/             # 业务逻辑层（按业务分子目录）
├── utils/                   # 通用工具（错误、重试、加解密、时间等）
├── vars/                    # 全局变量（版本、常量、错误码等）
├── Dockerfile               # 多阶段构建镜像
├── Makefile                 # 构建、测试、清理
└── go.mod / go.sum
```

各层中的 `module.go`（或 `moudle.go`）为 fx 模块定义，用于依赖注入，保持结构清晰。

---

## 6. 业务开发流程

按「Model → Repository → Service → Controller」顺序扩展：

1. **Model（数据模型）**
   - MySQL：`internal/models/do/mysql/{业务目录}/{表名}.go`
   - Mongo：`internal/models/do/mongo/{业务目录}/{文档名}.go`
   - 请求/响应与校验：`internal/models/vo/`

2. **Repository（数据访问）**
   - MySQL：`internal/repository/mysql/{业务目录}/{资源}.go`
   - Mongo：`internal/repository/mongo/{资源}.go`

3. **Service（业务逻辑）**
   - `internal/service/{业务目录}_srv/{资源}.go`

4. **Controller（HTTP 接口）**
   - `internal/controller/{资源}_handler.go`
   - 在 `internal/http/server.go` 或对应路由文件中注册路由

5. **注册到 fx**  
   在对应层的 `module.go`（或 `moudle.go`）中 `Provide` 新构造函数，确保被 `app/cmd/http.go` 的 `inject()` 所引用。

参考现有示例：`audit_cluster`（model → repository → service → controller）与 `task_list`（Mongo）。

---

## 7. 常用命令

```shell
# 生成 Swagger 文档
swag init -g app/main.go

# 编译
make build

# 清理
make clean

# 运行测试
go test ./...
```

---

## 8. Docker 构建与运行

```shell
# 构建镜像（使用根目录 Dockerfile，内部执行 make build）
docker build -t go-starter:latest .

# 运行（挂载配置或使用镜像内默认 config）
docker run -p 8080:8080 -v $(pwd)/config:/app/config go-starter:latest
```

镜像内默认执行：`/app/go_starter http -c /app/config/config.yml`（config 由 Dockerfile 从 `config_docker.yml` 复制）。

---

## 9. 其他

- 更多技术栈与 CI/CD 说明见 [README.md](README.md)。
- 依赖注入与命令结构见 `app/cmd/http.go` 中的 `inject()` 与 Cobra 命令定义。
