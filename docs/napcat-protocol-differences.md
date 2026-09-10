# NapCat 与 OneBot 11 标准的差异备忘

本文记录**实测到的** NapCat 行为与 OneBot 11 标准/其它实现不一致的地方, 以及 Easyonebot 与
消费方 (NothingBot_v4) 各自需要承担什么。

原则:

- 能在 Easyonebot 内部归一化的就归一化, 让消费方只看到一套语义
- 归一化不了的 (实现能力差异) 记在这里, 由消费方按 `get_version_info` 的 `app_name` 分支

实测环境: NapCat 4.18.19 (Docker), QQ 2393827810, 2026-09。

## 1 poke 通知的字段语义 (已归一化)

**差异**: OneBot 11 标准里 `user_id` 是**发起** poke 的人, `target_id` 是被戳的人;
NapCat 里 `user_id` 与 `target_id` **都是被戳方**, 发起方放在标准里不存在的 `sender_id`。

bot 主动戳好友时 NapCat 回推的原始事件:

```json
{
  "post_type": "notice", "notice_type": "notify", "sub_type": "poke",
  "self_id": 2393827810,
  "target_id": 982809597,   // 被戳方
  "user_id": 982809597,     // 也是被戳方 (标准里这里应是发起方)
  "sender_id": 2393827810,  // 真正发起的人 = bot
  "raw_info": [ ... ]
}
```

| 字段 | OneBot 11 标准 | NapCat |
|---|---|---|
| `user_id` | 发起方 | 被戳方 |
| `target_id` | 被戳方 | 被戳方 |
| `sender_id` | 不存在 | 发起方 |

**后果**: 如果按标准判断"谁戳了我"(`user_id == self_id`), NapCat 下永远不成立; 而
"`target_id == self_id`"虽然多数场景歪打正着, 但 **bot 自己发起的 poke 也会被 NapCat 回推**,
用 `target_id` 判断会把"我戳别人"误判成"别人戳我"。

**Easyonebot 的处理**:

- `event.NoticeNotifyPoke_Nc` 增加 `SenderId` 字段。注意必须加在 `Event` 直接嵌入的结构里 ——
  事件解码是 `mapstructure.Decode(m, event)` (conn.go) 再经 `utils.AnyCopy` 复制, `copier`
  不认嵌套层级的字段。
- `handler.go` 的 poke 分支把 `SenderId` 归一化: 缺失时回退成 `UserId`。因此消费方**只读
  `SenderId`** 即可, 不用关心背后是 NapCat 还是标准实现。

```go
onebot.OnNoticeNotifyPoke(func(p *event.NoticeNotifyPoke) {
    if p.SenderId == p.SelfId {   // 排除自己发起的 poke (NapCat 会回推)
        return
    }
    if p.TargetId != p.SelfId {   // 只处理戳向 bot 的
        return
    }
    ...
})
```

**回归测试**: `T_poke_test.go` 的 `TestNormalizePoke` (字段语义三种场景) 与
`TestClassifyPokeEvent` (走真实分发路径, 断言回调收到的发起方)。

## 2 合并转发节点获取方式 (消费方需分支)

同一个 `get_forward_msg`, 两种实现的返回形状不同:

| 实现 | `app_name` | 行为 |
|---|---|---|
| NapCat | `NapCat.Onebot` | 消息段里直接带 `content`, 就是 event 数组 |
| Lagrange | `Lagrange.OneBot` | 只给 `id`, 要再调一次 `get_forward_msg` 拿节点 |
| LLBot | `LLOneBot` | 与 NapCat 同语义 (返回的 `messages[].content` 即消息段数组) |
| SnowLuma | `SnowLuma` | 与 NapCat 同语义 |

**这不是 Easyonebot 归一化得了的**: 合并转发在 OneBot 11 标准里本身就有两套 (go-cqhttp 的
`messages[].content` 与 Lagrange 的 `id` 二次获取)。

