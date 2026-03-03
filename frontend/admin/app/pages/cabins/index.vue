<template>
  <div class="min-h-screen bg-slate-50 p-4 md:p-6">
    <div class="mx-auto max-w-7xl">
      <div class="mb-4 flex items-center justify-between">
        <h1 class="text-xl font-semibold text-slate-900">舱位商品管理</h1>
        <NuxtLink to="/cabins/new" class="rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-500">新建舱位</NuxtLink>
      </div>

      <div class="mb-4 rounded-lg border border-slate-200 bg-white p-3 shadow-sm">
        <div class="flex flex-wrap items-center gap-3">
          <select v-model.number="filters.cruiseId" data-test="filter-cruise" class="h-10 min-w-44 rounded-md border border-slate-200 px-3 text-sm outline-none ring-indigo-500 focus:ring-2" @change="syncRouteOptions">
            <option :value="0">邮轮</option>
            <option v-for="cruise in cruises" :key="cruise.id" :value="Number(cruise.id)">{{ cruise.name || `邮轮 #${cruise.id}` }}</option>
          </select>
          <select v-model.number="filters.routeId" data-test="filter-route" class="h-10 min-w-44 rounded-md border border-slate-200 px-3 text-sm outline-none ring-indigo-500 focus:ring-2">
            <option :value="0">航线</option>
            <option v-for="route in routeOptions" :key="route.id" :value="Number(route.id)">{{ route.name || `航线 #${route.id}` }}</option>
          </select>
          <select v-model.number="filters.cabinTypeId" data-test="filter-cabin-type" class="h-10 min-w-44 rounded-md border border-slate-200 px-3 text-sm outline-none ring-indigo-500 focus:ring-2">
            <option :value="0">舱型</option>
            <option v-for="type in cabinTypes" :key="type.id" :value="Number(type.id)">{{ type.name || `舱型 #${type.id}` }}</option>
          </select>
          <input v-model="filters.keyword" data-test="filter-keyword" placeholder="搜索舱位编号" class="h-10 w-48 rounded-md border border-slate-200 px-3 text-sm outline-none ring-indigo-500 focus:ring-2" @keyup.enter="loadItems" />
          <select v-model.number="filters.status" data-test="filter-status" class="h-10 min-w-32 rounded-md border border-slate-200 px-3 text-sm outline-none ring-indigo-500 focus:ring-2">
            <option :value="0">全部状态</option>
            <option :value="1">上架</option>
            <option :value="2">维护中</option>
            <option :value="-1">下架</option>
          </select>
          <button type="button" class="h-10 rounded-md border border-slate-200 px-3 text-sm text-slate-700 hover:bg-slate-50" @click="loadItems">筛选</button>
        </div>
      </div>

      <div class="rounded-lg border border-slate-200 bg-white shadow-sm">
        <table class="w-full text-sm">
          <thead class="bg-slate-50 text-left text-slate-600">
            <tr>
              <th class="w-10 p-3">
                <input type="checkbox" :checked="allChecked" @change="toggleAll(($event.target as HTMLInputElement).checked)" />
              </th>
              <th class="p-3">编号</th>
              <th class="p-3">邮轮</th>
              <th class="p-3">航线</th>
              <th class="p-3">舱型</th>
              <th class="p-3">面积</th>
              <th class="p-3">价格区间</th>
              <th class="p-3">库存</th>
              <th class="p-3">状态</th>
              <th class="p-3">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading"><td class="p-3" colspan="10">加载中...</td></tr>
            <tr v-else-if="error"><td class="p-3 text-rose-500" colspan="10">{{ error }}</td></tr>
            <tr v-else-if="items.length === 0"><td class="p-3" colspan="10">暂无数据</td></tr>
            <tr v-for="(item, idx) in items" v-else :key="item.id" :class="idx % 2 === 1 ? 'bg-slate-50' : ''">
              <td class="p-3"><input type="checkbox" :checked="selectedIds.has(Number(item.id))" @change="toggleSingle(Number(item.id), ($event.target as HTMLInputElement).checked)" /></td>
              <td class="p-3 font-medium text-slate-900">{{ item.code || `CABIN-${item.id}` }}</td>
              <td class="p-3 text-slate-600">{{ item.cruise_name || '-' }}</td>
              <td class="p-3 text-slate-600">{{ item.route_name || '-' }}</td>
              <td class="p-3 text-slate-600">{{ item.cabin_type_name || '-' }}</td>
              <td class="p-3 text-slate-600">{{ Number(item.area || 0) > 0 ? `${item.area} m2` : '-' }}</td>
              <td class="p-3 font-mono text-slate-700">{{ priceRangeText(item) }}</td>
              <td class="p-3">
                <div class="w-24">
                  <p class="mb-1 text-xs text-slate-600">{{ available(item) }}/{{ total(item) }}</p>
                  <div class="h-1.5 rounded bg-slate-200">
                    <div class="h-1.5 rounded" :class="inventoryBarClass(item)" :style="{ width: `${inventoryPercent(item)}%` }" />
                  </div>
                </div>
              </td>
              <td class="p-3"><span :class="statusClass(item.status)">{{ statusText(item.status) }}</span></td>
              <td class="p-3">
                <NuxtLink :to="`/cabins/${item.id}`" class="text-indigo-600 hover:text-indigo-500">编辑</NuxtLink>
                <NuxtLink :to="`/cabins/inventory?skuId=${item.id}`" class="ml-2 text-slate-600 hover:text-slate-500">库存</NuxtLink>
                <NuxtLink :to="`/cabins/pricing?skuId=${item.id}`" class="ml-2 text-slate-600 hover:text-slate-500">价格</NuxtLink>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="selectedIds.size > 0" data-test="batch-action" class="fixed bottom-0 left-0 right-0 flex items-center justify-center gap-3 bg-indigo-600 px-4 py-3 text-sm text-white">
        <span>已选 {{ selectedIds.size }} 项</span>
        <button type="button" class="rounded bg-white/20 px-3 py-1.5 hover:bg-white/30" @click="batchUpdateStatus(1)">批量上架</button>
        <button type="button" class="rounded bg-white/20 px-3 py-1.5 hover:bg-white/30" @click="batchUpdateStatus(-1)">批量下架</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

