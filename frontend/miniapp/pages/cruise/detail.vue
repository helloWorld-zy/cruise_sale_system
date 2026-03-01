<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { request } from '../../src/utils/request'

const props = defineProps<{ cruiseId?: number | string }>()

const loading = ref(false)
const error = ref('')
const detail = ref<Record<string, any> | null>(null)
const cabinTypes = ref<Record<string, any>[]>([])
const facilities = ref<Record<string, any>[]>([])
const routes = ref<Record<string, any>[]>([])
const activeFacilityTab = ref(0)

const gallery = computed(() => {
  if (!detail.value) return ['https://picsum.photos/seed/mini-detail-default/1000/700']
  const images = Array.isArray(detail.value.images) ? detail.value.images : []
  if (images.length > 0) return images.map((item: any) => item.url || item)
  const id = resolveCruiseId()
  return [
    `https://picsum.photos/seed/mini-detail-${id}-1/1000/700`,
    `https://picsum.photos/seed/mini-detail-${id}-2/1000/700`,
    `https://picsum.photos/seed/mini-detail-${id}-3/1000/700`,
  ]
})

const facilityTabs = computed(() => {
  const all = [{ id: 0, name: '全部' }]
  const unique = new Map<number, string>()
  facilities.value.forEach((item) => {
    const cid = Number(item.category_id || 0)
    if (cid > 0 && !unique.has(cid)) unique.set(cid, item.category_name || `分类${cid}`)
  })
  unique.forEach((name, id) => all.push({ id, name }))
  return all
})

const filteredFacilities = computed(() => {
  if (activeFacilityTab.value === 0) return facilities.value
  return facilities.value.filter((item) => Number(item.category_id) === activeFacilityTab.value)
})

function resolveCruiseId() {
  return Number(props.cruiseId ?? 0)
}

async function loadAll() {
  const id = resolveCruiseId()
  if (!id) {
    error.value = '缺少 cruiseId 参数'
    return
  }
  loading.value = true
  error.value = ''
  try {
    const [detailRes, typeRes, facilityRes, routeRes] = await Promise.all([
      request(`/cruises/${id}`),
      request(`/cabin-types?cruise_id=${id}&page=1&page_size=20`),
      request(`/facilities?cruise_id=${id}`),
      request('/routes'),
    ])

    detail.value = (detailRes as any)?.data ?? detailRes ?? null
    const typePayload = (typeRes as any)?.data ?? typeRes ?? {}
    cabinTypes.value = Array.isArray(typePayload) ? typePayload : typePayload?.list ?? []
    const facilityPayload = (facilityRes as any)?.data ?? facilityRes ?? []
    facilities.value = Array.isArray(facilityPayload) ? facilityPayload : facilityPayload?.list ?? []
    const routePayload = (routeRes as any)?.data ?? routeRes ?? []
    const source = Array.isArray(routePayload) ? routePayload : routePayload?.list ?? []
    routes.value = source.slice(0, 6).map((item: Record<string, any>) => ({
      id: item.id || Math.random().toString(36).substr(2, 9),
      date: item.departure_date || item.date || '-',
      name: item.name || item.route_name || '-',
      price: Math.round(Number(item.min_price_cents || item.price_cents || 0) / 100) || '-',
    }))
  } catch (e: any) {
    error.value = e?.message ?? '加载邮轮详情失败'
  } finally {
    loading.value = false
  }
}

onMounted(loadAll)
</script>

<template>
  <view class="page">
    <view v-if="loading" class="state">Loading...</view>
    <view v-else-if="error" class="state error">{{ error }}</view>
    <view v-else-if="detail" class="content">
      <swiper class="hero" indicator-dots circular>
        <swiper-item v-for="img in gallery" :key="img">
          <image class="hero-img" :src="img" mode="aspectFill" />
        </swiper-item>
      </swiper>

      <view class="section card">
        <text class="title">{{ detail.name || '-' }}</text>
        <scroll-view scroll-x class="metric-strip">
          <view class="metric-card">⚓<text>{{ detail.tonnage || '-' }}吨</text></view>
          <view class="metric-card">👤<text>{{ detail.passenger_capacity || '-' }}人</text></view>
          <view class="metric-card">🧭<text>{{ detail.length || '-' }}m</text></view>
          <view class="metric-card">🛳<text>{{ detail.deck_count || '-' }}层</text></view>
        </scroll-view>
      </view>

      <view class="section card">
        <text class="section-title">舱房类型</text>
        <view v-for="item in cabinTypes" :key="item.id" class="row">
          <image class="thumb" mode="aspectFill" :src="`https://picsum.photos/seed/mini-cabin-${item.id}/200/200`" />
          <view class="mid">
            <text class="row-title">{{ item.name || '-' }}</text>
            <text class="row-meta">{{ item.area_min || item.area || '-' }}m2 · {{ item.max_capacity || item.capacity || '-' }}人</text>
          </view>
          <text class="price">¥{{ Math.round(Number(item.min_price_cents || item.price_cents || 0) / 100) || '-' }}</text>
        </view>
      </view>

      <view class="section card">
        <text class="section-title">设施导览</text>
        <scroll-view scroll-x class="tabs">
          <text
            v-for="tab in facilityTabs"
            :key="tab.id"
            class="tab"
            :class="activeFacilityTab === tab.id ? 'tab-active' : ''"
            @click="activeFacilityTab = tab.id"
          >
            {{ tab.name }}
          </text>
        </scroll-view>
        <view class="facility-grid">
          <view v-for="fac in filteredFacilities" :key="fac.id" class="facility-item">
            <text class="facility-name">{{ fac.name }}</text>
            <text class="facility-tag" :class="fac.extra_charge ? 'tag-charge' : 'tag-free'">{{ fac.extra_charge ? '收费' : '免费' }}</text>
          </view>
        </view>
      </view>

      <view class="section card bottom-space">
        <text class="section-title">关联航线</text>
        <view v-for="(item, idx) in routes" :key="item.id || idx" class="route-row">
          <view>
            <text class="route-date">{{ item.date }}</text>
            <text class="route-name">{{ item.name }}</text>
          </view>
          <text class="route-price">¥{{ item.price }}</text>
        </view>
      </view>

      <view class="bottom-bar">
        <text class="icon-btn">♡ 收藏</text>
        <text class="icon-btn">↗ 分享</text>
        <button class="cta">查看航线</button>
      </view>
    </view>
  </view>
