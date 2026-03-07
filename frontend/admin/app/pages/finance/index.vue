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
  <div class="admin-page">
    <AdminPageHeader title="财务总览" subtitle="展示核心财务统计指标" />
    <h1 class="sr-only">Finance</h1>

    <AdminDataCard flush>
      <div class="overflow-x-auto">
        <p v-if="loading" class="p-3 text-sm text-slate-600">Loading...</p>
        <p v-else-if="error" class="p-3 text-sm text-rose-500">{{ error }}</p>
        <p v-else-if="rows.length === 0" class="p-3 text-sm text-slate-600">No data</p>
        <table v-else class="w-full min-w-[680px] text-sm">
          <thead class="bg-slate-50 text-left text-slate-600">
            <tr>
              <th class="p-3">指标</th>
              <th class="p-3">数值</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="row in rows" :key="row.label">
              <td class="p-3 text-slate-600">{{ row.label }}</td>
              <td class="p-3 text-slate-900">{{ row.value }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </AdminDataCard>
  </div>
</template>
