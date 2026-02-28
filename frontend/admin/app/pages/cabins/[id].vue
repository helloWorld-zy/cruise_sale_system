<!-- admin/app/pages/cabins/[id].vue -->
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
  voyage_id: 0,
  cabin_type_id: 0,
  total: 0,
  available: 0,
})

async function loadDetail() {
  loading.value = true
  error.value = null
  try {
    const res = await request(`/cabins/${id}`)
    const data = res?.data ?? res ?? {}
    form.value = {
      voyage_id: Number(data.voyage_id ?? 0),
      cabin_type_id: Number(data.cabin_type_id ?? 0),
      total: Number(data.total ?? 0),
      available: Number(data.available ?? 0),
    }
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load cabin detail'
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  if (saving.value) return
  saving.value = true
  error.value = null
  try {
    await request(`/cabins/${id}`, {
      method: 'PUT',
      body: {
        voyage_id: Number(form.value.voyage_id),
        cabin_type_id: Number(form.value.cabin_type_id),
        total: Number(form.value.total),
        available: Number(form.value.available),
      },
    })
    await navigateTo('/cabins')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to update cabin'
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  if (deleting.value) return
  if (!confirm(`确认删除舱房 #${id} 吗？`)) return
  deleting.value = true
  error.value = null
  try {
    await request(`/cabins/${id}`, { method: 'DELETE' })
    await navigateTo('/cabins')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to delete cabin'
  } finally {
    deleting.value = false
  }
}

onMounted(loadDetail)
</script>

<template>
  <div class="page">
    <h1>编辑舱房 #{{ id }}</h1>
    <p v-if="loading">Loading...</p>
    <form v-else @submit.prevent="handleSave">
      <input v-model.number="form.voyage_id" type="number" min="1" placeholder="Voyage ID" :disabled="saving || deleting" />
      <input v-model.number="form.cabin_type_id" type="number" min="1" placeholder="Cabin Type ID" :disabled="saving || deleting" />
      <input v-model.number="form.total" type="number" min="0" placeholder="Total" :disabled="saving || deleting" />
      <input v-model.number="form.available" type="number" min="0" placeholder="Available" :disabled="saving || deleting" />
      <p v-if="error" class="error">{{ error }}</p>
      <div>
        <button type="submit" :disabled="saving || deleting">{{ saving ? '保存中...' : '保存' }}</button>
        <button type="button" style="margin-left:8px" :disabled="saving || deleting" @click="handleDelete">{{ deleting ? '删除中...' : '删除' }}</button>
      </div>
    </form>
  </div>
</template>
