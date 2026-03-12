<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const { request } = useApi()
const router = useRouter()

// ===== 数据模型 =====
interface Route {
  id: number
  name: string
  code: string
  departure_port: string
  arrival_port: string
  stops?: string
  description?: string
  status: number
}

interface Voyage {
  id: number
  route_id?: number
  cruise_id?: number
  code: string
  image_url?: string
  brief_info?: string
  depart_date: string
  return_date: string
  status: number
  first_stop_city?: string
  itinerary_days?: number
}

// ===== 状态管理 =====
const allRoutes = ref<Route[]>([])
const allVoyages = ref<Voyage[]>([])
const filteredRoutes = ref<Route[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

// ===== 搜索过滤器 =====
const filters = ref({
  destination: '',
  departure: '',
  date: ''
})

// ===== 计算聚合数据供下拉框使用 =====
const uniqueDestinations = computed(() => {
  const ports = allRoutes.value.map(r => r.arrival_port).filter(Boolean)
  const voyagePorts = allVoyages.value.map(v => v.first_stop_city).filter(Boolean)
  return Array.from(new Set([...ports, ...voyagePorts])).sort() as string[]
})

const uniqueDepartures = computed(() => {
  const ports = allRoutes.value.map(r => r.departure_port).filter(Boolean)
  return Array.from(new Set(ports)).sort() as string[]
})

const uniqueDates = computed(() => {
  const datesMap = new Map<string, string>()
  allVoyages.value.forEach(v => {
    if (v.depart_date) {
      const d = new Date(v.depart_date)
      if (!isNaN(d.getTime())) {
        const yearMonth = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`
        const label = `${d.getFullYear()}年${d.getMonth() + 1}月`
        datesMap.set(yearMonth, label)
      }
    }
  })
  return Array.from(datesMap.entries())
    .map(([val, lbl]) => ({ value: val, label: lbl }))
    .sort((a, b) => a.value.localeCompare(b.value))
})

// ===== 数据获取逻辑 =====
async function loadSearchData() {
  loading.value = true
  error.value = null
  try {
    // 尝试获取 Routes (如果公共 API 失败则降级到 admin)
    let routesRes = await request('/routes').catch(() => request('/admin/routes').catch(() => null))
    let routesData = routesRes?.data ?? routesRes?.list ?? routesRes ?? []
    
    // 尝试获取 Voyages
    let voyagesRes = await request('/voyages').catch(() => request('/admin/voyages').catch(() => null))
    let voyagesData = voyagesRes?.data ?? voyagesRes?.list ?? voyagesRes ?? []

    if (!Array.isArray(routesData)) routesData = []
    if (!Array.isArray(voyagesData)) voyagesData = []

    allRoutes.value = routesData
    allVoyages.value = voyagesData
    applyFilters()
  } catch (err: any) {
    console.error('Failed to load search data:', err)
    error.value = '加载航线数据失败，请重试。'
  } finally {
    loading.value = false
  }
}

// ===== 过滤逻辑 =====
function applyFilters() {
  // 找出符合日期筛选项的相关 Voyages，从而找到匹配的 routes
  const validRouteOrCruiseIds = new Set<number>()
  if (filters.value.date) {
    allVoyages.value.forEach(v => {
      if (v.depart_date && v.depart_date.startsWith(filters.value.date)) {
        if (v.route_id) validRouteOrCruiseIds.add(v.route_id)
        if (v.cruise_id) validRouteOrCruiseIds.add(v.cruise_id)
      }
    })
  }

  filteredRoutes.value = allRoutes.value.filter(r => {
    // 1. Destination 目的地检查
    if (filters.value.destination) {
       // 检查 route.arrival_port 或者关联 voyage.first_stop_city
       if (r.arrival_port !== filters.value.destination) {
           const hasMatchingStop = allVoyages.value.some(v => 
              (v.route_id === r.id || v.cruise_id === r.id) && v.first_stop_city === filters.value.destination
           )
           if (!hasMatchingStop) return false
       }
    }

    // 2. Departure 出发地检查
    if (filters.value.departure && r.departure_port !== filters.value.departure) {
      return false
    }

    // 3. Dates 日期检查
    if (filters.value.date) {
      // 只有同时获取到 Voyages 和 Routes 时才严格过滤，避免空过滤
      if (allVoyages.value.length > 0 && !validRouteOrCruiseIds.has(r.id)) {
        return false
      }
    }

    // Check status (only show active routes)
    if (r.status !== undefined && r.status === 0) return false

    return true
  })
}

function clearFilters() {
  filters.value = { destination: '', departure: '', date: '' }
  applyFilters()
}

function goRouteDetail(route: Route) {
  // 前端以前的逻辑是点击进入具体的舱型或详情，例如 `/cabins?cruise_id=xxx`
  // NCL 官网上搜索出 routes 后会点击进入航线详情选舱，这里我们利用系统内现有的 `/cabins` 并带上参数
  router.push({ path: '/cabins', query: { cruise_id: route.id, source: 'search' } })
}

onMounted(() => {
  loadSearchData()
})
</script>

<template>
  <div class="search-page w-full min-h-screen pb-16">
    <!-- Hero Banner with Large Image Background -->
    <section class="relative w-full h-[650px] overflow-visible bg-[var(--color-primary)] flex flex-col items-center justify-center text-center text-white mb-24">
      <img src="https://images.unsplash.com/photo-1599640842225-85d111c60e6b?q=80&w=2000&auto=format&fit=crop" class="absolute inset-0 w-full h-full object-cover opacity-70" alt="邮轮背景" />
      
      <!-- Overlay text -->
      <div class="relative z-10 px-4 max-w-4xl mx-auto -mt-16">
        <div class="uppercase tracking-widest text-sm md:text-md mb-2 font-bold drop-shadow-md">开启您的完美假期</div>
        <h1 class="text-4xl md:text-6xl font-bold mb-6 drop-shadow-xl text-white">探寻下一个航海传奇</h1>
        <p class="text-lg md:text-2xl drop-shadow-md font-bold text-white max-w-2xl mx-auto">
          仅限本周末！预订指定航次立享高达 <span class="text-[var(--color-cta)]">$500</span> 船上消费金。
        </p>
      </div>

      <!-- Floating Search Bar -->
      <div class="absolute -bottom-12 md:-bottom-10 left-1/2 -translate-x-1/2 w-[92%] max-w-6xl bg-white rounded-lg shadow-2xl z-20 flex flex-col md:flex-row border border-[var(--web-border)]">
        
        <!-- Destination / Where -->
        <div class="flex-1 p-3 md:p-5 relative text-left border-b md:border-b-0 md:border-r border-[var(--web-border)] hover:bg-[#F8FAFC] transition">
          <label class="block text-[11px] font-bold uppercase text-[var(--web-muted)] mb-1">目的地 Where</label>
          <select v-model="filters.destination" class="w-full text-[var(--color-primary)] font-bold bg-transparent outline-none cursor-pointer appearance-none text-lg">
            <option value="">不限 Any Destination</option>
            <option v-for="dest in uniqueDestinations" :key="dest" :value="dest">{{ dest }}</option>
          </select>
          <div class="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none">
            <svg class="w-5 h-5 text-[var(--web-muted)]" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
          </div>
        </div>

        <!-- Departure -->
        <div class="flex-1 p-3 md:p-5 relative text-left border-b md:border-b-0 md:border-r border-[var(--web-border)] hover:bg-[#F8FAFC] transition">
          <label class="block text-[11px] font-bold uppercase text-[var(--web-muted)] mb-1">出发地 Departure</label>
          <select v-model="filters.departure" class="w-full text-[var(--color-primary)] font-bold bg-transparent outline-none cursor-pointer appearance-none text-lg">
            <option value="">不限 Any Port</option>
            <option v-for="port in uniqueDepartures" :key="port" :value="port">{{ port }}</option>
          </select>
          <div class="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none">
            <svg class="w-5 h-5 text-[var(--web-muted)]" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
          </div>
        </div>

        <!-- Dates -->
        <div class="flex-1 p-3 md:p-5 relative text-left hover:bg-[#F8FAFC] transition">
          <label class="block text-[11px] font-bold uppercase text-[var(--web-muted)] mb-1">日期 Dates</label>
          <select v-model="filters.date" class="w-full text-[var(--color-primary)] font-bold bg-transparent outline-none cursor-pointer appearance-none text-lg">
            <option value="">不限 Any Dates</option>
            <option v-for="date in uniqueDates" :key="date.value" :value="date.value">{{ date.label }}</option>
          </select>
          <div class="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none">
            <svg class="w-5 h-5 text-[var(--web-muted)]" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
          </div>
        </div>

        <!-- Action Button -->
        <button @click="applyFilters" class="bg-[var(--color-action)] text-white font-bold uppercase text-lg px-8 py-5 md:py-0 hover:bg-pink-700 transition flex items-center justify-center min-w-[200px] md:rounded-r-lg shadow-inner">
          <span class="mr-2">查找航线</span>
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path></svg>
        </button>
      </div>
    </section>

    <!-- Results Section -->
    <section class="max-w-7xl mx-auto px-6 py-6 mt-8 md:mt-12">
      <div class="flex justify-between items-end mb-10 border-b-2 border-slate-100 pb-4">
        <div>
          <h2 class="text-3xl font-bold text-[var(--color-primary)] heading-font m-0">精选航线推荐</h2>
          <p class="text-sm text-[var(--web-muted)] mt-1 mb-0">发现世界最美丽的海岸线，体验尊贵服务。</p>
        </div>
        <span class="text-sm font-bold bg-[#F1F5F9] px-3 py-1 rounded text-slate-600 hidden md:block">{{ filteredRoutes.length }} 条航线可选</span>
      </div>

      <!-- Loading / Empty States -->
      <div v-if="loading" class="text-center py-24">
        <div class="inline-block animate-spin rounded-full h-12 w-12 border-4 border-[var(--color-primary)] border-t-transparent mb-4"></div>
        <div class="text-lg font-bold text-[var(--color-primary)]">正在获取最新航线...</div>
      </div>
      
      <div v-else-if="error" class="text-center py-20">
        <p class="text-2xl font-bold text-[var(--web-danger)] mb-2">{{ error }}</p>
        <button @click="loadSearchData" class="btn-primary px-6 py-2 rounded shadow-md mt-4">重试刷新</button>
      </div>

      <div v-else-if="filteredRoutes.length === 0" class="text-center py-24 bg-white rounded-xl shadow-sm border border-slate-100">
        <svg class="w-16 h-16 mx-auto text-slate-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
        <p class="text-2xl font-bold text-[var(--color-primary)] mb-2">未找到匹配航线</p>
        <p class="text-[var(--web-muted)] mb-8">请尝试更改搜索条件或浏览所有航线。</p>
        <button @click="clearFilters" class="px-8 py-3 rounded border border-[var(--color-primary)] text-[var(--color-primary)] font-bold hover:bg-[var(--color-primary)] hover:text-white transition duration-300 uppercase text-sm">清除过滤条件</button>
      </div>
      
      <!-- Route Cards Grid -->
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <div v-for="route in filteredRoutes" :key="route.id" @click="goRouteDetail(route)" class="group flex flex-col bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.06)] hover:shadow-[0_12px_30px_rgba(0,0,0,0.12)] transition-all duration-300 cursor-pointer overflow-hidden border border-slate-100">
          
          <!-- Image Header -->
          <div class="relative h-60 overflow-hidden bg-slate-100">
            <img :src="'https://images.unsplash.com/photo-1548574505-5e239809ee19?q=80&w=800&auto=format&fit=crop'" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-[800ms] ease-out" alt="航线配图" />
            <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/10 to-transparent opacity-80 group-hover:opacity-90 transition-opacity"></div>
            
            <div class="absolute top-4 left-4 flex flex-col gap-2">
              <span class="bg-white/95 backdrop-blur text-[var(--color-primary)] text-xs font-bold py-1 px-2.5 rounded shadow-sm border border-white/40 tracking-wide">
                {{ route.departure_port || '随机港口' }} 出发
              </span>
            </div>
            
            <!-- Bottom stats overlay -->
            <div class="absolute bottom-4 left-4 right-4 text-white">
              <div class="text-xs uppercase tracking-widest font-bold opacity-80 mb-1 flex items-center gap-1">
                <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"></path></svg>
                目的地: {{ route.arrival_port || '未知' }}
              </div>
            </div>
          </div>
          
          <!-- Route Content -->
          <div class="p-6 flex-1 flex flex-col">
            <h3 class="text-2xl font-bold text-[var(--color-primary)] mb-3 leading-snug group-hover:text-[var(--color-action)] transition-colors">{{ route.name || '经典航线' }}</h3>
            
            <div class="flex items-center gap-2 mb-4">
              <span class="block w-6 h-[2px] bg-[var(--color-cta)]"></span>
              <span class="text-xs font-bold text-[var(--web-muted)] uppercase tracking-wider">航线代码: {{ route.code }}</span>
            </div>

            <p class="text-[0.95rem] text-slate-600 line-clamp-2 mt-1 mb-6 flex-1">{{ route.description || route.stops || '加入我们的梦想之旅，探索绝美目的地的难忘体验。' }}</p>
            
            <div class="flex items-end justify-between pt-5 border-t border-slate-100">
               <div>
                  <span class="block text-[11px] uppercase font-bold text-slate-400 mb-1 tracking-widest">起价</span>
                  <div class="flex items-baseline text-[var(--color-primary)] line-height-none">
                    <span class="text-lg font-bold mr-1">¥</span>
                    <span class="text-3xl font-extrabold heading-font">3,999</span>
                  </div>
               </div>
               <div class="w-10 h-10 rounded-full bg-slate-50 flex items-center justify-center text-[var(--color-primary)] group-hover:bg-[var(--color-action)] group-hover:text-white transition-colors">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path></svg>
               </div>
            </div>
          </div>

        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
/* Scoped overrides if necessary */
select option {
  color: var(--color-primary);
  font-weight: 500;
}
</style>
