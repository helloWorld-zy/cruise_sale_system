<!-- admin/app/pages/voyages/new.vue — 新建航次页面 -->
<!-- 提供航次编码和名称的输入表单 -->
<script setup lang="ts">
import { ref } from 'vue'
import { useApi } from '../../composables/useApi'
// 表单数据模型（与后端 domain.Voyage 字段对齐）
const form = ref({ code: '', route_id: 0, cruise_id: 1, depart_date: '', return_date: '' })
const loading = ref(false)
const error = ref<string | null>(null)
const { request } = useApi()

function toRFC3339(v: string) {
  if (!v) return v
  return v.length === 16 ? `${v}:00Z` : v
}

// handleSubmit 提交新建航次请求并在成功后跳转列表页。
async function handleSubmit() {
  if (loading.value) return
  loading.value = true
  error.value = null
  try {
    await request('/voyages', {
      method: 'POST',
      body: {
        code: form.value.code,
        route_id: Number(form.value.route_id),
        cruise_id: Number(form.value.cruise_id),
        depart_date: toRFC3339(form.value.depart_date),
        return_date: toRFC3339(form.value.return_date),
      },
    })
    await navigateTo('/voyages')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to create voyage'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="page">
    <h1>New Voyage</h1>
    <form @submit.prevent="handleSubmit">
      <!-- 航次编码输入 -->
      <input v-model="form.code" placeholder="Code" :disabled="loading" />
      <input v-model.number="form.route_id" type="number" min="1" placeholder="Route ID" :disabled="loading" />
      <input v-model.number="form.cruise_id" type="number" min="1" placeholder="Cruise ID" :disabled="loading" />
      <input v-model="form.depart_date" type="datetime-local" :disabled="loading" />
      <input v-model="form.return_date" type="datetime-local" :disabled="loading" />
      <p v-if="error" class="text-red-500">{{ error }}</p>
      <button type="button" :disabled="loading" @click="handleSubmit">{{ loading ? 'Submitting...' : 'Submit' }}</button>
    </form>
  </div>
</template>

