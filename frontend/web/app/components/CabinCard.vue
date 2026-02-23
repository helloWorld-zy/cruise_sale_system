<script setup lang="ts">
defineProps<{ cabin: { id: number; name: string; price: number } }>()
</script>

<template>
  <div>{{ cabin.name }} - {{ cabin.price }}</div>
</template>
