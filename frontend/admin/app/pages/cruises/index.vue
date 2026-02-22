<template>
  <div>
    <div class="flex items-center justify-between mb-4">
      <h1 class="text-xl font-semibold">邮轮管理</h1>
      <UButton to="/cruises/create" color="primary">新建邮轮</UButton>
    </div>
    <table class="w-full border">
      <thead>
        <tr class="bg-gray-50">
          <th class="text-left p-2">名称</th>
          <th class="text-left p-2">状态</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td class="p-2">示例邮轮</td>
          <td class="p-2">上架</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
</script>
