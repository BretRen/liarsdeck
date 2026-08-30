import { ref, watchEffect } from 'vue';

const dict = {
  zh: {
    app_title: "Liar's Deck",
    app_subtitle: '心理博弈 · 轮盘对决 · 卡牌狂欢',
    lobby_create_title: '🎲 创建对战房间',
    lobby_join_title: '🔑 凭码加入房间',
    lobby_spec_title: '👀 观战实时对局',
    nickname: '玩家昵称',
    nickname_ph: '输入你的霸气昵称...',
    room_code: '6位房间码',
    room_code_ph: '例如：ABC888',
    create_btn: '创建房间',
    join_btn: '加入对战',
    spectate_btn: '进入观战',
    back: '返回上一步',
    rules_btn: '查看规则',
    
    // Status badges
    status_waiting: '等待就绪',
    status_playing: '激烈对局中',
    status_game_over: '对局结束',
    
    // Player seat
    host_tag: '房主',
    spec_tag: '观众',
    me_tag: '我',
    dead_tag: '💀 阵亡',
    watching_tag: '观战中',
    bullets_label: '左轮弹仓',
    hand_count_label: '手牌',
    ready_btn: '准备',
    unready_btn: '取消准备',
    ready_status: '已就绪',
    unready_status: '未准备',
    kick_btn: '请离',
    kick_confirm: '确定要踢出该玩家吗？',
    
    // Header
    invite_btn: '复制邀请',
    copied_toast: '✅ 邀请链接已复制到剪贴板',
    audio_on: '🔊 音效开启',
    audio_off: '🔇 静音',
    timeout_warn: '秒后系统自动托管',
    
    // Table
    table_card_label: '本轮真牌 (TRUE CARD)',
    wild_card_tip: '万能牌 2 可作为任何真牌打出',
    cards_on_table: '张背面向下底牌',
    winner_label: '👑 最终胜利者',
    play_again_btn: '再来一局',
    wait_host_reset: '等待房主重置对局...',
    start_game_btn: '开始游戏',
    need_more_players: '需至少 2 名准备好的玩家',
    all_ready_needed: '等待所有玩家准备',
    
    // Hand
    my_hand_title: '我的手牌',
    selected_count: '已选',
    max_cards_tip: '每次最多出 1~3 张',
    
    // Actions
    play_cards_btn: '打出所选牌',
    call_liar_btn: '🚨 质疑说谎 (CALL LIAR)',
    
    // Events
    event_liar_alert: '🚨 质疑说谎！',
    event_calls_out: '怀疑',
    event_liar_claim: '在说谎！',
    event_cards_revealed: '翻开底牌验证',
    event_bang_title: '💥 致命枪击！',
    event_bang_sub: '扣中致命子弹，遗憾出局！',
    event_click_title: '💨 咔哒... 空包弹！',
    event_click_sub: '逃过一劫，轮盘继续！',
    
    // Spectator Notice
    spectator_banner: '你当前处于观战模式，可以实时观看桌面对决。',
    
    // Rules Modal
    rules_title: '🃏 Liar\'s Deck 游戏规则手册',
    rule_goal_title: '🎯 获胜目标',
    rule_goal_desc: '成为全场最后一名存活玩家，或者最先出完手牌并挺过下家的质疑！',
    rule_deck_title: '🎴 牌组构成 (24张)',
    rule_deck_desc: 'K (King)、Q (Queen)、A (Ace)、2 (万能 Wild)，每种各 6 张。2 可代替任何真牌。',
    rule_turn_title: '🔁 回合流程',
    rule_turn_1: '每轮随机翻出一张“真牌”，存活玩家每人分发 5 张手牌。',
    rule_turn_2: '轮到你时，选择 1~3 张手牌背面朝下打出，并声称全部为“真牌”。',
    rule_turn_3: '你可以打真牌，也可以尽情撒谎！',
    rule_liar_title: '🚨 质疑机制 (Call Liar)',
    rule_liar_1: '当下家轮到操作时，可以选择信任接牌，或直接【质疑说谎】。',
    rule_liar_2: '翻开底牌：若存在任何一张既不是真牌也不是 2 的牌 → 质疑成功，出牌者开枪！',
    rule_liar_3: '若底牌全为真牌或 2 → 质疑失败，质疑者自己开枪！',
    rule_gun_title: '🔫 左轮手枪轮盘赌',
    rule_gun_1: '每人一把 6 发弹仓左轮（5 发空包弹，1 发致命弹，随机乱序）。',
    rule_gun_2: '每次触发开枪扣动扳机一次，命中致命弹直接淘汰！',
    rule_got_it: '明白，进入对局！',
    
    // Logs
    battle_log_title: '📜 战局记录',
    
    // Errors
    err_enter_nickname: '请输入玩家昵称',
    err_enter_code: '请输入6位房间码',
    reconnecting: '连接断开，正在尝试自动重连...',
  },
  en: {
    app_title: "Liar's Deck",
    app_subtitle: 'Bluffing · Russian Roulette · Psychological Warfare',
    lobby_create_title: '🎲 Create Match Room',
    lobby_join_title: '🔑 Join with Code',
    lobby_spec_title: '👀 Spectate Live Game',
    nickname: 'Nickname',
    nickname_ph: 'Enter your cool name...',
    room_code: '6-Digit Room Code',
    room_code_ph: 'e.g. ABC888',
    create_btn: 'Create Room',
    join_btn: 'Join Match',
    spectate_btn: 'Spectate',
    back: 'Go Back',
    rules_btn: 'Game Rules',
    
    // Status badges
    status_waiting: 'Waiting',
    status_playing: 'In Battle',
    status_game_over: 'Game Over',
    
    // Player seat
    host_tag: 'Host',
    spec_tag: 'Spectator',
    me_tag: 'You',
    dead_tag: '💀 Dead',
    watching_tag: 'Watching',
    bullets_label: 'Chambers',
    hand_count_label: 'Cards',
    ready_btn: 'Ready',
    unready_btn: 'Cancel Ready',
    ready_status: 'Ready',
    unready_status: 'Not Ready',
    kick_btn: 'Kick',
    kick_confirm: 'Are you sure you want to kick this player?',
    
    // Header
    invite_btn: 'Copy Invite',
    copied_toast: '✅ Invite link copied to clipboard',
    audio_on: '🔊 Sound ON',
    audio_off: '🔇 Muted',
    timeout_warn: 's auto timeout',
    
    // Table
    table_card_label: 'TRUE CARD FOR THIS ROUND',
    wild_card_tip: 'Wild Card 2 counts as any True Card',
    cards_on_table: 'cards played face-down',
    winner_label: '👑 Champion',
    play_again_btn: 'Play Again',
    wait_host_reset: 'Waiting for host to restart...',
    start_game_btn: 'Start Game',
    need_more_players: 'Needs 2+ ready players',
    all_ready_needed: 'Waiting for all ready',
    
    // Hand
    my_hand_title: 'Your Hand',
    selected_count: 'Selected',
    max_cards_tip: 'Select 1 to 3 cards',
    
    // Actions
    play_cards_btn: 'Play Cards',
    call_liar_btn: '🚨 CALL LIAR!',
    
    // Events
    event_liar_alert: '🚨 CALL LIAR!',
    event_calls_out: 'suspects',
    event_liar_claim: 'is bluffing!',
    event_cards_revealed: 'Cards Verification',
    event_bang_title: '💥 BANG! FATAL SHOT!',
    event_bang_sub: 'Hit the fatal bullet and got eliminated!',
    event_click_title: '💨 CLICK... BLANK!',
    event_click_sub: 'Lucky escape! The roulette continues!',
    
    // Spectator Notice
    spectator_banner: 'You are in Spectator Mode watching the match live.',
    
    // Rules Modal
    rules_title: '🃏 Liar\'s Deck Rulebook',
    rule_goal_title: '🎯 Objective',
    rule_goal_desc: 'Be the last player standing or empty all your hand cards safely!',
    rule_deck_title: '🎴 The Deck (24 Cards)',
    rule_deck_desc: 'K (King), Q (Queen), A (Ace), 2 (Wild), 6 of each. 2 can match any True Card.',
    rule_turn_title: '🔁 Turn Sequence',
    rule_turn_1: 'A True Card is revealed each round. Each alive player receives 5 cards.',
    rule_turn_2: 'On your turn, play 1-3 cards face-down and claim they match the True Card.',
    rule_turn_3: 'You may play honestly or bluff boldly!',
    rule_liar_title: '🚨 Calling Liar',
    rule_liar_1: 'The next player may accept and play, or call out the previous player.',
    rule_liar_2: 'Cards are checked: if any card is not the True Card and not 2 → Liar caught, bluffer shoots!',
    rule_liar_3: 'If all cards were valid → False call, caller shoots!',
    rule_gun_title: '🔫 Russian Roulette',
    rule_gun_1: 'Each player has a 6-chamber revolver (5 blanks, 1 fatal bullet, shuffled).',
    rule_gun_2: 'Pull the trigger when penalized. Fatal bullet means elimination!',
    rule_got_it: 'Got it, let\'s play!',
    
    // Logs
    battle_log_title: '📜 Battle Logs',
    
    // Errors
    err_enter_nickname: 'Please enter a nickname',
    err_enter_code: 'Please enter 6-digit room code',
    reconnecting: 'Connection lost, trying to reconnect...',
  }
};

export function useI18n() {
  const lang = ref(localStorage.getItem('liarsdeck_lang') || 'zh');

  watchEffect(() => {
    localStorage.setItem('liarsdeck_lang', lang.value);
  });

  function t(key) {
    return (dict[lang.value] || dict.zh)[key] || key;
  }

  function toggleLang() {
    lang.value = lang.value === 'zh' ? 'en' : 'zh';
  }

  return { lang, t, toggleLang };
}
