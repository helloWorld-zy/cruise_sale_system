<!-- admin/app/components/PricingRow.vue — 价格行组件 -->
<!-- 展示单行价格信息：日期 / 入住人数 / 价格 -->
<script setup lang="ts">
// 组件 props：接收一行价格数据
defineProps<{ row: { date: string; occupancy: number; price: number } }>()
</script>

<template>
  <div>
    <!-- 格式：日期 / 入住人数 / 价格（分） -->
    {{ row.date }} / {{ row.occupancy }} / {{ row.price }}
  </div>
</template>
