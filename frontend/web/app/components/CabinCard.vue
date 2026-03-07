<!-- web/app/components/CabinCard.vue — 舱房卡片组件 -->
<!-- 展示舱房简要信息：名称和价格 -->
<script setup lang="ts">
// 组件 props：接收舱房数据
defineProps<{ cabin: { id: number; name: string; price: number } }>()
</script>

<template>
  <!-- 舱房卡片：显示名称和价格 -->
  <article class="card">
    <h3 class="heading-font text-[var(--color-primary)]">{{ cabin.name }}</h3>
    <p class="text-sm text-[var(--web-muted)]">舱房编号：#{{ cabin.id }}</p>
    <strong class="text-[var(--color-cta)] heading-font text-xl">￥{{ (cabin.price / 100).toFixed(2) }}</strong>
  </article>
</template>

<style scoped>
/* Scoped styles can be minimal since we use Tailwind / global CSS vars on card */
h3 {
  margin: 0 0 8px;
  font-size: 1.25rem;
}

p {
  margin: 0 0 8px;
}
</style>
