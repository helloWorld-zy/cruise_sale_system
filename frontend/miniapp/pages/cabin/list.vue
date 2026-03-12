<template>
  <view class="min-h-screen bg-background pb-8 overflow-x-hidden">
    <!-- Header -->
    <div class="relative z-0 bg-primary/95 pt-12 pb-16 px-6 shadow-md" style="border-bottom-left-radius: 40px; border-bottom-right-radius: 40px;">
      <h1 class="text-3xl font-heading font-bold text-white mb-2 relative z-10">精选舱房</h1>
      <p class="text-sm text-white/90 relative z-10">按航线、日期与港口筛选，快速对比热门舱型。</p>
      <!-- Decorative circles -->
      <div class="absolute top-0 right-0 w-32 h-32 bg-white/10 rounded-full blur-2xl -mr-10 -mt-10"></div>
    </div>

    <view class="px-4 -mt-10 relative z-10 flex flex-col gap-4">
      <view class="bg-white rounded-2xl p-4 shadow-md border border-gray-50 flex flex-col gap-3">
        <input v-model="filters.route" class="bg-gray-50 rounded-xl px-4 py-2.5 text-[14px] w-full outline-none focus:ring-2 focus:ring-secondary transition-smooth border border-transparent focus:border-secondary" placeholder="按航线筛选" />
        <input v-model="filters.date" class="bg-gray-50 rounded-xl px-4 py-2.5 text-[14px] w-full outline-none focus:ring-2 focus:ring-secondary transition-smooth border border-transparent focus:border-secondary" placeholder="按日期筛选(YYYY-MM-DD)" />
        <input v-model="filters.port" class="bg-gray-50 rounded-xl px-4 py-2.5 text-[14px] w-full outline-none focus:ring-2 focus:ring-secondary transition-smooth border border-transparent focus:border-secondary" placeholder="按出发港筛选" />
        <view class="flex flex-wrap gap-2 mt-1">
          <button class="px-3.5 py-1.5 border-0 rounded-full text-[13px] font-medium transition-smooth cursor-pointer" :class="filters.sortBy === 'default' ? 'bg-primary text-white shadow-sm hover:-translate-y-0.5' : 'bg-gray-100 text-gray-500 hover:bg-gray-200'" @click="filters.sortBy = 'default'">默认</button>
          <button class="px-3.5 py-1.5 border-0 rounded-full text-[13px] font-medium transition-smooth cursor-pointer" :class="filters.sortBy === 'price_asc' ? 'bg-primary text-white shadow-sm hover:-translate-y-0.5' : 'bg-gray-100 text-gray-500 hover:bg-gray-200'" @click="filters.sortBy = 'price_asc'">价格升序</button>
          <button class="px-3.5 py-1.5 border-0 rounded-full text-[13px] font-medium transition-smooth cursor-pointer" :class="filters.sortBy === 'price_desc' ? 'bg-primary text-white shadow-sm hover:-translate-y-0.5' : 'bg-gray-100 text-gray-500 hover:bg-gray-200'" @click="filters.sortBy = 'price_desc'">价格降序</button>
          <button class="px-3.5 py-1.5 border-0 rounded-full text-[13px] font-medium transition-smooth cursor-pointer" :class="filters.sortBy === 'area_desc' ? 'bg-primary text-white shadow-sm hover:-translate-y-0.5' : 'bg-gray-100 text-gray-500 hover:bg-gray-200'" @click="filters.sortBy = 'area_desc'">面积降序</button>
        </view>
      </view>

      <text v-if="loading" class="text-center text-text mt-12 block font-medium">Loading...</text>
      <text v-else-if="error" class="text-center text-red-500 mt-12 block">{{ error }}</text>
      <text v-else-if="visibleCabins.length === 0" class="text-center text-gray-400 mt-12 block">暂无可预订舱位</text>

      <view v-else class="flex flex-col gap-4 mt-2 pb-6">
        <CabinCard v-for="(item, index) in visibleCabins" :key="String(item.id || item.code || `row-${index}`)" :item="item" />
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import CabinCard from '../../components/CabinCard.vue'
import { request } from '../../src/utils/request'

type CabinListItem = {
  id?: number
  code?: string
  name?: string
  cabin_type_name?: string
  bed_type?: string
  area?: number
  amenities?: string
  price_cents?: number
  price?: number
  route_name?: string
  departure_port?: string
  travel_date?: string
  total?: number
  locked?: number
  sold?: number
  available?: number
}

const cabins = ref<CabinListItem[]>([])
const props = defineProps<{ preview?: boolean }>()
const loading = ref(false)
const error = ref('')
const filters = ref({
  route: '',
  date: '',
  port: '',
  sortBy: 'default' as 'default' | 'price_asc' | 'price_desc' | 'area_desc',
})

const visibleCabins = computed(() => {
  let rows = [...cabins.value]
  const route = filters.value.route.trim().toLowerCase()
  const date = filters.value.date.trim().toLowerCase()
  const port = filters.value.port.trim().toLowerCase()

  if (route) {
    rows = rows.filter((item) => String(item.route_name || '').toLowerCase().includes(route))
  }
  if (date) {
    rows = rows.filter((item) => String(item.travel_date || '').toLowerCase().includes(date))
  }
  if (port) {
    rows = rows.filter((item) => String(item.departure_port || '').toLowerCase().includes(port))
  }

  if (filters.value.sortBy === 'price_asc') {
    rows.sort((a, b) => Number(a.price_cents || a.price || 0) - Number(b.price_cents || b.price || 0))
  } else if (filters.value.sortBy === 'price_desc') {
    rows.sort((a, b) => Number(b.price_cents || b.price || 0) - Number(a.price_cents || a.price || 0))
  } else if (filters.value.sortBy === 'area_desc') {
    rows.sort((a, b) => Number(b.area || 0) - Number(a.area || 0))
  }

  return rows
})

async function loadCabins() {
  if (props.preview) {
    cabins.value = [
      {
        id: 1001,
        name: '云海阳台房',
        cabin_type_name: 'Balcony',
        bed_type: '大床',
        area: 32,
        amenities: '海景阳台,独立卫浴,迷你吧,WiFi',
        price_cents: 236000,
        route_name: '上海-济州-福冈',
        departure_port: '上海',
        travel_date: '2026-03-05',
        total: 20,
        sold: 14,
        locked: 3,
        available: 3,
      },
      {
        id: 1002,
        name: '星辰海景房',
        cabin_type_name: 'Ocean View',
        bed_type: '双床',
        area: 28,
        amenities: '景观窗,沙发区,WiFi',
        price_cents: 198000,
        route_name: '深圳-岘港-芽庄',
        departure_port: '深圳',
        travel_date: '2026-03-11',
        total: 24,
        sold: 12,
        locked: 2,
        available: 10,
        price: undefined,
      },
    ]
    return
  }
  loading.value = true
  error.value = ''
  try {
    const res: any = await request('/cabins', {
      data: {
        page: 1,
        page_size: 20,
      },
    })
    const payload = res?.data ?? res
    cabins.value = Array.isArray(payload?.list) ? payload.list : Array.isArray(payload) ? payload : []
  } catch (e: any) {
    error.value = e?.message || '舱房加载失败，请稍后重试'
    cabins.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadCabins()
})
</script>
