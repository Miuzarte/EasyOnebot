# Lagrange.Core V2 分析

- 仓库: <https://github.com/LagrangeDev/Lagrange.Core>
- 本地快照: `~/git/Lagrange.Core`, HEAD `20c2ba0` (2026-09-08), 浅克隆 depth=200, 当前分支为 V2
- 采集时间: 2026-09-10
- 定位: 纯 C# 的 NTQQ 协议实现, V1 已 sunset (移到 `v1` 分支)

## 0 一句话结论

**V2 不能接 OneBot 11, 不建议为 NothingBot_v4 换过去。** OneBot 11 是被明确放弃, 不是"暂未支持"。

## 1 OneBot 11 支持: 无

- 官方文档原文: "由于 OneBot 11 协议的各种历史遗留问题, 我们最终决定, LagrangeV2 不再支持 OneBot 11 协议, V1 的 OneBot 11 协议实现也一并不再维护" (<https://lagrangedev.github.io/Lagrange.Doc/v2/>)。
- 全仓 `grep -ri onebot` 仅 2 处命中: `README.md:29` 的 Milky 说明, `Lagrange.Core.Runner/QrCodeHelper.cs:7` 的代码出处注释。**零实现代码**; GitHub topic 残留 `onebot11` 只是标签未清。
- 没有任何 roadmap / issue 承诺恢复。
- 现在提供的是 **Milky 协议** (`Lagrange.Milky` 项目)。

中间层选项: `molanp/milky2onebot` (Python, 0 star, README 自称"小孩子不懂事写着玩的"、"api太多了, 我不知道有没有缺失api"), 不可用于生产。可行替代是自写适配层, 或者直接用同时支持 OneBot 11 / Milky 的 LuckyLilliaBot。

## 2 Milky 协议

- API: `POST /api/<name>` (JSON, `Authorization: Bearer`)。
- 事件: WebSocket / SSE / WebHook 三种, Lagrange.Milky 三种都已实现 (`Lagrange.Milky/Extensions/HostApplicationBuilderExtension.cs:98-113`)。
- 生态: NoneBot / Koishi / Karin / Go / Rust / .NET SDK 齐全 (<https://milky.ntqqrev.org/awesome>)。
- 文档: <https://milky.ntqqrev.org/>, 实现文档 <https://lagrangedev.github.io/Lagrange.Milky.Document>
- **实现落后于协议**: 本地实现自报 `MilkyVersion = "1.2"` (`Lagrange.Milky/Constants.cs:8`), 协议已到 1.3 (2026-08-07 发布, <https://milky.ntqqrev.org/changelog>)。

## 3 部署方式

5 步:

1. 拉 `ghcr.io/lagrangedev/lagrange.milky:main` (13.0 MiB, amd64 + arm64, `.github/workflows/milky-docker.yaml`) 或 `dotnet publish` 自编译。
2. 首次启动生成 `appsettings.json`。
3. 填 `Lagrange.Login.Uin` 与 `Lagrange.Protocol.Signer.Token`。
4. 扫码或密码登录 (生成 `<uin>.ks`)。
5. 连接 `ws://ip:3000/event`。

**不需要 QQ 客户端**, 但**签名是硬依赖且白名单制**: 官方 `sign.lagrangecore.org` 需用 GPG 私钥向 QQ 机器人自助申请 token (每 GitHub 账号最多绑 3 个 QQ, 见 `LagrangeDev/SignApiGuide`), 代码走 `HttpSigner` (`Lagrange.Milky/Signing/HttpSigner.cs:16-39`)。

发布物问题: GitHub 无 release, CI 产物只有 Actions artifacts (会过期); 官方下载页仍写 "WIP 当前 Lagrange.Milky 处于未完成状态", 且链接指向已 archived 的 LagrangeV2 仓库。

## 4 更新稳定性: 偏弱

- 无 version tag / release: tags API 只有 `nightly`, releases 只有 2025-08-02 的 Nightly。
- 提交节奏 (`git log --date=format:%Y-%m` 直方图): 2026-05 共 11 次 -> **2026-06 零提交** -> 2026-07 共 39 次 -> 2026-08 共 7 次 -> 2026-09 仅 2 次。
- 核心维护者 3 人 (2026 年 DarkRRb 30 / NoirHare 25 / Linwenxuan04 22 次提交)。
- 半年内破坏性变更: 时间字段改 Unix 时间戳 (`eb29489`)、Milky Refactor (`#908`)、appsettings 结构调整。
- V1 分支末次提交 2025-10-10 "[Chore] Update README.md termination notice"; V2 已并回 master。

## 5 资源占用与性能

- **未找到官方 benchmark** (仓库里只有 `Lagrange.Proto.Benchmark`, 比的是 protobuf 序列化)。
- 可确认的客观事实: Native AOT 单文件 + 13 MiB 镜像 + 完全不含 Electron/Node, 相比寄生 QQNT 的实现资源占用优势**是真实存在的**, 但没有量化数据。

## 6 兼容性映射 (若真要接 Milky)

现用点: `utilsBot.go:62,74`, `MO_GroupNotice.go:29,32`, 另有 `T_forwardMsg_test.go:35`。

| 现用 (EasyOnebot `.Lgr`) | Milky 对应 | 迁移改动 |
|---|---|---|
| `send_group_forward_msg` | `send_group_message` + `forward` 消息段 (`data.messages=[{user_id, sender_name, segments}]`) | 节点结构重写 (现在是 uin/name/content) |
| `send_private_forward_msg` | `send_private_message` + `forward` 段 | 同上 |
| `friend_poke` | `send_friend_nudge` (user_id, is_self) | 改名 + 参数 |
| `group_poke` | `send_group_nudge` (group_id, user_id) | 改名 |
| `get_forward_msg` (入站合并转发切片) | `get_forwarded_messages`, **Lagrange.Milky 未实现** (README 未勾选) | 功能直接缺失 |

真正的成本不在这 4 个端点, 而在 EasyOnebot 与全部 matcher/事件/消息段都是 OneBot 11 语义 (EasyOnebot 99 个端点) —— 换 Milky 等于重写协议层与事件模型。

## 7 代码质量

- 分层清晰 (`Lagrange.Core` / `Lagrange.Proto` / `Lagrange.Milky` / `Lagrange.Core.NativeAPI` + 源生成器)。
- 依赖极少 (`Lagrange.Milky.csproj` 只有 3 个 `PackageReference`)。
- 190 个 `[Test]` 但全在 Core/Proto/加密算法上; **Milky 零测试, CI 不跑 `dotnet test`**。
- 仓库根目录**无 LICENSE 文件**, NuGet `Lagrange.Core` 2.1.1 (2026-07-28) 的 licenseExpression 为空, 仅 `Lagrange.Proto` 标 GPL-3.0-or-later。

## 8 AI 痕迹与风险

- 未发现 AI 痕迹: 无 claude / copilot / chatgpt 字样, 代码风格统一, 手工 proto 库 + 源生成器。
- README 有标准免责声明。
- 长期风险: 签名白名单由维护者单点控制; 封号灰产属性; 无 release / 版本号; V1 已经 sunset 过一次, "能不能长期用"不乐观。

## 9 给 NothingBot_v4 的建议

不要迁移。若只是想要低资源, 换 LuckyLilliaBot 比换协议划算; 在没有稳定 release、且签名依赖白名单服务的前提下, 不建议把机器人绑到 Lagrange V2 上。
