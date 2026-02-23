<!-- admin/app/pages/routes/new.vue — 新建航线页面 -->
<!-- 提供航线编码和名称的输入表单 -->
<script setup lang="ts">
import { ref } from 'vue'
// 表单数据模型
const form = ref({ code: '', name: '' })
</script>

<template>
  <div class="page">
    <h1>New Route</h1>
    <!-- 航线编码输入 -->
    <input v-model="form.code" placeholder="Code" />
    <!-- 航线名称输入 -->
    <input v-model="form.name" placeholder="Name" />
  </div>
</template>