const { request } = useApi()
const loading = ref(false)
const error = ref<string | null>(null)
const items = ref<Record<string, any>[]>([])
const cruises = ref<Record<string, any>[]>([])
const routes = ref<Record<string, any>[]>([])
const cabinTypes = ref<Record<string, any>[]>([])
const selectedIds = ref<Set<number>>(new Set())

const filters = ref({
  cruiseId: 0,
  routeId: 0,
  cabinTypeId: 0,
  keyword: '',
  status: 0,
})

const routeOptions = computed(() => {
  if (!filters.value.cruiseId) return routes.value
  return routes.value.filter((item) => Number(item.cruise_id) === Number(filters.value.cruiseId))
})

const allChecked = computed(() => items.value.length > 0 && items.value.every((it) => selectedIds.value.has(Number(it.id))))

async function loadOptions() {
  try {
    const [cruiseRes, routeRes, typeRes] = await Promise.all([
      request('/cruises', { query: { page: 1, page_size: 100 } }),
      request('/routes'),
      request('/cabin-types', { query: { cruise_id: 1, page: 1, page_size: 200 } }),
    ])
    const cruisePayload = cruiseRes?.data ?? cruiseRes ?? {}
    cruises.value = Array.isArray(cruisePayload) ? cruisePayload : cruisePayload?.list ?? []
    const routePayload = routeRes?.data ?? routeRes ?? []
    routes.value = Array.isArray(routePayload) ? routePayload : routePayload?.list ?? []
    const typePayload = typeRes?.data ?? typeRes ?? {}
    cabinTypes.value = Array.isArray(typePayload) ? typePayload : typePayload?.list ?? []
  } catch {
    cruises.value = []
    routes.value = []
    cabinTypes.value = []
  }
}

