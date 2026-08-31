<template>
  <div v-if="isOpen" class="admin-modal-overlay" @click.self="$emit('close')">
    <div class="admin-modal glass-panel">
      <!-- Modal Header -->
      <div class="modal-header">
        <div class="header-left">
          <span class="admin-badge">ADMIN CONSOLE</span>
          <h2 class="font-serif">{{ t('admin_title') }}</h2>
        </div>
        <div class="header-right">
          <button class="btn-icon lang-btn" @click="toggleLang">
            {{ lang.toUpperCase() }}
          </button>
          <button class="btn-icon" @click="$emit('close')">✕</button>
        </div>
      </div>

      <!-- State 1: Password Prompt -->
      <div v-if="!isAuthenticated" class="modal-content auth-content">
        <div class="auth-box">
          <div class="auth-icon-circle">🔐</div>
          <h3>{{ t('admin_auth_title') }}</h3>
          <p class="auth-desc">{{ t('admin_auth_desc') }}</p>

          <form class="auth-form" @submit.prevent="onLogin">
            <input
              v-model="secretInput"
              type="password"
              :placeholder="t('admin_auth_ph')"
              autofocus
              class="secret-input"
            />
            <button type="submit" class="btn-primary full-btn" :disabled="isCheckingAuth">
              {{ isCheckingAuth ? '...' : t('admin_unlock_btn') }}
            </button>
          </form>
        </div>
      </div>

      <!-- State 2: Admin Dashboard -->
      <div v-else class="modal-content dashboard-content">
        <!-- 1. Version & Update Section -->
        <section class="admin-card">
          <div class="card-header">
            <h4>{{ t('admin_version_card') }}</h4>
            <span class="curr-tag">{{ currentVersion }}</span>
          </div>

          <div class="version-meta-grid">
            <div class="meta-item">
              <span class="meta-label">{{ t('admin_curr_version') }}:</span>
              <span class="meta-val font-serif">{{ currentVersion }}</span>
            </div>
            <div class="meta-item">
              <span class="meta-label">{{ t('admin_latest_version') }}:</span>
              <span class="meta-val font-serif" :class="{ 'new-version': hasUpdate }">
                {{ latestVersion || '—' }}
              </span>
            </div>
          </div>

          <div v-if="releaseNotes" class="release-notes-box">
            <div class="notes-title">{{ releaseName }}</div>
            <pre class="notes-body">{{ releaseNotes }}</pre>
          </div>

          <div class="update-actions">
            <button
              class="btn-secondary flex-1"
              :disabled="isCheckingUpdate || isUpdating"
              @click="checkUpdate"
            >
              {{ isCheckingUpdate ? '...' : t('admin_check_btn') }}
            </button>
            <button
              class="btn-primary flex-1 btn-update-trigger"
              :disabled="isUpdating"
              @click="triggerUpdate"
            >
              {{ isUpdating ? t('admin_updating') : t('admin_update_btn') }}
            </button>
          </div>

          <div v-if="updateLog" class="update-log-terminal">
            {{ updateLog }}
          </div>
        </section>

        <!-- 2. Server Global Broadcast Section -->
        <section class="admin-card">
          <div class="card-header">
            <h4>{{ t('admin_broadcast_card') }}</h4>
            <span class="broadcast-badge">LIVE</span>
          </div>

          <div class="broadcast-form">
            <textarea
              v-model="broadcastMsg"
              rows="2"
              :placeholder="t('admin_broadcast_ph')"
              class="broadcast-textarea"
              :disabled="isBroadcasting"
            ></textarea>

            <div class="broadcast-presets">
              <button
                type="button"
                class="preset-btn"
                @click="broadcastMsg = t('admin_broadcast_preset_1')"
              >
                🛠️ {{ lang === 'zh' ? '维护预设' : 'Maintenance' }}
              </button>
              <button
                type="button"
                class="preset-btn"
                @click="broadcastMsg = t('admin_broadcast_preset_2')"
              >
                🎮 {{ lang === 'zh' ? '欢迎预设' : 'Welcome' }}
              </button>
            </div>

            <button
              class="btn-primary full-btn btn-broadcast-send"
              :disabled="isBroadcasting || !broadcastMsg.trim()"
              @click="sendBroadcast"
            >
              {{ isBroadcasting ? t('admin_broadcast_sending') : t('admin_broadcast_send_btn') }}
            </button>
          </div>
        </section>

        <!-- 3. Metrics & Stats Section -->
        <section class="admin-card">
          <div class="card-header">
            <h4>{{ t('admin_stats_card') }}</h4>
            <button class="btn-icon btn-refresh" @click="fetchStats">⟳</button>
          </div>

          <div class="stats-grid">
            <div class="stat-pill">
              <span class="stat-num font-serif">{{ stats.total_rooms }}</span>
              <span class="stat-lbl">{{ t('admin_active_rooms') }}</span>
            </div>
            <div class="stat-pill">
              <span class="stat-num font-serif">{{ stats.total_players }}</span>
              <span class="stat-lbl">{{ t('admin_active_players') }}</span>
            </div>
          </div>
        </section>
      </div>

      <!-- Modal Footer -->
      <div class="modal-footer">
        <button class="btn-secondary" @click="$emit('close')">
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

