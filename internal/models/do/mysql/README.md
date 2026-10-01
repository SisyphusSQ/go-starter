# MySQL 模型

按业务域建立 `<domain>_do/` 子包，禁止 MySQL DO 平铺在 mysql 根包。子包内一个数据库表对应一个 .go 文件，包含该表 DO、明确的列映射和 TableName()；例如 user.go 对应 users。禁止将多张表合并进 models.go/schema.go，也不拆散同一表的 DO 与映射方法。

表结构与 docs/sqls/schema 的完整 SQL 保持一致，发布归档遵守 docs/sqls/README.md。共享存储定义的边界见 [业务子包约定](../../../../docs/design/architecture/packages.md)，类型责任见 [模型分层](../../../../docs/design/architecture/models.md)。
