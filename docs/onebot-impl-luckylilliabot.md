# LuckyLilliaBot (LLBot) 分析

- 仓库: <https://github.com/LLOneBot/LuckyLilliaBot>
- 本地快照: `~/git/LuckyLilliaBot`, HEAD `b18cbb3` (2026-09-07), 浅克隆 depth=200
- 采集时间: 2026-09-10
- 定位: OneBot 11 / Satori / Milky 服务端, Node.js (Node >= 24)

## 1 它是什么

不是"重写", 而是**同一仓库改名续作**: `github.com/LLOneBot/LLOneBot` 现 301 到 `LuckyLilliaBot`, 仓库 id 709730860, created_at **2023-10-25**。旧形态是 LiteLoaderQQNT 插件 (`doc/更新日志.txt`: "支持加载类似 LiteLoaderQQNT 的 js 插件")。

现为独立 Node.js 程序, 两条注入路径:

| 模式 | 机制 | 需要 QQ 客户端 |
|---|---|---|
| Direct (直连) | 纯代码复刻 QQ NT 协议 (native Rust sign-proxy + TCP + SSO), 模式由 argv `--pmhq-port=` 判定 | 否 |
| PMHQ (有头) | 真 QQ 客户端跑在独立 `linyuchen/pmhq` 容器, 注入 dll, llbot 经 WS/HTTP 连它 | 是 |

证据: `src/main/qqProtocol/direct.ts`, `src/main/qqProtocol/pmhq.ts`, `src/common/utils/environment.ts`, `docs/docker.md`。

## 2 部署方式

| 方式 | 内容 |
|---|---|
| Docker (推荐) | 交互向导 `LLBot-Docker.sh` -> 生成 compose -> 起容器 -> WebUI 扫码 (`script/install-llbot-docker.sh`, `docker/Dockerfile`, `docker/startup.sh`) |
| 二进制 | `LLBot-CLI-{win-x64,linux-x64,linux-arm64,macos-arm64}`, `LLBot-Desktop-{win,macos}`, `LLBot.zip` (releases v8.1.10) |
| 源码 | `yarn install` / `yarn build` / `node llbot.js` |

镜像 `linyuchen/llbot`: 57 tags, `latest` = 8.1.10 与 release **同日推送** (2026-09-04), amd64 + arm64, 压缩 153MB/arch, release 事件自动构建 (`.github/workflows/docker.yml`), 官方维护活跃。

对比 NapCat: 直连模式单容器、无 QQ 客户端、不追 QQ 版本, 更简单; 但**多一步必填 Auth Token** —— 安装脚本第一项就是 `Auth Token` (`script/install-llbot-docker.sh:6-23`), 去 `auth.luckylillia.com` 获取, 两种模式都强制要求。

## 3 更新稳定性

- 发布节奏: 8.1.3 (07-30) -> 8.1.10 (09-04) 共 8 个 patch, 约每周一版, 稳定可预期。
- 近 838 提交 (2025-12 至 2026-09) 月度 37/14/87/99/213/66/85/187/43/7。
- 维护者: idranme 454 + linyuchen 370 = 98%, 窗口内共 10 名贡献者。
- issue 修复周期 1-4 天 (`#855` 09-02 至 09-04, `#834` 08-05 至 08-07, `#835` 08-06 至 08-07, `#836` 08-10 至 08-14)。
- 破坏性改动少: 1030 行 changelog 仅 3 处 (v7.7.0 `pinned` int->bool, v7.8.3 移除 `groupAll`, v7.2.0 zip/js 改名)。
- **现存风险**: `#857` (2026-09-10, open) "docker 反复掉线", 扫码成功但登录不上, 必须重启容器。

## 4 OneBot 11 支持

连接方式 4 类: ws, ws-reverse, http, http-post (脚本里叫 WebHook) —— `src/onebot11/adapter.ts`, `src/common/types.ts:1-37`。无独立 OneBot webhook 类型 (webhook 仅 Milky 有)。ActionName 共 120 个 (`src/onebot11/action/types.ts`), notice 事件 18 个。文档: <https://www.luckylillia.com/guide/develop#onebot11-协议>

本次关注端点全部存在:

| 端点 | 位置 | 实现 |
|---|---|---|
| `send_group_forward_msg` | `action/types.ts:106` | `go-cqhttp/SendForwardMsg.ts` |
| `send_private_forward_msg` | `107` | 同上 |
| `friend_poke` | `48` | `llbot/user/FriendPoke.ts` |
| `group_poke` | `47` | `llbot/group/GroupPoke.ts` |
| `delete_msg` | `82` | `msg/DeleteMsg.ts` |
| `get_msg` | `78` | `msg/GetMsg.ts` |
| `get_forward_msg` | `115` | `go-cqhttp/GetForwardMsg.ts` |

