<!-- admin/app/pages/voyages/new.vue — 新建航次页面 -->
<!-- 提供航次编码和名称的输入表单 -->
<script setup lang="ts">
import { ref } from 'vue'
// 表单数据模型
const form = ref({ code: '', name: '' })
</script>

<template>
  <div class="page">
    <h1>New Voyage</h1>
    <!-- 航次编码输入 -->
    <input v-model="form.code" placeholder="Code" />
    <!-- 航次名称输入 -->
    <input v-model="form.name" placeholder="Name" />
  </div>
</template>
