<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  available: number
  total: number
}>()

const percent = computed(() => {
  if (props.total <= 0) return 0
  return Math.max(0, Math.min(100, Math.round((props.available / props.total) * 100)))
})

const label = computed(() => {
  if (props.available <= 0) return '已售罄'
  if (percent.value < 20) return '即将售罄'
  if (percent.value < 50) return '库存紧张'
  return '库存充足'
})

const badgeClass = computed(() => {
  if (props.available <= 0) return 'badge soldout'
  if (percent.value < 20) return 'badge danger'
  if (percent.value < 50) return 'badge warning'
  return 'badge ok'
})
</script>

<template>
  <span :class="badgeClass">{{ label }}</span>
</template>

<style scoped>
.badge {
  border-radius: 9999px;
  padding: 2px 8px;
  font-size: 12px;
  font-weight: 600;
}

.ok {
  background: #ecfdf5;
  color: #047857;
}

.warning {
  background: #fffbeb;
  color: #b45309;
}

.danger {
  background: #fef2f2;
  color: #b91c1c;
}

.soldout {
  background: #f1f5f9;
  color: #475569;
}
</style>
