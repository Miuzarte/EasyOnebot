# NapCat 的 OpenAPI 与 EasyOnebot 对齐指南

- 采集时间: 2026-09-10
- 目标: 说明 NapCat 是否有机器可读的 API 定义, 以及 EasyOnebot 重构时如何以它为准对齐
- 相关文档: [端点全量索引](./onebot-napcat-endpoints.md) | [四家协议端选型对比](./onebot-impl-README.md) | [NapCat 与 OB11 标准的差异备忘](./napcat-protocol-differences.md)

## 一 结论: 有, 而且是官方自动生成并随版本发布的

NapCat 仓库里有一个专门的包 `packages/napcat-schema`, 它把所有 action 的 TypeBox schema 聚合成一份 **OpenAPI 3.1.0** 文档。

生成链路 (`packages/napcat-schema/index.ts`):

1. `getAllHandlers()` + `AutoRegisterRouter` 遍历全部 action (`packages/napcat-onebot/action/index.ts`)。
2. 逐个读取 `OneBotAction` 的 `payload` / `return` schema, 以及 `action/schemas.ts` 与 `types/message` 里的公共 schema。
3. 组装成 OpenAPI: `openapi: '3.1.0'`, 每个 action 一个 `paths['/<action_name>']`, 响应统一 `allOf: [BaseResponse, {data: ...}]` (`index.ts:544`、`appendActionPaths` at `:617`)。
4. 写出 `packages/napcat-schema/dist/openapi.json`, 并生成 `missing_props.log` 报告哪些 action 缺 schema。

发布链路 (`.github/workflows/release-downstream.yml:57-83`): 每次发版执行 `pnpm run build:openapi`, 然后把产物拷进 **`NapNeko/NapCatDocs`** 仓库的 `src/api/<版本>/openapi.json` 并自动提交。

也就是说: **每发一个 NapCat 版本, 就有一份对应版本的 OpenAPI 文档进版本控制**, 从 `4.12.0` 一直到最新的 `4.18.19`, 共 84 个版本目录。

## 二 怎么拿到

| 方式 | 地址 | 说明 |
|---|---|---|
| 原始文件 (推荐) | `https://raw.githubusercontent.com/NapNeko/NapCatDocs/main/src/api/<版本>/openapi.json` | 实测 200, 648KB, 脚本化拉取最方便 |
| 按版本列目录 | `https://api.github.com/repos/NapNeko/NapCatDocs/contents/src/api` | 配合 `sort -V` 取最新 |
| 文档站 | <https://napneko.github.io/> | 在线浏览, 但不要猜 `/api/<版本>/openapi.json` 这种路径 (实测 404) |
| 自己生成 | `cd packages/napcat-schema && pnpm run build:openapi` | 需要 pnpm + Node, 依赖 typebox |

```bash
# 拉最新一版
ver=$(curl -s https://api.github.com/repos/NapNeko/NapCatDocs/contents/src/api \
  | python3 -c "import json,sys;print(sorted((x['name'] for x in json.load(sys.stdin)), key=lambda s:[int(i) for i in s.split('.')])[-1])")
curl -sL "https://raw.githubusercontent.com/NapNeko/NapCatDocs/main/src/api/$ver/openapi.json" -o napcat-openapi.json
```

## 三 这份 spec 长什么样 (以 4.18.19 为准)

| 项 | 值 |
|---|---|
| `openapi` | 3.1.0 |
| `info` | `NapCat OneBot 11 HTTP API`, version `4.18.19` |
| paths | **175** 个, 每个都是一条 `POST /<action>` |
| `components.schemas` | 39 个 (含 `BaseResponse`、各种 `OB11Message*` 段定义) |
| tags | 17 类: 核心/系统/消息/群组/用户/文件/频道接口, 以及对应的 `*扩展`, 另有 `Go-CQHTTP`、`流式接口`、`流式传输扩展`、`AI 扩展` |
| 响应包装 | `allOf: [{ $ref: BaseResponse }, { required: [data], properties: { data: ... } }]` |
| 文档字段 | summary / description / payloadExample / returnExample / errorExamples 都有, 中文描述完整 |