**现状**: NothingBot_v4 的 `M_ForwardSlice.go` 按 `GetVersionInfoCache().AppName` 分流,
目前只认两个名字, 换实现时要补 case:

```go
switch vi.AppName {
case "NapCat.Onebot", "LLOneBot", "SnowLuma": // 同一语义
    return forwardGetNodes_nc(forward)
case "Lagrange.OneBot":
    return forwardGetNodes_lgr(forward)
default:
    return nil, fmt.Errorf("unsupported onebot implementation: %s", vi.AppName)
}
```

**建议**: 这段分流更适合收进 Easyonebot (按 `app_name` 暴露一个"取合并转发节点"的统一方法),
下一步可以迁。

## 3 get_msg 的 message / sender (已在 spec 层修正)

NapCat 实际返回:

```json
{
  "message": [ {"type": "text", "data": {"text": "..."}} ],   // 消息段数组
  "sender":  { "user_id": ..., "nickname": "...", "card": "...", "role": "..." }
}
```

上游 OpenAPI spec 里这两个字段只写了 `{"type": "object"}`, 没有结构。`api/napcat` 的同步脚本
用覆盖表把它们改成 `[]OB11MessageData` 与具名对象 (`sync_spec.py` 的 `SCHEMA_OVERRIDES`),
依据就是上面这条实测。

`get_forward_msg` 的 `messages` 在上游 spec 里是 `[]interface{}`, 没有结构, 所以
`Ctx.GetForwardMsg` 目前仍用 `api.GetForwardMsgResp` (它解成了 `message.SegmentArray`)。

## 4 各项 ID / 时间戳的类型 (已在 spec 层修正)

上游把 ID 与时间戳一律声明成 `type: number` (无 format), oapi-codegen 会映射成 **float32** ——
QQ 号、消息 id、时间戳都远超 float32 的 24 位有效精度。实测这些字段返回的都是整数, 所以同步脚本
用 `INT64_FIELDS` 覆盖表把它们钉成 `int64` (62 处)。

反过来, **请求**里不少 ID 字段 spec 声明成 `string` (例如 `send_group_msg.group_id`), 消费方
传 int 时要显式转换 —— 这也是 NothingBot 迁移时 `int -> string` 改动的来源。

## 5 无 echo 的主动推送 (消费方需注意)

NapCat 会主动推**不带 `echo`** 的消息:

- `meta_event` 的 `lifecycle` / `heartbeat`
- 各类 `notice` 事件 (包括上面第 1 条里 bot 自己发起的 poke 回推)

所以任何自己实现的 API 调用层**必须按 echo 关联响应**, 顺序读会把推送当响应。
`api/napcat/live_test.go` 的 poster 就是按这个实现的, 可以直接参考。

## 6 实现识别

`get_version_info` 的 `app_name` 是判断实现类型的唯一可靠依据 (NapCat 的文档示例里写的
`NapCatQQ` 是过时示例, 源码里是 `NapCat.Onebot`):

| 实现 | `app_name` |
|---|---|
| NapCat | `NapCat.Onebot` |
| Lagrange | `Lagrange.OneBot` |
| LLBot | `LLOneBot` |
| SnowLuma | `SnowLuma` |

## 复现命令

```bash
# 抓所有 notice 推送 (含 poke), 观察真实字段
# 见 api/napcat/live_test.go 的 wsPoster 实现

# poke 字段语义的回归测试
go test -run 'TestNormalizePoke|TestClassifyPokeEvent' -v .

# 真机联调 (会发消息并撤回)
NAPCAT_WS_URL=ws://127.0.0.1:8081 NAPCAT_TEST_GROUP=<测试群> NAPCAT_TEST_USER=<测试好友> \
    go test ./api/napcat/ -run TestLive -v
```

## 相关文档

- [NapCat 的 OpenAPI 与 EasyOnebot 对齐指南](./onebot-impl-napcat-openapi.md)
- [NapCat OneBot 11 端点全量索引](./onebot-napcat-endpoints.md)
- [协议端选型对比](./onebot-impl-README.md)