function syncRouteOptions() {
  if (filters.value.routeId && !routeOptions.value.some((item) => Number(item.id) === Number(filters.value.routeId))) {
    filters.value.routeId = 0
  }
}

async function loadItems() {
  loading.value = true
  error.value = null
  try {
    const query: Record<string, any> = { page: 1, page_size: 50 }
    if (filters.value.routeId > 0) query.voyage_id = filters.value.routeId
    if (filters.value.cabinTypeId > 0) query.cabin_type_id = filters.value.cabinTypeId
    if (filters.value.keyword.trim()) query.keyword = filters.value.keyword.trim()
    if (filters.value.status !== 0) query.status = filters.value.status === -1 ? 0 : filters.value.status
    const res = await request('/cabins', { query })
    const payload = res?.data ?? res ?? {}
    items.value = Array.isArray(payload) ? payload : payload?.list ?? []
    selectedIds.value = new Set()
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load cabins'
  } finally {
    loading.value = false
  }
}

function total(item: Record<string, any>) {
  return Math.max(Number(item.total || 0), 0)
}

function available(item: Record<string, any>) {
  if (item.available !== undefined) return Math.max(Number(item.available || 0), 0)
  return Math.max(total(item) - Number(item.locked || 0) - Number(item.sold || 0), 0)
}

function inventoryPercent(item: Record<string, any>) {
  const totalValue = total(item)
  if (totalValue <= 0) return 0
  return Math.min(100, Math.max(0, Math.round((available(item) / totalValue) * 100)))
}

function inventoryBarClass(item: Record<string, any>) {
  const percent = inventoryPercent(item)
  if (percent > 50) return 'bg-emerald-500'
  if (percent >= 20) return 'bg-amber-500'
  return 'bg-rose-500'
}

function priceRangeText(item: Record<string, any>) {
  const min = Number(item.min_price_cents || item.min_price || 0)
  const max = Number(item.max_price_cents || item.max_price || 0)
  if (min > 0 && max > 0) return `${Math.round(min / 100)}-${Math.round(max / 100)}`
  const single = Number(item.price_cents || 0)
  return single > 0 ? `${Math.round(single / 100)}` : '-'
}

function statusText(statusRaw: unknown) {
  const status = Number(statusRaw)
  if (status === 1) return '上架'
  if (status === 2) return '维护中'
  return '下架'
}

function statusClass(statusRaw: unknown) {
  const status = Number(statusRaw)
  if (status === 1) return 'rounded-full bg-emerald-50 px-2.5 py-0.5 text-xs font-medium text-emerald-700'
  if (status === 2) return 'rounded-full bg-amber-50 px-2.5 py-0.5 text-xs font-medium text-amber-700'
  return 'rounded-full bg-slate-100 px-2.5 py-0.5 text-xs font-medium text-slate-600'
}

function toggleSingle(id: number, checked: boolean) {
  const next = new Set(selectedIds.value)
  if (checked) next.add(id)
  else next.delete(id)
  selectedIds.value = next
}

function toggleAll(checked: boolean) {
  if (!checked) {
    selectedIds.value = new Set()
    return
  }
  selectedIds.value = new Set(items.value.map((item) => Number(item.id)).filter((id) => Number.isFinite(id) && id > 0))
}

async function batchUpdateStatus(status: number) {
  if (selectedIds.value.size === 0) return
  try {
    await request('/cabins/batch-status', {
      method: 'PUT',
      body: { ids: Array.from(selectedIds.value), status },
    })
    await loadItems()
  } catch (e: any) {
    error.value = e?.message ?? 'failed to batch update status'
  }
}

onMounted(async () => {
  await loadOptions()
  await loadItems()
})
</script>

