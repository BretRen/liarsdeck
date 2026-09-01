import { ref, computed } from 'vue';
import { useAudio } from './useAudio';
import confetti from 'canvas-confetti';

const connected = ref(false);
const hasJoinedRoom = ref(false);
const isDisconnected = ref(false);
const isReconnecting = ref(false);
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
  pause_deadline: 0,
  paused_player: '',
  remaining_turn_seconds: 0,
  winner: '',
  room_code: '',
});

const currentUnix = ref(Math.floor(Date.now() / 1000));
setInterval(() => {
  currentUnix.value = Math.floor(Date.now() / 1000);
}, 250);

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
const globalBroadcast = ref({ message: '', timestamp: 0, visible: false });
const myPing = ref(0);
const isScreenShaking = ref(false);
const isFatalShotActive = ref(false);
const inspectedCard = ref('');
const showEagleEyeModal = ref(false);
const itemUsedEvent = ref(null); // { nickname, item } — cleared after 2.4s
let itemUsedTimer = null;

let toastTimer = null;
let reconnectTimer = null;
let broadcastTimer = null;
let pingInterval = null;

const SESSION_KEY = 'liarsdeck_active_session';

function saveSession(roomCode, token, nickname) {
  if (roomCode && token && nickname) {
    localStorage.setItem(
      SESSION_KEY,
      JSON.stringify({ roomCode, token, nickname, savedAt: Date.now() })
    );
  }
}

function clearSession() {
  localStorage.removeItem(SESSION_KEY);
}

