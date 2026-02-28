<!-- admin/app/pages/bookings/[id].vue — 订单详情页面 -->
<!-- 根据路由参数加载并展示单个订单的详细信息 -->
<script setup lang="ts">
import { onMounted, ref } from 'vue'

declare const useApi: any
declare const navigateTo: any

// 根据路由参数中的订单 ID 加载并展示订单详情
const route = useRoute()
const { request } = useApi()

const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const error = ref<string | null>(null)
const booking = ref<{ id: number; status: string; total_cents: number } | null>(null)
const status = ref('')

onMounted(async () => {
  loading.value = true
  error.value = null
  try {
    const id = String(route.params.id ?? '')
    const res = await request(`/bookings/${id}`)
    booking.value = res?.data ?? res ?? null
    status.value = booking.value?.status ?? ''
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load booking detail'
  } finally {
    loading.value = false
  }
})

async function handleSave() {
  if (!booking.value || saving.value) return
  saving.value = true
  error.value = null
  try {
    await request(`/bookings/${booking.value.id}`, {
      method: 'PUT',
      body: {
        status: status.value,
      },
    })
    booking.value.status = status.value
  } catch (e: any) {
    error.value = e?.message ?? 'failed to update booking'
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  if (!booking.value || deleting.value) return
  if (!confirm(`确认删除订单 #${booking.value.id} 吗？`)) return
  deleting.value = true
  error.value = null
  try {
    await request(`/bookings/${booking.value.id}`, { method: 'DELETE' })
    await navigateTo('/bookings')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to delete booking'
  } finally {
    deleting.value = false
  }
}
</script>

<template>
  <div class="page">
    <h1>Booking Detail {{ route.params.id }}</h1>
    <p v-if="loading">Loading...</p>
    <p v-else-if="error" class="error">{{ error }}</p>
    <div v-else-if="booking">
      <p>ID: {{ booking.id }}</p>
      <p>Status: {{ booking.status }}</p>
      <p>Total: {{ booking.total_cents }}</p>
      <input v-model="status" placeholder="新状态" :disabled="saving || deleting" />
      <div>
        <button type="button" :disabled="saving || deleting" @click="handleSave">{{ saving ? '保存中...' : '保存状态' }}</button>
        <button type="button" style="margin-left:8px" :disabled="saving || deleting" @click="handleDelete">{{ deleting ? '删除中...' : '删除订单' }}</button>
      </div>
    </div>
  </div>
</template>
