<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    type?: 'success' | 'warning' | 'danger' | 'info'
    text: string
  }>(),
  {
    type: 'info',
  },
)

const classes = computed(() => `admin-status-tag admin-status-tag--${props.type}`)
</script>

<template>
  <span :class="classes">{{ text }}</span>
</template>
