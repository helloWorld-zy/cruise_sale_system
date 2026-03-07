<script setup lang="ts">
defineProps<{
  flush?: boolean
}>()
</script>

<template>
  <section class="admin-data-card" :class="{ 'admin-data-card--flush': flush }">
    <slot />
  </section>
</template>
