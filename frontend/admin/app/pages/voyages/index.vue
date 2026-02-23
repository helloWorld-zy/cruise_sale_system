<!-- admin/app/pages/voyages/index.vue — 航次列表页面 -->
<!-- 展示航次列表数据 -->
<script setup lang="ts">
// 示例数据（待对接 API）
const items = [{ id: 1, code: 'V1', name: 'Voyage 1' }]
</script>

<template>
  <div class="page">
    <h1>Voyages</h1>
    <!-- 航次数据表格 -->
    <table>
      <tbody>
        <tr v-for="r in items" :key="r.id">
          <td>{{ r.code }}</td>
          <td>{{ r.name }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
