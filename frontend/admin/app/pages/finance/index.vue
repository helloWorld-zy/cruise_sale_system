<!-- admin/app/pages/finance/index.vue — 财务管理页面 -->
<!-- 展示财务概览数据表格，使用统计摘要接口 -->
<script setup lang="ts">
import { onMounted, ref } from 'vue'

const { request } = useApi()
const loading = ref(false)
const error = ref<string | null>(null)

interface FinanceRow {
  label: string
  value: number | string
}
const rows = ref<FinanceRow[]>([])

onMounted(async () => {
  loading.value = true
  error.value = null
  try {
    const summary = await request('/admin/analytics/summary')
    const payload = summary?.data ?? summary
    rows.value = [
      { label: '今日销售额（分）', value: Number(payload?.today_sales ?? 0) },
      { label: '今日订单数', value: Number(payload?.today_orders ?? 0) },
      { label: '总收入（分）', value: Number(payload?.total_revenue ?? 0) },
      { label: '总订单数', value: Number(payload?.total_orders ?? 0) },
    ]
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load finance data'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="page">
    <h1>Finance</h1>
    <p v-if="loading">Loading...</p>
    <p v-else-if="error" class="text-red-500">{{ error }}</p>
    <p v-else-if="rows.length === 0">No data</p>
    <!-- 财务概览表格 -->
    <table v-else>
      <thead>
        <tr>
          <th>指标</th>
          <th>数值</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="row in rows" :key="row.label">
          <td>{{ row.label }}</td>
          <td>{{ row.value }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
