<script setup lang="ts">
defineProps<{ title: string; value: number | string }>()
</script>

<template>
  <div class="stat-card">
    <h3>{{ title }}</h3>
    <div>{{ value }}</div>
  </div>
</template>
