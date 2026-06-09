# LiarsDeck 大更新计划

## 1. 后端修复 (main.go)

### 1.1 Bug 修复
- `FireGun` 和 `CallLiar` 中 nil client 空指针崩溃
- `play_cards` 没有检查玩家打完手牌获胜的逻辑
- 移除废弃的 `rand.Seed`

### 1.2 新增功能
- 自动生成房间码 (6位字母)
- `IsHost` 房主身份 (第一个加入的人)
- `IsSpectator` 观战模式
- Host 可以移除玩家
- 创建/加入房间的 action 分离

## 2. 前端重写 (index.html)

### 2.1 UI 重设计
- 暗色卡牌主题，游戏氛围
- 所有功能集成：房间创建/加入 → 大厅 → 游戏 → 结束

### 2.2 新功能
- 房间自动生成 + 邀请链接展示
- 房主管理 (踢人)
- 观战模式
- 规则介绍弹窗
- 更好的牌面动画和视觉

## 3. 编译部署
- `go build -o liarsbar-web`
- `systemctl restart pdnode-liarsdeck`
