import { ref, watchEffect } from 'vue';

const dict = {
  zh: {
    app_title: "Liar's Deck",
    app_subtitle: '感谢Hermes (Model: Deepseek V4 Flash) & Gemini 3.7 Flash!',
    lobby_create_title: '创建新房间',
    lobby_join_title: '输入房间码加入',
    lobby_spec_title: '观战模式',
    nickname: '玩家昵称',
    nickname_ph: '输入你的称呼',
    room_code: '房间码',
    room_code_ph: '6位大写字母/数字',
    create_btn: '创建房间',
    join_btn: '加入对战',
    spectate_btn: '进入观战',
    back: '返回',
    rules_btn: '规则说明',

    // Status
    status_waiting: '等待其他人加入中',
    status_playing: '对局进行中',
    status_paused: '对局暂停中',
    status_game_over: '游戏结束',

    // Seat
    host_tag: '房主',
    spec_tag: '旁观',
    me_tag: '自己',
    dead_tag: '已出局',
    offline_tag: '断线中',
    watching_tag: '观战中',
    bullets_label: '剩余子弹',
    hand_count_label: '手牌',
    ready_btn: '准备就绪',
    unready_btn: '取消准备',
    ready_status: '已就绪',
    unready_status: '未准备',
    kick_btn: '踢出房间',

    // Header
    invite_btn: '复制邀请链接',
    copied_toast: '邀请链接已复制到剪贴板',
    audio_on: '声音开启',
    audio_off: '已静音',
    timeout_warn: '操作倒计时',
    pause_warn: '重连倒计时',

    // Table
    table_card_label: '本轮真牌',
    wild_card_tip: '2 为万能牌，可作为任何真牌打出',
    cards_on_table: '张暗牌在桌面',
    winner_label: '胜出者',
    play_again_btn: '再战一局',
    wait_host_reset: '等待房主重新开始',
    start_game_btn: '开始对局',
    need_more_players: '需至少 2 名已准备玩家',
    all_ready_needed: '等待所有玩家准备',

    // Hand
    my_hand_title: '我的手牌',
    selected_count: '已选',
    max_cards_tip: '可选 1~3 张',

    // Actions
    play_cards_btn: '出牌',
    call_liar_btn: '质疑 (CALL LIAR)',

    // Events
    event_liar_alert: '发起质疑！',
    event_calls_out: '怀疑',
    event_liar_claim: '在说谎！',
    event_cards_revealed: '翻开底牌验牌',
    event_bang_title: '致命枪响！',
    event_bang_sub: '击中实弹，遗憾淘汰出局',
    event_click_title: '空包弹！',
    event_click_sub: '扣动扳机为纯空枪，逃过一劫',

    // Spectator
    spectator_banner: '你正在旁观本场对局',

    // Rules Modal
    rules_title: '游戏规则与牌理',
    rule_goal_title: '对局目标',
    rule_goal_desc: '存活至全桌最后一人，或打光所有手牌并在下家质疑中证明清白。',
    rule_deck_title: '牌组构成 (共24张)',
    rule_deck_desc: '包含 K、Q、A 及万能牌 2，每种各 6 张。2 可代替本轮任何真牌。',
    rule_turn_title: '轮流出牌',
    rule_turn_1: '每轮随机翻开一张“真牌”，存活玩家各分发 5 张手牌。',
    rule_turn_2: '到你的回合，选 1~3 张牌背面朝下打出，声称其全为本轮真牌。',
    rule_turn_3: '你可以诚实出牌，也可以虚张声势（撒谎）。',
    rule_liar_title: '质疑与裁决',
    rule_liar_1: '下家可选择继续出牌，或对上家发起【质疑说谎】。',
    rule_liar_2: '翻开底牌：若存在非真牌且非 2 的牌 → 抓谎成功，出牌者开枪惩罚！',
    rule_liar_3: '若全为真牌或 2 → 质疑失败，质疑者自己开枪惩罚！',
    rule_gun_title: '左轮轮盘手枪与断线规则',
    rule_gun_1: '每位玩家持有一把 6 发弹仓手枪（5 发空包弹，1 发致命实弹，随机装填）。',
    rule_gun_2: '局中玩家首次断线将暂停游戏 30 秒等待其重连；若超时未归或二次断线，将直接被判定出局！',
    rule_got_it: '我已了解',

    // Disconnect & Pause Modals
    dc_title: '网络连接已中断',
    dc_desc: '与酒馆牌桌的连接已断开，正在尝试自动重新连接...',
    dc_reconnecting: '正在重新连接中',
    dc_retry_btn: '立即重试连接',
    dc_exit_btn: '退出到大厅',
    dc_success: '重新连接成功！',
    dc_grace_tip: '你有 30 秒时间重新连接，超时将被判定直接淘汰（每人仅限 1 次机会）。',

    pause_modal_title: '对局已暂停',
    pause_modal_desc: '玩家【{name}】断线，等待其重新连接中...',
    pause_countdown_label: '重连剩余时间',
    pause_tip: '30 秒内重连将恢复对局；若超时仍未归来，系统将直接处死该玩家。',

    // Admin Modal
    admin_title: '管理控制台',
    admin_auth_title: '管理员身份验证 (Ctrl + X)',
    admin_auth_desc: '请输入管理员密钥以解锁服务端维护与热更新操作：',
    admin_auth_ph: '输入 ADMIN_SECRET 管理密钥',
    admin_unlock_btn: '解锁管理面板',
    admin_version_card: '系统版本与 GitHub Releases 维护',
    admin_curr_version: '当前版本',
    admin_latest_version: '最新版本',
    admin_check_btn: '检查 GitHub Releases 更新',
    admin_update_btn: '更新并重启服务端',
    admin_updating: '正在拉起 update.go 后台程序...',
    admin_stats_card: '实时运行状态',
    admin_active_rooms: '当前活跃房间',
    admin_active_players: '在线玩家总数',
    admin_close_btn: '关闭控制台',
    admin_broadcast_card: '全服实时广播',
    admin_broadcast_ph: '输入要向全服所有在线房间广播的内容...',
    admin_broadcast_send_btn: '发送全服广播',
    admin_broadcast_sending: '正在发送广播...',
    admin_broadcast_preset_1: '服务器即将在 5 分钟后进行更新维护，请尽快完成当前对局。',
    admin_broadcast_preset_2: '欢迎来到 Liar\'s Deck 骗子酒馆！祝各位玩家游戏愉快！',
    // Login Modal
    login_modal_title: "Liar's Deck · 酒馆大门",
    login_modal_desc: '酒馆正在营业中。请出示您的身份凭证以入座牌桌。',
    login_btn: '登录 / Login',
    login_loading_btn: '正在跳转授权中心...',
    login_secure_tip: '使用 pdnode 统一认证中心 · PKCE 传输加密保护',
    logout_btn: '登出',
    logout_confirm: '确定要退出当前账号并登出吗？',
    logged_in_as: '已登录账号',

    // Logs
    battle_log_title: '对局动态',

    // Errors
    err_enter_nickname: '请填写玩家昵称',
    err_enter_code: '请输入6位房间码',
    reconnecting: '连接中断，正在重连...',
  },
  en: {
    app_title: "Liar's Deck",
    app_subtitle: 'Thanks to Hermes (Model: Deepseek V4 Flash) & Gemini 3.7 Flash',
    lobby_create_title: 'Host a Room',
    lobby_join_title: 'Join with Code',
    lobby_spec_title: 'Spectate Table',
    nickname: 'Handle',
    nickname_ph: 'Your player name',
    room_code: 'Room Code',
    room_code_ph: '6-character code',
    create_btn: 'Create Room',
    join_btn: 'Join Table',
    spectate_btn: 'Spectate',
    back: 'Back',
    rules_btn: 'House Rules',

    // Status
    status_waiting: 'Waiting for Players',
    status_playing: 'Hand in Progress',
    status_paused: 'Match Paused',
    status_game_over: 'Showdown Over',

    // Seat
    host_tag: 'Host',
    spec_tag: 'Spectator',
    me_tag: 'You',
    dead_tag: 'Eliminated',
    offline_tag: 'Offline',
    watching_tag: 'Watching',
    bullets_label: 'Chamber',
    hand_count_label: 'Cards',
    ready_btn: 'Ready Up',
    unready_btn: 'Unready',
    ready_status: 'Ready',
    unready_status: 'Not Ready',
    kick_btn: 'Kick',

    // Header
    invite_btn: 'Copy Invite',
    copied_toast: 'Invite link copied',
    audio_on: 'Sound On',
    audio_off: 'Muted',
    timeout_warn: 'Turn Timer',
    pause_warn: 'Pause Timer',

    // Table
    table_card_label: 'Table Card',
    wild_card_tip: 'Card 2 is Wild and counts as the Table Card',
    cards_on_table: 'cards played face-down',
    winner_label: 'Sole Survivor',
    play_again_btn: 'Deal Again',
    wait_host_reset: 'Waiting for host to deal...',
    start_game_btn: 'Deal Cards',
    need_more_players: 'Requires 2+ ready players',
    all_ready_needed: 'Waiting for ready status',

    // Hand
    my_hand_title: 'Your Hand',
    selected_count: 'Selected',
    max_cards_tip: 'Pick 1 to 3 cards',

    // Actions
    play_cards_btn: 'Play Face-Down',
    call_liar_btn: 'CALL LIAR',

    // Events
    event_liar_alert: 'CALLING LIAR!',
    event_calls_out: 'challenges',
    event_liar_claim: 'for bluffing!',
    event_cards_revealed: 'Verifying Cards',
    event_bang_title: 'FATAL SHOT!',
    event_bang_sub: 'Hit live round and eliminated!',
    event_click_title: 'BLANK!',
    event_click_sub: 'Dry fire! Lucky escape!',

    // Spectator
    spectator_banner: 'You are currently observing the table',

    // Rules Modal
    rules_title: 'Tavern Rules & Mechanics',
    rule_goal_title: 'The Objective',
    rule_goal_desc: 'Be the last one standing, or discard all cards and survive the final challenge.',
    rule_deck_title: 'Deck (24 Cards)',
    rule_deck_desc: 'K, Q, A and Wild 2s (6 each). 2s match any active Table Card.',
    rule_turn_title: 'Playing a Turn',
    rule_turn_1: 'A Table Card is flipped face-up. Each alive player draws 5 cards.',
    rule_turn_2: 'On your turn, play 1-3 cards face-down and claim they match the Table Card.',
    rule_turn_3: 'Play honestly or pull a daring bluff.',
    rule_liar_title: 'The Challenge',
    rule_liar_1: 'The next player may play their turn or call out the previous player.',
    rule_liar_2: 'If any card fails to match the Table Card and is not a 2 -> Bluffer pulls the trigger!',
    rule_liar_3: 'If all cards were valid -> Challenger pulls the trigger!',
    rule_gun_title: 'Russian Roulette & Disconnect Rules',
    rule_gun_1: 'Every player holds a 6-chamber cylinder with 1 fatal live round.',
    rule_gun_2: 'First disconnect pauses the game for 30s. If the player fails to reconnect or disconnects again, they are eliminated!',
    rule_got_it: 'Understood',

    // Disconnect & Pause Modals
    dc_title: 'Connection Interrupted',
    dc_desc: 'Connection to the table was lost. Attempting to reconnect automatically...',
    dc_reconnecting: 'Reconnecting...',
    dc_retry_btn: 'Retry Now',
    dc_exit_btn: 'Exit to Lobby',
    dc_success: 'Reconnected successfully!',
    dc_grace_tip: 'You have 30 seconds to reconnect before elimination (1-time chance only).',

    pause_modal_title: 'Match Paused',
    pause_modal_desc: 'Player [{name}] disconnected. Waiting for reconnection...',
    pause_countdown_label: 'Grace Time Remaining',
    pause_tip: 'Match will resume if reconnected in 30s. Player will be eliminated upon timeout.',

    // Admin Modal
    admin_title: 'Admin Console',
    admin_auth_title: 'Admin Authentication (Ctrl + X)',
    admin_auth_desc: 'Enter your admin secret key to unlock server maintenance & hot updates:',
    admin_auth_ph: 'Enter ADMIN_SECRET',
    admin_unlock_btn: 'Unlock Console',
    admin_version_card: 'Version & GitHub Releases Maintenance',
    admin_curr_version: 'Current Version',
    admin_latest_version: 'Latest Release',
    admin_check_btn: 'Check GitHub Releases',
    admin_update_btn: 'Update & Restart Server',
    admin_updating: 'Spawning update.go daemon process...',
    admin_stats_card: 'Realtime Server Metrics',
    admin_active_rooms: 'Active Rooms',
    admin_active_players: 'Online Players',
    admin_close_btn: 'Close Console',
    admin_broadcast_card: 'Server Global Broadcast',
    admin_broadcast_ph: 'Enter announcement message to broadcast to all active rooms...',
    admin_broadcast_send_btn: 'Send Global Broadcast',
    admin_broadcast_sending: 'Broadcasting...',
    admin_broadcast_preset_1: 'Server will undergo maintenance in 5 minutes. Please conclude current matches.',
    admin_broadcast_preset_2: 'Welcome to Liar\'s Deck! Best of luck at the tavern tables!',
    // Login Modal
    login_modal_title: "Liar's Deck · Tavern Entrance",
    login_modal_desc: 'The tavern is open for business. Please verify your identity to take a seat.',
    login_btn: 'Login with pdnode ID',
    login_loading_btn: 'Redirecting to Auth Server...',
    login_secure_tip: 'Secured by pdnode ID · Protected with OAuth2 PKCE',
    logout_btn: 'Logout',
    logout_confirm: 'Are you sure you want to sign out?',
    logged_in_as: 'Signed in as',

    // Logs
    battle_log_title: 'Action Log',

    // Errors
    err_enter_nickname: 'Please enter a name',
    err_enter_code: 'Please enter a 6-digit code',
    reconnecting: 'Connection interrupted, reconnecting...',
  }
};

// 全局单例语言响应式状态，确保全站组件双向同步切换
const currentLang = ref(localStorage.getItem('liarsdeck_lang') || 'zh');

watchEffect(() => {
  localStorage.setItem('liarsdeck_lang', currentLang.value);
});

export function useI18n() {
  function t(key, params = {}) {
    let str = (dict[currentLang.value] || dict.zh)[key] || key;
    for (const [k, v] of Object.entries(params)) {
      str = str.replace(new RegExp(`\\{${k}\\}`, 'g'), v);
    }
    return str;
  }

  function toggleLang() {
    currentLang.value = currentLang.value === 'zh' ? 'en' : 'zh';
  }

  return { lang: currentLang, t, toggleLang };
}
