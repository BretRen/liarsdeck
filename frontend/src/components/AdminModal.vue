<template>
  <div v-if="isOpen" class="fixed inset-0 z-[2000] flex items-center justify-center p-4 bg-slate-950/85 backdrop-blur-md animate-in fade-in duration-200" @click.self="$emit('close')">
    <div class="card w-full max-w-xl bg-slate-900/95 border border-slate-700/80 rounded-2xl shadow-2xl shadow-black/90 flex flex-col max-h-[90vh] overflow-hidden">
      <!-- Modal Header -->
      <div class="flex items-center justify-between px-6 py-4 bg-slate-950/80 border-b border-slate-800">
        <div class="flex items-center gap-2">
          <span class="badge badge-primary badge-xs font-black tracking-widest uppercase">ADMIN</span>
          <h2 class="text-lg font-bold font-serif text-slate-100">{{ t('admin_title') }}</h2>
        </div>
        <div class="flex items-center gap-1.5">
          <button class="btn btn-ghost btn-xs text-slate-400 hover:text-white" @click="toggleLang">
            {{ lang.toUpperCase() }}
          </button>
          <button class="btn btn-ghost btn-xs w-7 h-7 rounded-full text-slate-400 hover:text-white p-0" @click="$emit('close')">✕</button>
        </div>
      </div>

      <!-- State 1: Password Prompt -->
      <div v-if="!isAuthenticated" class="p-8 flex flex-col items-center text-center">
        <div class="w-16 h-16 rounded-2xl bg-indigo-950/50 border border-indigo-500/30 flex items-center justify-center text-3xl mb-4 shadow-lg shadow-indigo-500/10">
          🔐
        </div>
        <h3 class="text-base font-bold text-slate-100 mb-1">{{ t('admin_auth_title') }}</h3>
        <p class="text-xs text-slate-400 mb-6 max-w-xs">{{ t('admin_auth_desc') }}</p>

        <form class="w-full max-w-sm flex flex-col gap-3" @submit.prevent="onLogin">
          <input
            v-model="secretInput"
            type="password"
            :placeholder="t('admin_auth_ph')"
            autofocus
            class="input input-bordered w-full h-11 bg-slate-950/80 border-slate-800 text-slate-100 focus:border-indigo-500 text-sm"
          />
          <button type="submit" class="btn btn-primary w-full h-11 font-bold shadow-lg shadow-indigo-600/30 bg-gradient-to-r from-indigo-600 to-indigo-700 hover:from-indigo-500 hover:to-indigo-600 border-none text-white" :disabled="isCheckingAuth">
            {{ isCheckingAuth ? '...' : t('admin_unlock_btn') }}
          </button>
        </form>
      </div>

      <!-- State 2: Admin Dashboard -->
      <div v-else class="p-6 overflow-y-auto flex flex-col gap-5 text-left text-xs text-slate-300">
        <!-- 1. Version & Update Section -->
        <section class="card bg-slate-950/70 border border-slate-800/90 p-4 flex flex-col gap-3">
          <div class="flex items-center justify-between border-b border-slate-800/80 pb-2">
            <h4 class="font-bold text-slate-100 text-sm">{{ t('admin_version_card') }}</h4>
            <span class="badge badge-neutral badge-sm font-mono font-bold text-indigo-300">{{ currentVersion }}</span>
          </div>

          <div class="grid grid-cols-2 gap-3 p-3 bg-slate-900/80 rounded-xl border border-slate-800">
            <div class="flex flex-col gap-0.5">
              <span class="text-slate-400 text-[11px]">{{ t('admin_curr_version') }}:</span>
              <span class="font-mono font-bold text-slate-200 text-sm">{{ currentVersion }}</span>
            </div>
            <div class="flex flex-col gap-0.5">
              <span class="text-slate-400 text-[11px]">{{ t('admin_latest_version') }}:</span>
              <span class="font-mono font-bold text-sm" :class="hasUpdate ? 'text-emerald-400 animate-pulse' : 'text-slate-200'">
                {{ latestVersion || '—' }}
              </span>
            </div>
          </div>

          <div v-if="releaseNotes" class="p-3 bg-slate-900/60 rounded-xl border border-slate-800 flex flex-col gap-1">
            <div class="font-bold text-indigo-300 text-xs">{{ releaseName }}</div>
            <pre class="font-mono text-[11px] text-slate-400 whitespace-pre-wrap max-h-32 overflow-y-auto">{{ releaseNotes }}</pre>
          </div>

          <div class="flex gap-2">
            <button
              class="btn btn-neutral flex-1 h-10 min-h-0 bg-slate-800 hover:bg-slate-700 border-slate-700 text-slate-200 font-semibold"
              :disabled="isCheckingUpdate || isUpdating"
              @click="checkUpdate"
            >
              {{ isCheckingUpdate ? '...' : t('admin_check_btn') }}
            </button>
            <button
              class="btn btn-primary flex-1 h-10 min-h-0 font-bold shadow-lg shadow-indigo-600/30 bg-gradient-to-r from-indigo-600 to-indigo-700 hover:from-indigo-500 hover:to-indigo-600 border-none text-white disabled:opacity-40"
              :disabled="isUpdating"
              @click="triggerUpdate"
            >
              {{ isUpdating ? t('admin_updating') : t('admin_update_btn') }}
            </button>
          </div>

          <div v-if="updateLog" class="p-3 bg-black/80 rounded-lg border border-slate-800 font-mono text-[11px] text-emerald-400 whitespace-pre-wrap max-h-40 overflow-y-auto">
            {{ updateLog }}
          </div>
        </section>

        <!-- 2. Server Global Broadcast Section -->
        <section class="card bg-slate-950/70 border border-slate-800/90 p-4 flex flex-col gap-3">
          <div class="flex items-center justify-between border-b border-slate-800/80 pb-2">
            <h4 class="font-bold text-slate-100 text-sm">{{ t('admin_broadcast_card') }}</h4>
            <span class="badge badge-error badge-xs font-bold text-white uppercase">LIVE</span>
          </div>

          <div class="flex flex-col gap-2.5">
            <textarea
              v-model="broadcastMsg"
              rows="2"
              :placeholder="t('admin_broadcast_ph')"
              class="textarea textarea-bordered w-full bg-slate-900/80 border-slate-800 text-slate-100 focus:border-indigo-500 text-xs resize-none"
              :disabled="isBroadcasting"
            ></textarea>

            <div class="flex gap-2">
              <button
                type="button"
                class="btn btn-ghost btn-xs bg-slate-800/70 hover:bg-slate-800 text-slate-300 border border-slate-700/60"
                @click="broadcastMsg = t('admin_broadcast_preset_1')"
              >
                🛠️ {{ lang === 'zh' ? '维护预设' : 'Maintenance' }}
              </button>
              <button
                type="button"
                class="btn btn-ghost btn-xs bg-slate-800/70 hover:bg-slate-800 text-slate-300 border border-slate-700/60"
                @click="broadcastMsg = t('admin_broadcast_preset_2')"
              >
                🎮 {{ lang === 'zh' ? '欢迎预设' : 'Welcome' }}
              </button>
            </div>

            <button
              class="btn btn-primary w-full h-10 min-h-0 font-bold shadow-lg shadow-indigo-600/30 bg-gradient-to-r from-indigo-600 to-indigo-700 hover:from-indigo-500 hover:to-indigo-600 border-none text-white disabled:opacity-40"
              :disabled="isBroadcasting || !broadcastMsg.trim()"
              @click="sendBroadcast"
            >
              {{ isBroadcasting ? t('admin_broadcast_sending') : t('admin_broadcast_send_btn') }}
            </button>
          </div>
        </section>

        <!-- 3. Metrics & Stats Section -->
        <section class="card bg-slate-950/70 border border-slate-800/90 p-4 flex flex-col gap-3">
          <div class="flex items-center justify-between border-b border-slate-800/80 pb-2">
            <h4 class="font-bold text-slate-100 text-sm">{{ t('admin_stats_card') }}</h4>
            <button class="btn btn-ghost btn-xs text-slate-400 hover:text-white p-0 h-5 w-5" @click="fetchStats">⟳</button>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div class="p-3 bg-slate-900/80 rounded-xl border border-slate-800 flex flex-col items-center gap-0.5">
              <span class="text-2xl font-black font-mono text-indigo-400">{{ stats.total_rooms }}</span>
              <span class="text-[11px] text-slate-400">{{ t('admin_active_rooms') }}</span>
            </div>
            <div class="p-3 bg-slate-900/80 rounded-xl border border-slate-800 flex flex-col items-center gap-0.5">
              <span class="text-2xl font-black font-mono text-sky-400">{{ stats.total_players }}</span>
              <span class="text-[11px] text-slate-400">{{ t('admin_active_players') }}</span>
            </div>
          </div>
        </section>
      </div>

      <!-- Modal Footer -->
      <div class="p-4 bg-slate-950/70 border-t border-slate-800 flex justify-end">
        <button class="btn btn-neutral btn-sm bg-slate-800 hover:bg-slate-700 border-slate-700 text-slate-300 font-semibold" @click="$emit('close')">
          {{ t('admin_close_btn') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch } from 'vue';
import { useI18n } from '../composables/useI18n';
import { useGameStore } from '../composables/useGameStore';

const props = defineProps({
  isOpen: { type: Boolean, default: false },
});

defineEmits(['close']);
const { t, lang, toggleLang } = useI18n();
const { showToast } = useGameStore();

const isAuthenticated = ref(false);
const isCheckingAuth = ref(false);
const secretInput = ref('');
const savedSecret = ref('');

const currentVersion = ref(typeof __APP_VERSION__ !== 'undefined' ? __APP_VERSION__ : 'v2.3.0');
const latestVersion = ref('');
const hasUpdate = ref(false);
const releaseName = ref('');
const releaseNotes = ref('');
const isCheckingUpdate = ref(false);
const isUpdating = ref(false);
const updateLog = ref('');

const broadcastMsg = ref('');
const isBroadcasting = ref(false);

const stats = reactive({
  total_rooms: 0,
  total_players: 0,
});

async function fetchServerVersion() {
  try {
    const res = await fetch('/api/version');
    if (res.ok) {
      const data = await res.json();
      if (data.version) {
        currentVersion.value = data.version;
      }
    }
  } catch (_) {}
}

watch(
  () => props.isOpen,
  (val) => {
    if (val) {
      fetchServerVersion();
      if (isAuthenticated.value) {
        fetchStats();
      }
    }
  },
  { immediate: true }
);

async function onLogin() {
  if (!secretInput.value.trim()) return;
  isCheckingAuth.value = true;
  try {
    const res = await fetch('/api/admin/auth', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ secret: secretInput.value.trim() }),
    });
    const data = await res.json();
    if (res.ok && data.authenticated) {
      isAuthenticated.value = true;
      savedSecret.value = secretInput.value.trim();
      currentVersion.value = data.version || 'v2.2.0';
      showToast('管理员身份验证通过');
      fetchStats();
    } else {
      showToast(data.error || '密钥错误');
    }
  } catch (e) {
    showToast('验证请求失败: ' + e.message);
  } finally {
    isCheckingAuth.value = false;
  }
}

