<!-- admin/app/pages/voyages/[id].vue -->
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
const form = ref({
  code: '',
  name: '',
  route_id: 0,
  departure_date: '',
  arrival_date: '',
})

async function loadDetail() {
  loading.value = true
  error.value = null
  try {
    const res = await request(`/voyages/${id}`)
    const data = res?.data ?? res ?? {}
    form.value = {
      code: data.code ?? '',
      name: data.name ?? '',
      route_id: Number(data.route_id ?? 0),
      departure_date: data.departure_date?.slice?.(0, 10) ?? '',
      arrival_date: data.arrival_date?.slice?.(0, 10) ?? '',
    }
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load voyage detail'
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  if (saving.value) return
  saving.value = true
  error.value = null
  try {
    await request(`/voyages/${id}`, {
      method: 'PUT',
      body: {
        code: form.value.code,
        name: form.value.name,
        route_id: Number(form.value.route_id),
        departure_date: form.value.departure_date,
        arrival_date: form.value.arrival_date,
      },
    })
    await navigateTo('/voyages')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to update voyage'
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  if (deleting.value) return
  if (!confirm(`确认删除航次 #${id} 吗？`)) return
  deleting.value = true
  error.value = null
  try {
    await request(`/voyages/${id}`, { method: 'DELETE' })
    await navigateTo('/voyages')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to delete voyage'
  } finally {
    deleting.value = false
  }
}

onMounted(loadDetail)
</script>

<template>
  <div class="page">
    <h1>编辑航次 #{{ id }}</h1>
    <p v-if="loading">Loading...</p>
    <form v-else @submit.prevent="handleSave">
      <input v-model="form.code" placeholder="Code" :disabled="saving || deleting" />
      <input v-model="form.name" placeholder="Name" :disabled="saving || deleting" />
      <input v-model.number="form.route_id" type="number" min="1" placeholder="Route ID" :disabled="saving || deleting" />
      <input v-model="form.departure_date" type="date" :disabled="saving || deleting" />
      <input v-model="form.arrival_date" type="date" :disabled="saving || deleting" />
      <p v-if="error" class="error">{{ error }}</p>
      <div>
        <button type="submit" :disabled="saving || deleting">{{ saving ? '保存中...' : '保存' }}</button>
        <button type="button" style="margin-left:8px" :disabled="saving || deleting" @click="handleDelete">{{ deleting ? '删除中...' : '删除' }}</button>
      </div>
    </form>
  </div>
</template>
