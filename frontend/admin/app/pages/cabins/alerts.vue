<template>
  <div class="min-h-screen bg-slate-50 p-4 md:p-6">
    <div class="mx-auto max-w-6xl space-y-4">
      <div class="rounded-lg border border-rose-200 bg-rose-50 p-4">
        <h1 class="text-lg font-semibold text-rose-700">库存预警</h1>
        <p class="text-sm text-rose-600">以下舱位可用库存已低于阈值，请及时补仓或调整销售策略。</p>
      </div>

      <div class="rounded-lg border border-slate-200 bg-white shadow-sm">
        <table class="w-full text-sm">
          <thead class="bg-slate-50 text-left text-slate-600">
            <tr>
              <th class="p-3">舱位编号</th>
              <th class="p-3">可用库存</th>
              <th class="p-3">预警阈值</th>
              <th class="p-3">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading"><td class="p-3" colspan="4">加载中...</td></tr>
            <tr v-else-if="error"><td class="p-3 text-rose-500" colspan="4">{{ error }}</td></tr>
            <tr v-else-if="alerts.length === 0"><td class="p-3" colspan="4">暂无预警</td></tr>
            <tr v-for="item in alerts" v-else :key="item.cabin_sku_id">
              <td class="p-3 font-medium text-slate-900">{{ item.cabin_code || `SKU-${item.cabin_sku_id}` }}</td>
              <td class="p-3 text-base font-bold text-rose-600">{{ item.available }}</td>
              <td class="p-3 text-slate-600">{{ item.alert_threshold }}</td>
              <td class="p-3">
                <NuxtLink :to="`/cabins/inventory?skuId=${item.cabin_sku_id}`" class="text-indigo-600 hover:text-indigo-500">调整阈值</NuxtLink>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'

const { request } = useApi()
const loading = ref(false)
const error = ref<string | null>(null)
const alerts = ref<Record<string, any>[]>([])

async function loadAlerts() {
  loading.value = true
  error.value = null
  try {
    const res = await request('/cabins/alerts')
    const payload = res?.data ?? res ?? []
    alerts.value = Array.isArray(payload) ? payload : payload?.list ?? []
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load alerts'
  } finally {
    loading.value = false
  }
}

onMounted(loadAlerts)
</script>
