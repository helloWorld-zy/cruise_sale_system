<template>
  <div class="min-h-screen bg-slate-50 p-4 md:p-6">
    <div class="mx-auto max-w-7xl">
      <div class="mb-4 flex items-center justify-between">
        <h1 class="text-xl font-semibold text-slate-900">舱房类型管理</h1>
        <NuxtLink to="/cabin-types/new" class="rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-500">新建舱房类型</NuxtLink>
      </div>

      <div class="mb-4 rounded-lg border border-slate-200 bg-white p-3 shadow-sm">
        <div class="flex flex-wrap items-center gap-3">
          <select
            v-model.number="filters.cruiseId"
            data-test="cruise-filter"
            class="h-10 min-w-64 rounded-md border border-slate-200 px-3 text-sm text-slate-700 outline-none ring-indigo-500 focus:ring-2"
          >
            <option :value="0">选择邮轮</option>
            <option v-for="cruise in cruises" :key="cruise.id" :value="Number(cruise.id)">{{ cruise.name || `邮轮 #${cruise.id}` }}</option>
          </select>
          <button type="button" class="h-10 rounded-md border border-slate-200 px-3 text-sm text-slate-700 hover:bg-slate-50" @click="loadItems">筛选</button>
        </div>
      </div>

      <div class="rounded-lg border border-slate-200 bg-white shadow-sm">
        <table class="w-full text-sm">
          <thead class="bg-slate-50 text-left text-slate-600">
            <tr>
              <th class="p-3">名称</th>
              <th class="p-3">代码</th>
              <th class="p-3">面积范围</th>
              <th class="p-3">容量</th>
              <th class="p-3">状态</th>
              <th class="p-3">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td class="p-3" colspan="6">加载中...</td>
            </tr>
            <tr v-else-if="error">
              <td class="p-3 text-rose-500" colspan="6">{{ error }}</td>
            </tr>
            <tr v-else-if="items.length === 0">
              <td class="p-3" colspan="6">暂无数据</td>
            </tr>
            <tr v-for="(item, idx) in items" v-else :key="item.id" :class="idx % 2 === 1 ? 'bg-slate-50' : ''">
              <td class="p-3 font-medium text-slate-900">{{ item.name || '-' }}</td>
              <td class="p-3 text-slate-600">{{ item.code || '-' }}</td>
              <td class="p-3 text-slate-600">{{ areaText(item) }}</td>
              <td class="p-3 text-slate-600">{{ item.max_capacity || item.capacity || '-' }}</td>
              <td class="p-3">
                <span :class="statusClass(item.status)">{{ statusText(item.status) }}</span>
              </td>
              <td class="p-3">
                <NuxtLink :to="`/cabin-types/${item.id}`" class="text-indigo-600 hover:text-indigo-500">编辑</NuxtLink>
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
const cruises = ref<Record<string, any>[]>([])
const items = ref<Record<string, any>[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

const filters = ref({
  cruiseId: 0,
})

async function loadCruises() {
  try {
    const res = await request('/cruises', { query: { page: 1, page_size: 100 } })
    const payload = res?.data ?? res ?? {}
    cruises.value = Array.isArray(payload) ? payload : payload?.list ?? []
    if (filters.value.cruiseId === 0 && cruises.value.length > 0) {
      filters.value.cruiseId = Number(cruises.value[0].id) || 0
    }
  } catch {
    cruises.value = []
  }
}

async function loadItems() {
  if (!filters.value.cruiseId) {
    items.value = []
    return
  }
  loading.value = true
  error.value = null
  try {
    const res = await request('/cabin-types', {
      query: {
        cruise_id: filters.value.cruiseId,
        page: 1,
        page_size: 50,
      },
    })
    const payload = res?.data ?? res ?? {}
    items.value = Array.isArray(payload) ? payload : payload?.list ?? []
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load cabin types'
  } finally {
    loading.value = false
  }
}

function statusText(statusRaw: unknown) {
  const status = Number(statusRaw)
  return status === 1 ? '上架' : '下架'
}

function statusClass(statusRaw: unknown) {
  const status = Number(statusRaw)
  if (status === 1) return 'rounded-full bg-emerald-50 px-2.5 py-0.5 text-xs font-medium text-emerald-700'
  return 'rounded-full bg-slate-100 px-2.5 py-0.5 text-xs font-medium text-slate-600'
}

function areaText(item: Record<string, any>) {
  const min = Number(item.area_min || 0)
  const max = Number(item.area_max || 0)
  if (min > 0 && max > 0) return `${min}-${max} m2`
  const area = Number(item.area || 0)
  return area > 0 ? `${area} m2` : '-'
}

onMounted(async () => {
  await loadCruises()
  await loadItems()
})
</script>
