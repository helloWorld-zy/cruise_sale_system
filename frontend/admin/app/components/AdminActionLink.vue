<script setup lang="ts">
import { computed } from 'vue'

type ActionVariant = 'primary' | 'secondary' | 'danger'
type ActionSize = 'md' | 'sm'

const props = withDefaults(
  defineProps<{
    to: string
    variant?: ActionVariant
    size?: ActionSize
  }>(),
  {
    variant: 'secondary',
    size: 'sm',
  },
)

const linkClass = computed(() => [
  'admin-btn',
  props.variant === 'danger' ? 'admin-btn--danger' : '',
  props.size === 'sm' ? 'admin-btn--sm' : '',
].join(' '))
</script>

<template>
  <NuxtLink :to="to" :class="linkClass">
    <slot />
  </NuxtLink>
</template>