async function checkUpdate() {
  isCheckingUpdate.value = true;
  updateLog.value = '';
  try {
    const res = await fetch('/api/admin/check-update', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ secret: savedSecret.value }),
    });
    const data = await res.json();
    if (res.ok) {
      currentVersion.value = data.current_version;
      latestVersion.value = data.latest_version;
      hasUpdate.value = data.has_update;
      releaseName.value = data.release_name || '';
      releaseNotes.value = data.release_body || '';
      if (!data.has_update) {
        showToast('当前已是最新版本');
      }
    } else {
      showToast(data.error || '检查更新失败');
    }
  } catch (e) {
    showToast('检查更新异常: ' + e.message);
  } finally {
    isCheckingUpdate.value = false;
  }
}

async function triggerUpdate() {
  if (isUpdating.value) return;
  if (!confirm('确定要拉起自动更新程序并重启服务端吗？此操作将热替换程序二进制并重启。')) {
    return;
  }

  isUpdating.value = true;
  updateLog.value = '⏳ 正在向服务端发送更新指令...';

  try {
    const res = await fetch('/api/admin/trigger-update', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ secret: savedSecret.value }),
    });
    const data = await res.json();
    if (res.ok && data.success) {
      updateLog.value = `🚀 ${data.message}\n⏳ 服务端正在执行热替换并重启，请在 5-10 秒后刷新网页！`;
      showToast('更新程序已启动！');
    } else {
      updateLog.value = `❌ 更新失败: ${data.error || '未知错误'}`;
      showToast(data.error || '触发更新失败');
      isUpdating.value = false;
    }
  } catch (e) {
    updateLog.value = `❌ 网络异常: ${e.message}\nℹ️ 服务端可能已经重启完成，请尝试刷新页面。`;
    isUpdating.value = false;
  }
}

async function fetchStats() {
  try {
    const res = await fetch('/api/admin/stats', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ secret: savedSecret.value }),
    });
    const data = await res.json();
    if (res.ok) {
      stats.total_rooms = data.total_rooms || 0;
      stats.total_players = data.total_players || 0;
    }
  } catch (_) {}
}

async function sendBroadcast() {
  if (!broadcastMsg.value.trim() || isBroadcasting.value) return;
  isBroadcasting.value = true;
  try {
    const res = await fetch('/api/admin/broadcast', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        secret: savedSecret.value,
        message: broadcastMsg.value.trim(),
      }),
    });
    const data = await res.json();
    if (res.ok && data.success) {
      showToast(`广播已发送给 ${data.broadcast_count} 个在线房间！`);
      broadcastMsg.value = '';
    } else {
      showToast(data.error || '广播发送失败');
    }
  } catch (e) {
    showToast('发送广播异常: ' + e.message);
  } finally {
    isBroadcasting.value = false;
  }
}
</script>
