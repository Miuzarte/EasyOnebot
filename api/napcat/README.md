# api/napcat —— NapCat 的 OneBot 11 类型定义

这个包把 NapCat 官方的 OpenAPI spec 拉进仓库, 并用 oapi-codegen 生成 Go 类型。
目的是让 EasyOnebot 手写的端点名 / 参数 / 响应结构有一个**机器可校验的事实源**,
而不是对着 onebot11 文档手抄。

## 文件

| 文件 | 说明 |
|---|---|
| `spec/openapi-raw-<版本>.json` | 上游原样产物, **不要手改** |
| `spec/openapi-<版本>.json` | 注入 operationId 与去重后的实际生成输入 |
| `spec/version.json` | 版本指纹: 来源 URL、端点/schema 数量、注入与改动记录 |
| `sync_spec.py` | 拉取上游 spec 并做上述处理 |
| `cfg.yaml` | oapi-codegen 配置 (只生成 models) |
| `openapi.gen.go` | 生成结果, **不要手改** |
| `spec.go` | 读 spec 的小工具: `SpecEndpoints()` / `SpecOperationIDs()` |
| `spec_test.go` | 一致性测试, 见下 |

## 类型命名

| 端点 | 请求体 | 响应 data |
|---|---|---|
| `send_group_forward_msg` | `SendGroupForwardMsgJSONBody` | `SendGroupForwardMsgData` |
| `get_login_info` | `GetLoginInfoJSONBody` | `OB11User` (spec 里就是 `$ref`) |
| `delete_msg` | `DeleteMsgJSONBody` | `DeleteMsgData` |

请求体字段大多是可空指针 + `omitempty`, 数字类 ID 在 spec 里是 `string`
(例如 `GroupID *string`), 调用方传 int 时需要显式转换。

取指针用 Go 1.26+ 的 `new(expr)` (例如 `GroupID: new("123")`), 不再需要自定义 `Ptr` 辅助函数。

## 刷新流程

```bash
cd api/napcat
python3 sync_spec.py             # 同步上游最新版本 (也可 sync_spec.py 4.18.20)
python3 sync_spec.py --check     # 只看本地是否落后 (退出码 1 表示有新版)
# 同步新版本后要改两处常量:
#   generate.go 里的 go:generate 指令
#   spec.go 里的 specFile
go generate ./...                # 重新生成 openapi.gen.go
go test ./api/napcat/            # 跑一致性测试
```

## 为什么需要预处理

上游 spec 是 NapCat 从 TypeBox schema 自动生成的, 直接拿来喂 oapi-codegen 有两处问题:

1. **175 个端点全部没有 operationId**。不注入的话, 生成的方法/类型名只能由 HTTP 方法 + path
   推断 (`PostSendGroupMsg`), 又长又难用。`sync_spec.py` 把 path 当作端点名直接注入
   (`/send_group_forward_msg` -> `SendGroupForwardMsg`), 生成的类型就是
   `SendGroupForwardMsgJSONBody` / `SendGroupForwardMsgJSONRequestBody`。
2. **响应 data 是匿名内联 schema** (`allOf[BaseResponse, {data: {...}}]`), 只生成 models 时
   不会产出 Go 类型, 调用方没法强类型解析。脚本给每个 data 打上 `x-go-type-name`,
   生成 `<OperationID>Data` (如 `SendGroupForwardMsgData`); data 本身是 `$ref` 的端点
   (如 `get_login_info` -> `OB11User`) 保持原样。
3. **少量字段下划线与驼峰双写** (`category_id` 与 `categoryId`, `reverse_order` 与 `reverseOrder`),
   会生成同名字段导致编译失败。脚本只在两个字段的 schema 除 description 外完全一致时丢弃
   非下划线形式, 并把丢弃记录写进 `version.json`; schema 不一致会直接报错要求人工确认。

另外两个带前导点的旧路由 (`/.ocr_image`、`/.handle_quick_operation`) 在 Go 名字规范化时会与
公开端点撞名, 脚本里显式改名为 `OcrImageLegacy` / `HandleQuickOperationLegacy`。

## 只生成 models 的原因

NapCat 的 spec 描述的是 HTTP `POST /<action>` 接口, 而 EasyOnebot 走 WebSocket 收发同样形状的
JSON。生成的 HTTP client (含 351 个方法与请求构造) 一行都用不上, 且会让生成文件从 1.4 万行涨到
6 万行。所以 `cfg.yaml` 里只开 `models`, 由 EasyOnebot 自己在 WS 传输层上按端点名分发。

## 一致性测试

| 测试 | 覆盖的问题 |
|---|---|
| `TestHandwrittenEndpointsExistInSpec` | 手写 `NewReq("...")` 的端点名是否都存在于 spec (拼写错误 / NapCat 已移除); `legacyEndpoints` 列表是否过期 |
| `TestSpecOperationIDs` | operationId 是否注入完整、是否重名 (生成类型一一对应的前提) |
| `TestGeneratedCodeIsFresh` | 生成代码是否覆盖全部端点 (忘记 `go generate` 会失败) |
| `TestVersionMetaMatchesSpec` | `version.json` 指纹与当前 spec、`specFile` 是否一致 |

## 已知的端点缺口

`spec_test.go` 的 `legacyEndpoints` 记录了 EasyOnebot 手写层里、NapCat 已经不提供的端点
(共 10 个, 含 `delete_group_file_folder` 这种仅命名差异的)。同步新 spec 后如果某个端点回到
spec 里, 测试会失败并提醒把这个条目删掉。
