<!-- admin/components/StatCard.vue — 统计卡片组件 -->
<!-- 展示单个统计指标（标题 + 数值） -->
<script setup lang="ts">
// 组件 props：接收统计指标的标题和数值
defineProps<{ title: string; value: number | string }>()
</script>

<template>
  <div class="stat-card">
    <h3>{{ title }}</h3>
    <div>{{ value }}</div>
  </div>
</template>
