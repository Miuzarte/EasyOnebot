# OneBot 11 协议端选型对比

本文对比四个 QQ 协议端实现, 面向的消费方是 **NothingBot_v4** (Go + EasyOnebot, OneBot 11 over WebSocket)。
四份单仓库详报见同目录:

- [Lagrange.Core V2](./onebot-impl-lagrange-core-v2.md)
- [NapCatQQ](./onebot-impl-napcatqq.md)
- [LuckyLilliaBot (LLBot)](./onebot-impl-luckylilliabot.md)
- [SnowLuma](./onebot-impl-snowluma.md)

数据采集时间: **2026-09-10** (本地浅克隆 depth=200 + GitHub REST API + Docker Hub API)。

## 一 总览

| | NapCatQQ | Lagrange.Core V2 | LuckyLilliaBot | SnowLuma |
|---|---|---|---|---|
| 仓库 | NapNeko/NapCatQQ | LagrangeDev/Lagrange.Core | LLOneBot/LuckyLilliaBot | SnowLuma/SnowLuma |
| 语言 | TypeScript | C# (.NET) | TypeScript (Node >=24) | TypeScript (Node >=22) |
| 项目年龄 | 2024-03 起 | V1 已 sunset, V2 现役 | 2023-10 起 (原 LLOneBot 改名) | 代码首提交 2026-05, 仅约 4 个月 |
| star / fork | 10.5k / 818 | 3.0k / 352 | 3.6k / 281 | 1.1k / 73 |
| OneBot 11 | 支持 (184 个 action) | **已明确放弃** | 支持 (120 个 action) | 支持 (约 187 个 action) |
| 默认协议 | OneBot 11 | Milky | OneBot 11 + Satori + Milky | OneBot 11 |
| 注入方式 | 补丁 QQNT (Electron) | 纯协议复刻, 不要客户端 | 直连复刻 或 真客户端容器 | ptrace 注入运行中的 QQNT |
| Docker 官方镜像 | 有 (NapCat-Docker) | ghcr.io 有 (milky:main, 13 MiB) | 有 (linyuchen/llbot, 57 tags) | 有 (284 tags) |
| Docker 部署难度 | 一条 `docker run` + WebUI 扫码 | 5 步 + 白名单签名 token | 向导脚本 + 必填 Auth Token | 需 `SYS_PTRACE`/`seccomp=unconfined`/`shm 1g` + noVNC 扫码 |
| 需要自装 QQ | 否 (镜像内置 LinuxQQ) | 否 | 直连否 / PMHQ 模式要 | 否 (镜像内置) |
| 稳定 release | v4.18.19 (2026-08-14) 1-3 周/版 | **无 release**, 只有 2025-08 nightly | v8.1.10 (2026-09-04) 约每周 | v1.14.15 (2026-08-30) 约 1.2 天/版 |
| 近半年提交 | 90 | 93 (其中 2026-06 整月为 0) | 670 | 1126 |
| 核心维护者 | 活跃, 多人 | 3 人 | 2 人 (占 98%) | 1 人 (占 78%) |
| 客户端 app_name | `NapCat.Onebot` | 无 OneBot | `LLOneBot` | `SnowLuma` |
| 许可证 | 自定义限制性 (禁商用) | 根目录无 LICENSE, NuGet licenseExpression 为空 | 未见特殊限制 | 源码可见非商业 + EULA 5.4 禁并入第三方镜像 |

## 二 与 NothingBot_v4 的兼容矩阵

机器人侧真正会踩的只有两处:

**1 合并转发节点获取必须按实现分发**

`M_ForwardSlice.go:96-114` 只认两个 app_name:

```go
const (
    LAGRANGE_APPNAME = "Lagrange.OneBot"
    NAPCAT_APPNAME   = "NapCat.Onebot"
)
switch vi.AppName {
case NAPCAT_APPNAME:   return forwardGetNodes_nc(forward)  // content 直接是 event 数组
case LAGRANGE_APPNAME: return forwardGetNodes_lgr(forward) // 拿 id 再调 get_forward_msg
default:               return nil, fmt.Errorf("unsupported onebot implementation: %s", vi.AppName)
}
```

换到 LLBot 或 SnowLuma 时这一段会直接报 unsupported, 需要新增分支。好消息是两者的 `get_forward_msg` 语义与 NapCat 一致 (返回的 `messages[].content` 就是消息段数组, 不需要二次取节点), 因此可以并进 NapCat 分支:

```go
case "NapCat.Onebot", "LLOneBot", "SnowLuma":
    return forwardGetNodes_nc(forward)
```

**2 扩展端点 (EasyOnebot `.Lgr.*`) 是否被实现**

机器人用到的非标准端点只有 4 处 (`utilsBot.go:62,74`, `MO_GroupNotice.go:29,32`, 另有 `T_forwardMsg_test.go:35`):