另有 `send_poke` (65)、`send_forward_msg` (43)、`forward_friend/group_single_msg`。

## 5 app_name (对 NothingBot_v4 关键)

**`'LLOneBot'`**, 硬编码 (`src/onebot11/action/system/GetVersionInfo.ts:11`, 类型 `src/onebot11/types.ts:404`); `protocol_version: 'v11'`, `app_version` = `src/version.ts`。

影响: `GetVersionInfoCache().AppName` 会走 LLOneBot 分支, 而 NothingBot_v4 `M_ForwardSlice.go:108` 的 switch 目前只认 `NapCat.Onebot` / `Lagrange.OneBot`, 会报 `unsupported onebot implementation`。好消息是它的 `get_forward_msg` 同时接受 `id` 和 `message_id` (gocq resId), 返回的 `messages[].content` 直接是消息段数组, 与 NapCat 语义一致, 不需要像 Lagrange 那样二次取节点 (`action/go-cqhttp/GetForwardMsg.ts`), 因此可以并进 NapCat 分支处理。

## 6 资源占用

- 单 Node.js 进程 (`src/main/main.ts`, cordis 插件容器)。
- Direct 模式**不含 QQ 客户端**; PMHQ 模式额外一个真 QQ 容器 (compose 双服务, `docs/docker.md`)。
- 镜像约 540MB 磁盘 / 153MB 压缩。**未找到实测内存/CPU 数据**, 仓库仅定性描述直连"省内存"、PMHQ"吃内存"。

## 7 AI 开发占比的客观核查

| 项 | 结果 | 证据 |
|---|---|---|
| AI 协作/提示词文件 | 仅 `CLAUDE.md` **398 字节**, 3 行文档索引; `.gitignore` 排除 `.claude/`; 未找到 AGENTS.md / .cursorrules / copilot-instructions.md | `CLAUDE.md`, `.gitignore` |
| 提交标题 AI 痕迹 | `codex\|claude\|ai-\|copilot` 仅 **2/838 = 0.24%**; `Co-Authored-By` **0 条**; "Generated with/Claude Code" **0 条** | `git log --oneline \| grep -ciE`, `git log --format=%B \| grep -ci co-authored-by` |
| LLM 风格注释/文档 | 注释 1384 行 / 42848 TS 行 = 3.2%; 连续注释块 (>=5 行) 23 处, 最长 31 行 | awk 统计 |
| revert 频率 | **5/838 = 0.6%**, 4 条是 CI/Dockerfile/changelog, 仅 1 条行为回滚 (`6db4643`) | `git log --oneline \| grep -i revert` |
| 唯一 codex 痕迹 | PR `#856` `codex/fix-direct-session-auth` (外部贡献者 HLfromZ): 带 24 条回归测试、在旧代码上复现 7 项预期失败、被 Sourcery 指出竞态后补 15 项测试, 质量**高于**平均 PR | GitHub API issue 856 |
| 真实质量缺口 | main 上存在既有类型错误 `src/ntqqapi/helper/messageBuilding.ts:74` TS2322, 而 CI (`.github/workflows/test.yml`) 只跑 unit/webui 测试与 build, **无 tsc 门禁** | `test.yml`, PR `#856` 正文 |

结论: "AI 开发占比太多"在可核查的 git 元数据上不成立; 但元数据会低估 AI 参与 (无 trailer 提交无法区分), 所以只能说查不到证据, 不能反证。已发现的工程质量问题 (CI 无类型检查、遗留类型错误、9 月掉线 bug) 与 AI 参与之间没有因果关联。

## 8 结论

**优点**: 2023 年起同仓库持续演进 (不是新项目); 发布与 Docker 镜像同日、多架构、每周 patch; 2 名核心维护者, issue 1-4 天响应; 直连模式单容器零 QQ 客户端, 部署比 NapCat 少一层 QQ 版本适配; 120 个 OneBot 11 端点, 本项目需要的全在; 破坏性改动极少。

**风险**:

1. 直连模式强依赖作者中心化服务 `api-auth.luckylillia.com` 与 Auth Token, 且 token 有"可用 QQ 数量上限" (`src/main/config/index.ts:355`, `src/main/qqProtocol/direct.ts:472`) —— 服务停或 token 失效则完全不可用, 这是 NapCat 没有的单点。
2. `#857` Docker 反复掉线未修。
3. CI 无类型检查, main 带既有类型错误。
4. 生态、文档、社区体量小于 NapCat, 核心仅 2 人。

**建议**: 若迁移, 直连模式最省资源但把命脉交给了第三方 auth/sign 服务; 不能接受该单点就继续用 NapCat Docker。若试用, 用 PMHQ 模式 (真 QQ 客户端) 可绕开签名服务, 代价是双容器与更多内存。