const currentVersion = ref('v2.0.6');
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

watch(
  () => props.isOpen,
  (val) => {
    if (val && isAuthenticated.value) {
      fetchStats();
    }
  }
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
      currentVersion.value = data.version || 'v2.0.0';
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
  if (!confirm('确定要触发服务端自更新并热重启吗？')) return;
  isUpdating.value = true;
  updateLog.value = '⏳ 正在向服务端发送热更新指令...\n';
  try {
    const res = await fetch('/api/admin/trigger-update', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ secret: savedSecret.value }),
    });
    const data = await res.json();
    if (res.ok) {
      updateLog.value += `✅ ${data.message}\n请等待约 3~5 秒后刷新网页连接新版本。`;
      showToast('更新程序已启动');
    } else {
      updateLog.value += `❌ ${data.error}`;
      showToast(data.error || '更新失败');
    }
  } catch (e) {
    updateLog.value += `❌ 请求异常: ${e.message}`;
  } finally {
    isUpdating.value = false;
  }
}

async function sendBroadcast() {
  if (!broadcastMsg.value.trim()) return;
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
    if (res.ok) {
      showToast(data.message || '广播发送成功');
      broadcastMsg.value = '';
    } else {
      showToast(data.error || '广播发送失败');
    }
  } catch (e) {
    showToast('发送异常: ' + e.message);
  } finally {
    isBroadcasting.value = false;
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
      stats.total_rooms = data.total_rooms;
      stats.total_players = data.total_players;
    }
  } catch (_) {}
}
</script>

<style scoped>
.admin-modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.85);
  backdrop-filter: blur(10px);
  z-index: 4000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  animation: fadeIn 0.2s ease;
}

