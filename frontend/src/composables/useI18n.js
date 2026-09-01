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
    room_code_ph: '请输入6位房间码',
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
    rule_section_basic: '基础规则',
    rule_section_items: '道具模式',
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

    // Items Mode Rules
    rule_items_how_title: '道具的获取',
    rule_items_how_desc: '每次你被开枪但抽到空包弹幸存后，系统自动随机补发一个道具（上限 2 个）。道具只能在你的回合中使用。',
    rule_items_list_title: '道具效果说明',
    rule_item_eagle_eye_detail: '查看上家最后打出的一组牌中随机一张的真实点数，仅自己可见。桌面无牌时不可用。',
    rule_item_sawed_off_detail: '激活后，下一次你被迫开枪惩罚时连扣两次扳机，每次独立判定是否命中实弹。',
    rule_item_hard_liquor_detail: '把手牌中最多 2 张非目标牌（假牌）替换为新的随机牌，帮你更新手牌结构。',
    rule_item_kevlar_armor_detail: '装备后吸收下一次致命实弹，护甲碎掉但你存活。不能挡空包弹，护甲不叠加。',
    rule_item_fate_shift_detail: '将桌面目标牌随机替换为另一张（K/Q/A 三选一但不与当前相同），迫使所有人重新评估手牌。',

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
    admin_auth_desc: '请输入管理员密钥以解锁管理控制台：',
    admin_auth_ph: '输入管理密钥',
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
    global_broadcast_title: '全服系统公告',

    // Login Modal
    login_modal_title: "Liar's Deck",
    login_modal_desc: '请登录后游玩',
    login_btn: '登录 / Login',
    login_loading_btn: '正在跳转...',
    login_secure_tip: '使用 Pdnode Auth 提供保护和身份提供',
    logout_btn: '登出',
    logout_confirm: '确定要退出当前账号并登出吗？',
    logged_in_as: '已登录账号',

    // Changelog
    changelog_btn: '更新日志',
    changelog_title: '版本更新日志',
    changelog_subtitle: '查看各版本迭代与特性更新记录',
    changelog_versions: '历史版本',
    changelog_loading: '正在加载更新日志...',

    // Volume & Ping
    volume_label: '音量调节',
    ping_label: '延迟',

    // Public Rooms
    public_rooms_title: '公开对局大厅',
    public_rooms_subtitle: '在线牌桌实时列表，点击即可免码快速加入',
    public_rooms_empty: '当前暂无开放牌桌，快来创建第一间吧！',
    quick_join_btn: '快速入座',
    refresh_btn: '刷新',
    room_host: '房主',
    table_status_waiting: '等待中',
    table_status_playing: '对局中',
    table_players: '人数',

    // Create Room Modal
    create_room_modal_title: '创建房间',
    create_room_modal_subtitle: '选择人数与模式',
    select_players_label: '对局人数',
    players_2: '2 人',
    players_3: '3 人',
    players_4: '4 人',
    select_mode_label: '模式',
    mode_classic_title: '普通',
    mode_classic_desc: '经典暗牌欺诈与俄罗斯轮盘',
    mode_items_title: '道具',
    mode_items_desc: '可使用道具改变战局',
    items_mode_warning_title: '测试版本',
    items_mode_warning_desc: '道具模式目前处于测试阶段，可能存在不平衡或未完成的内容，欢迎反馈。',
    create_room_confirm_lock: '继续 ({s}s)',
    create_room_confirm_items: '继续',
    create_room_confirm_classic: '继续',
    cancel: '取消',

    // Items
    item_eagle_eye_name: '放大镜',
    item_eagle_eye_desc: '查看桌面上最后打出的一张牌',
    item_sawed_off_name: '猎枪',
    item_sawed_off_desc: '下一次受罚开枪连扣两次扳机',
    item_hard_liquor_name: '啤酒',
    item_hard_liquor_desc: '替换手中的假牌',
    item_kevlar_armor_name: '防弹衣',
    item_kevlar_armor_desc: '抵消一次致命实弹',
    item_fate_shift_name: '骰子',
    item_fate_shift_desc: '重置桌面上的目标牌',
    item_used_toast: '已使用【{name}】',
    item_used_by: '{user} 使用了道具',
    item_use_btn: '使用',
    item_slot_empty: '空',
    item_disabled_not_turn: '只能在你的回合使用',
    item_disabled_empty_table: '桌面上暂无出牌',
    eagle_eye_modal_title: '放大镜',
    eagle_eye_modal_desc: '上家打出的其中一张牌为：',
    eagle_eye_close: '继续',
    double_damage_banner: '猎枪生效中：下一次开枪判定连开两枪',
    armor_equipped_tag: '防弹衣',
    items_inventory_label: '道具',
    no_items_tip: '暂无道具',

    // Events shot details
    event_shot_badge_fatal: '💥 致命实弹',
    event_shot_badge_blank: '🛡️ 空包弹',
    event_shot_badge_double: '⚡ 猎枪连开两枪',
    event_shot_badge_armor: '🛡️ 防弹衣抵消实弹',
    event_shot_armor_sub: '防弹衣吸收了致命实弹！已重新装填新轮盘',
    event_shot_double_sub: '猎枪生效！承受了连扣两次扳机判定',

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
    room_code_ph: 'Enter 6-character room code',
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
    rule_section_basic: 'Basic Rules',
    rule_section_items: 'Items Mode',
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

    // Items Mode Rules
    rule_items_how_title: 'How to Get Items',
    rule_items_how_desc: 'Each time you survive a blank shot, the server randomly grants you one item (max 2 held at once). Items can only be used on your own turn.',
    rule_items_list_title: 'Item Effects',
    rule_item_eagle_eye_detail: 'Peek at one random card from the last played group. Only you can see the result. Disabled when no cards are on the table.',
    rule_item_sawed_off_detail: 'When activated, your next gunshot penalty fires twice. Each shot is independently rolled for a live or blank round.',
    rule_item_hard_liquor_detail: 'Replace up to 2 non-table cards (bluff cards) in your hand with new random draws. Helps rebuild a bad hand.',
    rule_item_kevlar_armor_detail: 'Absorbs the next fatal live round hit. The armor breaks after use but you survive. Does not stack.',
    rule_item_fate_shift_detail: 'Reroll the table card to a different one (K/Q/A, never the same as current). Forces everyone to re-evaluate their hands.',

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
    admin_auth_desc: 'Enter your admin secret key to unlock the admin console:',
    admin_auth_ph: 'Enter admin secret',
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
    global_broadcast_title: 'SERVER ANNOUNCEMENT',

    // Login Modal
    login_modal_title: "Liar's Deck",
    login_modal_desc: 'Please sign in to play',
    login_btn: 'Login / Sign In',
    login_loading_btn: 'Redirecting...',
    login_secure_tip: 'Protected and identity provided by Pdnode Auth',
    logout_btn: 'Logout',
    logout_confirm: 'Are you sure you want to sign out?',
    logged_in_as: 'Signed in as',

    // Changelog
    changelog_btn: 'Changelog',
    changelog_title: 'Release Changelog',
    changelog_subtitle: 'Browse features and release history',
    changelog_versions: 'VERSIONS',
    changelog_loading: 'Loading changelog...',

    // Volume & Ping
    volume_label: 'Volume',
    ping_label: 'Ping',

    // Public Rooms
    public_rooms_title: 'Public Tables',
    public_rooms_subtitle: 'Browse active tables and jump in instantly',
    public_rooms_empty: 'No active tables yet. Host one now!',
    quick_join_btn: 'Join Table',
    refresh_btn: 'Refresh',
    room_host: 'Host',
    table_status_waiting: 'Waiting',
    table_status_playing: 'In Game',
    table_players: 'Players',

    // Create Room Modal
    create_room_modal_title: 'Create Room',
    create_room_modal_subtitle: 'Configure player count and mode',
    select_players_label: 'Player Count',
    players_2: '2 Players',
    players_3: '3 Players',
    players_4: '4 Players',
    select_mode_label: 'Mode',
    mode_classic_title: 'Normal',
    mode_classic_desc: 'Classic bluffing and Russian roulette',
    mode_items_title: 'Items',
    mode_items_desc: 'Use items to change the game',
    items_mode_warning_title: 'Beta Feature',
    items_mode_warning_desc: 'Items mode is currently in testing. Balance and content may be incomplete. Feedback welcome.',
    create_room_confirm_lock: 'Continue ({s}s)',
    create_room_confirm_items: 'Continue',
    create_room_confirm_classic: 'Continue',
    cancel: 'Cancel',

    // Items
    item_eagle_eye_name: 'Magnifier',
    item_eagle_eye_desc: 'View one of the last played cards',
    item_sawed_off_name: 'Shotgun',
    item_sawed_off_desc: 'Next gunshot penalty fires twice',
    item_hard_liquor_name: 'Beer',
    item_hard_liquor_desc: 'Replace false cards in hand',
    item_kevlar_armor_name: 'Vest',
    item_kevlar_armor_desc: 'Block one fatal bullet',
    item_fate_shift_name: 'Dice',
    item_fate_shift_desc: 'Reroll the table card',
    item_used_toast: 'Used [{name}]',
    item_used_by: '{user} used an item',
    item_use_btn: 'Use',
    item_slot_empty: 'Empty',
    item_disabled_not_turn: 'Can only use on your turn',
    item_disabled_empty_table: 'No cards on table yet',
    eagle_eye_modal_title: 'Magnifier',
    eagle_eye_modal_desc: 'One of the cards played by the previous player:',
    eagle_eye_close: 'Continue',
    double_damage_banner: 'Shotgun active: next gunshot fires twice',
    armor_equipped_tag: 'Vest',
    items_inventory_label: 'Items',
    no_items_tip: 'No items',

    // Events shot details
    event_shot_badge_fatal: '💥 FATAL ROUND',
    event_shot_badge_blank: '🛡️ DRY FIRE',
    event_shot_badge_double: '⚡ SHOTGUN DOUBLE SHOT',
    event_shot_badge_armor: '🛡️ VEST DEFLECTED FATAL SHOT',
    event_shot_armor_sub: 'Vest absorbed the fatal bullet! Reloaded new cylinder.',
    event_shot_double_sub: 'Shotgun active! Two trigger pulls were executed.',

    // Logs
    battle_log_title: 'Action Log',

    // Errors
    err_enter_nickname: 'Please enter a name',
    err_enter_code: 'Please enter a 6-digit code',
    reconnecting: 'Connection interrupted, reconnecting...',
  }
};

// 全局单例语言响应式状态，默认英文
const currentLang = ref(localStorage.getItem('liarsdeck_lang') || 'en');

watchEffect(() => {
  localStorage.setItem('liarsdeck_lang', currentLang.value);
});

export function useI18n() {
  function t(key, params = {}) {
    let str = (dict[currentLang.value] || dict.en || dict.zh)[key] || (dict.en || dict.zh)[key] || key;
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
