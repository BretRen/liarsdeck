<template>
  <div class="fixed inset-0 z-[2100] flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-md animate-in fade-in duration-200" @click.self="$emit('close')">
    <div class="card w-full max-w-3xl bg-slate-900/95 border border-slate-700/80 rounded-2xl shadow-2xl shadow-black/80 flex flex-col max-h-[88vh] overflow-hidden">
      <!-- Header -->
      <div class="px-5 py-4 border-b border-slate-800 flex items-center justify-between bg-slate-950/50">
        <div class="flex items-center gap-2.5">
          <div class="w-8 h-8 rounded-lg bg-indigo-600/20 border border-indigo-500/30 flex items-center justify-center text-indigo-400 text-base">
            📜
          </div>
          <div>
            <h2 class="text-base md:text-lg font-bold font-serif text-slate-100 tracking-wide">{{ t('changelog_title') }}</h2>
            <p class="text-[11px] text-slate-400">{{ t('changelog_subtitle') }}</p>
          </div>
        </div>
        <button class="btn btn-ghost btn-xs w-7 h-7 rounded-full text-slate-400 hover:text-white p-0" @click="$emit('close')">✕</button>
      </div>

      <!-- Main Layout: Sidebar Version List + Main Markdown Viewer -->
      <div class="flex flex-col md:flex-row flex-1 overflow-hidden">
        <!-- Sidebar Version List -->
        <div class="w-full md:w-48 bg-slate-950/60 border-b md:border-b-0 md:border-r border-slate-800 p-2.5 flex md:flex-col gap-1.5 overflow-x-auto md:overflow-y-auto shrink-0">
          <div class="hidden md:block px-2 py-1 text-[10px] font-bold text-slate-500 uppercase tracking-wider">
            {{ t('changelog_versions') }}
          </div>
          <button
            v-for="ver in versions"
            :key="ver"
            class="btn btn-xs justify-between font-mono rounded-lg transition-all"
            :class="selectedVersion === ver ? 'btn-primary text-white font-bold shadow-md shadow-indigo-600/30' : 'btn-ghost text-slate-400 hover:text-slate-200 hover:bg-slate-800/60'"
            @click="selectVersion(ver)"
          >
            <span>{{ ver }}</span>
            <span v-if="ver === latestVersion" class="badge badge-xs bg-indigo-500/20 border-indigo-500/40 text-indigo-300 text-[8px]">LATEST</span>
          </button>
          <div v-if="versions.length === 0 && !loading" class="text-center py-4 text-xs text-slate-500">
            暂无版本记录
          </div>
        </div>

        <!-- Markdown Content Viewer -->
        <div class="flex-1 p-5 overflow-y-auto bg-slate-900/40 flex flex-col justify-between">
          <div v-if="loading" class="flex flex-col items-center justify-center py-16 gap-3 text-slate-400">
            <span class="loading loading-spinner loading-md text-indigo-500"></span>
            <span class="text-xs">{{ t('changelog_loading') }}</span>
          </div>

          <div v-else-if="error" class="alert alert-error bg-rose-950/40 border-rose-800 text-rose-300 text-xs my-4">
            <span>{{ error }}</span>
          </div>

          <div v-else class="markdown-body select-text" v-html="renderedHtml"></div>

          <!-- Footer Metadata -->
          <div class="mt-6 pt-3 border-t border-slate-800/80 flex items-center justify-between text-[11px] text-slate-500 font-mono">
            <span>Version: {{ selectedVersion }}</span>
            <span>Liar's Deck Core Engine</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { marked } from 'marked';
import { useI18n } from '../composables/useI18n';

defineEmits(['close']);
const { t } = useI18n();

const versions = ref([]);
const selectedVersion = ref('');
const markdownContent = ref('');
const loading = ref(false);
const error = ref('');

const latestVersion = computed(() => (versions.value.length > 0 ? versions.value[0] : ''));

const renderedHtml = computed(() => {
  if (!markdownContent.value) return '';
  return marked.parse(markdownContent.value);
});

async function fetchVersions() {
  loading.value = true;
  error.value = '';
  try {
    const res = await fetch('/api/changelogs');
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    const data = await res.json();
    if (data.success && Array.isArray(data.versions)) {
      versions.value = data.versions;
      if (versions.value.length > 0) {
        selectVersion(versions.value[0]);
      }
    }
  } catch (err) {
    error.value = '无法加载更新日志列表: ' + (err.message || '未知错误');
  } finally {
    loading.value = false;
  }
}

async function selectVersion(ver) {
  selectedVersion.value = ver;
  loading.value = true;
  error.value = '';
  try {
    const res = await fetch(`/api/changelogs/${encodeURIComponent(ver)}`);
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    const data = await res.json();
    if (data.success && data.content) {
      markdownContent.value = data.content;
    } else {
      markdownContent.value = '该版本暂无详细更新记录。';
    }
  } catch (err) {
    error.value = `加载 ${ver} 更新日志失败: ` + (err.message || '未知错误');
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  fetchVersions();
});
</script>
