<template>
  <view class="card">
    <image :src="cover" class="cover" mode="aspectFill" />
    <view class="body">
      <text class="name">{{ title }}</text>
      <text class="desc">{{ description }}</text>
      <view class="meta-row">
        <text class="price">{{ priceLabel }}</text>
        <text class="stock" :class="stockClass">{{ stockLabel }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

type CabinCardItem = {
  id?: number
  code?: string
  name?: string
  title?: string
  cabin_type_name?: string
  cabin_type_code?: string
  bed_type?: string
  area?: number
  amenities?: string
  price_cents?: number
  price?: number
  total?: number
  locked?: number
  sold?: number
  available?: number
}

const props = defineProps<{ item: CabinCardItem }>()

const title = computed(() => {
  return (
    props.item.name ||
    props.item.title ||
    props.item.cabin_type_name ||
    props.item.cabin_type_code ||
    props.item.code ||
    `舱位#${props.item.id ?? '-'}`
  )
})

const description = computed(() => {
  const feature = props.item.bed_type ? `床型 ${props.item.bed_type}` : ''
  const area = props.item.area ? `${props.item.area}㎡` : ''
  const amenity = props.item.amenities ? String(props.item.amenities).split(',')[0] : ''
  return [feature, area, amenity].filter(Boolean).join(' · ') || '配置待补充'
})

const priceLabel = computed(() => {
  const cents = Number(props.item.price_cents || props.item.price || 0)
  if (cents <= 0) return '价格待定'
  return `￥${Math.round(cents / 100)} / 晚`
})

const cover = computed(() => `https://picsum.photos/seed/mini-cabin-card-${props.item.id || props.item.code || 'x'}/900/520`)

const stockCount = computed(() => {
  if (props.item.available !== undefined) return Number(props.item.available)
  return Number(props.item.total || 0) - Number(props.item.locked || 0) - Number(props.item.sold || 0)
})

const stockPercent = computed(() => {
  const total = Number(props.item.total || 0)
  if (total <= 0) return 0
  return Math.max(0, Math.min(100, Math.round((stockCount.value / total) * 100)))
})

const stockLabel = computed(() => {
  if (stockCount.value <= 0) return '已售罄'
  if (stockPercent.value < 20) return '即将售罄'
  if (stockPercent.value < 50) return '库存紧张'
  return '库存充足'
})

const stockClass = computed(() => {
  if (stockCount.value <= 0) return 'stock-soldout'
  if (stockPercent.value < 20) return 'stock-danger'
  if (stockPercent.value < 50) return 'stock-warning'
  return 'stock-ok'
})
</script>

<style scoped>
.card {
  border-radius: 32rpx;
  overflow: hidden;
  background: #fff;
  border: none;
  box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.05);
}

.cover {
  width: 100%;
  height: 280rpx;
}

.body {
  display: flex;
  flex-direction: column;
  gap: 12rpx;
  padding: 30rpx;
}

.name {
  font-size: 36rpx;
  font-weight: 700;
  color: #222;
}

.desc {
  font-size: 26rpx;
  color: #888;
}

.price {
  font-size: 38rpx;
  font-weight: 800;
  color: #ff6b6b;
}

.meta-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12rpx;
  margin-top: 12rpx;
}

.stock {
  font-size: 22rpx;
  font-weight: 600;
  padding: 6rpx 16rpx;
  border-radius: 999rpx;
}

.stock-ok {
  background: #ecfdf3;
  color: #027a48;
}

.stock-warning {
  background: #fff7ed;
  color: #c4320a;
}

.stock-danger {
  background: #fef3f2;
  color: #b42318;
}

.stock-soldout {
  background: #f1f5f9;
  color: #64748b;
}
</style>
