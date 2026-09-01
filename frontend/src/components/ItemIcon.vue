<template>
  <!-- Eagle Eye: Magnifier -->
  <svg
    v-if="item === 'eagle_eye'"
    class="w-4 h-4"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="2"
    stroke-linecap="round"
    stroke-linejoin="round"
  >
    <circle cx="11" cy="11" r="8" />
    <line x1="21" y1="21" x2="16.65" y2="16.65" />
  </svg>

  <!-- Sawed Off: Lightning Bolt / Double Trigger -->
  <svg
    v-else-if="item === 'sawed_off'"
    class="w-4 h-4"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="2"
    stroke-linecap="round"
    stroke-linejoin="round"
  >
    <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2" />
  </svg>

  <!-- Beer / Hard Liquor: Refresh Cycle -->
  <svg
    v-else-if="item === 'hard_liquor'"
    class="w-4 h-4"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="2"
    stroke-linecap="round"
    stroke-linejoin="round"
  >
    <path d="M21.5 2v6h-6M21.34 15.57a10 10 0 1 1-.57-8.38l5.67-5.67" />
  </svg>

  <!-- Kevlar Armor: Shield -->
  <svg
    v-else-if="item === 'kevlar_armor'"
    class="w-4 h-4"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="2"
    stroke-linecap="round"
    stroke-linejoin="round"
  >
    <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
  </svg>

  <!-- Fate Shift: Dice -->
  <svg
    v-else-if="item === 'fate_shift'"
    class="w-4 h-4"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="2"
    stroke-linecap="round"
    stroke-linejoin="round"
  >
    <rect x="3" y="3" width="18" height="18" rx="3" ry="3" />
    <circle cx="8.5" cy="8.5" r="1.2" fill="currentColor" />
    <circle cx="15.5" cy="15.5" r="1.2" fill="currentColor" />
    <circle cx="12" cy="12" r="1.2" fill="currentColor" />
  </svg>

  <!-- Fallback Box Icon -->
  <svg
    v-else
    class="w-4 h-4"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="2"
    stroke-linecap="round"
    stroke-linejoin="round"
  >
    <rect x="3" y="3" width="18" height="18" rx="2" />
  </svg>
</template>

<script setup>
defineProps({
  item: { type: String, default: '' },
});
</script>