</template>

<style scoped>
.page {
  min-height: 100vh;
  background: #f6f7fb;
}
.hero {
  height: 480rpx;
}
.hero-img {
  width: 100%;
  height: 480rpx;
}
.content {
  padding: 20rpx;
}
.card {
  background: #fff;
  border-radius: 16rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.06);
  padding: 20rpx;
}
.section {
  margin-top: 16rpx;
}
.title {
  display: block;
  font-size: 34rpx;
  font-weight: 700;
  color: #172533;
}
.metric-strip {
  margin-top: 14rpx;
  white-space: nowrap;
}
.metric-card {
  display: inline-flex;
  flex-direction: column;
  width: 140rpx;
  margin-right: 12rpx;
  padding: 12rpx;
  border-radius: 12rpx;
  background: #f8f9fa;
  font-size: 22rpx;
  color: #4c5d70;
}
.section-title {
  display: block;
  margin-bottom: 12rpx;
  font-size: 28rpx;
  font-weight: 600;
}
.row {
  display: flex;
  align-items: center;
  margin-bottom: 14rpx;
}
.thumb {
  width: 100rpx;
  height: 100rpx;
  border-radius: 12rpx;
  margin-right: 12rpx;
}
.mid {
  flex: 1;
}
.row-title {
  display: block;
  font-size: 26rpx;
  color: #1f2d3a;
}
.row-meta {
  display: block;
  margin-top: 4rpx;
  font-size: 22rpx;
  color: #7a8796;
}
.price {
  color: #ff6b6b;
  font-size: 28rpx;
  font-weight: 700;
}
.tabs {
  white-space: nowrap;
  margin-bottom: 10rpx;
}
.tab {
  display: inline-block;
  margin-right: 20rpx;
  padding-bottom: 8rpx;
  font-size: 24rpx;
  color: #6b7c8f;
}
.tab-active {
  color: #ff6b6b;
  border-bottom: 4rpx solid #ff6b6b;
}
.facility-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10rpx;
}
.facility-item {
  border-radius: 12rpx;
  background: #f8f9fa;
  padding: 12rpx;
}
.facility-name {
  display: block;
  font-size: 24rpx;
  color: #25384c;
}
.facility-tag {
  display: inline-block;
  margin-top: 6rpx;
  border-radius: 999rpx;
  padding: 4rpx 12rpx;
  font-size: 20rpx;
}
.tag-free {
  background: #ecfdf3;
  color: #0f8f4b;
}
.tag-charge {
  background: #fff7ed;
  color: #c66b06;
}
.route-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10rpx 0;
}
.route-date {
  display: block;
  font-size: 22rpx;
  color: #7a8796;
}
.route-name {
  display: block;
  margin-top: 2rpx;
  font-size: 24rpx;
  color: #1f2d3a;
}
.route-price {
  color: #ff6b6b;
  font-size: 28rpx;
  font-weight: 700;
}
.bottom-space {
  margin-bottom: 120rpx;
}
.bottom-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  gap: 16rpx;
  background: #fff;
  padding: 16rpx 20rpx;
  padding-bottom: calc(16rpx + env(safe-area-inset-bottom));
  box-shadow: 0 -2rpx 12rpx rgba(0, 0, 0, 0.08);
}
.icon-btn {
  font-size: 24rpx;
  color: #5e6f82;
}
.cta {
  flex: 1;
  border-radius: 44rpx;
  height: 88rpx;
  line-height: 88rpx;
  background: #ff6b6b;
  color: #fff;
  font-size: 28rpx;
}
.state {
  padding: 30rpx;
  color: #6b7c8f;
}
.error {
  color: #d7415f;
}
</style>
