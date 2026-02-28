<!-- web/app/components/CabinCard.vue — 舱房卡片组件 -->
<!-- 展示舱房简要信息：名称和价格 -->
<script setup lang="ts">
// 组件 props：接收舱房数据
defineProps<{ cabin: { id: number; name: string; price: number } }>()
</script>

<template>
  <!-- 舱房卡片：显示名称和价格 -->
  <article class="cabin-card">
    <h3>{{ cabin.name }}</h3>
    <p>舱房编号：#{{ cabin.id }}</p>
    <strong>￥{{ (cabin.price / 100).toFixed(2) }}</strong>
  </article>
</template>

<style scoped>
.cabin-card {
  border: 1px solid #f1d9c9;
  border-radius: 14px;
  padding: 14px;
  margin-bottom: 12px;
  background: #fff;
}

h3 {
  margin: 0 0 8px;
  font-size: 18px;
}

p {
  margin: 0 0 8px;
}
</style>
