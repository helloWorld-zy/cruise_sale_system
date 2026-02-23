<!-- admin/app/pages/routes/index.vue — 航线列表页面 -->
<!-- 展示所有航线数据表格 -->
<script setup lang="ts">
import RouteTable from '../../components/RouteTable.vue'
// 示例数据（待对接 API）
const items = [{ id: 1, code: 'R1', name: 'Route 1' }]
</script>

<template>
  <div class="page">
    <h1>Routes</h1>
    <!-- 航线数据表格组件 -->
    <RouteTable :items="items" />
  </div>
</template>
