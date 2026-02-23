<!-- admin/app/pages/cabins/pricing.vue — 定价管理页面 -->
<!-- 展示价格日历矩阵，每行对应一个日期/入住人数的价格 -->
<script setup lang="ts">
import PricingRow from '../../components/PricingRow.vue'
// 示例价格数据（待对接 API）
const rows = [{ date: '2026-05-01', occupancy: 2, price: 19900 }]
</script>

<template>
  <div class="page">
    <h1>Pricing Matrix</h1>
    <!-- 渲染每行价格数据 -->
    <PricingRow v-for="r in rows" :key="r.date" :row="r" />
  </div>
</template>