**三个必须提前知道的坑**:

1. **数字参数在 spec 里大多声明为 `string`**。统计到约 188 处"疑似数字却写成 string"的字段, 例如 `send_group_forward_msg` 的 `user_id` / `group_id`, `friend_poke` 的 `group_id` / `user_id` / `target_id`。这是历史兼容 (OneBot 11 允许字符串 ID) 的结果, 直接拿它生成 Go 客户端会到处要 `strconv`。
2. **枚举/联合类型用了 `null`、`allOf`、甚至少量非 JSON Schema 的 `type` 值** (统计里能看到 `type: "stream"`、`"private"`、`"bili"` 之类, 来自 TypeBox 的 `x-schema-id` 与嵌套定义)。生成器要能容忍脏值, 别假设严格合法。
3. **只覆盖 action, 不覆盖事件**。事件结构在 `packages/napcat-onebot/event` 里, 另有 `napcat-types` 包; 如果想要事件也机器可读, 得自己从 `packages/napcat-onebot/types` 抽, 或者对着 `NapCatDocs` 的事件章节写。

## 四 EasyOnebot 现状 vs spec 差距

用 `grep -oE 'NewReq\("([a-z0-9_]+)"' api/*.go` 提取现有封装:

| 项 | 数量 |
|---|---|
| EasyOnebot 已封装端点 (去重) | **77** |
| NapCat 4.18.19 端点 | **175** |
| 交集 | 68 |
| 只有 EasyOnebot 有 (NapCat 没有或改名) | **9** |
| NapCat 有但 EasyOnebot 没封装 | **107** |

只有 EasyOnebot 有的 9 个 (见 `api_lgr_*.go`), 其中多数是 **NapCat 已改名或移除**, 例如:

| EasyOnebot 里的名字 | NapCat 现状 |
|---|---|
| `delete_group_file_folder` | NapCat 叫 `delete_group_folder` (纯命名差异) |
| `rename_group_file_folder` | 未找到 |
| `fetch_mface_key` | 未找到 |
| `set_group_reaction` | 未找到 |
| `set_group_bot_status` | 未找到 |
| `send_group_bot_callback` | 未找到 |
| `upload_image` | NapCat 只有 `upload_image_to_qun_album` |
| `set_group_anonymous` / `set_group_anonymous_ban` | 未找到 (NapCat 已不再暴露匿名相关) |

这也是重构时最该处理的一类技术债: **命名对齐 + 对不存在的端点做显式标注**, 而不是留着让调用方在运行期才发现。

现有 `api/` 包本身的结构问题 (重构要动的):

- `api_std.go` (40KB) 一个文件塞了 OneBot 标准里几乎所有东西; `api_lgr_*.go` 再把 Lagrange 扩展摊在 6 个文件里。
- 命名空间靠**类型**区分: `MixCaller{Std, Lgr, Nc}`, 但同一个语义在不同实现下同名 (`Std.SendGroupMsg` vs `Lgr.SendGroupForwardMsg`), 调用方要自己判断该用哪个; `Nc` 目前只有一个 `api_nc_message.go`, 实际项目里大量 NapCat 能力被塞进了 `Lgr` (因为 NapCat 兼容 Lagrange 命名), 语义是歪的。
- 参数用 `map[string]any` 手写 (`NewReq("send_private_msg", map[string]any{...})`), 没有结构体、没有编译期校验, 也没有和任何 spec 对齐, 所以 endpoint 拼写错误只能在运行期暴露。
- 响应类型散落各处且命名不统一 (`SendPrivateMsgResp = SendAnyMsgResp`、`SendGroupForwardMsgResp` vs `SendPrivateMsgResp` 互相指向对方的诡异别名, 见 `ctxOp_lgr.go:28-35`)。
- `validate` tag 与 struct 校验是手写的 (`api_lgr_ability.go` 一类), 与 spec 无关联。

## 五 对齐方案建议

分四步, 每步都能独立验收, 不必一次性大爆炸。

