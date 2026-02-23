<!-- admin/app/pages/cruises/index.vue — 邮轮管理列表页面 -->
<!-- 展示邮轮列表，提供新建邮轮按钮入口 -->
<template>
  <div>
    <!-- 页面标题栏 -->
    <div class="flex items-center justify-between mb-4">
      <h1 class="text-xl font-semibold">邮轮管理</h1>
      <UButton to="/cruises/create" color="primary">新建邮轮</UButton>
    </div>
    <!-- 邮轮数据表格 -->
    <table class="w-full border">
      <thead>
        <tr class="bg-gray-50">
          <th class="text-left p-2">名称</th>
          <th class="text-left p-2">状态</th>
        </tr>
      </thead>
      <tbody>
        <!-- TODO: 对接 API 加载真实数据 -->
        <tr>
          <td class="p-2">示例邮轮</td>
          <td class="p-2">上架</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
/* 邮轮列表页面：待对接后端 API */
</script>
