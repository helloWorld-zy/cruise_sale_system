<template>
  <view class="page">
    <text class="title">精选舱房</text>
    <text class="subtitle">按航线、日期与港口筛选，快速对比热门舱型。</text>

    <view class="filters">
      <input v-model="filters.route" class="filter-input" placeholder="按航线筛选" />
      <input v-model="filters.date" class="filter-input" placeholder="按日期筛选(YYYY-MM-DD)" />
      <input v-model="filters.port" class="filter-input" placeholder="按出发港筛选" />
      <view class="sort-row">
        <text class="sort-item" :class="filters.sortBy === 'default' ? 'sort-active' : ''" @click="filters.sortBy = 'default'">默认</text>
        <text class="sort-item" :class="filters.sortBy === 'price_asc' ? 'sort-active' : ''" @click="filters.sortBy = 'price_asc'">价格升序</text>
        <text class="sort-item" :class="filters.sortBy === 'price_desc' ? 'sort-active' : ''" @click="filters.sortBy = 'price_desc'">价格降序</text>
        <text class="sort-item" :class="filters.sortBy === 'area_desc' ? 'sort-active' : ''" @click="filters.sortBy = 'area_desc'">面积降序</text>
      </view>
    </view>

    <text v-if="loading" class="hint">Loading...</text>
    <text v-else-if="error" class="error">{{ error }}</text>
    <text v-else-if="visibleCabins.length === 0" class="hint">暂无可预订舱位</text>

    <view v-else class="list">
      <CabinCard v-for="(item, index) in visibleCabins" :key="String(item.id || item.code || `row-${index}`)" :item="item" />
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

<style scoped>
.page {
  min-height: 100vh;
  padding: 32rpx;
  background: #f5f7fa;
  position: relative;
}

.page::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 400rpx;
  background: linear-gradient(135deg, #0cebeb 0%, #20e3b2 50%, #29ffc6 100%);
  border-bottom-left-radius: 40rpx;
  border-bottom-right-radius: 40rpx;
  z-index: 0;
}

.title {
  position: relative;
  z-index: 1;
  display: block;
  margin-bottom: 8rpx;
  font-size: 48rpx;
  font-weight: 800;
  color: #fff;
}

.subtitle {
  position: relative;
  z-index: 1;
  display: block;
  margin-bottom: 24rpx;
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.9);
}

.filters {
  position: relative;
  z-index: 1;
  margin-bottom: 24rpx;
  display: flex;
  flex-direction: column;
  gap: 16rpx;
  padding: 24rpx;
  border-radius: 32rpx;
  background: #fff;
  border: none;
  box-shadow: 0 12rpx 32rpx rgba(0, 0, 0, 0.05);
}

.filter-input {
  border: none;
  background: #f5f7fa;
  border-radius: 20rpx;
  padding: 16rpx 24rpx;
  font-size: 26rpx;
}

.sort-row {
  display: flex;
  gap: 12rpx;
  flex-wrap: wrap;
  margin-top: 8rpx;
}

.sort-item {
  border-radius: 999rpx;
  border: none;
  background: #f5f7fa;
  padding: 10rpx 20rpx;
  font-size: 24rpx;
  color: #666;
  font-weight: 600;
}

.sort-active {
  background: #0cebeb;
  color: #fff;
}

.list {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  gap: 24rpx;
  padding-bottom: 40rpx;
}

.hint {
  position: relative;
  z-index: 1;
  color: #fff;
  font-size: 28rpx;
  text-align: center;
  display: block;
  margin-top: 40rpx;
}

.error {
  position: relative;
  z-index: 1;
  color: #ffeeee;
  font-size: 28rpx;
  text-align: center;
  display: block;
  margin-top: 40rpx;
}
</style>
