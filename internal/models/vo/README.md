# VO：HTTP 边界对象

沿用本仓既有约定：HTTP 请求与响应均放 VO。请求使用 <Action><Resource>Req，响应使用 <Resource>Resp、<Resource>View 或明确的列表类型。

业务请求/响应必须放 `<domain>_vo/` 子包，单文件业务也不得平铺。根包仅保留统一响应、绑定、通用校验及通用分页能力，不导入业务 VO 子包。

请求实现 Validate() error 并由 controller 调用 BindAndValidate；service 保留业务规则与资源授权。Response 和响应 helper 为统一出口，VO 不包含 gorm/bson 标签或密码哈希。

规则见 [模型分层](../../../docs/design/architecture/models.md) 和 [业务子包约定](../../../docs/design/architecture/packages.md)。
