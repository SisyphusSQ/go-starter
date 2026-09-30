# SQL 文件规范与最终范围

日期：2026-09-30；两仓位于 suqing/starter-modernization，未提交或推送。

## 用户确认的范围

- 基底和 examples 不附带业务 DDL、建表 SQL 或种子数据。
- 统一在 docs/sqls 保存规则：schema 为当前完整结构，unreleased 为待交付 SQL，releases/vX.Y.Z 为正式版本交付档案。
- 默认只生成 README，不创建 SQL 文件或虚构的版本目录；后续实际业务需要时再补 SQL。
- 新环境使用匹配代码版本的完整 schema；应用启动、构建和文件归档不自动执行 SQL。
- 文件命名、顺序、前置条件、验证、版本归档、不可变历史和独立执行证据已写入 docs/sqls/README.md，并链接根 AGENTS 与模型/开发规范。

## 参考与取舍

只读检查了 SQMC03 的 BaaS 和 DBIO 项目 docs/sqls，采用其“完整结构、未发布区、按真实版本归档、归档不代替执行证据”原则；当前基底简化为 schema/unreleased/releases 三类目录，不引入执行框架。

## 验证

| 项目 | 结果 | 范围 |
|---|---|---|
| 生成器与 CLI 单元回归 | Pass | go test ./internal/scaf_fold ./cmd |
| 受影响生成组合 | Pass | standalone、mysql-example、mysql-only、full-example：tidy、race、vet、build，集成标签仅编译 |
| 参考工程检查 | Pass | make verify lint；integration 标签仅编译；lint 0 issues |
| CLI 命令 | Pass | --help 只保留现有 HTTP、version 和 Cobra 帮助入口 |
| 两仓六平台构建 | Pass | make -j2 build-all；12 个产物目标架构核对 |
| 文档与分发边界 | Pass | 127 个本地链接解析，参考清单不含 SQL 文件或已删除入口 |
| 参考仓同步 | Pass | reference changes: 0 |
| 数据库执行与实机平台运行 | Not Run | 本轮不以编译测试声称真实 SQL 执行或六平台运行已通过 |

原始证据在生成器本地 .agents/runs/sql-layout-*.log。SQL 测试仅在显式隔离环境里创建随机命名的测试表并清理；本轮没有连接数据库执行 SQL。
