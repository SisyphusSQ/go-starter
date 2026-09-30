# 六平台构建、模型规范与文档分层

日期：2026-09-30；工作区 SQMC04；两仓分支 suqing/starter-modernization。追加改造未提交或推送。

## 修改范围

- Makefile 保留 build/release，并新增 build-all、release-all 和 build/release-<os>-<arch>。覆盖 windows/darwin/linux × amd64/arm64，Windows 自动使用 .exe；独立输出 bin/<os>-<arch>/，默认 CGO_ENABLED=0。
- CI 增加相同的六平台构建入口；构建目标不运行测试、不签名、不打 tag、不发布。
- docs/README.md 只做导航，设计移入 docs/design/architecture 和 docs/design/details/<topic>；参考 tuyu-studio 的架构/细节分层、codex-pulse 的按职责索引方式。
- 根 AGENTS 指定必读规范，models/AGENTS 明确依赖约束，models/do、dto、vo 的 README 始终生成。没有 examples 也保留开发入口；不创建无业务含义的空 struct。
- 沿用本仓模型语义：DO 是持久化模型，DTO 是跨服务/外部系统传输对象，VO 是 HTTP 请求和响应；业务转换由 service 完成，service 不定义应归入 models 的业务结构体。
- Go 规范涵盖 import、命名、类型归属、注释、错误、日志、context、Fx 生命周期、权限和验证；新增模块有明确实施步骤。

## 验证证据

| 验证 | 结果 | 范围 |
|---|---|---|
| go-web-starter make -j2 build-all | Pass | 六个平台实际交叉编译 |
| go-starter make -j2 build-all | Pass | 六个平台实际交叉编译 |
| 12 个产物格式与架构 | Pass | file + go version -m；Windows PE x86-64/Aarch64，macOS Mach-O x86_64/arm64，Linux ELF x86-64/aarch64；GOOS/GOARCH/CGO_ENABLED=0 匹配 |
| release-all 目标展开 | Pass | 两仓 make -n release-all；确认同一六平台构建配方加 -trimpath，本轮没有额外重复编译 release 产物 |
| 生成器与 CLI 回归 | Pass | go test ./internal/scaf_fold ./cmd；无 examples 时断言模型目录、规范与分层索引存在 |
| Markdown 本地链接 | Pass | 两仓和模板的 104 个本地链接可解析（生成模板按输出路径校验） |
| Workflow YAML | Pass | 两仓 workflow 可解析；不代表远端 CI 已运行 |
| 参考工程同步 | Pass | reference 比较 changes: 0 |
| 六平台实机运行 | Not Run | 本轮验收为交叉编译和产物元数据，未宣称 Windows/Linux/macOS 两架构均完成实机运行 |
| 远端 CI | Not Run | 未推送 |

原始日志保存在生成器的 .agents/runs/platforms-generator.log、platforms-reference.log、platform-artifacts.log、platform-release-plan.log 和 docs-generation-tests.log，不提交二进制或原始运行日志。

当前参考模板摘要：ba207f925c6b1e6bdcd81e45ed7487426a64eaffafbcc6c425b145a57001144a。所有默认生成的模型说明和规范均进入 .starter.json 管理。

## 安全与兼容性

本轮修改涉及构建和协作规范，没有新增业务网络入口、数据库操作或秘密配置。产物在 Git 忽略的 bin 目录；交叉构建关闭 CGO，未来引入 CGO 依赖时需明确各平台工具链并重新验证。

Makefile 使用 GNU Make/POSIX shell，Windows 主机需 Git Bash/MSYS2 与 GNU Make 或 WSL；Windows 目标二进制可由 macOS/Linux 交叉构建。代码规范已经纳入无业务工程，后续 agent 从根 AGENTS 与 docs/design 进入。
