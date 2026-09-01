import { ref } from 'vue';

const isMuted = ref(localStorage.getItem('liarsdeck_muted') === 'true');
const masterVolume = ref(parseFloat(localStorage.getItem('liarsdeck_volume') || '0.8'));

let audioCtx = null;

function getAudioContext() {
  if (!audioCtx) {
    const AudioContext = window.AudioContext || window.webkitAudioContext;
    if (AudioContext) {
      audioCtx = new AudioContext();
    }
  }
  if (audioCtx && audioCtx.state === 'suspended') {
    audioCtx.resume();
  }
  return audioCtx;
}

function calcGain(baseGain) {
  if (isMuted.value) return 0;
  return baseGain * Math.max(0, Math.min(1, masterVolume.value));
}

export function useAudio() {
  function toggleMute() {
    isMuted.value = !isMuted.value;
    localStorage.setItem('liarsdeck_muted', isMuted.value ? 'true' : 'false');
  }

  function setVolume(val) {
    const num = Math.max(0, Math.min(1, parseFloat(val) || 0));
    masterVolume.value = num;
    localStorage.setItem('liarsdeck_volume', num.toString());
    if (num > 0 && isMuted.value) {
      isMuted.value = false;
      localStorage.setItem('liarsdeck_muted', 'false');
    }
  }

  // 牌张打出/发牌音效（厚实木桌扑克摩擦声）
  function playCardDeal() {
    if (isMuted.value || masterVolume.value <= 0) return;
    const ctx = getAudioContext();
    if (!ctx) return;

    const osc = ctx.createOscillator();
    const gain = ctx.createGain();

    osc.type = 'triangle';
    osc.frequency.setValueAtTime(320, ctx.currentTime);
    osc.frequency.exponentialRampToValueAtTime(80, ctx.currentTime + 0.09);

    const peak = calcGain(0.28);
    gain.gain.setValueAtTime(peak, ctx.currentTime);
    gain.gain.exponentialRampToValueAtTime(0.005, ctx.currentTime + 0.09);

    osc.connect(gain);
    gain.connect(ctx.destination);

    osc.start();
    osc.stop(ctx.currentTime + 0.09);
  }

  // 选牌敲击声（清脆筹码/卡牌轻抬）
  function playCardSelect() {
    if (isMuted.value || masterVolume.value <= 0) return;
    const ctx = getAudioContext();
    if (!ctx) return;

    const osc = ctx.createOscillator();
    const gain = ctx.createGain();

    osc.type = 'sine';
    osc.frequency.setValueAtTime(720, ctx.currentTime);
    osc.frequency.exponentialRampToValueAtTime(1100, ctx.currentTime + 0.04);

    const peak = calcGain(0.15);
    gain.gain.setValueAtTime(peak, ctx.currentTime);
    gain.gain.exponentialRampToValueAtTime(0.005, ctx.currentTime + 0.04);

    osc.connect(gain);
    gain.connect(ctx.destination);

    osc.start();
    osc.stop(ctx.currentTime + 0.04);
  }

  // 揭牌 3D 翻面飞掠声
  function playCardFlip() {
    if (isMuted.value || masterVolume.value <= 0) return;
    const ctx = getAudioContext();
    if (!ctx) return;

    const osc = ctx.createOscillator();
    const gain = ctx.createGain();

    osc.type = 'sine';
    osc.frequency.setValueAtTime(280, ctx.currentTime);
    osc.frequency.exponentialRampToValueAtTime(650, ctx.currentTime + 0.08);

    const peak = calcGain(0.22);
    gain.gain.setValueAtTime(peak, ctx.currentTime);
    gain.gain.exponentialRampToValueAtTime(0.005, ctx.currentTime + 0.08);

    osc.connect(gain);
    gain.connect(ctx.destination);

    osc.start();
    osc.stop(ctx.currentTime + 0.08);
  }

  // 左轮手枪空包弹（纯机械金属撞针清脆咔哒声）
  function playGunClick() {
    if (isMuted.value || masterVolume.value <= 0) return;
    const ctx = getAudioContext();
    if (!ctx) return;

    // 撞针敲击
    const osc = ctx.createOscillator();
    const gain = ctx.createGain();
    osc.type = 'square';
    osc.frequency.setValueAtTime(2200, ctx.currentTime);
    osc.frequency.exponentialRampToValueAtTime(450, ctx.currentTime + 0.035);

    const peak = calcGain(0.4);
    gain.gain.setValueAtTime(peak, ctx.currentTime);
    gain.gain.exponentialRampToValueAtTime(0.005, ctx.currentTime + 0.035);

    osc.connect(gain);
    gain.connect(ctx.destination);

    osc.start();
    osc.stop(ctx.currentTime + 0.04);

    // 弹仓余震金属泛音
    const osc2 = ctx.createOscillator();
    const gain2 = ctx.createGain();
    osc2.type = 'sine';
    osc2.frequency.setValueAtTime(850, ctx.currentTime + 0.02);
    osc2.frequency.exponentialRampToValueAtTime(300, ctx.currentTime + 0.1);

    gain2.gain.setValueAtTime(calcGain(0.1), ctx.currentTime + 0.02);
    gain2.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.1);

    osc2.connect(gain2);
    gain2.connect(ctx.destination);
    osc2.start(ctx.currentTime + 0.02);
    osc2.stop(ctx.currentTime + 0.1);
  }

  // 致命实弹枪声（深度低频冲击 + 瞬态白噪爆破 + 混响衰减）
  function playGunshot() {
    if (isMuted.value || masterVolume.value <= 0) return;
    const ctx = getAudioContext();
    if (!ctx) return;

    const t = ctx.currentTime;

    // 1. 低频胸腔震撼冲击波
    const kick = ctx.createOscillator();
    const kickGain = ctx.createGain();
    kick.type = 'sine';
    kick.frequency.setValueAtTime(220, t);
    kick.frequency.exponentialRampToValueAtTime(24, t + 0.7);

    kickGain.gain.setValueAtTime(calcGain(1.0), t);
    kickGain.gain.exponentialRampToValueAtTime(0.005, t + 0.7);
    kick.connect(kickGain);
    kickGain.connect(ctx.destination);
    kick.start(t);
    kick.stop(t + 0.7);

    // 2. 枪口火光破空瞬态
    const bufferSize = ctx.sampleRate * 0.9;
    const buffer = ctx.createBuffer(1, bufferSize, ctx.sampleRate);
    const data = buffer.getChannelData(0);
    for (let i = 0; i < bufferSize; i++) {
      data[i] = Math.random() * 2 - 1;
    }

    const noise = ctx.createBufferSource();
    noise.buffer = buffer;

    const filter = ctx.createBiquadFilter();
    filter.type = 'lowpass';
    filter.frequency.setValueAtTime(4500, t);
    filter.frequency.exponentialRampToValueAtTime(350, t + 0.7);

    const noiseGain = ctx.createGain();
    noiseGain.gain.setValueAtTime(calcGain(1.0), t);
    noiseGain.gain.exponentialRampToValueAtTime(0.005, t + 0.85);

    noise.connect(filter);
    filter.connect(noiseGain);
    noiseGain.connect(ctx.destination);

    noise.start(t);
    noise.stop(t + 0.9);
  }

  // 紧张心跳悬念音（扣动扳机/质疑前夕）
  function playHeartbeat() {
    if (isMuted.value || masterVolume.value <= 0) return;
    const ctx = getAudioContext();
    if (!ctx) return;

    const t = ctx.currentTime;
    const osc = ctx.createOscillator();
    const gain = ctx.createGain();

    osc.type = 'sine';
    osc.frequency.setValueAtTime(58, t);
    osc.frequency.exponentialRampToValueAtTime(36, t + 0.25);

    gain.gain.setValueAtTime(calcGain(0.55), t);
    gain.gain.exponentialRampToValueAtTime(0.005, t + 0.25);

    osc.connect(gain);
    gain.connect(ctx.destination);

    osc.start(t);
    osc.stop(t + 0.25);
  }

  // 质疑警报（重低音氛围低吼 + 危险音阶）
  function playLiarAlert() {
    if (isMuted.value || masterVolume.value <= 0) return;
    const ctx = getAudioContext();
    if (!ctx) return;

    const t = ctx.currentTime;

    const osc = ctx.createOscillator();
    const gain = ctx.createGain();
    osc.type = 'sawtooth';
    osc.frequency.setValueAtTime(120, t);
    osc.frequency.linearRampToValueAtTime(160, t + 0.2);
    osc.frequency.linearRampToValueAtTime(110, t + 0.4);

    gain.gain.setValueAtTime(calcGain(0.45), t);
    gain.gain.exponentialRampToValueAtTime(0.005, t + 0.6);

    const filter = ctx.createBiquadFilter();
    filter.type = 'lowpass';
    filter.frequency.value = 600;

    osc.connect(filter);
    filter.connect(gain);
    gain.connect(ctx.destination);

    osc.start(t);
    osc.stop(t + 0.6);
  }

  // 胜利号角和弦（大气古典管弦和弦）
  function playVictory() {
    if (isMuted.value || masterVolume.value <= 0) return;
    const ctx = getAudioContext();
    if (!ctx) return;

    const chords = [
      { freq: 261.63, delay: 0 },    // C4
      { freq: 329.63, delay: 0.08 }, // E4
      { freq: 392.00, delay: 0.16 }, // G4
      { freq: 523.25, delay: 0.26 }, // C5
    ];

    chords.forEach(({ freq, delay }) => {
      const osc = ctx.createOscillator();
      const gain = ctx.createGain();

      osc.type = 'triangle';
      osc.frequency.value = freq;

      const st = ctx.currentTime + delay;
      gain.gain.setValueAtTime(0, st);
      gain.gain.linearRampToValueAtTime(calcGain(0.3), st + 0.06);
      gain.gain.exponentialRampToValueAtTime(0.001, st + 0.9);

      osc.connect(gain);
      gain.connect(ctx.destination);

      osc.start(st);
      osc.stop(st + 0.95);
    });
  }

  return {
    isMuted,
    masterVolume,
    toggleMute,
    setVolume,
    playCardDeal,
    playCardSelect,
    playCardFlip,
    playGunClick,
    playGunshot,
    playHeartbeat,
    playLiarAlert,
    playVictory,
  };
}
