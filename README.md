# go-starter

[![GitHub stars](https://img.shields.io/github/stars/pinkhello/go-starter?color=0088ff)](https://github.com/pinkhello/go-starter) [![GitHub forks](https://img.shields.io/github/forks/pinkhello/go-starter?color=0088ff)](https://github.com/pinkhello/go-starter) [![GitHub issues](https://img.shields.io/github/issues/pinkhello/go-starter?color=0088ff)](https://github.com/pinkhello/go-starter)

基于 Go 的轻量级 HTTP API 脚手架，分层清晰、依赖注入，便于快速扩展业务。详细使用步骤见 [Guide.md](Guide.md)。

---

## 技术栈

| 类别       | 组件 | 说明 |
|------------|------|------|
| Web        | [Echo](https://github.com/labstack/echo) | HTTP 框架 |
| CLI        | [Cobra](https://github.com/spf13/cobra) | 命令行 |
| 配置       | [Viper](https://github.com/spf13/viper) | YAML/环境变量 |
| 依赖注入   | [fx](https://github.com/uber-go/fx) | Uber 出品 DI/IOC |
| ORM        | [GORM](https://github.com/go-gorm/gorm) | MySQL 等 |
| 文档库     | [qmgo](https://github.com/qiniu/qmgo) | MongoDB |
| 缓存       | [go-redis](https://github.com/redis/go-redis) | Redis 客户端 |
| 定时任务   | [cron](https://github.com/robfig/cron/v3) | 定时任务 |
| 日志       | [zap](https://github.com/uber-go/zap) + [logrus](https://github.com/sirupsen/logrus) | 日志与 Lark 通知 |
| 监控       | [Prometheus](https://github.com/prometheus/client_golang) | 指标 |
| API 文档   | [swag](https://github.com/swaggo/swag) + [echo-swagger](https://github.com/swaggo/echo-swagger) | Swagger 2.0 |
| 其他       | golangci-lint、Github Actions 等 | 代码检查与 CI |

---

## 项目结构

```
go-starter/
├── app/                    # 应用入口
│   ├── main.go
│   └── cmd/                # Cobra 子命令
│       ├── root.go
│       ├── http.go         # HTTP 服务启动与 fx 注入
│       └── version.go
├── config/                 # 配置
│   ├── config.go
│   ├── config.yml
│   └── config_docker.yml
├── docs/                   # Swag 生成的 Swagger 文档
├── internal/               # 核心业务
│   ├── controller/         # HTTP 控制器
│   ├── cron/               # 定时任务
│   ├── http/               # Echo 服务与中间件
│   ├── lib/                # 基础设施
│   │   ├── gorm/           # GORM MySQL
│   │   ├── mongodb/        # MongoDB
│   │   ├── redis/          # Redis
│   │   └── log/            # 日志（含 Lark）
│   ├── models/             # 数据模型
│   │   ├── do/             # MySQL / Mongo 实体
│   │   └── vo/             # 请求/响应与校验
│   ├── repository/         # 数据访问层
│   │   ├── mysql/
│   │   └── mongo/
│   └── service/            # 业务逻辑层
├── utils/                  # 通用工具
├── vars/                   # 全局变量、错误码等
├── Dockerfile              # 多阶段构建
├── Makefile
├── go.mod / go.sum
├── Guide.md                # 使用指南
└── README.md
```

---

## 快速开始

### 环境要求

- Go 1.24+
- 按需准备：MySQL、Redis、MongoDB（见 `config/config.yml`）

### 本地运行

```shell
git clone git@github.com:pinkhello/go-starter.git
cd go-starter
go mod tidy
make build
./bin/go_starter http -c config/config.yml
```

### Swagger

```shell
swag init -g app/main.go
```

访问：`http://{IP}:{PORT}/swagger/index.html`

### 健康检查

```text
GET http://{IP}:{PORT}/
```

---

## 构建与部署

### 本地构建

```shell
make build
# 产物：bin/go_starter
```

### Docker 构建与运行

```shell
# 构建镜像（Dockerfile 在项目根目录）
docker build -t go-starter:latest .

# 运行（可按需挂载 config）
docker run -p 8080:8080 go-starter:latest
```

镜像内使用 `config_docker.yml` 复制得到的 `config.yml`，需保证 MySQL、Redis 等对容器可达（同机或网络互通）。

### GitHub Actions

若使用 CI 推送镜像，在仓库 Secrets 中配置 `ACCESS_USERNAME`（如 Docker Hub 用户名）等所需变量。

---

## 依赖注入（fx）

HTTP 服务通过 `app/cmd/http.go` 中的 `inject()` 组装：

```go
func inject() fx.Option {
	return fx.Options(
		fx.Provide(
			config.NewConfig,
			utils.NewTimeoutContext,
		),
		libs.GlobalModule,   // GORM / MongoDB / Redis / Log
		repository.Module,
		service.Module,
		cron.Module,
		controller.Module,
		http.Module,
	)
}
```

新增数据源或业务模块时，在对应 `internal/*/module.go` 中 `Provide`，并在上述 `inject()` 中挂载对应 Module 即可。

---

## 其他

- **使用与开发流程**：见 [Guide.md](Guide.md)。
- **License**：见 [LICENSE](LICENSE)。
