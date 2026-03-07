<!-- admin/app/pages/dashboard/index.vue — 仪表盘页面 -->
<!-- 展示销售总额和订单数等关键统计指标 -->
<script setup lang="ts">
import { onMounted, ref } from 'vue'
import StatCard from '../../../components/StatCard.vue'

// 仪表盘统计数据，改为调用后端分析接口。
const { request } = useApi()
const loading = ref(false)
const error = ref<string | null>(null)
const summary = ref({ today_sales: 0, weekly_trend: [] as number[], today_orders: 0 })
const empty = ref(false)

onMounted(async () => {
  loading.value = true
  error.value = null
  empty.value = false
  try {
    const res = await request('/admin/analytics/summary')
    const payload = res?.data ?? res ?? null
    if (!payload || Object.keys(payload).length === 0) {
      empty.value = true
      return
    }
    summary.value = payload
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load dashboard'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="admin-page">
    <AdminPageHeader title="仪表盘" subtitle="查看销售与订单核心指标" />
    <h1 class="sr-only">Dashboard</h1>

    <AdminDataCard>
      <p v-if="loading" class="p-3 text-sm text-slate-600">Loading...</p>
      <p v-else-if="error" class="p-3 text-sm text-rose-500">{{ error }}</p>
      <p v-else-if="empty" data-test="empty" class="p-3 text-sm text-slate-600">暂无数据</p>
      <div v-else class="grid gap-4 md:grid-cols-2">
        <!-- 销售额统计卡片 -->
        <StatCard title="Sales" :value="summary.today_sales" />
        <!-- 订单数统计卡片 -->
        <StatCard title="Orders" :value="summary.today_orders" />
      </div>
    </AdminDataCard>
  </div>
</template>
