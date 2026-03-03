<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { request } from '../../src/utils/request'

// 通过 props 接收页面参数，便于运行时与测试场景复用。
const props = defineProps<{ cabinSkuId?: number | string; preview?: boolean }>()
const loading = ref(false)
const error = ref('')
const detail = ref<any>(null)
const prices = ref<any[]>([])
const inventory = ref<any>(null)
const selectedDate = ref('')
const currentImage = ref(0)

function resolveSkuId() {
  return Number(props.cabinSkuId ?? 0)
}

async function loadDetail() {
  if (props.preview) {
    detail.value = {
      id: resolveSkuId() || 1,
      name: '天际阳台房',
      bed_type: '大床',
      position: '中层甲板',
      orientation: '海景',
      has_window: true,
      has_balcony: true,
      amenities: '独立卫浴,迷你吧,高速WiFi,观景阳台',
      price_cents: 268000,
    }
    prices.value = [
      { date: '2026-03-05', price_cents: 268000 },
      { date: '2026-03-06', price_cents: 288000 },
      { date: '2026-03-07', price_cents: 276000 },
    ]
    inventory.value = { total: 12, locked: 2, sold: 6, available: 4 }
    return
  }
  const skuId = resolveSkuId()
  if (!skuId) {
    error.value = '缺少 cabinSkuId 参数'
    return
  }
  loading.value = true
  error.value = ''
  try {
    const detailRes = await request(`/cabins/${skuId}`)
    detail.value = detailRes?.data ?? detailRes
    const pricesRes = await request(`/cabins/${skuId}/prices`)
    prices.value = pricesRes?.data ?? pricesRes ?? []
    const inventoryRes = await request(`/cabins/${skuId}/inventory`)
    inventory.value = inventoryRes?.data ?? inventoryRes ?? null
  } catch (e: any) {
    error.value = e?.message ?? '加载舱房失败'
  } finally {
    loading.value = false
  }
}

const gallery = computed(() => {
  const images = Array.isArray(detail.value?.images) ? detail.value.images : []
  if (images.length > 0) return images.map((item: any) => item.url || item)
  const skuId = resolveSkuId()
  return [
    `https://picsum.photos/seed/mini-cabin-${skuId}-1/1000/700`,
    `https://picsum.photos/seed/mini-cabin-${skuId}-2/1000/700`,
    `https://picsum.photos/seed/mini-cabin-${skuId}-3/1000/700`,
  ]
})

const tags = computed(() => {
  const source = [
    detail.value?.bed_type && `床型:${detail.value.bed_type}`,
    detail.value?.position && `位置:${detail.value.position}`,
    detail.value?.orientation && `朝向:${detail.value.orientation}`,
    `窗户:${detail.value?.has_window ? '有' : '无'}`,
    `阳台:${detail.value?.has_balcony ? '有' : '无'}`,
  ]
  return source.filter(Boolean)
})

function normalizedDate(raw: string) {
  return raw ? raw.slice(0, 10) : ''
}

const dateItems = computed(() => {
  return prices.value.map((item) => normalizedDate(item.date)).filter(Boolean)
})

const selectedPrice = computed(() => {
  const date = selectedDate.value
  const row = prices.value.find((item) => normalizedDate(item.date) === date)
  return row ? Number(row.price_cents || row.price || 0) : 0
})

const availableLabel = computed(() => {
  const total = Number(inventory.value?.total || 0)
  const locked = Number(inventory.value?.locked || 0)
  const sold = Number(inventory.value?.sold || 0)
  const available = inventory.value?.available !== undefined ? Number(inventory.value.available) : total - locked - sold
  if (available <= 3) return `仅剩${available}间`
  return '库存充足'
})

const availableClass = computed(() => (availableLabel.value.includes('仅剩') ? 'stock-low' : 'stock-ok'))

onMounted(() => {
  loadDetail().then(() => {
    if (dateItems.value.length > 0) selectedDate.value = dateItems.value[0]
  })
})
</script>

