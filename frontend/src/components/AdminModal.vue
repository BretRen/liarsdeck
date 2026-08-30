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

        <!-- 2. Metrics & Stats Section -->
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

const currentVersion = ref('v2.0.5');
const latestVersion = ref('');
const hasUpdate = ref(false);
const releaseName = ref('');
const releaseNotes = ref('');
const isCheckingUpdate = ref(false);
const isUpdating = ref(false);
const updateLog = ref('');

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
      if (data.has_update) {
        showToast(`发现新版本 ${data.latest_version}`);
      } else {
        showToast('当前已是最新版本');
      }
    } else {
      showToast(data.error || '检查失败');
    }
  } catch (e) {
    showToast('网络错误: ' + e.message);
  } finally {
    isCheckingUpdate.value = false;
  }
}

async function triggerUpdate() {
  if (!confirm('确认要在后台下载并替换当前服务端进程吗？')) return;
  isUpdating.value = true;
  updateLog.value = '正在触发 update.go 进程...\n';
  try {
    const res = await fetch('/api/admin/trigger-update', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ secret: savedSecret.value }),
    });
    const data = await res.json();
    if (res.ok) {
      updateLog.value += `✅ ${data.message}\n服务正在后台平滑热替换并重启，请稍候 3~5 秒自动重连...`;
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
  max-height: 85vh;
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
  font-size: 9.5px;
  font-weight: 800;
  letter-spacing: 1px;
  padding: 2px 6px;
  background: #3b1414;
  border: 1px solid #742a2a;
  color: #ff8787;
  border-radius: 3px;
}

.modal-header h2 {
  font-size: 1.15rem;
  color: var(--gold-accent);
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.lang-btn {
  font-weight: 700;
}

.modal-content {
  padding: 20px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* Auth View */
.auth-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 12px;
  padding: 20px 10px;
}

.auth-icon-circle {
  font-size: 36px;
}

.auth-desc {
  font-size: 13px;
  color: var(--text-secondary);
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 10px;
  width: 100%;
  max-width: 320px;
}

.secret-input {
  text-align: center;
  letter-spacing: 2px;
}

/* Dashboard View */
.admin-card {
  background: #0f121a;
  border: 1px solid var(--border-brass);
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
}

.card-header h4 {
  font-size: 0.92rem;
  color: var(--text-primary);
  margin: 0;
}

.curr-tag {
  font-family: var(--font-serif);
  font-size: 11px;
  font-weight: 700;
  color: var(--gold-accent);
  background: #1a1e28;
  padding: 2px 8px;
  border-radius: 3px;
  border: 1px solid var(--border-brass);
}

.version-meta-grid {
  display: flex;
  gap: 20px;
  font-size: 13px;
}

.meta-item {
  display: flex;
  align-items: baseline;
  gap: 6px;
}

.meta-label {
  color: var(--text-muted);
}

.meta-val {
  font-weight: 700;
  color: var(--text-primary);
}

.meta-val.new-version {
  color: #51cf66;
}

.release-notes-box {
  background: #141720;
  border: 1px solid var(--border-subtle);
  border-radius: 4px;
  padding: 8px 12px;
  max-height: 120px;
  overflow-y: auto;
}

.notes-title {
  font-weight: 700;
  font-size: 12px;
  color: var(--gold-accent);
  margin-bottom: 4px;
}

.notes-body {
  font-family: var(--font-sans);
  font-size: 11.5px;
  color: var(--text-secondary);
  white-space: pre-wrap;
  margin: 0;
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
