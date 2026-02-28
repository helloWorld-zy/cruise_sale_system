<!-- admin/app/pages/routes/[id].vue -->
<script setup lang="ts">
import { onMounted, ref } from 'vue'

declare const useApi: any
declare const navigateTo: any

const route = useRoute()
const { request } = useApi()
const id = Number(route.params.id)

const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const error = ref<string | null>(null)
const form = ref({ code: '', name: '' })

async function loadDetail() {
  loading.value = true
  error.value = null
  try {
    const res = await request(`/routes/${id}`)
    const data = res?.data ?? res ?? {}
    form.value = { code: data.code ?? '', name: data.name ?? '' }
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load route detail'
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  if (saving.value) return
  saving.value = true
  error.value = null
  try {
    await request(`/routes/${id}`, {
      method: 'PUT',
      body: { code: form.value.code, name: form.value.name },
    })
    await navigateTo('/routes')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to update route'
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  if (deleting.value) return
  if (!confirm(`确认删除航线 #${id} 吗？`)) return
  deleting.value = true
  error.value = null
  try {
    await request(`/routes/${id}`, { method: 'DELETE' })
    await navigateTo('/routes')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to delete route'
  } finally {
    deleting.value = false
  }
}

onMounted(loadDetail)
</script>

<template>
  <div class="page">
    <h1>编辑航线 #{{ id }}</h1>
    <p v-if="loading">Loading...</p>
    <form v-else @submit.prevent="handleSave">
      <input v-model="form.code" placeholder="Code" :disabled="saving || deleting" />
      <input v-model="form.name" placeholder="Name" :disabled="saving || deleting" />
      <p v-if="error" class="error">{{ error }}</p>
      <div>
        <button type="submit" :disabled="saving || deleting">{{ saving ? '保存中...' : '保存' }}</button>
        <button type="button" style="margin-left:8px" :disabled="saving || deleting" @click="handleDelete">{{ deleting ? '删除中...' : '删除' }}</button>
      </div>
    </form>
  </div>
</template>
