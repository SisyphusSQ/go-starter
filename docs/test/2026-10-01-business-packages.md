# v2.0.1 业务子包规范开发验证

日期：2026-10-01。范围为 go-web-starter 模板及固定 go-starter 参考工程；本记录只描述本轮开发阶段实际完成的验证。

## 变更与验收

- controller/service/repository 及业务 DO/DTO/VO 按业务域分子包；单文件业务同样分包，根包仅承担装配和明确公共能力。
- MySQL DO 业务子包内一表一文件，DO、TableName() 和列映射同文件；共享类型例外不允许汇集多表 DO。
- 示例用户 VO 迁至 models/vo/example_vo；Lark DTO 迁至 models/dto/lark_dto。同步调用方、生成开关和 reference 清单，无 examples/无 Lark 时不残留对应业务子包。
- 更新架构与开发索引、根/各层 AGENTS、模型目录说明与新增模块步骤，固定参考版本和生成来源为 v2.0.1。

## 已完成的开发验证

| 入口 | 实际结果与范围 |
|---|---|
| go-web-starter：make test | Pass；生成器、CLI 与同步保护测试，含 race |
| go-web-starter：make verify lint | Pass；格式、vet、构建；golangci-lint 为 0 issues |
| go-web-starter：make integration | Pass；10 个生成组合，对产物运行 race/test/vet/build，并编译 integration 标签 |
| go-starter：make verify lint | Pass；参考工程格式、单元测试/race、vet、构建；golangci-lint 为 0 issues |
| 本地文档链接读取检查 | Pass；125 个相对链接均有目标 |

生成矩阵包括 standalone、redis-only、jwt-redis、optional-sdk、mysql-example、mongo-example、full-example、mysql-only、mongodb-only、mysql-and-mongodb。新增断言检查业务 VO/DTO 的子包位置，以及关闭 examples/Lark 后不存在旧平铺文件和空业务子包。

开发原始日志保留在两个仓库各自的 .agents/runs/v2.0.1-* 文件中，不纳入 Git。reference --write 已按固定配置同步生成文件；最终 reference 比较结果为 changes: 0，参考工程与模板一致。

## 验证边界与迁移

- 未执行 Docker、真实 MySQL/MongoDB/Redis 联调或六平台实机验收；integration 标签编译不等于真实数据库验收。
- 六平台发布构建、远端 CI、PR 合并、tag 和 Release 状态以实际发布平台记录为准，本开发记录不预先宣称成功。
- 提交、PR 和发版收尾按用户约定不重复执行测试，复用上表已有输出。
- 已有业务工程不会自动迁移。先读取业务子包规范，再移动文件、调整 package/import 和 Fx 注册；保持 HTTP、JSON、表名、认证、授权及事务契约。
- 本版本无 SQL 变更，不创建 SQL 版本目录或执行数据库操作。

安全自查：此次代码 diff 仅改变模型 package/import、生成过滤和版本来源；未新增接口、外部请求、持久化操作或日志输出，保留原有认证、授权和资源管理路径。