**第 1 步 建立单点事实源 (spec 进仓)**

把 NapCat 的 `openapi.json` 作为 dev 资产纳入 EasyOnebot: 例如 `docs/napcat/openapi-4.18.19.json`, 并写一个 `script/sync-napcat-spec.sh` 负责从 `NapCatDocs` 拉取、按版本落盘、更新 `versions.json`。这样离线可查、diff 可审、不需要联网构建。

**第 2 步 生成端点清单与常量, 先拿到"编译期对齐"**

不必立刻做全量客户端生成, 先把最大的一笔收益吃掉:

- 从 spec 生成 `api/endpoints_gen.go`, 内容是全部 175 个端点名常量 (`const ActionSendGroupMsg = "send_group_msg"`)。
- 把 `api/*.go` 里所有 `NewReq("字面量")` 改成 `NewReq(ActionXxx)`, 拼写错误从此变成编译错误。
- 加一个测试: 遍历 spec 的 paths, 断言每个 Action 常量都存在; 反向断言每个常量都在 spec 里 (允许白名单放那 9 个 NapCat 没有的扩展)。

**第 3 步 参数与响应用生成代码兜底**

- 由 spec 生成 payload struct (Go 结构体 + json tag), 数字字段按 spec 声明为 `string` 的, 在生成时统一转成 `int` 或保留 `string` 并加显式转换辅助函数, 不要就地 `strconv` 满地撒。
- 由 spec 生成响应 struct, 复用一个统一的 `BaseResponse` + `data` 泛型包装 (Go 1.27 可以用泛型 `Resp[T]`), 顺手把现有那些互相指向的别名删掉。
- 生成器要能容忍第三节说的三类脏值; 只生成 `paths`/`components` 里存在的端点, 缺 schema 的 action 记进 `missing.log` 供人工补。

**第 4 步 重新划命名空间与调用入口**

- `Std` / `Lgr` / `Nc` 三层保留, 但明确规则: **一个端点只归属一层**, 属于 NapCat 原生能力的放 `Nc`, 只有 Lagrange 也提供的才允许留 `Lgr`, 并加注释标注 NapCat 侧的真实端点名。
- `MixCaller` 改成语义入口 (如 `Call().Message.SendGroupForwardMsg(...)`) 或按能力聚合, 让调用方不再需要知道 "这个端点其实来自哪个实现"。
- `M_ForwardSlice.go` 那种按 app_name 分支的逻辑, 收拢成实现能力表 (`napcat` / `llonebot` / `snowluma` 支持同一种 forward 语义), 而不是散在业务代码里 `switch vi.AppName`。

**验收标准**: `go build` + `go vet` + 新增的 spec 一致性测试全绿; 业务侧 (NothingBot_v4) 只需为一处 app_name 分支加 case, 其余调用保持编译通过。

## 六 其他可选事实源

- `packages/napcat-onebot/action/schemas.ts` 与 `packages/napcat-onebot/types/**`: spec 的上游, 想要更细的类型 (message 段、事件) 时直接看这里。
- `packages/napcat-schema/index.ts` 里的 `missing_props.log`: 能直接告诉你哪些 action 还没有完整 schema, 生成客户端时应视为"字段不全"。
- `NapCatDocs` 仓库的 `src/api/[version].md`: 文档站渲染入口, 想知道官方怎么展示某个 action 时看它。

## 七 复现命令

```bash
# 1 取 spec
git clone --depth 1 https://github.com/NapNeko/NapCatDocs /tmp/napcatdocs
ls /tmp/napcatdocs/src/api | sort -V | tail -1        # 4.18.19

# 2 基本统计
python3 - <<'PY'
import json
d = json.load(open('/tmp/napcatdocs/src/api/4.18.19/openapi.json'))
print(d['openapi'], d['info']['version'], len(d['paths']), len(d['components']['schemas']))
PY

# 3 与 EasyOnebot 对比
cd ~/git/EasyOnebot
grep -ohE 'NewReq\("[a-z0-9_]+"' api/*.go | sort -u | wc -l   # 77
```
