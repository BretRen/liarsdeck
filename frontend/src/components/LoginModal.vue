<template>
  <div v-if="!isAuthenticated" class="login-modal-overlay">
    <div class="login-modal glass-panel">
      <!-- Top Bar with Language Selector -->
      <div class="modal-top">
        <span class="auth-tag">AUTHENTICATION</span>
        <button class="btn-icon lang-btn" @click="toggleLang">
          {{ lang.toUpperCase() }}
        </button>
      </div>

      <!-- Main Login Container -->
      <div class="login-body">
        <div class="tavern-emblem">
          <div class="emblem-glow"></div>
          <div class="emblem-icon">🃏</div>
        </div>

        <h2 class="tavern-title font-serif">{{ t('login_modal_title') }}</h2>
        <p class="tavern-subtitle">{{ t('login_modal_desc') }}</p>

        <div v-if="authError" class="auth-error-banner">
          ⚠️ {{ authError }}
        </div>

        <!-- Single Action Button -->
        <button
          class="btn-primary full-btn btn-login-action"
          :disabled="isLoggingIn"
          @click="login"
        >
          <span v-if="isLoggingIn" class="spinner"></span>
          <span>{{ isLoggingIn ? t('login_loading_btn') : t('login_btn') }}</span>
        </button>
      </div>

      <!-- Footer Note -->
      <div class="modal-bottom">
        <span class="auth-footer-text">
          🛡️ {{ t('login_secure_tip') }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from '../composables/useI18n';
import { useAuth } from '../composables/useAuth';

const { t, lang, toggleLang } = useI18n();
const { isAuthenticated, isLoggingIn, authError, login } = useAuth();
</script>

<style scoped>
.login-modal-overlay {
  position: fixed;
  inset: 0;
  background: radial-gradient(circle at center, rgba(16, 19, 27, 0.96) 0%, rgba(5, 6, 8, 0.99) 100%);
  backdrop-filter: blur(16px);
  z-index: 99999;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
}

.login-modal {
  width: 100%;
  max-width: 440px;
  background: #12151d;
  border: 1.5px solid var(--accent-gold);
  border-radius: var(--radius-lg);
  box-shadow: 0 24px 64px rgba(0, 0, 0, 0.95), 0 0 40px rgba(212, 175, 55, 0.25);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  animation: modal-enter 0.35s cubic-bezier(0.34, 1.56, 0.64, 1);
}

@keyframes modal-enter {
  from {
    opacity: 0;
    transform: scale(0.92) translateY(10px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

.modal-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 20px;
  background: rgba(10, 12, 18, 0.8);
  border-bottom: 1px solid var(--border-brass);
}

.auth-tag {
  font-size: 10.5px;
  font-weight: 800;
  letter-spacing: 1.5px;
  color: var(--accent-gold);
}

.lang-btn {
  font-size: 11px;
  padding: 3px 8px;
}

.login-body {
  padding: 36px 28px 28px;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.tavern-emblem {
  position: relative;
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(212, 175, 55, 0.2) 0%, rgba(20, 24, 34, 0.8) 100%);
  border: 1.5px solid var(--accent-gold);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
  box-shadow: 0 0 24px rgba(212, 175, 55, 0.3);
}

.emblem-icon {
  font-size: 38px;
  filter: drop-shadow(0 2px 8px rgba(0, 0, 0, 0.6));
}

.tavern-title {
  font-size: 1.6rem;
  color: var(--text-primary);
  margin-bottom: 10px;
  font-weight: 700;
  letter-spacing: 0.5px;
}

.tavern-subtitle {
  font-size: 13.5px;
  color: var(--text-secondary);
  line-height: 1.55;
  max-width: 320px;
  margin-bottom: 26px;
}

.auth-error-banner {
  width: 100%;
  background: rgba(239, 68, 68, 0.15);
  border: 1px solid var(--accent-crimson);
  color: #fca5a5;
  font-size: 12.5px;
  padding: 8px 12px;
  border-radius: var(--radius-md);
  margin-bottom: 20px;
  text-align: left;
  line-height: 1.4;
}

.btn-login-action {
  width: 100%;
  padding: 14px 20px;
  font-size: 15.5px;
  font-weight: 800;
  letter-spacing: 1px;
  background: linear-gradient(180deg, #e5a93c 0%, #b8811b 100%);
  border: 1px solid #f3c267;
  color: #0b0d13;
  border-radius: var(--radius-md);
  cursor: pointer;
  box-shadow: 0 6px 20px rgba(212, 175, 55, 0.35);
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
}

.btn-login-action:hover:not(:disabled) {
  background: linear-gradient(180deg, #f5bc52 0%, #c98e22 100%);
  transform: translateY(-1px);
  box-shadow: 0 8px 26px rgba(212, 175, 55, 0.5);
}

.btn-login-action:active:not(:disabled) {
  transform: translateY(1px);
}

.btn-login-action:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(0, 0, 0, 0.3);
  border-top-color: #000;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.modal-bottom {
  padding: 12px 20px;
  background: rgba(10, 12, 18, 0.8);
  border-top: 1px solid var(--border-brass);
  text-align: center;
}

.auth-footer-text {
  font-size: 11.5px;
  color: var(--text-muted);
}
</style>