| 端点 | NapCat | LLBot | SnowLuma | Lagrange V2 |
|---|---|---|---|---|
| `send_group_forward_msg` | 有 | 有 (`action/types.ts:106`) | 有 (`extended.ts:2477`) | 无 |
| `send_private_forward_msg` | 有 | 有 (`107`) | 有 (`2494`) | 无 |
| `friend_poke` | 有 (`router.ts:173`) | 有 (`48`) | 有 (`402`) | 无 |
| `group_poke` | 有 (`router.ts:174`) | 有 (`47`) | 有 (`415`) | 无 |

也就是说 NapCat / LLBot / SnowLuma 三家都不需要改机器人代码就能继续用扩展端点, 只有 app_name 分发那一处要补。

## 三 结论与建议

优先级前提是 **部署简单 + 更新稳定**, 不是性能。

1. **继续用 NapCat (Docker)** —— 综合最优, 没有换的必要。官方镜像一条命令、免装 QQ、生态最大 (10.5k star)、issue 响应活跃; 官方 Docker 镜像把 LinuxQQ 版本钉在构建时, 正好对冲它最大的被动升级风险。主要缺点是被 QQ 客户端版本牵着走 (`#2043` PacketBackend 不支持新 QQ 版本)、无 CHANGELOG 且 release note 由 LLM 生成、有内存泄漏/CPU 打满类 issue (`#1047` `#1218` `#1771`)。
2. **想要更省资源可试 LuckyLilliaBot 的 PMHQ 模式** —— 唯一"老项目 + 显式多协议 + 同时提供 OneBot 11 / Milky"的替代品, 2023 年至今同仓库演进, release 与 Docker 镜像同日推送。注意两点: 直连模式强依赖作者中心化 `auth.luckylillia.com` 的 Auth Token (token 还有可用 QQ 数量上限), 属于 NapCat 没有的单点; 另有未修的 Docker 反复掉线 issue (`#857`)。要绕开签名单点就用 PMHQ (真客户端) 模式, 代价是双容器与更多内存。
3. **Lagrange V2 不适合本项目** —— 官方文档明确"不再支持 OneBot 11 协议", 全仓 `grep -ri onebot` 零实现代码; 现在只有 Milky 协议, 唯一中间层 milky2onebot 是 0 star 的玩具项目。C# AOT 的性能优势是真的 (13 MiB 镜像、无 Electron), 但代价是重写整个 OneBot 协议层, 而它自己连正式 release 都没有, 签名靠白名单制 GPG 申请 —— V1 已经 sunset 过一次, 长期可用性不乐观。
4. **SnowLuma 暂不推荐 (可观察)** —— OneBot 覆盖与工程规范都不错, issue 中位关闭 0.68 天, 但项目实际只有 4 个月、单人占 78% 提交、99 个 tag 约 1.2 天/版, 且路线图已写「V2.0.0 前计划」预示破坏性变更; Docker 还要 `SYS_PTRACE` + noVNC, 比 NapCat 繁琐, 许可证也限制自建镜像。适合愿意高频跟进更新、看重 OneBot 覆盖面的用户。

## 四 对 AI 开发占比的客观核查

用户对后两个新项目的直觉 ("AI 开发占比太多") 在可核查的 git 元数据上**不成立**:

| 指标 | LLBot | SnowLuma |
|---|---|---|
| AI 协作/提示词文件 | 仅 `CLAUDE.md` 398 字节 (3 行文档索引) | 无 CLAUDE.md / AGENTS.md / .cursor |
| 提交标题含 codex/claude/ai-/copilot | 2/838 = 0.24% | 0/1126 |
| AI co-author | 0 | 6/1126 = 0.5% |
| revert 比例 | 5/838 = 0.6% (4 条是 CI/Dockerfile) | 6 次 |
| 注释率 | 3.2% | 9.7% |
| 真实质量缺口 | CI 无 tsc 门禁, main 带既有类型错误 | 协议逆向类 bug 多 (`#417` 修后复现 `#433`) |

需要说明: 元数据会系统性低估 AI 参与 (AI 辅助常以人类作者身份无 trailer 提交), 所以只能说"查不到证据", 不能反证。反过来, 已观察到的质量问题 (CI 无类型检查、协议漂移 bug) 与 AI 参与之间没有因果关系 —— 四家都有协议逆向类 bug, 只是 NapCat 生态大、暴露面广。

## 五 复现命令

```bash
# 仓库克隆 (浅克隆够用)
git clone --depth 200 https://github.com/NapNeko/NapCatQQ        ~/git/NapCatQQ
git clone --depth 200 https://github.com/LagrangeDev/Lagrange.Core ~/git/Lagrange.Core
git clone --depth 200 https://github.com/LLOneBot/LuckyLilliaBot  ~/git/LuckyLilliaBot
git clone --depth 200 https://github.com/SnowLuma/SnowLuma        ~/git/SnowLuma

# 关键事实核查
grep -ri onebot ~/git/Lagrange.Core | head            # V2 零 OneBot 实现
grep -rhoE 'app_name.*' ~/git/*/src --include=*.ts    # 各家 app_name
curl -s https://api.github.com/repos/NapNeko/NapCatQQ       # star / created / pushed
curl -s "https://api.github.com/repos/SnowLuma/SnowLuma/releases?per_page=8"
```
