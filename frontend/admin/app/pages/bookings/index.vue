<!-- admin/app/pages/bookings/index.vue — 订单列表页面 -->
<!-- H-02 修复：使用 useApi 获取真实订单数据，去掉硬编码模拟数据 -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'

declare const useApi: any

// 使用 useApi 请求后端订单列表数据
const { request } = useApi()
const items = ref<{ id: number; status: string; total: number }[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

onMounted(async () => {
  loading.value = true
  try {
    const res = await request('/bookings')
    items.value = res?.data ?? []
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load bookings'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="page">
    <h1>Bookings</h1>
    <p v-if="loading">Loading...</p>
    <p v-else-if="error" class="error">{{ error }}</p>
    <BookingRow v-for="b in items" :key="b.id" :booking="b" />
  </div>
</template>
