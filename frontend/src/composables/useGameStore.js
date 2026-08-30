import { ref, computed } from 'vue';
import { useAudio } from './useAudio';
import confetti from 'canvas-confetti';

const connected = ref(false);
const ws = ref(null);
const state = ref({
  status: 'waiting',
  players: [],
  current_turn: -1,
  table_card: '',
  last_player: -1,
  last_played_cnt: 0,
  logs: [],
  deadline: 0,
  winner: '',
  room_code: '',
});

const currentUnix = ref(Math.floor(Date.now() / 1000));
setInterval(() => {
  currentUnix.value = Math.floor(Date.now() / 1000);
}, 1000);

const selectedIndexes = ref([]);
const eventQueue = ref([]);
const processingEvent = ref(false);
const currentStep = ref('');
const currentStepData = ref({});
const myNickname = ref('');
const myRoomCode = ref('');
const myPlayerToken = ref('');
const toast = ref('');
const errorMsg = ref('');
const pendingState = ref(null);

let toastTimer = null;

export function useGameStore() {
  const audio = useAudio();

  function showToast(msg, dur = 2200) {
    toast.value = msg;
    if (toastTimer) clearTimeout(toastTimer);
    toastTimer = setTimeout(() => {
      toast.value = '';
    }, dur);
  }

  // ── Computed ──
  const myPlayer = computed(() => {
    if (myPlayerToken.value) {
      const byId = state.value.players.find((p) => p.id === myPlayerToken.value);
      if (byId) return byId;
    }
    return state.value.players.find((p) => p.nickname === myNickname.value);
  });

  const myId = computed(() => (myPlayer.value ? myPlayer.value.id : ''));
  const myHand = computed(() => (myPlayer.value ? myPlayer.value.hand || [] : []));
  const isPlayer = computed(() => myPlayer.value && !myPlayer.value.is_spectator);
  const isSpectator = computed(() => myPlayer.value && myPlayer.value.is_spectator);
  const amHost = computed(() => myPlayer.value && myPlayer.value.is_host);

  const isMyTurn = computed(() => {
    if (!state.value.players || state.value.current_turn === -1 || !myPlayer.value) return false;
    return state.value.players[state.value.current_turn]?.id === myId.value;
  });

  const canPlay = computed(() => {
    if (!isMyTurn.value) return false;
    // 如果上家打光了手牌，下家必须质疑
    if (state.value.last_player !== -1 && state.value.last_player < state.value.players.length) {
      const lastP = state.value.players[state.value.last_player];
      if (lastP.hand && lastP.hand.length === 0) return false;
    }
    return selectedIndexes.value.length >= 1 && selectedIndexes.value.length <= 3;
  });

  const canCallLiar = computed(() => {
    return (
      isMyTurn.value &&
      state.value.last_player !== -1 &&
      state.value.last_player !== state.value.current_turn
    );
  });

  const canStart = computed(() => {
    if (!amHost.value) return false;
    const players = state.value.players.filter((p) => !p.is_spectator);
    return players.length >= 2 && players.every((p) => p.is_ready);
  });

  const allReady = computed(() => {
    const players = state.value.players.filter((p) => !p.is_spectator);
    return players.length >= 2 && players.every((p) => p.is_ready);
  });

  // ── Event Queue ──
  function processQueue() {
    if (processingEvent.value || eventQueue.value.length === 0) return;
    processingEvent.value = true;
    const ev = eventQueue.value.shift();
    currentStep.value = ev.type;
    currentStepData.value = ev.data;

    if (ev.type === 'liar_call') {
      audio.playLiarAlert();
    } else if (ev.type === 'shot') {
      if (ev.data.fatal) {
        audio.playGunshot();
      } else {
        audio.playGunClick();
      }
    } else if (ev.type === 'reveal') {
      audio.playCardDeal();
    }

    const dur = { liar_call: 1800, reveal: 2600, shot: 2200 };
    setTimeout(() => {
      currentStep.value = '';
      processingEvent.value = false;

      if (eventQueue.value.length === 0 && pendingState.value) {
        state.value = pendingState.value;
        pendingState.value = null;
        if (!isMyTurn.value) selectedIndexes.value = [];
      }
      processQueue();
    }, dur[ev.type] || 2000);
  }

  // ── WebSocket Connect ──
  function connect(action, roomCode, name) {
    errorMsg.value = '';
    myNickname.value = name;
    if (roomCode) myRoomCode.value = roomCode;

    const proto = window.location.protocol === 'https:' ? 'wss' : 'ws';
    const host = window.location.host;
    let url = `${proto}://${host}/ws?action=${action}&name=${encodeURIComponent(name)}`;
    if (roomCode) url += `&code=${encodeURIComponent(roomCode)}`;

    doConnect(url);
  }

  function doConnect(url) {
    if (ws.value) {
      ws.value.close();
    }
    const socket = new WebSocket(url);
    ws.value = socket;

    socket.onopen = () => {
      connected.value = true;
      errorMsg.value = '';
    };

    socket.onclose = () => {
      if (connected.value) {
        showToast('连接已断开，正在尝试重连...');
        if (myPlayerToken.value && myRoomCode.value) {
          setTimeout(tryReconnect, 1500);
        } else {
          setTimeout(() => {
            connected.value = false;
          }, 1500);
        }
      }
    };

    socket.onmessage = handleMessage;
  }

  function tryReconnect() {
    if (!myRoomCode.value || !myPlayerToken.value || !myNickname.value) return;
    const proto = window.location.protocol === 'https:' ? 'wss' : 'ws';
    const host = window.location.host;
    const url = `${proto}://${host}/ws?action=reconnect&code=${encodeURIComponent(
      myRoomCode.value
    )}&token=${encodeURIComponent(myPlayerToken.value)}&name=${encodeURIComponent(myNickname.value)}`;

    const newWs = new WebSocket(url);
    newWs.onopen = () => {
      showToast('🔗 重新连接成功！');
      ws.value = newWs;
      connected.value = true;
      newWs.onclose = () => {
        if (connected.value) {
          showToast('连接已断开');
          setTimeout(tryReconnect, 3000);
        }
      };
      newWs.onmessage = handleMessage;
    };
    newWs.onclose = () => {
      setTimeout(tryReconnect, 3000);
    };
  }

  function handleMessage(e) {
    try {
      const msg = JSON.parse(e.data);
      if (msg.error) {
        errorMsg.value = msg.error;
        showToast(msg.error);
        return;
      }

      if (msg.type === 'game_state') {
        const prevStatus = state.value.status;
        const prevTurn = state.value.current_turn;

        if (msg.data.players && myNickname.value) {
          const me = msg.data.players.find((p) => p.nickname === myNickname.value);
          if (me && me.id) myPlayerToken.value = me.id;
        }

        if (msg.data.room_code) {
          myRoomCode.value = msg.data.room_code;
        }

        // 胜出特效
        if (msg.data.status === 'game_over' && prevStatus !== 'game_over' && msg.data.winner) {
          audio.playVictory();
          confetti({
            particleCount: 100,
            spread: 70,
            origin: { y: 0.6 },
          });
        }

        if (currentStep.value !== '') {
          pendingState.value = msg.data;
        } else {
          state.value = msg.data;
          if (!isMyTurn.value) selectedIndexes.value = [];
        }
      } else if (['liar_call', 'reveal', 'shot'].includes(msg.type)) {
        eventQueue.value.push({ type: msg.type, data: msg.data });
        processQueue();
      }
    } catch (_) {}
  }

  // ── Actions ──
  function sendAction(action, payload = {}) {
    if (ws.value && ws.value.readyState === WebSocket.OPEN) {
      ws.value.send(JSON.stringify({ action, payload }));
    }
  }

  function toggleCardSelect(idx) {
    if (!isMyTurn.value) return;
    const pos = selectedIndexes.value.indexOf(idx);
    if (pos > -1) {
      selectedIndexes.value.splice(pos, 1);
    } else if (selectedIndexes.value.length < 3) {
      selectedIndexes.value.push(idx);
      audio.playCardSelect();
    }
  }

  function clearSelection() {
    selectedIndexes.value = [];
  }

  function playCards() {
    if (!canPlay.value) return;
    const cards = selectedIndexes.value.map((i) => myHand.value[i]);
    sendAction('play_cards', { cards });
    clearSelection();
    audio.playCardDeal();
  }

  function callLiar() {
    if (!canCallLiar.value) return;
    sendAction('call_liar', {});
  }

  function toggleReady() {
    sendAction('ready', {});
  }

  function startGame() {
    if (!canStart.value) return;
    sendAction('start', {});
  }

  function resetGame() {
    if (!amHost.value) return;
    sendAction('reset', {});
  }

  function kickPlayer(targetId) {
    if (!amHost.value) return;
    sendAction('remove_player', { target_id: targetId });
  }

  function copyInvite() {
    const link = `${window.location.origin}?room=${state.value.room_code || myRoomCode.value}`;
    navigator.clipboard
      .writeText(link)
      .then(() => showToast('✅ 邀请链接已复制'))
      .catch(() => showToast(`📋 ${link}`));
  }

  function disconnect() {
    if (ws.value) {
      ws.value.close();
      ws.value = null;
    }
    connected.value = false;
  }

  return {
    connected,
    state,
    currentUnix,
    selectedIndexes,
    currentStep,
    currentStepData,
    myNickname,
    myId,
    myRoomCode,
    toast,
    errorMsg,
    myPlayer,
    myHand,
    isPlayer,
    isSpectator,
    amHost,
    isMyTurn,
    canPlay,
    canCallLiar,
    canStart,
    allReady,
    connect,
    sendAction,
    toggleCardSelect,
    clearSelection,
    playCards,
    callLiar,
    toggleReady,
    startGame,
    resetGame,
    kickPlayer,
    copyInvite,
    disconnect,
    showToast,
  };
}
