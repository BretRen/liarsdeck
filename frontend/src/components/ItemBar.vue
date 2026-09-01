<template>
  <div class="p-3 mb-3.5 bg-slate-900/80 border border-slate-700/60 rounded-2xl shadow-xl shadow-black/50 backdrop-blur-xl w-full">
    <!-- 顶栏小标签 -->
    <div class="flex items-center justify-between mb-2 px-0.5 text-[11px] font-bold text-slate-400">
      <div class="flex items-center gap-1.5">
        <span class="tracking-wider text-slate-300">{{ t('items_inventory_label') }}</span>
        <span class="badge badge-xs bg-slate-800 text-slate-300 border-slate-700 font-mono">
          {{ items.length }}/1
        </span>
      </div>
      <span v-if="items.length === 0" class="text-[10px] text-slate-500">
        {{ t('no_items_tip') }}
      </span>
      <span v-else-if="!isMyTurn" class="text-[10px] text-slate-500">
        {{ t('item_disabled_not_turn') }}
      </span>
    </div>

    <!-- 1 格道具槽位 -->
    <div class="w-full">
      <div
        class="h-14 rounded-xl border transition-all duration-200 relative flex items-center p-3"
        :class="getSlotClass(items[0])"
      >
        <!-- 槽位有道具 -->
        <template v-if="items[0]">
          <div class="flex items-center justify-between w-full gap-3">
            <div class="flex items-center gap-2.5 overflow-hidden">
              <div class="w-8 h-8 rounded-lg bg-amber-500/10 border border-amber-500/30 flex items-center justify-center text-amber-400 shrink-0">
                <ItemIcon :item="items[0]" class="w-4 h-4" />
              </div>
              <div class="flex flex-col text-left overflow-hidden">
                <span class="font-bold text-xs text-slate-100 truncate">
                  {{ getItemName(items[0]) }}
                </span>
                <span class="text-[10.5px] text-slate-400 truncate max-w-[200px]" :title="getItemDesc(items[0])">
                  {{ getItemDesc(items[0]) }}
                </span>
              </div>
            </div>

            <!-- 使用按钮 -->
            <button
              type="button"
              class="btn btn-sm shrink-0 font-bold px-4 transition-all"
              :class="getButtonClass(items[0])"
              :disabled="isItemDisabled(items[0])"
              @click="onUse(items[0])"
            >
              {{ t('item_use_btn') }}
            </button>
          </div>
        </template>

        <!-- 空槽位 -->
        <template v-else>
          <div class="flex items-center justify-center w-full text-slate-600 text-xs font-mono py-1">
            <span class="text-[11px] font-semibold tracking-wide">{{ t('item_slot_empty') }}</span>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from '../composables/useI18n';
import ItemIcon from './ItemIcon.vue';

const props = defineProps({
  items: {
    type: Array,
    default: () => [],
  },
  isMyTurn: {
    type: Boolean,
    default: false,
  },
  tableHasCards: {
    type: Boolean,
    default: false,
  },
  hasArmor: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['use-item']);
const { t } = useI18n();

function getItemName(item) {
  const map = {
    eagle_eye: t('item_eagle_eye_name'),
    sawed_off: t('item_sawed_off_name'),
    hard_liquor: t('item_hard_liquor_name'),
    kevlar_armor: t('item_kevlar_armor_name'),
    fate_shift: t('item_fate_shift_name'),
  };
  return map[item] || item;
}

function getItemDesc(item) {
  if (item === 'kevlar_armor' && props.hasArmor) {
    return t('armor_already_equipped') || t('armor_equipped_tag');
  }
  const map = {
    eagle_eye: t('item_eagle_eye_desc'),
    sawed_off: t('item_sawed_off_desc'),
    hard_liquor: t('item_hard_liquor_desc'),
    kevlar_armor: t('item_kevlar_armor_desc'),
    fate_shift: t('item_fate_shift_desc'),
  };
  return map[item] || '';
}

function isItemDisabled(item) {
  if (!props.isMyTurn) return true;
  if (item === 'eagle_eye' && !props.tableHasCards) return true;
  if (item === 'kevlar_armor' && props.hasArmor) return true;
  return false;
}

function getSlotClass(item) {
  if (!item) {
    return 'bg-slate-950/40 border-dashed border-slate-800/80';
  }
  if (!props.isMyTurn) {
    return 'bg-slate-900/60 border-slate-800 opacity-80';
  }
  return 'bg-slate-900/90 border-indigo-700/50 shadow-sm shadow-indigo-950/30';
}

function getButtonClass(item) {
  if (isItemDisabled(item)) {
    return 'btn-disabled bg-slate-800 border-slate-700 text-slate-500 cursor-not-allowed';
  }
  return 'btn-primary bg-indigo-600 hover:bg-indigo-500 border-none text-white shadow-sm active:scale-95';
}

function onUse(item) {
  if (isItemDisabled(item)) return;
  emit('use-item', item);
}
</script>
