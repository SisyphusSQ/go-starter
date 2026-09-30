# Go 基底现代化：实现与验证记录

日期：2026-09-30。当前开发机器 SQMC04，两仓均在 `suqing/starter-modernization`；变更尚未 commit、push 或发布。

- 生成器：`go-web-starter`，基线 `bcbc29c`，模板为权威来源。
- 参考工程：`go-starter`，基线 `5496e63`，固定包含全部可选组件，默认关闭外部连接且不带业务示例。
- Go 统一为 1.27.1；继续使用 Uber Fx，不引入 Wire。
- Harness：v0.7.0，来源提交 `b20e5e8ece6a529c7d74aa0a8b1de77bd06c374c`。
- 首轮改造模板摘要：`5a11d9e113a36b9833949e6618e6525fa88117657e34f1fc3311d54de7a1b3f9`。`.starter.json` 保存生成选项与受管理文件 SHA-256。

## 已实现范围

1. 默认 `--db none`，显式选择 MySQL/MongoDB；`--with` 选择 Redis、cron、Lark、Prometheus query、JWT，JWT 要求 Redis。
2. `--examples` 才生成 User CRUD 代码；默认无示例业务、无预置账号、无演示任务。参考工程包含组件能力，运行时需显式 enabled。
3. Fx 构造阶段不连接外部服务；OnStart 验证连接和监听，失败清理、退出释放；默认工程可独立启动。
4. 请求 ID、结构化日志、统一错误与参数校验；健康/就绪区分，metrics 沿用鉴权且正确统计错误与鉴权拒绝。
5. JWT 校验、命名空间隔离、独立令牌 ID 和 Redis 撤销；SQL 同库事务。
6. docs、AGENTS、计划/state/runs 模板、MIT 来源、Compose、CI 和固定版本检查工具已预置。
7. reference 默认比较，显式 --write 同步；本地修改/删除、未管理文件碰撞及符号链接均受到保护。

## 依赖

全部保留的 Go module 版本按实施时查询结果固定，并检查生成后的间接依赖；普通生成不联网获取 latest。代表版本：Echo v5.4.0、MongoDB Driver v2.9.1、Redis v9.22.0、Lark SDK v3.12.0、Fx v1.24.0、Zap v1.28.0、Prometheus client_golang v1.24.1。移除 qmgo、golib 和旧 module path，完整版本以 go.mod 和模板为准。

## 验证结果

| 范围 | 结果 | 证据与边界 |
|---|---|---|
| 生成器格式、race、vet、build、lint | Pass | `make verify test lint`，lint 0 issues |
| 参考工程格式、race、vet、build、lint | Pass | `make verify lint`，lint 0 issues；审查修复后针对受影响代码复验 |
| 十种生成组合 | Pass | standalone、redis-only、jwt-redis、optional-sdk、mysql-example、mongo-example、full-example、mysql-only、mongodb-only、mysql-and-mongodb；产物 tidy/格式/race/vet/build，integration 标签仅编译 |
| 审查修复回归 | Pass | 事务跨库上下文、CORS 通配域名、MySQL 原值更新、metrics 状态记录先复现再修复；受影响生成组合复验 |
| 真实 MySQL 8.0.46 | Pass | SQMC03 独立临时实例，事务回滚、原值更新匹配、就绪与关闭；同架构测试二进制在 SQMC04 编译后执行，未在远端重新解析 module |
| 真实 Redis 8.6.2 | Pass | SQMC04 独立临时实例，真实读写、就绪、JWT 生效与撤销、关闭 |
| 临时资源清理 | Pass | 上述自建进程、数据和凭据已清理；未访问业务数据库 |
| 模板与参考仓一致性 | Pass | reference 比较退出 0，changes: 0；独立 ChangeLog、LICENSE 和未管理文档保留 |
| Workflow/Compose YAML 语法 | Pass | Ruby Psych 解析；不等同容器执行 |
| Docker 镜像/Compose 运行 | Not Run | 已检查机器无可用 Docker；非 root 镜像、健康检查与 CI 入口已配置 |
| 真实 MongoDB / 软删除集成用例 | Not Run | 无可用 MongoDB 实例；Driver v2 编译和生命周期构造测试通过，真实读写/软删除用例已纳入隔离测试 |
| GitHub Actions | Not Run | 尚未 push，不能把 workflow 文件当作已通过的 CI |
| 飞书发送 / 业务环境验收 | Not Run | 本次未调用真实消息接口或业务服务 |

原始日志保存在生成器本地 `.agents/runs/`，被 Git 忽略；关键记录为 `generator-final.log`、`reference-final.log`、`matrix-final.log`、`matrix-review-fixes.log`、`matrix-metrics-fix.log`、`mysql-final-live.log`、`redis-jwt-live.log`、`metrics-green.log` 和 `reference-check.log`。

## 漏洞检查与安全自查

生成器 govulncheck 未发现漏洞。参考工程为 0 个可达漏洞、0 个已导入包漏洞，另有 1 个 module 级告警：[GO-2026-5932](https://pkg.go.dev/vuln/GO-2026-5932)，指向 `golang.org/x/crypto/openpgp`；当前项目未导入该包，数据库未提供修复版本。不能将此结果描述为整个依赖图零告警。

本次审查已覆盖鉴权和探针边界、CORS、外部请求的可信配置约束、秘密与日志、文件写入范围、资源释放及事务语义。修复了通配 CORS、跨库事务逃逸、原值更新误判、metrics 失真和软删除后的更新问题。示例只提供认证层保护，不是正式资源授权或多租户系统；正式业务须在服务/查询层实现实际授权。

## 兼容性与剩余验证

- cron 为进程内调度与防重入；多副本协调需明确设计，已移除无实现的 lockTTL 配置。
- 后续具备容器环境时，按 docs/test/integration.md 运行 Docker、MongoDB 和 MySQL 8.4 验证，并读取实际 CI 结果。