function getSavedSession() {
  try {
    const raw = localStorage.getItem(SESSION_KEY);
    if (!raw) return null;
    const sess = JSON.parse(raw);
    if (sess && sess.roomCode && sess.token && sess.nickname) {
      if (Date.now() - (sess.savedAt || 0) < 3600000) {
        return sess;
      }
    }
  } catch (_) {}
  return null;
}

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
    if (state.value.status !== 'playing' || !state.value.players || state.value.current_turn === -1) return false;
    const turnPlayer = state.value.players[state.value.current_turn];
    if (!turnPlayer) return false;
    return myPlayerToken.value !== '' && turnPlayer.id === myPlayerToken.value;
  });

  const remainingSeconds = computed(() => {
    if (state.value.status !== 'playing' || !state.value.deadline) return 0;
    return Math.max(0, state.value.deadline - currentUnix.value);
  });

  const pauseRemainingSeconds = computed(() => {
    if (state.value.status !== 'paused' || !state.value.pause_deadline) return 0;
    return Math.max(0, state.value.pause_deadline - currentUnix.value);
  });

  const canPlay = computed(() => {
    if (state.value.status !== 'playing') return false;
    if (!isMyTurn.value) return false;
    if (state.value.last_player !== -1 && state.value.last_player < state.value.players.length) {
      const lastP = state.value.players[state.value.last_player];
      if (lastP.hand && lastP.hand.length === 0) return false;
    }
    return selectedIndexes.value.length >= 1 && selectedIndexes.value.length <= 3;
  });

  const canCallLiar = computed(() => {
    if (state.value.status !== 'playing') return false;
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

  // ── Event Queue & Delayed Victory Animation ──
  function processQueue() {
    if (processingEvent.value || eventQueue.value.length === 0) return;
    processingEvent.value = true;
    const ev = eventQueue.value.shift();
    currentStep.value = ev.type;
    currentStepData.value = ev.data;

    if (ev.type === 'liar_call') {
      audio.playLiarAlert();
      audio.playHeartbeat();
    } else if (ev.type === 'shot') {
      if (ev.data.fatal) {
        audio.playGunshot();
        isFatalShotActive.value = true;
        isScreenShaking.value = true;
        setTimeout(() => {
          isScreenShaking.value = false;
        }, 750);
        setTimeout(() => {
          isFatalShotActive.value = false;
        }, 2200);
      } else {
        audio.playGunClick();
      }
    } else if (ev.type === 'reveal') {
      audio.playCardFlip();
    }

    const dur = { liar_call: 1900, reveal: 2400, shot: 2200 };
    setTimeout(() => {
      currentStep.value = '';
      processingEvent.value = false;

      // When ALL event queue animations are finished, commit pending state
      if (eventQueue.value.length === 0 && pendingState.value) {
        const isWinningTransition =
          pendingState.value.status === 'game_over' &&
          state.value.status !== 'game_over' &&
          pendingState.value.winner;

        state.value = pendingState.value;
        pendingState.value = null;
        if (!isMyTurn.value) selectedIndexes.value = [];

        // Trigger victory celebration ONLY after all animations have ended
        if (isWinningTransition) {
          setTimeout(() => {
            audio.playVictory();
            confetti({
              particleCount: 120,
              spread: 80,
              origin: { y: 0.6 },
            });
          }, 300);
        }
      }
      processQueue();
    }, dur[ev.type] || 2000);
  }

  // ── WebSocket Connect ──
  function connect(action, roomCode, name, token = '', options = {}) {
    errorMsg.value = '';
    myNickname.value = name;
    if (roomCode) myRoomCode.value = roomCode;
    if (token) myPlayerToken.value = token;

    const proto = window.location.protocol === 'https:' ? 'wss' : 'ws';
    const host = window.location.host;
    let url = `${proto}://${host}/ws?action=${action}&name=${encodeURIComponent(name)}`;
    if (roomCode) url += `&code=${encodeURIComponent(roomCode)}`;
    if (token) url += `&token=${encodeURIComponent(token)}`;
    if (options.mode) url += `&mode=${encodeURIComponent(options.mode)}`;
    if (options.maxPlayers) url += `&max_players=${encodeURIComponent(options.maxPlayers)}`;

    doConnect(url);
  }

  function doConnect(url) {
    if (ws.value) {
      ws.value.close();
    }
    const socket = new WebSocket(url);
    ws.value = socket;

    socket.onopen = () => {
      errorMsg.value = '';
      if (reconnectTimer) {
        clearTimeout(reconnectTimer);
        reconnectTimer = null;
      }
      if (pingInterval) clearInterval(pingInterval);
      pingInterval = setInterval(() => {
        if (ws.value && ws.value.readyState === WebSocket.OPEN) {
          ws.value.send(JSON.stringify({ action: 'ping', payload: { client_time: Date.now() } }));
        }
      }, 3000);
      // 立即触发一次
      ws.value.send(JSON.stringify({ action: 'ping', payload: { client_time: Date.now() } }));
    };

    socket.onclose = () => {
      if (pingInterval) {
        clearInterval(pingInterval);
        pingInterval = null;
      }
      // 只有在已成功加入房间且非主动退出的情况下，网络意外断开才弹窗重连
      if (hasJoinedRoom.value) {
        isDisconnected.value = true;
        isReconnecting.value = true;
        if (myPlayerToken.value && myRoomCode.value) {
          scheduleReconnect();
        }
      } else {
        connected.value = false;
        isDisconnected.value = false;
        isReconnecting.value = false;
      }
    };

    socket.onmessage = handleMessage;
  }

  function scheduleReconnect() {
    if (reconnectTimer) clearTimeout(reconnectTimer);
    reconnectTimer = setTimeout(() => {
      if (isDisconnected.value && hasJoinedRoom.value) {
        tryReconnect();
      }
    }, 2000);
  }

  function tryReconnect() {
    if (!myRoomCode.value || !myPlayerToken.value || !myNickname.value) return;
    isReconnecting.value = true;
    const proto = window.location.protocol === 'https:' ? 'wss' : 'ws';
    const host = window.location.host;
    const url = `${proto}://${host}/ws?action=reconnect&code=${encodeURIComponent(
      myRoomCode.value
    )}&token=${encodeURIComponent(myPlayerToken.value)}&name=${encodeURIComponent(myNickname.value)}`;

    const newWs = new WebSocket(url);
    newWs.onopen = () => {
      ws.value = newWs;
      if (reconnectTimer) {
        clearTimeout(reconnectTimer);
        reconnectTimer = null;
      }
      if (pingInterval) clearInterval(pingInterval);
      pingInterval = setInterval(() => {
        if (ws.value && ws.value.readyState === WebSocket.OPEN) {
          ws.value.send(JSON.stringify({ action: 'ping', payload: { client_time: Date.now() } }));
        }
      }, 3000);
      ws.value.send(JSON.stringify({ action: 'ping', payload: { client_time: Date.now() } }));

      newWs.onclose = () => {
        if (pingInterval) {
          clearInterval(pingInterval);
          pingInterval = null;
        }
        if (hasJoinedRoom.value) {
          isDisconnected.value = true;
          isReconnecting.value = true;
          scheduleReconnect();
        }
      };
      newWs.onmessage = handleMessage;
    };
    newWs.onclose = () => {
      if (isDisconnected.value && hasJoinedRoom.value) {
        scheduleReconnect();
      }
    };
  }

  function handleMessage(e) {
    try {
      const msg = JSON.parse(e.data);
      if (msg.error) {
        errorMsg.value = msg.error;
        showToast(msg.error);
        connected.value = false;
        hasJoinedRoom.value = false;
        isDisconnected.value = false;
        isReconnecting.value = false;
        clearSession();
        if (ws.value) {
          ws.value.close();
          ws.value = null;
        }
        return;
      }

      if (msg.type === 'pong') {
        if (msg.data && msg.data.client_time) {
          myPing.value = Math.max(1, Date.now() - msg.data.client_time);
        }
        return;
      }

      if (msg.type === 'game_state') {
        const wasReconnecting = isReconnecting.value || isDisconnected.value;

        hasJoinedRoom.value = true;
        connected.value = true;
        isDisconnected.value = false;
        isReconnecting.value = false;

        const prevStatus = state.value.status;

        if (msg.data.players) {
          let me = null;
          if (myPlayerToken.value) {
            me = msg.data.players.find((p) => p.id === myPlayerToken.value);
          }
          if (!me && myNickname.value) {
            me = msg.data.players.find((p) => p.nickname === myNickname.value);
          }
          if (me && me.id) {
            myPlayerToken.value = me.id;
            saveSession(msg.data.room_code || myRoomCode.value, me.id, myNickname.value);
          }
        }

        if (msg.data.room_code) {
          myRoomCode.value = msg.data.room_code;
        }

        // On reconnect: force-flush any stale animation queue so the new state
        // from the server is applied immediately and isMyTurn is correct.
        if (wasReconnecting) {
          eventQueue.value = [];
          pendingState.value = null;
          currentStep.value = '';
          processingEvent.value = false;
        }

        // If there are ongoing or queued animations, delay state commit and victory fanfare
        if (currentStep.value !== '' || eventQueue.value.length > 0) {
          pendingState.value = msg.data;
        } else {
          const isWinningTransition =
            msg.data.status === 'game_over' &&
            prevStatus !== 'game_over' &&
            msg.data.winner;

          state.value = msg.data;
          if (!isMyTurn.value) selectedIndexes.value = [];

          if (isWinningTransition) {
            audio.playVictory();
            confetti({
              particleCount: 120,
              spread: 80,
              origin: { y: 0.6 },
            });
          }
        }
      } else if (['liar_call', 'reveal', 'shot'].includes(msg.type)) {
        eventQueue.value.push({ type: msg.type, data: msg.data });
        processQueue();
      } else if (msg.type === 'eagle_eye_result') {
        if (msg.data && msg.data.card) {
          inspectedCard.value = msg.data.card;
          showEagleEyeModal.value = true;
          audio.playCardFlip();
        }
      } else if (msg.type === 'item_used') {
        if (msg.data && msg.data.item) {
          // Play sound
          if (msg.data.item === 'sawed_off') {
            audio.playGunClick();
          } else if (msg.data.item === 'hard_liquor' || msg.data.item === 'fate_shift') {
            audio.playCardFlip();
          } else {
            audio.playLiarAlert();
          }
          // Trigger visual notification
          if (itemUsedTimer) clearTimeout(itemUsedTimer);
          itemUsedEvent.value = { nickname: msg.data.nickname || '?', item: msg.data.item };
          itemUsedTimer = setTimeout(() => { itemUsedEvent.value = null; }, 2400);
        }
      } else if (msg.type === 'server_broadcast') {
        if (msg.data && msg.data.message) {
          showBroadcast(msg.data.message);
          audio.playLiarAlert();
        }
      }
    } catch (_) {}
  }

  function showBroadcast(msgText) {
    globalBroadcast.value = {
      message: msgText,
      timestamp: Date.now(),
      visible: true,
    };
    if (broadcastTimer) clearTimeout(broadcastTimer);
    broadcastTimer = setTimeout(() => {
      globalBroadcast.value.visible = false;
    }, 12000);
  }

  function dismissBroadcast() {
    if (broadcastTimer) clearTimeout(broadcastTimer);
    globalBroadcast.value.visible = false;
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
    } else {
      if (selectedIndexes.value.length >= 3) {
        showToast('每次最多只能选择 3 张暗牌');
        return;
      }
      selectedIndexes.value.push(idx);
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
      .then(() => showToast('邀请链接已复制'))
      .catch(() => showToast(link));
  }

  function exitToLobby() {
    clearSession();
    if (reconnectTimer) {
      clearTimeout(reconnectTimer);
      reconnectTimer = null;
    }
    if (ws.value) {
      ws.value.close();
      ws.value = null;
    }
    hasJoinedRoom.value = false;
    isDisconnected.value = false;
    isReconnecting.value = false;
    connected.value = false;
    myRoomCode.value = '';
    myPlayerToken.value = '';
  }

  function useItem(itemType) {
    sendAction('use_item', { item: itemType });
  }

  function closeEagleEyeModal() {
    showEagleEyeModal.value = false;
    inspectedCard.value = '';
  }

  function disconnect() {
    exitToLobby();
  }

  return {
    connected,
    hasJoinedRoom,
    isDisconnected,
    isReconnecting,
    state,
    currentUnix,
    remainingSeconds,
    pauseRemainingSeconds,
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
    tryReconnect,
    exitToLobby,
    disconnect,
    showToast,
    getSavedSession,
    clearSession,
    globalBroadcast,
    dismissBroadcast,
    myPing,
    isScreenShaking,
    isFatalShotActive,
    inspectedCard,
    showEagleEyeModal,
    closeEagleEyeModal,
    useItem,
    itemUsedEvent,
  };
}
