<template>
  <div class="min-h-screen bg-[#fcfbf9] pb-24 text-slate-900 font-sans">
    <div class="bg-[#0f3d5c] text-[#fcfbf9] py-16 px-6 mb-10 relative overflow-hidden">
      <div class="absolute inset-0 opacity-10 bg-[url('https://www.transparenttextures.com/patterns/stardust.png')]"></div>
      <div class="max-w-6xl mx-auto relative z-10">
        <p class="text-[10px] tracking-[0.3em] font-bold uppercase mb-4 text-[#c9a96e]">Cabin Selection</p>
        <h1 class="font-['Playfair_Display','Georgia',serif] text-4xl md:text-5xl text-white">搜索舱位</h1>
      </div>
    </div>

    <main class="mx-auto max-w-6xl px-6">
      <div class="mb-12 flex flex-wrap gap-4 -mt-24 relative z-20 bg-white/90 backdrop-blur-md p-6 border border-[#eadfcb]/50 shadow-2xl">
        <input 
          v-model="filters.keyword"
          type="text"
          placeholder="搜索舱位编号"
          class="h-12 flex-1 min-w-48 border-b-2 border-slate-200 bg-transparent px-3 text-sm outline-none focus:border-[#c9a96e] transition-colors rounded-none"
          @keyup.enter="handleSearch"
        />
        
        <select v-model.number="filters.routeId" class="h-12 min-w-40 border-b-2 border-slate-200 bg-transparent px-3 text-sm outline-none focus:border-[#c9a96e] transition-colors rounded-none appearance-none">
          <option :value="0">选择航线</option>
          <option v-for="route in routes" :key="route.id" :value="Number(route.id)">{{ route.name || `航线 #${route.id}` }}</option>
        </select>
        
        <select v-model.number="filters.cabinTypeId" class="h-12 min-w-40 border-b-2 border-slate-200 bg-transparent px-3 text-sm outline-none focus:border-[#c9a96e] transition-colors rounded-none appearance-none">
          <option :value="0">选择舱型</option>
          <option v-for="ct in cabinTypes" :key="ct.id" :value="Number(ct.id)">{{ ct.name || `舱型 #${ct.id}` }}</option>
        </select>

        <button type="button" class="h-12 bg-[#0f3d5c] px-8 text-xs tracking-widest text-white hover:bg-[#1a3a5c] transition-colors uppercase font-serif" @click="handleSearch">
          搜索
        </button>
      </div>

      <div v-if="loading" class="py-12 text-center text-slate-500 font-serif italic">
        加载中...
      </div>

      <div v-else-if="error" class="py-12 text-center text-rose-800">
        {{ error }}
      </div>

      <div v-else-if="results.length === 0" class="py-12 text-center text-slate-500 font-serif italic">
        暂无搜索结果
      </div>

      <div v-else class="grid grid-cols-1 gap-10 md:grid-cols-2 lg:grid-cols-3">
        <NuxtLink
          v-for="cabin in results"
          :key="cabin.id"
          :to="`/cabins/${cabin.id}`"
          class="group overflow-hidden bg-white shadow-sm transition-all duration-500 hover:-translate-y-2 hover:shadow-2xl border border-transparent hover:border-[#eadfcb]"
        >
          <div class="relative aspect-[16/10] overflow-hidden bg-slate-100">
            <img 
              :src="cabinImage(cabin)" 
              :alt="cabin.code || '舱位'"
              class="h-full w-full object-cover transition-transform duration-700 group-hover:scale-110"
            />
            <div class="absolute inset-0 bg-gradient-to-t from-black/40 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
          </div>
          <div class="flex flex-col space-y-4 p-6 bg-white relative">
            <h3 class="font-['Playfair_Display','Georgia',serif] text-2xl text-[#0f3d5c]">
              {{ cabin.code || `舱位 #${cabin.id}` }}
            </h3>
            <div class="flex flex-wrap gap-3 text-[11px] font-semibold tracking-wider text-slate-500">
              <span class="border-b border-[#eadfcb] pb-1">{{ cabin.cabin_type_name || '-' }}</span>
              <span class="border-b border-[#eadfcb] pb-1">{{ cabin.deck || '-' }}甲板</span>
              <span class="border-b border-[#eadfcb] pb-1">{{ cabin.area || '-' }}m²</span>
            </div>
            
            <div class="pt-4 mt-2 border-t border-slate-100 flex items-center justify-between">
              <div>
                <span class="font-['Playfair_Display','Georgia',serif] text-2xl text-[#c9a96e]">¥{{ displayPrice(cabin) }}</span>
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
import { computed, onMounted, ref } from 'vue'

