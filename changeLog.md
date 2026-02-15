## v1.0.0(20260215)
#### feature:
1. 新增用户接口响应体定义 `UserListResp`、`UserMongoListResp`、`UserIDResp`、`UserMongoIDResp`，统一列表与主键返回结构
2. 新增 MySQL 与 Mongo 用户 Service 直接组装响应体的返回逻辑，List/Update/Delete 由 Service 返回业务 Resp，Controller 仅做参数处理与透传
3. 新增 `docs/schema/users_example_insert.sql` 的 `password` 字段示例数据，支持账号登录联调

#### optimization:
1. 优化 MySQL 与 Mongo 用户 Handler，移除 `map[string]any` 形式返回，改为显式 VO 类型返回，提升接口契约清晰度
2. 优化相关单元测试与 Mock 签名，适配 Service 返回类型调整，保证测试用例与当前接口一致

