<template>
  <div class="fixed inset-0 z-[2000] flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-md animate-in fade-in duration-200" @click.self="$emit('close')">
    <div class="card w-full max-w-lg bg-slate-900/95 border border-slate-700/80 rounded-2xl shadow-2xl shadow-black/80 flex flex-col max-h-[85vh] overflow-hidden">
      <div class="flex items-center justify-between px-6 py-4 border-b border-slate-800 bg-slate-950/70">
        <h2 class="text-xl font-bold font-serif text-slate-100 tracking-wide">{{ t('rules_title') }}</h2>
        <button class="btn btn-ghost btn-xs w-7 h-7 rounded-full text-slate-400 hover:text-white p-0" @click="$emit('close')">✕</button>
      </div>

      <div class="p-6 overflow-y-auto flex flex-col gap-5 text-left text-sm text-slate-300">
        <section class="flex flex-col gap-1.5">
          <div class="flex items-center gap-2">
            <span class="text-xs font-mono font-black text-indigo-400 bg-indigo-950/80 border border-indigo-500/40 px-1.5 py-0.5 rounded">01</span>
            <h3 class="font-bold text-slate-100">{{ t('rule_goal_title') }}</h3>
          </div>
          <p class="text-slate-400 text-xs leading-relaxed pl-7">{{ t('rule_goal_desc') }}</p>
        </section>

        <section class="flex flex-col gap-1.5">
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

        <section class="flex flex-col gap-1.5">
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

        <section class="flex flex-col gap-1.5">
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

        <section class="flex flex-col gap-1.5">
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

      <div class="p-4 bg-slate-950/70 border-t border-slate-800">
        <button class="btn btn-primary w-full font-bold shadow-lg shadow-indigo-600/30 bg-gradient-to-r from-indigo-600 to-indigo-700 hover:from-indigo-500 hover:to-indigo-600 border-none text-white" @click="$emit('close')">
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
