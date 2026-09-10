# SnowLuma 分析

- 仓库: <https://github.com/SnowLuma/SnowLuma>
- 本地快照: `~/git/SnowLuma`, HEAD `fb5f9b21` (2026-08-30), 浅克隆 depth=200
- 采集时间: 2026-09-10
- 定位: 面向 QQ 客户端的 TypeScript 互操作运行时, 把 QQ 原生会话桥接为 OneBot v11 动作与事件 (README.md:19)

## 1 定位与注入方式

- TypeScript + Node 22+, pnpm monorepo, 13 个包 (`packages/{bridge,common,core,mcp,onebot,proto-defs,protocol,proton,runtime,sdk}`)。
- 自研 native addon (`packages/runtime/native/snowluma-{linux-x64,linux-arm64,win32-x64}.node`) 用 **ptrace 注入运行中的 QQNT 进程**, 不同于 NapCat 的 Electron 补丁方式 (`packages/bridge/src/hook-manager.ts`)。
- 自我定位是 NapCat 兼容: `packages/onebot/src/streaming.ts:1` 注释 "Wire format mirrors NapCat exactly so its client scripts work unchanged"; README:96 说明 "项目参考了 LagrangeV2 的协议定义与 NapCatQQ 的实现思路"。

## 2 部署方式

| 平台 | 方式 |
|---|---|
| Linux | **仅官方 Docker**, 镜像自带 LinuxQQ 3.2.32 + Xvfb + VNC/noVNC + supervisord |
| Windows | 需用户自装桌面 NTQQ + 解压跑 `launcher.bat` (`docs/guide/deploy/windows.html`) |

官方镜像维护良好: Docker Hub 284 tags, 每次 release 自动派发构建, `latest` 当日同步 (`.github/workflows/release.yml` 的 dispatch-docker)。

**比 NapCat Docker 复杂**: 必须 `--cap-add=SYS_PTRACE --security-opt seccomp=unconfined --shm-size=1g`, 5 个端口, 用 noVNC 扫码; NapCat 只需一条 `docker run` + 6099 WebUI 内扫码, 无额外 capability (`NapCat-Docker` README)。

## 3 更新稳定性: 年轻且快

- GitHub created 2025-06-13, 但代码首提交 **2026-05-02**, 实际只有约 4 个月历史。
- 1126 提交核实无误 (月 390/282/274/180, 递减); 99 个 tag 约 1.2 天/版, 当前 `v1.14.15`。
- 无 BREAKING 声明但已有 6 次 revert; `RoadMap.md` 标题为「V2.0.0 前计划」, 预示 V2 会有破坏性变更。
- 单一主作者占 78% (876/1126)。
- issue: 254 个, 已关 245, 中位关闭 0.68 天, 17% not_planned。

## 4 OneBot 11 支持

- 连接方式: ws server/client + http server/post (`packages/onebot/src/config.ts`)。
- 文档称 193 个动作, 源码 grep 到 187 个名字。
- 本次关注端点全部存在:

| 端点 | 位置 |
|---|---|
| `send_group_forward_msg` | `packages/onebot/src/actions/extended.ts:2477` |
| `send_private_forward_msg` | `extended.ts:2494` |
| `friend_poke` | `extended.ts:402` |
| `group_poke` | `extended.ts:415` |
| `delete_msg` | `actions/message.ts:136` |
| `get_msg` | `actions/message.ts:115` |
| `get_forward_msg` | `extended.ts:2512` |

缺口: `.handle_quick_operation_async` 未注册 (issue `#434` 仍 open, 影响 NoneBot1)。

## 5 app_name (对 NothingBot_v4 关键)

`app_name = 'SnowLuma'`, `app_version = '<版本>-node'`, `protocol_version = 'v11'` (`packages/onebot/src/actions/info.ts:60-66`)。

影响: 客户端需要为 `SnowLuma` 新增 AppName 分支, 否则 NothingBot_v4 `M_ForwardSlice.go:108` 的 switch 会报 `unsupported onebot implementation`。

## 6 资源占用

- 每个账号一个 QQ (Chromium) 进程, 约 300-500MB; `shm >= 1G`; 默认关 GPU 以防 SwiftShader 泄漏 (Docker 文档 + Docker.Framework README)。
- SnowLuma 自身进程的内存**未找到官方基准**。

## 7 AI 开发占比的客观核查

| 项 | 结果 | 证据 |
|---|---|---|
| AI 协作/提示词文件 | 无 CLAUDE.md / AGENTS.md / .cursor / CONTEXT.md | `find`, `git ls-files` |
| 提交标题 AI 痕迹 | `codex\|claude\|ai-\|copilot` = **0/1126** | `git log --oneline \| grep -ciE` |
| AI co-author | 6/1126 = 0.5%; `codex/*` 分支 2/167 PR | `git log --format=%B`, GitHub API |
| 已装的 AI 工具 | OpenAI Codex GitHub App 与 Copilot Autofix (PR `#425`, 3 提交) | PR `#425` |
| 代码风格 | src 注释率 9.7%, TODO 仅 2 处, 无成片 LLM 注释/过度文档; 提交信息冗长统一 | awk 统计 |
| 实际缺陷分布 | 集中在协议逆向与 QQ 版本漂移, 非 AI 空转: 139 个 bug issue, `message_id`/reply id 修后复现 (`#417` -> `#433`), 同一 bug 重报 4 次 (`#340/341/342/370`), 群权限设置 3 次 (`#388/392/411`), QQ 崩溃/注入失败 (`#312/366/350`) | issue 列表 |
| 工程反面 | 369 个测试文件, 82k 测试行, CI 有 typecheck + 测试 + catalog 校验 | CI 配置 |

结论: 未见 AI 生成痕迹的统计证据, 且测试与 CI 反而比多数同类项目更完整; 它的真实短板是**项目年龄与协议漂移**, 不是代码来源。

## 8 许可证与风险

- EULA/LICENSE: 无 AI 条款; 限制性内容为源码可见的非商业许可, 禁商业使用与公开发布衍生版; 专有 `.node` 禁逆向。
- **EULA 5.4 明确禁止把专有组件并入第三方 Docker 镜像或自动化脚本部署** (需书面授权) —— 自建镜像/自写部署脚本前必须确认。
- hook 与 QQ 版本强绑定, 官方要求黑洞 `qqpatch.gtimg.cn` 冻结 QQ。
- 项目仅 4 个月、单人主导 (78%), 存在维护者停更风险。

## 9 结论

**优点**: OneBot 覆盖广 (187-193 动作), 官方 Docker 随版本自动发版, issue 响应极快 (中位 0.68 天), 工程规范 (369 测试文件 + CI typecheck)。

**风险**: 项目年轻 (4 个月) 且单人主导; hook 与 QQ 版本强绑定; Docker 需 `ptrace` + noVNC, 比 NapCat 繁琐; 许可证限制自建镜像; 路线图预示 V2 破坏性变更。

**建议**: 对"部署省心 + 更新稳定"这一诉求, NapCat 仍更省心; SnowLuma 适合愿意高频跟进更新、看重 OneBot 覆盖面的用户, 建议先观察一个 V2 周期再考虑。