.admin-modal {
  width: 100%;
  max-width: 540px;
  max-height: 88vh;
  display: flex;
  flex-direction: column;
  background: #141720;
  border: 1.5px solid #d4af37;
  border-radius: var(--radius-lg);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.9), 0 0 30px rgba(212, 175, 55, 0.25);
  overflow: hidden;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 20px;
  background: #10131a;
  border-bottom: 1px solid var(--border-brass);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.admin-badge {
  background: linear-gradient(135deg, #d4af37 0%, #aa820a 100%);
  color: #0b0d13;
  font-size: 10px;
  font-weight: 900;
  padding: 2px 6px;
  border-radius: 4px;
  letter-spacing: 0.5px;
}

.header-left h2 {
  font-size: 1.15rem;
  color: var(--text-primary);
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.modal-content {
  padding: 20px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* Auth Box */
.auth-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 20px 10px;
}

.auth-icon-circle {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: rgba(212, 175, 55, 0.1);
  border: 1px solid var(--border-brass);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  margin-bottom: 14px;
}

.auth-box h3 {
  font-size: 1.1rem;
  margin-bottom: 6px;
  color: var(--text-primary);
}

.auth-desc {
  font-size: 12.5px;
  color: var(--text-muted);
  max-width: 360px;
  margin-bottom: 18px;
  line-height: 1.4;
}

.auth-form {
  width: 100%;
  max-width: 320px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.secret-input {
  width: 100%;
  padding: 10px 14px;
  background: #0d0f15;
  border: 1px solid var(--border-brass);
  border-radius: var(--radius-md);
  color: #fff;
  font-size: 14px;
  text-align: center;
  letter-spacing: 1px;
}

.secret-input:focus {
  outline: none;
  border-color: var(--gold-accent);
  box-shadow: 0 0 10px rgba(212, 175, 55, 0.3);
}

/* Dashboard Cards */
.admin-card {
  background: #0d1017;
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  padding: 14px 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  padding-bottom: 8px;
}

.card-header h4 {
  font-size: 13.5px;
  color: var(--text-primary);
  margin: 0;
}

.curr-tag {
  font-size: 11px;
  font-family: monospace;
  background: rgba(255, 255, 255, 0.08);
  padding: 2px 6px;
  border-radius: 4px;
  color: var(--text-secondary);
}

.broadcast-badge {
  font-size: 10px;
  font-weight: 800;
  color: #f59e0b;
  background: rgba(245, 158, 11, 0.15);
  border: 1px solid rgba(245, 158, 11, 0.4);
  padding: 1px 6px;
  border-radius: 4px;
  letter-spacing: 0.5px;
}

.version-meta-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.meta-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
  background: #141720;
  padding: 8px 10px;
  border-radius: 4px;
}

.meta-label {
  font-size: 11px;
  color: var(--text-muted);
}

.meta-val {
  font-size: 14px;
  font-weight: 700;
  color: var(--text-primary);
}

.meta-val.new-version {
  color: #51cf66;
}

.release-notes-box {
  background: #080a0f;
  border: 1px solid var(--border-subtle);
  border-radius: 4px;
  padding: 8px 10px;
  max-height: 120px;
  overflow-y: auto;
}

.notes-title {
  font-size: 12px;
  font-weight: 700;
  color: var(--gold-accent);
  margin-bottom: 4px;
}

.notes-body {
  font-size: 11px;
  color: var(--text-secondary);
  white-space: pre-wrap;
  margin: 0;
  font-family: inherit;
}

.update-actions {
  display: flex;
  gap: 10px;
}

.flex-1 {
  flex: 1;
}

.btn-update-trigger {
  background: linear-gradient(180deg, #2b8a3e 0%, #1e662c 100%);
  border-color: #3ca852;
  color: #fff;
}
.btn-update-trigger:hover:not(:disabled) {
  background: linear-gradient(180deg, #329946 0%, #237533 100%);
}

.update-log-terminal {
  background: #050608;
  border: 1px solid var(--border-brass);
  border-radius: 4px;
  padding: 8px 10px;
  font-family: monospace;
  font-size: 11px;
  color: #51cf66;
  white-space: pre-wrap;
}

/* Broadcast Section */
.broadcast-form {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.broadcast-textarea {
  width: 100%;
  padding: 8px 12px;
  background: #080a0f;
  border: 1px solid var(--border-brass);
  border-radius: var(--radius-md);
  color: #fff;
  font-size: 13px;
  resize: vertical;
  font-family: inherit;
  line-height: 1.4;
}

.broadcast-textarea:focus {
  outline: none;
  border-color: var(--gold-accent);
}

.broadcast-presets {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.preset-btn {
  background: rgba(255, 255, 255, 0.05);
  border: 1px dashed var(--border-subtle);
  color: var(--text-muted);
  font-size: 11.5px;
  padding: 4px 8px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.preset-btn:hover {
  background: rgba(212, 175, 55, 0.1);
  border-color: var(--gold-accent);
  color: var(--text-primary);
}

.btn-broadcast-send {
  background: linear-gradient(180deg, #d97706 0%, #b45309 100%);
  border-color: #f59e0b;
  color: #fff;
  font-weight: 700;
  margin-top: 2px;
}

.btn-broadcast-send:hover:not(:disabled) {
  background: linear-gradient(180deg, #f59e0b 0%, #d97706 100%);
}

.stats-grid {
  display: flex;
  gap: 12px;
}

.stat-pill {
  flex: 1;
  background: #141720;
  border: 1px solid var(--border-subtle);
  border-radius: 6px;
  padding: 10px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.stat-num {
  font-size: 1.6rem;
  font-weight: 900;
  color: var(--gold-accent);
  line-height: 1;
}

.stat-lbl {
  font-size: 11px;
  color: var(--text-muted);
}

.btn-refresh {
  font-size: 12px;
  padding: 2px 6px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  padding: 12px 20px;
  background: #10131a;
  border-top: 1px solid var(--border-brass);
}

@keyframes fadeIn {
  from { opacity: 0; transform: scale(0.97); }
  to { opacity: 1; transform: scale(1); }
}
</style>
