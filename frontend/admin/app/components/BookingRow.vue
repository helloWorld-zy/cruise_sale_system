<script setup lang="ts">
// BookingRow 用于渲染订单行的核心摘要信息。
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