interface CabinItem {
  id: number
  code: string
  cabin_type_name: string
  deck: string
  area: number
  min_price_cents: number
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
const route = useRoute()

const loading = ref(false)
const error = ref<string | null>(null)
const results = ref<CabinItem[]>([])
const routes = ref<RouteItem[]>([])
const cabinTypes = ref<CabinTypeItem[]>([])

const filters = ref({
  keyword: (route.query.keyword as string) || '',
  routeId: Number(route.query.routeId) || 0,
  cabinTypeId: Number(route.query.cabinTypeId) || 0,
})

async function loadOptions() {
  try {
    const routeRes = await request('/cruises', { query: { page: 1, page_size: 100 } })
    const routePayload = routeRes?.data ?? routeRes ?? []
    routes.value = Array.isArray(routePayload) ? routePayload : routePayload?.list ?? []

    const selectedCruiseId = filters.value.routeId || Number(routes.value[0]?.id || 0)
    if (selectedCruiseId > 0) {
      const typeRes = await request('/cabin-types', { query: { cruise_id: selectedCruiseId, page: 1, page_size: 100 } })
      const typePayload = typeRes?.data ?? typeRes ?? {}
      cabinTypes.value = Array.isArray(typePayload) ? typePayload : typePayload?.list ?? []
      if (!filters.value.routeId) {
        filters.value.routeId = selectedCruiseId
      }
    } else {
      cabinTypes.value = []
    }
  } catch (e) {
    console.error('Failed to load options:', e)
  }
}

async function handleSearch() {
  loading.value = true
  error.value = null
  
  try {
    const query: Record<string, unknown> = { page: 1, page_size: 20 }
    
    if (filters.value.keyword.trim()) {
      query.keyword = filters.value.keyword.trim()
    }
    if (filters.value.routeId > 0) {
      query.voyage_id = filters.value.routeId
    }
    if (filters.value.cabinTypeId > 0) {
      query.cabin_type_id = filters.value.cabinTypeId
    }
    
    const cruiseId = Number(query.voyage_id || filters.value.routeId || 0)
    if (cruiseId <= 0) {
      results.value = []
      return
    }

    const res = await request('/cabin-types', {
      query: {
        cruise_id: cruiseId,
        page: query.page,
        page_size: query.page_size,
      },
    })
    const payload = res?.data ?? res ?? {}
    const list = Array.isArray(payload) ? payload : payload?.list ?? []
    const keyword = String(query.keyword || '').toLowerCase()
    results.value = list
      .filter((item: Record<string, any>) => {
        if (!keyword) return true
        const name = String(item.name || '').toLowerCase()
        return name.includes(keyword)
      })
      .map((item: Record<string, any>) => ({
        id: Number(item.id),
        code: item.name || `舱型 #${item.id}`,
        cabin_type_name: item.name || '-',
        deck: '-',
        area: Number(item.area_min || item.area || 0),
        min_price_cents: Number(item.min_price_cents || item.price_cents || 0),
        total: Number(item.inventory || 0),
        locked: 0,
        sold: 0,
        images: item.images,
      }))
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : '搜索失败'
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
  if (filters.value.keyword || filters.value.routeId || filters.value.cabinTypeId) {
    await handleSearch()
  }
})
</script>
