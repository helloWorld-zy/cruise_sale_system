<template>
  <div class="min-h-screen bg-[#f8f4ed] pb-24 text-slate-900">
    <main class="mx-auto max-w-6xl px-6 py-8">
      <h1 class="mb-6 font-['Playfair_Display','Georgia',serif] text-3xl text-[#12263a]">舱位浏览</h1>

      <div class="mb-6 flex flex-wrap gap-4 rounded-lg border border-[#eadfcb] bg-white p-4 shadow-sm">
        <select v-model.number="filters.routeId" class="h-10 min-w-40 rounded-md border border-slate-200 px-3 text-sm outline-none ring-indigo-500 focus:ring-2">
          <option :value="0">选择航线</option>
          <option v-for="route in routes" :key="route.id" :value="Number(route.id)">{{ route.name || `航线 #${route.id}` }}</option>
        </select>
        
        <select v-model.number="filters.cabinTypeId" class="h-10 min-w-40 rounded-md border border-slate-200 px-3 text-sm outline-none ring-indigo-500 focus:ring-2">
          <option :value="0">选择舱型</option>
          <option v-for="ct in cabinTypes" :key="ct.id" :value="Number(ct.id)">{{ ct.name || `舱型 #${ct.id}` }}</option>
        </select>

        <select v-model.number="filters.sortBy" class="h-10 min-w-32 rounded-md border border-slate-200 px-3 text-sm outline-none ring-indigo-500 focus:ring-2">
          <option :value="0">默认排序</option>
          <option :value="1">价格从低到高</option>
          <option :value="2">价格从高到低</option>
          <option :value="3">面积从大到小</option>
        </select>

        <button type="button" class="rounded-md bg-[#12263a] px-4 py-2 text-sm text-white hover:bg-[#1a3a5c]" @click="loadCabins">
          搜索
        </button>
      </div>

      <div v-if="loading" class="py-12 text-center text-slate-600">
        加载中...
      </div>

      <div v-else-if="error" class="py-12 text-center text-rose-600">
        {{ error }}
      </div>

      <div v-else-if="cabins.length === 0" class="py-12 text-center text-slate-600">
        暂无可用舱位
      </div>

      <div v-else class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
        <NuxtLink
          v-for="cabin in cabins"
          :key="cabin.id"
          :to="`/cabins/${cabin.id}`"
          class="group overflow-hidden rounded-2xl border border-[#eadfcb] bg-white shadow-sm transition-shadow hover:shadow-md"
        >
          <div class="aspect-[16/10] overflow-hidden">
            <img 
              :src="cabinImage(cabin)" 
              :alt="cabin.code || '舱位'"
              class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
            />
          </div>
          <div class="p-4">
            <h3 class="mb-2 font-['Playfair_Display','Georgia',serif] text-lg text-[#12263a]">
              {{ cabin.code || `舱位 #${cabin.id}` }}
            </h3>
            <p class="mb-3 text-sm text-slate-600">
              {{ cabin.cabin_type_name || '-' }} · {{ cabin.deck || '-' }}甲板 · {{ cabin.area || '-' }}m²
            </p>
            
            <div class="mb-3 flex flex-wrap gap-1">
              <span v-if="cabin.has_window" class="rounded-full bg-blue-50 px-2 py-0.5 text-xs text-blue-700">有窗</span>
              <span v-if="cabin.has_balcony" class="rounded-full bg-blue-50 px-2 py-0.5 text-xs text-blue-700">有阳台</span>
              <span v-if="cabin.bed_type" class="rounded-full bg-blue-50 px-2 py-0.5 text-xs text-blue-700">{{ cabin.bed_type }}</span>
            </div>

            <div class="flex items-center justify-between">
              <div>
                <span class="font-['Playfair_Display','Georgia',serif] text-xl text-[#12263a]">¥{{ displayPrice(cabin) }}</span>
                <span v-if="cabin.child_price_cents" class="ml-2 text-xs text-slate-500">儿童 ¥{{ Math.round(cabin.child_price_cents / 100) }}</span>
              </div>
              <span :class="inventoryBadgeClass(cabin)">{{ inventoryLabel(cabin) }}</span>
            </div>
          </div>
        </NuxtLink>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'

