## Unreleased

## v2.0.0(20260930)

#### feature:

1. 新增 docs/sqls/schema、unreleased、releases 的 SQL 文件规范，明确版本归档及执行证据要求。

2. 新增 Windows、macOS、Linux 的 amd64/arm64 六平台构建及独立产物目录；保留本机构建入口。
3. 文档按 docs/design/architecture 与 docs/design/details/<topic> 分层，固化 DO/DTO/VO、Go 代码规范和新增模块步骤；默认工程保留模型目录说明。

4. 作为 go-web-starter 的固定参考工程，通过 reference 命令生成、比较和同步，记录模板与 Harness 来源。
5. 预置 Harness v0.7.0、架构与开发文档、Compose 及隔离集成测试入口。
6. 新增受鉴权保护的 metrics、请求关联日志、统一错误响应、健康与就绪探针及同库事务入口。

#### optimization:


1. 升级至 Go 1.27.1、Echo v5 和官方 MongoDB Driver v2，保留 Uber Fx，更新全部保留的 Go 依赖。
2. 外部组件默认关闭，按配置装配和管理生命周期；开箱启动不要求数据库或 Redis。
3. 默认工程移除 User CRUD、示例数据和定时任务；完整示例通过生成器 --examples 显式获取。

#### note:

- 模块路径采用 /v2；version 与 --version 使用同一版本来源，正式产物固定注入 v2.0.0。

1. 本基底用于新项目初始化；配置、响应体和示例目录以当前文档为准。
2. JWT 要求启用 Redis 和独立 namespace。

## v1.0.0(20260215)

#### feature:

1. 新增用户接口响应体定义 `UserListResp`、`UserMongoListResp`、`UserIDResp`、`UserMongoIDResp`，统一列表与主键返回结构
2. 新增 MySQL 与 Mongo 用户 Service 直接组装响应体的返回逻辑，List/Update/Delete 由 Service 返回业务 Resp，Controller 仅做参数处理与透传
3. 新增 `docs/schema/users_example_insert.sql` 的 `password` 字段示例数据，支持账号登录联调

#### optimization:

1. 优化 MySQL 与 Mongo 用户 Handler，移除 `map[string]any` 形式返回，改为显式 VO 类型返回，提升接口契约清晰度
2. 优化相关单元测试与 Mock 签名，适配 Service 返回类型调整，保证测试用例与当前接口一致
