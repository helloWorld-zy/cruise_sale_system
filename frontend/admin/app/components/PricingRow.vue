<script setup lang="ts">
defineProps<{ row: { date: string; occupancy: number; price: number } }>()
</script>

<template>
  <div>
    {{ row.date }} / {{ row.occupancy }} / {{ row.price }}
  </div>
</template>
