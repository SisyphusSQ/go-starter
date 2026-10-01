# DTO：传输与跨层对象

按业务域或外部系统建立独立子包，例如 payment_dto/、lark_dto/；即使只有一个文件也不得放在 dto 根包。DTO 用于应用服务输入、命令、事件、外部 SDK 的请求和结果；单纯 HTTP 请求沿用本仓 VO 约定。

禁止把带 json/form/validate 标签的跨层结构体藏在 service 中，或用匿名结构体绕过模型分层。DTO 不依赖 service/controller，也不携带数据库连接、Echo context 或持久化标签。

规则见 [模型分层](../../../docs/design/architecture/models.md) 和 [业务子包约定](../../../docs/design/architecture/packages.md)。
