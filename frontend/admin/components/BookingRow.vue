<script setup lang="ts">
defineProps<{
  booking: {
    id: number
    status: string
    total: number
  }
}>()
</script>

<template>
  <div class="booking-row">
    {{ booking.id }} - {{ booking.status }} - {{ booking.total }}
  </div>
</template>
