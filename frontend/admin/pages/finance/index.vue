<!-- admin/pages/finance/index.vue — 财务管理页面 -->
<!-- 展示财务流水数据表格 -->
<script setup lang="ts">
import { onMounted, ref } from 'vue'
// 暂时复用分析摘要接口作为财务概览数据源。
// TODO: 替换为专用财务 API
const { request } = useApi()
const loading = ref(false)
const error = ref<string | null>(null)
const rows = ref<Array<{ date: string; amount: number }>>([])

onMounted(async () => {
  loading.value = true
  error.value = null
  try {
    const summary = await request('/admin/analytics/summary')
    const payload = summary?.data ?? summary
    rows.value = [
      { date: '今日销售', amount: Number(payload?.today_sales ?? 0) },
      { date: '今日订单', amount: Number(payload?.today_orders ?? 0) },
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
    <!-- 财务流水表格 -->
    <table v-else>
      <thead>
        <tr>
          <th>日期</th>
          <th>金额</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="row in rows" :key="row.date">
          <td>{{ row.date }}</td>
          <td>{{ row.amount }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

