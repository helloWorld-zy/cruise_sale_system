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
  border-radius: 22rpx;
  overflow: hidden;
  background: #fff;
  border: 1rpx solid #d5e1ea;
  box-shadow: 0 14rpx 32rpx rgba(20, 56, 84, 0.14);
}

.cover {
  width: 100%;
  height: 220rpx;
}

.body {
  display: flex;
  flex-direction: column;
  gap: 12rpx;
  padding: 20rpx;
}

.name {
  font-size: 34rpx;
  font-weight: 700;
  color: #15324b;
}

.desc {
  color: #5a7188;
}

.price {
  font-size: 32rpx;
  font-weight: 700;
  color: #0f3d5c;
}

.meta-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12rpx;
}

.stock {
  font-size: 22rpx;
  font-weight: 600;
}

.stock-ok {
  color: #0f8a60;
}

.stock-warning {
  color: #b27722;
}

.stock-danger {
  color: #c53f57;
}

.stock-soldout {
  color: #64748b;
}
</style>
