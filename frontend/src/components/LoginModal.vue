<template>
  <div v-if="!isAuthenticated" class="fixed inset-0 z-[99999] flex items-center justify-center p-4 bg-slate-950/90 backdrop-blur-xl">
    <div class="card w-full max-w-md bg-slate-900/90 border border-slate-700/70 shadow-2xl shadow-black/80 rounded-2xl overflow-hidden animate-in fade-in zoom-in duration-300">
      <!-- Top Bar with Language Selector -->
      <div class="flex items-center justify-between px-5 py-3.5 bg-slate-950/70 border-b border-slate-800">
        <span class="text-[11px] font-extrabold tracking-widest text-indigo-400 uppercase">AUTHENTICATION</span>
        <button class="btn btn-ghost btn-xs text-slate-300 hover:bg-slate-800" @click="toggleLang">
          {{ lang.toUpperCase() }}
        </button>
      </div>

      <!-- Main Login Container -->
      <div class="card-body p-8 flex flex-col items-center text-center">
        <div class="relative w-20 h-20 rounded-2xl bg-gradient-to-br from-indigo-500/20 to-slate-800 border border-indigo-500/30 flex items-center justify-center mb-5 shadow-lg shadow-indigo-500/10">
          <span class="text-4xl filter drop-shadow">🃏</span>
        </div>

        <h2 class="text-2xl font-bold font-serif text-slate-100 tracking-wide mb-2">{{ t('login_modal_title') }}</h2>
        <p class="text-sm text-slate-400 leading-relaxed max-w-xs mb-6">{{ t('login_modal_desc') }}</p>

        <div v-if="authError" class="alert alert-error text-xs py-2 px-3 mb-5 w-full text-left rounded-lg bg-rose-950/60 border border-rose-600/50 text-rose-200">
          <span>⚠️ {{ authError }}</span>
        </div>

        <!-- Single Action Button (DaisyUI btn-primary / indigo glow) -->
        <button
          class="btn btn-primary w-full h-12 text-base font-bold tracking-wider uppercase shadow-lg shadow-indigo-600/30 bg-gradient-to-r from-indigo-600 to-indigo-700 hover:from-indigo-500 hover:to-indigo-600 border-none text-white transition-all duration-200 active:scale-[0.98]"
          :disabled="isLoggingIn"
          @click="login"
        >
          <span v-if="isLoggingIn" class="loading loading-spinner loading-sm"></span>
          <span>{{ isLoggingIn ? t('login_loading_btn') : t('login_btn') }}</span>
        </button>
      </div>

      <!-- Footer Note -->
      <div class="px-5 py-3 bg-slate-950/70 border-t border-slate-800/80 text-center">
        <span class="text-xs text-slate-400">
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
