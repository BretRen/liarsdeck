<template>
  <div
    class="fixed inset-0 z-[2000] flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-md"
    @click.self="$emit('close')"
  >
    <div
      class="card w-full max-w-lg bg-slate-900/95 border border-slate-700/80 rounded-2xl shadow-2xl shadow-black/80 flex flex-col max-h-[88vh] overflow-hidden rules-modal-enter"
    >
      <!-- Header -->
      <div class="flex items-center justify-between px-6 py-4 border-b border-slate-800 bg-slate-950/70 shrink-0">
        <h2 class="text-xl font-bold font-serif text-slate-100 tracking-wide">{{ t('rules_title') }}</h2>
        <button class="btn btn-ghost btn-xs w-7 h-7 rounded-full text-slate-400 hover:text-white p-0" @click="$emit('close')">✕</button>
      </div>

      <!-- Scrollable content -->
      <div class="p-6 overflow-y-auto flex flex-col gap-0 text-left text-sm text-slate-300">

        <!-- ── 基础规则 ── -->
        <div class="flex items-center gap-2 mb-3">
          <div class="h-px flex-1 bg-slate-800"></div>
          <span class="text-[10px] font-bold text-slate-500 uppercase tracking-widest px-2">{{ t('rule_section_basic') }}</span>
          <div class="h-px flex-1 bg-slate-800"></div>
        </div>

        <div class="flex flex-col gap-4 mb-6">
          <section class="flex flex-col gap-1.5 rule-section" style="--i:0">
            <div class="flex items-center gap-2">
              <span class="text-xs font-mono font-black text-indigo-400 bg-indigo-950/80 border border-indigo-500/40 px-1.5 py-0.5 rounded">01</span>
              <h3 class="font-bold text-slate-100">{{ t('rule_goal_title') }}</h3>
            </div>
            <p class="text-slate-400 text-xs leading-relaxed pl-7">{{ t('rule_goal_desc') }}</p>
          </section>

          <section class="flex flex-col gap-1.5 rule-section" style="--i:1">
            <div class="flex items-center gap-2">
              <span class="text-xs font-mono font-black text-indigo-400 bg-indigo-950/80 border border-indigo-500/40 px-1.5 py-0.5 rounded">02</span>
              <h3 class="font-bold text-slate-100">{{ t('rule_deck_title') }}</h3>
            </div>
            <p class="text-slate-400 text-xs leading-relaxed pl-7 mb-2">{{ t('rule_deck_desc') }}</p>
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-2 pl-7">
              <div class="p-2 text-center rounded-lg bg-slate-950/80 border border-slate-800 text-xs font-bold text-slate-200">K (King)</div>
              <div class="p-2 text-center rounded-lg bg-slate-950/80 border border-slate-800 text-xs font-bold text-rose-400">Q (Queen)</div>
              <div class="p-2 text-center rounded-lg bg-slate-950/80 border border-slate-800 text-xs font-bold text-slate-200">A (Ace)</div>
              <div class="p-2 text-center rounded-lg bg-indigo-950/80 border border-indigo-500/50 text-xs font-bold text-indigo-300">2 (Wild)</div>
            </div>
          </section>

          <section class="flex flex-col gap-1.5 rule-section" style="--i:2">
            <div class="flex items-center gap-2">
              <span class="text-xs font-mono font-black text-indigo-400 bg-indigo-950/80 border border-indigo-500/40 px-1.5 py-0.5 rounded">03</span>
              <h3 class="font-bold text-slate-100">{{ t('rule_turn_title') }}</h3>
            </div>
            <ul class="list-disc list-inside text-slate-400 text-xs leading-relaxed pl-7 flex flex-col gap-1">
              <li>{{ t('rule_turn_1') }}</li>
              <li>{{ t('rule_turn_2') }}</li>
              <li>{{ t('rule_turn_3') }}</li>
            </ul>
          </section>

          <section class="flex flex-col gap-1.5 rule-section" style="--i:3">
            <div class="flex items-center gap-2">
              <span class="text-xs font-mono font-black text-indigo-400 bg-indigo-950/80 border border-indigo-500/40 px-1.5 py-0.5 rounded">04</span>
              <h3 class="font-bold text-slate-100">{{ t('rule_liar_title') }}</h3>
            </div>
            <ul class="list-disc list-inside text-slate-400 text-xs leading-relaxed pl-7 flex flex-col gap-1">
              <li>{{ t('rule_liar_1') }}</li>
              <li>{{ t('rule_liar_2') }}</li>
              <li>{{ t('rule_liar_3') }}</li>
            </ul>
          </section>

          <section class="flex flex-col gap-1.5 rule-section" style="--i:4">
            <div class="flex items-center gap-2">
              <span class="text-xs font-mono font-black text-indigo-400 bg-indigo-950/80 border border-indigo-500/40 px-1.5 py-0.5 rounded">05</span>
              <h3 class="font-bold text-slate-100">{{ t('rule_gun_title') }}</h3>
            </div>
            <ul class="list-disc list-inside text-slate-400 text-xs leading-relaxed pl-7 flex flex-col gap-1">
              <li>{{ t('rule_gun_1') }}</li>
              <li>{{ t('rule_gun_2') }}</li>
            </ul>
          </section>
        </div>

        <!-- ── 道具模式 ── -->
        <div class="flex items-center gap-2 mb-3">
          <div class="h-px flex-1 bg-slate-800"></div>
          <span class="text-[10px] font-bold text-amber-600 uppercase tracking-widest px-2">{{ t('rule_section_items') }}</span>
          <div class="h-px flex-1 bg-slate-800"></div>
        </div>

        <div class="flex flex-col gap-4">
          <section class="flex flex-col gap-1.5 rule-section" style="--i:5">
            <div class="flex items-center gap-2">
              <span class="text-xs font-mono font-black text-amber-500 bg-amber-950/80 border border-amber-600/40 px-1.5 py-0.5 rounded">06</span>
              <h3 class="font-bold text-slate-100">{{ t('rule_items_how_title') }}</h3>
            </div>
            <p class="text-slate-400 text-xs leading-relaxed pl-7">{{ t('rule_items_how_desc') }}</p>
          </section>

          <section class="flex flex-col gap-2 rule-section" style="--i:6">
            <div class="flex items-center gap-2">
              <span class="text-xs font-mono font-black text-amber-500 bg-amber-950/80 border border-amber-600/40 px-1.5 py-0.5 rounded">07</span>
              <h3 class="font-bold text-slate-100">{{ t('rule_items_list_title') }}</h3>
            </div>
            <div class="pl-7 flex flex-col gap-1.5">
              <div v-for="item in itemList" :key="item.key"
                class="flex items-start gap-2.5 p-2.5 rounded-lg bg-slate-950/60 border border-slate-800">
                <span class="text-xs font-bold text-amber-300 min-w-[48px] shrink-0">{{ item.name }}</span>
                <span class="text-[11px] text-slate-400 leading-relaxed">{{ item.desc }}</span>
              </div>
            </div>
          </section>
        </div>

      </div>

      <!-- Footer -->
      <div class="p-4 bg-slate-950/70 border-t border-slate-800 shrink-0">
        <button
          class="btn btn-primary w-full font-bold shadow-lg shadow-indigo-600/30 bg-gradient-to-r from-indigo-600 to-indigo-700 hover:from-indigo-500 hover:to-indigo-600 border-none text-white"
          @click="$emit('close')"
        >
          {{ t('rule_got_it') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useI18n } from '../composables/useI18n';

defineEmits(['close']);
const { t } = useI18n();

const itemList = computed(() => [
  { key: 'eagle_eye', name: t('item_eagle_eye_name'), desc: t('rule_item_eagle_eye_detail') },
  { key: 'sawed_off', name: t('item_sawed_off_name'), desc: t('rule_item_sawed_off_detail') },
  { key: 'hard_liquor', name: t('item_hard_liquor_name'), desc: t('rule_item_hard_liquor_detail') },
  { key: 'kevlar_armor', name: t('item_kevlar_armor_name'), desc: t('rule_item_kevlar_armor_detail') },
  { key: 'fate_shift', name: t('item_fate_shift_name'), desc: t('rule_item_fate_shift_detail') },
]);
</script>

<style scoped>
.rules-modal-enter {
  animation: modalIn 0.22s cubic-bezier(0.34, 1.56, 0.64, 1) both;
}

@keyframes modalIn {
  from { opacity: 0; transform: scale(0.95) translateY(8px); }
  to   { opacity: 1; transform: scale(1) translateY(0); }
}

.rule-section {
  animation: sectionIn 0.3s ease both;
  animation-delay: calc(var(--i) * 50ms + 80ms);
}

@keyframes sectionIn {
  from { opacity: 0; transform: translateX(-8px); }
  to   { opacity: 1; transform: translateX(0); }
}
</style>
