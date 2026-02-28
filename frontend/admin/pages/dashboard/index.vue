<!-- admin/pages/dashboard/index.vue — 仪表盘页面 -->
<!-- 展示销售总额和订单数等关键统计指标 -->
<script setup lang="ts">
import { onMounted, ref } from 'vue'
import StatCard from '../../components/StatCard.vue'
// 仪表盘统计数据，改为调用后端分析接口。
const { request } = useApi()
const loading = ref(false)
const error = ref<string | null>(null)
const summary = ref({ today_sales: 0, weekly_trend: [] as number[], today_orders: 0 })

onMounted(async () => {
  loading.value = true
  error.value = null
  try {
    const res = await request('/admin/analytics/summary')
    summary.value = res?.data ?? res
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load dashboard'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="page">
    <h1>Dashboard</h1>
    <p v-if="loading">Loading...</p>
    <p v-else-if="error" class="text-red-500">{{ error }}</p>
    <template v-else>
    <!-- 销售额统计卡片 -->
    <StatCard title="Sales" :value="summary.today_sales" />
    <!-- 订单数统计卡片 -->
    <StatCard title="Orders" :value="summary.today_orders" />
    </template>
  </div>
</template>

