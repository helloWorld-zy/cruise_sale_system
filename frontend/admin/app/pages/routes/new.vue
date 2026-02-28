<!-- admin/app/pages/routes/new.vue — 新建航线页面 -->
<!-- 提供航线编码和名称的输入表单 -->
<script setup lang="ts">
import { ref } from 'vue'
import { useApi } from '../../composables/useApi'
// 表单数据模型
const form = ref({ code: '', name: '' })
const loading = ref(false)
const error = ref<string | null>(null)
const { request } = useApi()

// handleSubmit 提交新建航线请求并在成功后跳转列表页。
async function handleSubmit() {
  if (loading.value) return
  loading.value = true
  error.value = null
  try {
    await request('/routes', {
      method: 'POST',
      body: {
        code: form.value.code,
        name: form.value.name,
      },
    })
    await navigateTo('/routes')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to create route'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="page">
    <h1>New Route</h1>
    <form @submit.prevent="handleSubmit">
      <!-- 航线编码输入 -->
      <input v-model="form.code" placeholder="Code" :disabled="loading" />
      <!-- 航线名称输入 -->
      <input v-model="form.name" placeholder="Name" :disabled="loading" />
      <p v-if="error" class="text-red-500">{{ error }}</p>
      <button type="button" :disabled="loading" @click="handleSubmit">{{ loading ? 'Submitting...' : 'Submit' }}</button>
    </form>
  </div>
</template>