interface CabinItem {
  id: number
  code: string
  cabin_type_name: string
  deck: string
  area: number
  has_window: boolean
  has_balcony: boolean
  bed_type: string
  min_price_cents: number
  child_price_cents: number
  total: number
  locked: number
  sold: number
  images?: Array<{ url: string }>
}

interface RouteItem {
  id: number
  name: string
}

interface CabinTypeItem {
  id: number
  name: string
}

const { request } = useApi()

const loading = ref(false)
const error = ref<string | null>(null)
const cabins = ref<CabinItem[]>([])
const routes = ref<RouteItem[]>([])
const cabinTypes = ref<CabinTypeItem[]>([])

const filters = ref({
  routeId: 0,
  cabinTypeId: 0,
  sortBy: 0,
})

const sortedCabins = computed(() => {
  const list = [...cabins.value]
  if (filters.value.sortBy === 1) {
    return list.sort((a, b) => (a.min_price_cents || 0) - (b.min_price_cents || 0))
  }
  if (filters.value.sortBy === 2) {
    return list.sort((a, b) => (b.min_price_cents || 0) - (a.min_price_cents || 0))
  }
  if (filters.value.sortBy === 3) {
    return list.sort((a, b) => (b.area || 0) - (a.area || 0))
  }
  return list
})

async function loadOptions() {
  try {
    const [routeRes, typeRes] = await Promise.all([
      request('/routes'),
      request('/cabin-types', { query: { page: 1, page_size: 100 } }),
    ])
    
    const routePayload = routeRes?.data ?? routeRes ?? []
    routes.value = Array.isArray(routePayload) ? routePayload : routePayload?.list ?? []
    
    const typePayload = typeRes?.data ?? typeRes ?? {}
    cabinTypes.value = Array.isArray(typePayload) ? typePayload : typePayload?.list ?? []
  } catch (e) {
    console.error('Failed to load options:', e)
  }
}

async function loadCabins() {
  loading.value = true
  error.value = null
  
  try {
    const query: Record<string, unknown> = { page: 1, page_size: 20 }
    if (filters.value.routeId > 0) query.voyage_id = filters.value.routeId
    if (filters.value.cabinTypeId > 0) query.cabin_type_id = filters.value.cabinTypeId
    
    const res = await request('/cabins', { query })
    const payload = res?.data ?? res ?? {}
    cabins.value = Array.isArray(payload) ? payload : payload?.list ?? []
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : '加载舱位失败'
  } finally {
    loading.value = false
  }
}

function cabinImage(cabin: CabinItem): string {
  if (cabin.images && cabin.images.length > 0) {
    return cabin.images[0].url
  }
  return `https://picsum.photos/seed/cabin-${cabin.id}/800/500`
}

function displayPrice(cabin: CabinItem): string {
  const price = cabin.min_price_cents
  return price ? Math.round(price / 100).toString() : '-'
}

function inventoryCount(cabin: CabinItem): number {
  const total = cabin.total || 0
  const locked = cabin.locked || 0
  const sold = cabin.sold || 0
  return Math.max(0, total - locked - sold)
}

function inventoryPercent(cabin: CabinItem): number {
  const total = cabin.total || 0
  if (total <= 0) return 0
  return Math.round((inventoryCount(cabin) / total) * 100)
}

function inventoryLabel(cabin: CabinItem): string {
  const percent = inventoryPercent(cabin)
  const count = inventoryCount(cabin)
  if (count === 0) return '已售罄'
  if (percent < 20) return '即将售罄'
  if (percent < 50) return '库存紧张'
  return '库存充足'
}

function inventoryBadgeClass(cabin: CabinItem): string {
  const percent = inventoryPercent(cabin)
  const count = inventoryCount(cabin)
  if (count === 0) return 'rounded-full bg-slate-100 px-2 py-0.5 text-xs font-medium text-slate-600'
  if (percent < 20) return 'rounded-full bg-rose-50 px-2 py-0.5 text-xs font-medium text-rose-700'
  if (percent < 50) return 'rounded-full bg-amber-50 px-2 py-0.5 text-xs font-medium text-amber-700'
  return 'rounded-full bg-emerald-50 px-2 py-0.5 text-xs font-medium text-emerald-700'
}

onMounted(async () => {
  await loadOptions()
  await loadCabins()
})
</script>
