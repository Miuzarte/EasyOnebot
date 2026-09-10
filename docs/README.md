# EasyOnebot 文档

本目录收集协议端实现的分析、选型对比, 以及实测到的协议差异。这些内容服务于
**EasyOnebot 与 NothingBot_v4**: 前者是对齐 OneBot 11 / NapCat 的库, 后者是消费方。

## 协议差异 (写代码前先看)

- [NapCat 与 OneBot 11 标准的差异备忘](./napcat-protocol-differences.md) —— 实测到的字段语义差异:
  poke 的 `sender_id`、合并转发节点获取、ID/时间戳类型、无 echo 的主动推送, 以及各自该由谁处理

## NapCat 对齐

- [NapCat 的 OpenAPI 与 EasyOnebot 对齐指南](./onebot-impl-napcat-openapi.md) —— NapCat 每版自动生成
  OpenAPI 3.1.0 (175 个端点), spec 在哪、怎么用 oapi-codegen 消费、有哪些坑
- [NapCat OneBot 11 端点全量索引](./onebot-napcat-endpoints.md) —— 175 个端点按 tag 分组, 标注 EasyOnebot
  是否已封装

生成流水线本身的说明在 `api/napcat/README.md` (与代码放在一起)。

## 协议端选型

- [选型对比总表](./onebot-impl-README.md) —— 四家 (NapCat / Lagrange V2 / LLBot / SnowLuma) 横向对比
  + 与 NothingBot_v4 的兼容矩阵 + 结论
- [NapCatQQ](./onebot-impl-napcatqq.md)
- [Lagrange.Core V2](./onebot-impl-lagrange-core-v2.md)
- [LuckyLilliaBot (LLBot)](./onebot-impl-luckylilliabot.md)
- [SnowLuma](./onebot-impl-snowluma.md)

## 维护约定

- 每条结论都标证据 (源码路径、命令或实测事件), 不写推测
- 实测数据的采集时间写在文首; 上游行为可能随版本变化, 复现命令都在文末
- 新增协议差异时, 同时在这里登记, 并在对应代码处加注释指回来
- 源码一律 LF 行尾 (与 NothingBot_v4 一致), 不要提交 CRLF; `gofmt -l` 应保持为空
