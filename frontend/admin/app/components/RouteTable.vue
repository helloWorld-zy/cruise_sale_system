<!-- admin/app/components/RouteTable.vue — 航线数据表格组件 -->
<!-- 可复用的表格组件，展示航线列表（编码和名称列） -->
<script setup lang="ts">
// 组件 props：接收航线数据列表
defineProps<{ items: { id: number; code: string; name: string }[] }>()
</script>

<template>
  <table class="route-table">
    <thead>
      <tr>
        <th>ID</th>
        <th>编码</th>
        <th>名称</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="r in items" :key="r.id">
        <td>{{ r.id }}</td>
        <td>{{ r.code }}</td>
        <td>{{ r.name }}</td>
      </tr>
    </tbody>
  </table>
</template>
