import { ref } from 'vue';

const isMuted = ref(localStorage.getItem('liarsdeck_muted') === 'true');

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

export function useAudio() {
  function toggleMute() {
    isMuted.value = !isMuted.value;
    localStorage.setItem('liarsdeck_muted', isMuted.value ? 'true' : 'false');
  }

  function playCardDeal() {
    if (isMuted.value) return;
    const ctx = getAudioContext();
    if (!ctx) return;

    const osc = ctx.createOscillator();
    const gain = ctx.createGain();

    osc.type = 'triangle';
    osc.frequency.setValueAtTime(440, ctx.currentTime);
    osc.frequency.exponentialRampToValueAtTime(120, ctx.currentTime + 0.08);

    gain.gain.setValueAtTime(0.2, ctx.currentTime);
    gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.08);

    osc.connect(gain);
    gain.connect(ctx.destination);

    osc.start();
    osc.stop(ctx.currentTime + 0.08);
  }

  function playCardSelect() {
    if (isMuted.value) return;
    const ctx = getAudioContext();
    if (!ctx) return;

    const osc = ctx.createOscillator();
    const gain = ctx.createGain();

    osc.type = 'sine';
    osc.frequency.setValueAtTime(580, ctx.currentTime);
    osc.frequency.exponentialRampToValueAtTime(880, ctx.currentTime + 0.05);

    gain.gain.setValueAtTime(0.15, ctx.currentTime);
    gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.05);

    osc.connect(gain);
    gain.connect(ctx.destination);

    osc.start();
    osc.stop(ctx.currentTime + 0.05);
  }

  function playGunClick() {
    if (isMuted.value) return;
    const ctx = getAudioContext();
    if (!ctx) return;

    // Metallic dry click
    const osc = ctx.createOscillator();
    const gain = ctx.createGain();

    osc.type = 'square';
    osc.frequency.setValueAtTime(1600, ctx.currentTime);
    osc.frequency.exponentialRampToValueAtTime(300, ctx.currentTime + 0.04);

    gain.gain.setValueAtTime(0.3, ctx.currentTime);
    gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.04);

    osc.connect(gain);
    gain.connect(ctx.destination);

    osc.start();
    osc.stop(ctx.currentTime + 0.04);
  }

  function playGunshot() {
    if (isMuted.value) return;
    const ctx = getAudioContext();
    if (!ctx) return;

    // 1. Low punch kick
    const kick = ctx.createOscillator();
    const kickGain = ctx.createGain();
    kick.type = 'sine';
    kick.frequency.setValueAtTime(150, ctx.currentTime);
    kick.frequency.exponentialRampToValueAtTime(25, ctx.currentTime + 0.5);

    kickGain.gain.setValueAtTime(0.8, ctx.currentTime);
    kickGain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.5);
    kick.connect(kickGain);
    kickGain.connect(ctx.destination);
    kick.start();
    kick.stop(ctx.currentTime + 0.5);

    // 2. White noise blast
    const bufferSize = ctx.sampleRate * 0.8;
    const buffer = ctx.createBuffer(1, bufferSize, ctx.sampleRate);
    const data = buffer.getChannelData(0);
    for (let i = 0; i < bufferSize; i++) {
      data[i] = Math.random() * 2 - 1;
    }

    const noise = ctx.createBufferSource();
    noise.buffer = buffer;

    const filter = ctx.createBiquadFilter();
    filter.type = 'lowpass';
    filter.frequency.setValueAtTime(3000, ctx.currentTime);
    filter.frequency.linearRampToValueAtTime(400, ctx.currentTime + 0.6);

    const noiseGain = ctx.createGain();
    noiseGain.gain.setValueAtTime(0.9, ctx.currentTime);
    noiseGain.gain.exponentialRampToValueAtTime(0.005, ctx.currentTime + 0.7);

    noise.connect(filter);
    filter.connect(noiseGain);
    noiseGain.connect(ctx.destination);

    noise.start();
    noise.stop(ctx.currentTime + 0.8);
  }

  function playLiarAlert() {
    if (isMuted.value) return;
    const ctx = getAudioContext();
    if (!ctx) return;

    const osc = ctx.createOscillator();
    const gain = ctx.createGain();

    osc.type = 'sawtooth';
    osc.frequency.setValueAtTime(750, ctx.currentTime);
    osc.frequency.linearRampToValueAtTime(950, ctx.currentTime + 0.15);
    osc.frequency.linearRampToValueAtTime(750, ctx.currentTime + 0.3);
    osc.frequency.linearRampToValueAtTime(950, ctx.currentTime + 0.45);

    gain.gain.setValueAtTime(0.3, ctx.currentTime);
    gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.55);

    osc.connect(gain);
    gain.connect(ctx.destination);

    osc.start();
    osc.stop(ctx.currentTime + 0.55);
  }

  function playVictory() {
    if (isMuted.value) return;
    const ctx = getAudioContext();
    if (!ctx) return;

    const notes = [440, 554.37, 659.25, 880];
    notes.forEach((freq, idx) => {
      const osc = ctx.createOscillator();
      const gain = ctx.createGain();

      osc.type = 'triangle';
      osc.frequency.value = freq;

      const startTime = ctx.currentTime + idx * 0.12;
      gain.gain.setValueAtTime(0, startTime);
      gain.gain.linearRampToValueAtTime(0.3, startTime + 0.05);
      gain.gain.exponentialRampToValueAtTime(0.01, startTime + 0.6);

      osc.connect(gain);
      gain.connect(ctx.destination);

      osc.start(startTime);
      osc.stop(startTime + 0.65);
    });
  }

  return {
    isMuted,
    toggleMute,
    playCardDeal,
    playCardSelect,
    playGunClick,
    playGunshot,
    playLiarAlert,
    playVictory,
  };
}
