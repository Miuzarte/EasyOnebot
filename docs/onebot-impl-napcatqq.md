# NapCatQQ 分析

- 仓库: <https://github.com/NapNeko/NapCatQQ>
- 本地快照: `~/git/NapCatQQ`, HEAD `109d0c1` (2026-09-09), 浅克隆 depth=200
- 采集时间: 2026-09-10
- 定位: 基于 NTQQ 的现代协议端框架, TypeScript

## 1 部署方式

| 方式 | 步骤 | 需自装 QQ | 官方 |
|---|---|---|---|
| Docker (推荐) | `docker run` 一条, 或 compose | 否, 镜像内置 LinuxQQ | 是, `NapNeko/NapCat-Docker` |
| 一键脚本 | `curl` + `bash install.sh`, 支持 `--tui` / `--docker` | 脚本处理 | 是, NapCat-Installer |
| 二进制 Shell | 下 `NapCat.Shell.zip` + 装 QQ + `launcher.bat` | 是 | 是 |
| 免安装整合包 | `NapCat.Shell.Windows.OneKey.zip` 解压 + 运行 Installer (内置 QQ) | 否 | 是, 仅 Windows |
| Node 整合包 | `NapCat.Shell.Windows.Node.zip` 解压 + `napcat.bat` (内置 node.exe) | 是 (要拷 QQ dll) | 是 |
| LiteLoaderQQNT 插件 | `NapCat.Framework.zip` 作为插件加载 | 是 | 是 |
| 其他 | AppImage / Termux / Mac / Nix / 宝塔 / 1Panel / Railway / Zeabur | 部分免 | 部分社区 |

Docker 镜像与主仓库同属 NapNeko 组织, 由主仓库 `release-downstream.yml:88-105` 自动触发发布, 属于官方维护, 但 Dockerfile 不在本仓库 (在 NapCat-Docker)。

## 2 更新稳定性

- 77 个 tag, 最新 `v4.18.19` = 2026-08-14; 节奏 1-3 周, 但 08-04 / 08-06 / 08-07 三天连发。
- **无 CHANGELOG 文件**; release note 由 LLM 生成 (`release-publish.yml` 里 `OPENAI_API_URL=openrouter`, `OPENAI_MODEL=google/gemma-4-26b-a4b-it:free`)。
- `package.json` 的 version 恒为 `0.0.1` (私有 monorepo), 真实版本来自 `VITE_NAPCAT_VERSION` (`napcat-common/src/version.ts`)。
- **被动升级压力真实存在**: 新 QQ 版本会直接不可用 —— `#2043` "PacketBackend 不支持当前QQ版本架构 9.9.35-52892", `#2042` 重启后登录卡死, `#2047` GPU/zygote 崩溃循环。release note 固定警告 "QQ 版本推荐 40768+, 最低 40768"。

## 3 资源占用

- 进程模型: 寄生在 QQNT (Electron/Chromium) 内, worker 用 fork 或 Electron `utilityProcess`; 独立 Node 路径也要 QQ 的 `wrapper.node` / `QQNT.dll` / `ncnn.dll` (`packages/napcat-shell/napcat.ts:58-260`, `packages/napcat-shell-loader/launcher.bat`)。Chromium 进程树始终存在。
- docs 自称 Shell 版"低内存", 但**未找到官方内存/CPU 基准数值**。
- 已知性能 issue: `#1218` 内存异常致服务器崩溃, `#1047` 疑似内存泄漏未释放 (not_planned), `#1771` 插件定时任务 CPU 从 <1% 到约 100%, `#2047` GPU 进程崩溃。

## 4 OneBot 11 支持

连接方式齐全 (`packages/napcat-onebot/config/config.ts:1-80`): websocket-server (反向 WS, 默认 3001), websocket-client (正向 WS), http-server, http-sse-server (SSE), http-client (webhook POST, 带 `x-signature` HMAC-SHA1, `network/http-client.ts:23-27`), plugin。

共 184 个 action (`action/router.ts`)。关键端点全部存在且非桩:

| 端点 | 位置 |
|---|---|
| `send_group_forward_msg` / `send_private_forward_msg` | `action/router.ts:103-104`, `go-cqhttp/SendForwardMsg.ts` |
| `friend_poke` / `group_poke` | `action/router.ts:173-174`, `packet/SendPoke.ts:40-45` |

与 EasyOnebot 的 43 个 `.Lgr` 端点比对: 36 个支持, 缺 7 个 —— `delete_group_file_folder` (NapCat 名为 `delete_group_folder`, 属命名差异)、`rename_group_file_folder`、`fetch_mface_key`、`set_group_reaction`、`set_group_bot_status`、`send_group_bot_callback`、`upload_image` (NapCat 只有 `upload_image_to_qun_album`)。

## 5 app_name

`app_name = 'NapCat.Onebot'`, `protocol_version = 'v11'`, `app_version = napCatVersion` (`packages/napcat-onebot/action/system/GetVersionInfo.ts:29-32`)。

注意文档示例里写的 `'NapCatQQ'` 是过时示例 (`example/SystemActionsExamples.ts:36`), 以源码为准。这个值正好命中 NothingBot_v4 `M_ForwardSlice.go` 的 `NAPCAT_APPNAME` 分支。

## 6 代码质量与 AI 痕迹

- TypeScript + pnpm monorepo, 27 个 package, 约 696 源文件 / 约 71k 行。
- 测试偏少: 10 个 test 文件 / 约 2981 行, 约 4% (`packages/napcat-test/*.test.ts`); CI 有 typecheck + vitest (`.github/workflows/build-ci.yml:12-34`)。
- 未找到 CLAUDE.md / AGENTS.md / .cursorrules / copilot 配置。
- 有明确 AI 协作但仅限发布流程: `.github/prompt/release_note_prompt.txt` 是给 LLM 的发布说明提示词, `default.md` 为兜底模板, CI 调 OpenRouter 生成 release note。
- 提交风格规范 (`fix:` / `feat:` / `chore:` / `style:`), 维护者活跃 (MliKiowa 8 天内关闭 `#2040`), dependabot 在跑, issue 普遍有回复。

## 7 风险

- **许可证**: LICENSE 是自定义 "Limited Redistribution License" (非 OSI), 禁止商业用途、禁止未授权基于 NapCat 代码开发项目; README:81-87 声明"本仓库仅用于提高易用性, 实现消息推送类功能"。涉及商业或二次分发需注意。
- **停摆风险**: 与腾讯 QQNT 高耦合, QQ 客户端升级就可能需要跟进 (packetBackend 白名单/偏移表, 见 `#2043`), 且旧 QQ 会被强制升级/登录失效, 无法长期锁版本; Docker 镜像把 QQ 版本钉在构建时 (`Dockerfile` 里 `linuxqq_3.2.30-50969`) 是缓解手段。
- **封号风险**: 未找到量化证据, README 仅泛指"使用请遵守当地法律法规"; 属于同类 NTQQ 协议端的共有风险。
- **版本纪律**: 无正式 CHANGELOG, 发布说明由 LLM 生成, 靠 release 页警告用户自行匹配 QQ 版本, 升级前需要人工确认对应 QQ 版本。
