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
    <strong>#{{ booking.id }}</strong>
    <span>状态：{{ booking.status }}</span>
    <span>总额：{{ booking.total }}</span>
  </div>
</template>

<style scoped>
.booking-row {
  display: flex;
  gap: 16px;
  align-items: center;
}

span {
  color: #38506f;
}
</style>
