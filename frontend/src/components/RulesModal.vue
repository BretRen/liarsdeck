<template>
  <div class="rules-modal-overlay" @click.self="$emit('close')">
    <div class="rules-modal glass-panel">
      <div class="modal-header">
        <h2 class="font-serif">{{ t('rules_title') }}</h2>
        <button class="btn-icon" @click="$emit('close')">✕</button>
      </div>

      <div class="modal-content">
        <section class="rule-block">
          <div class="rule-title">
            <span class="rule-num">01</span>
            <h3>{{ t('rule_goal_title') }}</h3>
          </div>
          <p>{{ t('rule_goal_desc') }}</p>
        </section>

        <section class="rule-block">
          <div class="rule-title">
            <span class="rule-num">02</span>
            <h3>{{ t('rule_deck_title') }}</h3>
          </div>
          <p>{{ t('rule_deck_desc') }}</p>
          <div class="cards-preview">
            <div class="preview-card rank-K">K (King)</div>
            <div class="preview-card rank-Q">Q (Queen)</div>
            <div class="preview-card rank-A">A (Ace)</div>
            <div class="preview-card rank-2">2 (Wild)</div>
          </div>
        </section>

        <section class="rule-block">
          <div class="rule-title">
            <span class="rule-num">03</span>
            <h3>{{ t('rule_turn_title') }}</h3>
          </div>
          <ul>
            <li>{{ t('rule_turn_1') }}</li>
            <li>{{ t('rule_turn_2') }}</li>
            <li>{{ t('rule_turn_3') }}</li>
          </ul>
        </section>

        <section class="rule-block">
          <div class="rule-title">
            <span class="rule-num">04</span>
            <h3>{{ t('rule_liar_title') }}</h3>
          </div>
          <ul>
            <li>{{ t('rule_liar_1') }}</li>
            <li>{{ t('rule_liar_2') }}</li>
            <li>{{ t('rule_liar_3') }}</li>
          </ul>
        </section>

        <section class="rule-block">
          <div class="rule-title">
            <span class="rule-num">05</span>
            <h3>{{ t('rule_gun_title') }}</h3>
          </div>
          <ul>
            <li>{{ t('rule_gun_1') }}</li>
            <li>{{ t('rule_gun_2') }}</li>
          </ul>
        </section>
      </div>

      <div class="modal-footer">
        <button class="btn-primary full-width" @click="$emit('close')">
          {{ t('rule_got_it') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from '../composables/useI18n';

defineEmits(['close']);
const { t } = useI18n();
</script>

<style scoped>
.rules-modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(8px);
  z-index: 2000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  animation: fadeIn 0.2s ease;
}

.rules-modal {
  width: 100%;
  max-width: 520px;
  max-height: 85vh;
  display: flex;
  flex-direction: column;
  background: #141720;
  border: 1px solid var(--border-brass);
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: 0 16px 40px rgba(0, 0, 0, 0.8);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-brass);
  background: #10131a;
}

.modal-header h2 {
  font-size: 1.1rem;
  color: var(--gold-accent);
  margin: 0;
}

.modal-content {
  padding: 18px 20px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.rule-block {
  background: #0f1218;
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  padding: 12px 16px;
}

.rule-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}

.rule-num {
  font-family: var(--font-serif);
  font-size: 12px;
  font-weight: 700;
  color: var(--gold-accent);
  border: 1px solid var(--border-brass);
  padding: 1px 5px;
  border-radius: 2px;
  background: #181d28;
}

.rule-title h3 {
  font-size: 0.92rem;
  font-weight: 700;
  color: var(--text-primary);
}

.rule-block p, .rule-block li {
  font-size: 0.86rem;
  color: var(--text-secondary);
  line-height: 1.6;
}

.rule-block ul {
  padding-left: 18px;
  margin-top: 4px;
}

.cards-preview {
  display: flex;
  gap: 8px;
  margin-top: 8px;
  flex-wrap: wrap;
}

.preview-card {
  padding: 3px 8px;
  border-radius: 4px;
  font-family: var(--font-serif);
  font-weight: 700;
  font-size: 12px;
  border: 1px solid;
}

.preview-card.rank-K { color: #1a1a1a; border-color: #d9d2c5; background: #fbf8f2; }
.preview-card.rank-Q { color: #991b1b; border-color: #fca5a5; background: #fbf8f2; }
.preview-card.rank-A { color: #1a1a1a; border-color: #d9d2c5; background: #fbf8f2; }
.preview-card.rank-2 { color: #065f46; border-color: #059669; background: #f0fdf4; }

.modal-footer {
  padding: 12px 20px;
  border-top: 1px solid var(--border-brass);
  background: #10131a;
}

.full-width {
  width: 100%;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
</style>
