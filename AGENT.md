# AGENT.md - Liar's Deck 项目规范与记忆库

本文档记录了 **Liar's Deck** 项目的核心架构、不可违背的开发规范、关键业务记忆以及未来规划。每次与 AI 协作时均需严格遵循本文档要求。

---

## 🏗️ 1. 项目架构与技术栈

- **后端核心**：Go 1.24 + Echo v4 Web 框架 + Gorilla WebSocket
- **随机与安全**：全系统洗牌、子弹装填、房间码生成均使用 `crypto/rand` 密码学安全随机数
- **前端核心**：Vue 3 (Composition API) + Vite
- **UI 设计系统**：Tailwind CSS v4 (`@tailwindcss/vite`) + DaisyUI v5 (`daisyui@5.7.22`)
- **认证体系**：OAuth2 Authorization Code Flow with PKCE，集成 pdnode / Zitadel OIDC 单点登录
- **部署与自更新**：纯 Go 原生解压与热替换机制，支持 `--wait-pid` 优雅释放端口重启与 GitHub Actions 跨平台 Release

---

## 🧠 2. 关键业务记忆与不可违背规则 (Important Memories & Rules)

### 📌 规则一：更新日志规范（发布前必提醒）
1. **存放路径**：每个 Tag 独立存放于 `changelogs/vX.Y.Z.md`（如 `changelogs/v2.3.4.md`）；
2. **发布流程约束**：
   - 每次准备发布新版本/Tag 时，**AI 必须主动提醒用户编写该版本的更新日志文件**；
   - 更新日志文件必须**由用户亲自编写/确认**，禁止 AI 擅自代写后续版本（历史版本 v1.0.0~v2.3.4 已由系统一次性初始化录入）；
   - 用户确认更新日志并提交后，方可打 Tag 与推送 Release。

### 📌 规则二：Git 提交严格使用 GPG 签名
1. 用户的开发环境已全局开启 GPG 签名（`commit.gpgsign=true`）；
2. 代码提交必须使用 `git commit -S -m "..."`；
3. 若出现 GPG 弹窗等待或签名阻塞，AI 需提示用户在本地处理，**严禁私自添加 `--no-gpg-sign`**。

### 📌 规则三：账号身份唯一绑定（以 `sub` 为唯一主键）
1. 玩家席位认领、断线重连、同账号防多开、游戏暂停/超时处死判定等全生命周期，**必须严格基于不可变的 Zitadel OIDC `sub` claim (`token` / `p.ID` / `PausedPlayerID`)**；
2. **严禁依赖 `nickname` 做任何唯一性或鉴权比对**（Zitadel 允许同名昵称，且用户可随时修改昵称）。

### 📌 规则四：断线保护为 30 秒累积时间池 (`DisconnectGraceRemaining`)
1. 玩家单局拥有 30 秒断线保护时间池；
2. 每次掉线暂停游戏并按实际离线时长精确扣减，重连后保留剩余秒数，彻底避免网络偶发抖动瞬断被直接淘汰的误杀；
3. 暂停/超时处死全生命周期严格绑定 `PausedPlayerID`。

### 📌 规则五：保持经典的每轮随机座次博弈机制
1. 在 `internal/game/game.go` 的 `StartRound()` 中，**必须保留每小轮发牌开局时的 `CryptoShuffle` 玩家数组打乱逻辑**；
2. 这一机制决定了每轮玩家出牌的上下家关系与【质疑 (Call Liar)】的博弈顺位。

### 📌 规则六：UI 设计系统规范（深邃暗黑黑曜石）
1. 严格使用 Tailwind CSS v4 + DaisyUI v5 构建；
2. 主题色调为深邃黑曜石（Slate-950 / Zinc-900 / 电光科技蓝 / 猩红等现代高对比配色），**严禁出现老旧金色**。

### 📌 规则七：多层级动态版本注入
1. Go 后端通过 `-ldflags` -> `runtime/debug.ReadBuildInfo()` -> `git describe --tags` 三重回退探测，通过 `/api/version` 实时暴露；
2. 前端通过 Vite 注入 `__APP_VERSION__`，并在管理面板中动态与后端及 GitHub API 比对，严禁硬编码旧版本号。

---

## 🗺️ 3. 项目路线图与未来规划 (Roadmap)

- [ ] **AI 电脑陪玩模式 (Bot Players)**：
  - 支持单人开房练习或房间人数不足 4 人时一键添加具有不同博弈心理性格（稳健型、激进抓谎型、虚张声势型）的 AI Bot。
- [ ] **多语言本地化扩展 (i18n Expansion)**：
  - 在现有中英双语基础上，扩展日语、韩语、西班牙语等语种。
- [ ] **音效与视效深度增强**：
  - 新增音量调节滑块组件；
  - 增加左轮手枪开火火花动效与击中屏幕震动粒子特效。
- [ ] **战绩与成就统计系统**：
  - 记录玩家历史对局数据（胜率、抓谎成功率、致命空枪生还率、荣誉称号等）。
