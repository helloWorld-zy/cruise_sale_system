<!-- web/app/components/CabinCard.vue — 舱房卡片组件 -->
<!-- 展示舱房简要信息：名称和价格 -->
<script setup lang="ts">
// 组件 props：接收舱房数据
defineProps<{ cabin: { id: number; name: string; price: number } }>()
</script>

<template>
  <!-- 舱房卡片：显示名称和价格 -->
  <div>{{ cabin.name }} - {{ cabin.price }}</div>
</template>
