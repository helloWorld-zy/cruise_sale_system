<!-- web/app/pages/search/index.vue — 舱房搜索页面 -->
<!-- 提供搜索筛选条件和搜索结果列表 -->
<script setup lang="ts">
import CabinCard from '../../components/CabinCard.vue'
// 示例搜索结果（待对接 API）
const results = [{ id: 1, name: 'Balcony Cabin', price: 19900 }]
</script>

<template>
  <div class="page">
    <h1>Search</h1>
    <!-- 搜索筛选条件区域 -->
    <div>Filter</div>
    <!-- 搜索结果列表 -->
    <CabinCard v-for="r in results" :key="r.id" :cabin="r" />
  </div>
</template>
