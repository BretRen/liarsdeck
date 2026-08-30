<template>
  <div class="rules-modal-overlay" @click.self="$emit('close')">
    <div class="rules-modal glass-panel">
      <div class="modal-header">
        <h2>{{ t('rules_title') }}</h2>
        <button class="btn-icon" @click="$emit('close')">✕</button>
      </div>

      <div class="modal-content">
        <section class="rule-block">
          <div class="rule-title">
            <span class="icon">🎯</span>
            <h3>{{ t('rule_goal_title') }}</h3>
          </div>
          <p>{{ t('rule_goal_desc') }}</p>
        </section>

        <section class="rule-block">
          <div class="rule-title">
            <span class="icon">🎴</span>
            <h3>{{ t('rule_deck_title') }}</h3>
          </div>
          <p>{{ t('rule_deck_desc') }}</p>
          <div class="cards-preview">
            <div class="preview-card rank-K">K</div>
            <div class="preview-card rank-Q">Q</div>
            <div class="preview-card rank-A">A</div>
            <div class="preview-card rank-2">2 (Wild)</div>
          </div>
        </section>

        <section class="rule-block">
          <div class="rule-title">
            <span class="icon">🔁</span>
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
            <span class="icon">🚨</span>
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
            <span class="icon">🔫</span>
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
  background: rgba(0, 0, 0, 0.75);
  backdrop-filter: blur(10px);
  z-index: 2000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  animation: fadeIn 0.25s ease;
}

.rules-modal {
  width: 100%;
  max-width: 540px;
  max-height: 85vh;
  display: flex;
  flex-direction: column;
  background: #141720;
  border: 1px solid var(--border-gold);
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: 0 16px 40px rgba(0, 0, 0, 0.8), 0 0 30px rgba(245, 158, 11, 0.15);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 18px 24px;
  border-bottom: 1px solid var(--border-subtle);
  background: rgba(245, 158, 11, 0.05);
}

.modal-header h2 {
  font-size: 1.15rem;
  color: var(--accent-gold);
  margin: 0;
}

.modal-content {
  padding: 20px 24px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.rule-block {
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  padding: 14px 16px;
}

.rule-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}

.rule-title h3 {
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--text-primary);
}

.rule-block p, .rule-block li {
  font-size: 0.88rem;
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
  margin-top: 10px;
}

.preview-card {
  padding: 4px 10px;
  border-radius: 6px;
  font-family: var(--font-heading);
  font-weight: 700;
  font-size: 13px;
  border: 1px solid;
}

.preview-card.rank-K { color: #fbbf24; border-color: #92400e; background: #2a1b0d; }
.preview-card.rank-Q { color: #f472b6; border-color: #831843; background: #260c18; }
.preview-card.rank-A { color: #60a5fa; border-color: #1e3a8a; background: #0c1a26; }
.preview-card.rank-2 { color: #34d399; border-color: #065f46; background: #0c261d; }

.modal-footer {
  padding: 14px 24px;
  border-top: 1px solid var(--border-subtle);
  background: rgba(0, 0, 0, 0.2);
}

.full-width {
  width: 100%;
}

@keyframes fadeIn {
  from { opacity: 0; transform: scale(0.95); }
  to { opacity: 1; transform: scale(1); }
}
</style>
