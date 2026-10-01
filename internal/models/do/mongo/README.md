# MongoDB 模型

按业务域建立 `<domain>_do/` 子包，禁止 DO 平铺在 mongo 根包；使用官方 Driver v2 的 bson 类型和标签。明确集合、索引与软删除契约；collection 在 Fx OnStart 后使用。见 [模型分层](../../../../docs/design/architecture/models.md) 和 [业务子包约定](../../../../docs/design/architecture/packages.md)。
