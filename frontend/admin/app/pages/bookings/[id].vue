<!-- admin/app/pages/bookings/[id].vue — 订单详情页面 -->
<!-- 根据路由参数加载并展示单个订单的详细信息 -->
<script setup lang="ts">
import { onMounted, ref } from 'vue'

declare const useApi: any

// 根据路由参数中的订单 ID 加载并展示订单详情
const route = useRoute()
const { request } = useApi()

const loading = ref(false)
const error = ref<string | null>(null)
const booking = ref<{ id: number; status: string; total: number } | null>(null)

onMounted(async () => {
  loading.value = true
  error.value = null
  try {
    const id = String(route.params.id ?? '')
    const res = await request(`/bookings/${id}`)
    booking.value = res?.data ?? null
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load booking detail'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="page">
    <h1>Booking Detail {{ route.params.id }}</h1>
    <p v-if="loading">Loading...</p>
    <p v-else-if="error" class="error">{{ error }}</p>
    <div v-else-if="booking">
      <p>ID: {{ booking.id }}</p>
      <p>Status: {{ booking.status }}</p>
      <p>Total: {{ booking.total }}</p>
    </div>
  </div>
</template>
