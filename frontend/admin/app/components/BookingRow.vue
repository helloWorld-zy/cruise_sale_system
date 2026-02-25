<!-- admin/app/components/BookingRow.vue — 订单行组件 -->
<!-- 渲染单行订单摘要信息：ID / 状态 / 金额 -->
<script setup lang="ts">
// 组件 props：接收一行订单数据
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