<template>
  <view class="page">
    <text v-if="loading" class="hint">Loading...</text>
    <text v-else-if="error" class="error">{{ error }}</text>
    <view v-else-if="detail" class="panel">
      <swiper class="hero" circular indicator-dots @change="(e: any) => (currentImage = e.detail.current)">
        <swiper-item v-for="img in gallery" :key="img">
          <image :src="img" class="hero-img" mode="aspectFill" />
        </swiper-item>
      </swiper>
      <text class="pager">{{ currentImage + 1 }}/{{ gallery.length }}</text>

      <view class="price-line">
        <text class="price">¥{{ Math.round((detail.price_cents || detail.price || selectedPrice || 0) / 100) || '-' }}起</text>
        <text class="tax">含税费</text>
      </view>

      <view class="tag-flow">
        <text v-for="tag in tags" :key="String(tag)" class="tag">{{ tag }}</text>
      </view>

      <view class="facility-grid">
        <text v-for="item in (detail.amenities || '').split(',').filter(Boolean)" :key="item" class="facility">{{ item }}</text>
      </view>

      <scroll-view scroll-x class="date-strip">
        <text
          v-for="date in dateItems"
          :key="date"
          class="date-item"
          :class="selectedDate === date ? 'date-active' : ''"
          @click="selectedDate = date"
        >
          {{ date.slice(5) }}
        </text>
      </scroll-view>

      <view class="price-panel">
        <text>日期：{{ selectedDate || '-' }}</text>
        <text>价格：¥{{ selectedPrice > 0 ? Math.round(selectedPrice / 100) : '-' }}</text>
      </view>

      <text class="stock" :class="availableClass">{{ availableLabel }}</text>

      <view class="action-bar">
        <text class="action-icon">☎ 客服</text>
        <button class="book-btn">立即预订</button>
      </view>
    </view>
  </view>
</template>

<style scoped>
.page {
  min-height: 100vh;
  background:
    radial-gradient(circle at 8% 0, #dbe9f5 0, transparent 30%),
    linear-gradient(180deg, #f3f8fb 0%, #edf3f7 100%);
  padding: 22rpx 20rpx 138rpx;
}
.panel {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
  padding: 20rpx;
  border-radius: 24rpx;
  background: #fff;
  border: 1rpx solid #d4e0ea;
  box-shadow: 0 16rpx 36rpx rgba(16, 47, 72, 0.12);
}
.hero {
  height: 420rpx;
  border-radius: 18rpx;
  overflow: hidden;
}
.hero-img {
  width: 100%;
  height: 420rpx;
}
.pager {
  align-self: flex-end;
  margin-top: -8rpx;
  font-size: 22rpx;
  color: #6b7c8f;
}
.price-line {
  display: flex;
  align-items: baseline;
  gap: 10rpx;
}
.price {
  color: #113d5c;
  font-size: 46rpx;
  font-weight: 700;
}
.tax {
  color: #6b8095;
  font-size: 22rpx;
}
.tag-flow {
  display: flex;
  flex-wrap: wrap;
  gap: 12rpx;
}
.tag {
  background: #eef5fb;
  color: #2f5f85;
  border-radius: 8rpx;
  padding: 8rpx 16rpx;
  font-size: 24rpx;
}
.facility-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10rpx;
}
.facility {
  background: #f5f9fc;
  border-radius: 10rpx;
  padding: 10rpx;
  font-size: 24rpx;
  color: #3f5569;
}
.date-strip {
  white-space: nowrap;
}
.date-item {
  display: inline-block;
  margin-right: 12rpx;
  border-radius: 999rpx;
  border: 1rpx solid #d7dfe8;
  padding: 8rpx 18rpx;
  color: #637689;
  font-size: 22rpx;
}
.date-active {
  border-color: #113d5c;
  background: #113d5c;
  color: #fff;
}
.price-panel {
  border-radius: 14rpx;
  background: #f4f8fb;
  padding: 12rpx;
  display: flex;
  flex-direction: column;
  gap: 6rpx;
  color: #425568;
  font-size: 24rpx;
}
.stock {
  font-size: 26rpx;
  font-weight: 600;
}
.stock-low {
  color: #c53f57;
}
.stock-ok {
  color: #0f8a60;
}
.action-bar {
  display: flex;
  align-items: center;
  gap: 14rpx;
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  background: #ffffff;
  padding: 14rpx 18rpx;
  padding-bottom: calc(14rpx + env(safe-area-inset-bottom));
  box-shadow: 0 -8rpx 26rpx rgba(16, 47, 72, 0.16);
}
.action-icon {
  font-size: 24rpx;
  color: #5e6f82;
}
.book-btn {
  flex: 1;
  border-radius: 44rpx;
  height: 88rpx;
  line-height: 88rpx;
  background: linear-gradient(135deg, #0f3d5c, #1f5f86);
  color: #fff;
  font-size: 28rpx;
}

.hint {
  color: #5a7190;
}

.error {
  color: #d13e5b;
}
</style>
