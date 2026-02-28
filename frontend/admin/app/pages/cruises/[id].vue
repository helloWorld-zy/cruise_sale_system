<!-- admin/app/pages/cruises/[id].vue -->
<script setup lang="ts">
import { onMounted, ref } from 'vue'

declare const useApi: any
declare const navigateTo: any

const route = useRoute()
const { request } = useApi()

const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const error = ref<string | null>(null)
const form = ref({
  id: 0,
  name: '',
  company_id: 0,
  status: '',
})

const id = Number(route.params.id)

async function loadDetail() {
  loading.value = true
  error.value = null
  try {
    const res = await request(`/cruises/${id}`)
    const data = res?.data ?? res ?? {}
    form.value = {
      id: Number(data.id ?? id),
      name: data.name ?? '',
      company_id: Number(data.company_id ?? 0),
      status: data.status ?? 'draft',
    }
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load cruise detail'
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  if (saving.value) return
  saving.value = true
  error.value = null
  try {
    await request(`/cruises/${id}`, {
      method: 'PUT',
      body: {
        name: form.value.name,
        company_id: Number(form.value.company_id),
        status: form.value.status,
      },
    })
    await navigateTo('/cruises')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to update cruise'
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  if (deleting.value) return
  if (!confirm(`确认删除邮轮 #${id} 吗？`)) return
  deleting.value = true
  error.value = null
  try {
    await request(`/cruises/${id}`, { method: 'DELETE' })
    await navigateTo('/cruises')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to delete cruise'
  } finally {
    deleting.value = false
  }
}

onMounted(loadDetail)
</script>

<template>
  <div class="page">
    <h1>编辑邮轮 #{{ id }}</h1>
    <p v-if="loading">Loading...</p>
    <form v-else @submit.prevent="handleSave">
      <input v-model="form.name" placeholder="邮轮名称" :disabled="saving || deleting" />
      <input v-model.number="form.company_id" type="number" min="1" placeholder="公司 ID" :disabled="saving || deleting" />
      <input v-model="form.status" placeholder="状态" :disabled="saving || deleting" />
      <p v-if="error" class="error">{{ error }}</p>
      <div>
        <button type="submit" :disabled="saving || deleting">{{ saving ? '保存中...' : '保存' }}</button>
        <button type="button" style="margin-left:8px" :disabled="saving || deleting" @click="handleDelete">{{ deleting ? '删除中...' : '删除' }}</button>
      </div>
    </form>
  </div>
</template>
