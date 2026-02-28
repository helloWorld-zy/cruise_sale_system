<!-- admin/app/pages/cruises/create.vue -->
<script setup lang="ts">
import { ref } from 'vue'
import { useApi } from '../../composables/useApi'
const { request } = useApi()
const loading = ref(false)
const error = ref<string | null>(null)
const form = ref({
  name: '',
  company_id: 0,
  status: 'draft',
})

async function handleSubmit() {
  if (loading.value) return
  loading.value = true
  error.value = null
  try {
    await request('/cruises', {
      method: 'POST',
      body: {
        name: form.value.name,
        company_id: Number(form.value.company_id),
        status: form.value.status,
      },
    })
    await navigateTo('/cruises')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to create cruise'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="page">
    <h1>新建邮轮</h1>
    <form @submit.prevent="handleSubmit">
      <input v-model="form.name" placeholder="邮轮名称" :disabled="loading" />
      <input v-model.number="form.company_id" type="number" min="1" placeholder="公司 ID" :disabled="loading" />
      <input v-model="form.status" placeholder="状态（draft/online/offline）" :disabled="loading" />
      <p v-if="error" class="error">{{ error }}</p>
      <button type="button" :disabled="loading" @click="handleSubmit">{{ loading ? '提交中...' : '创建' }}</button>
    </form>
  </div>
</template>

