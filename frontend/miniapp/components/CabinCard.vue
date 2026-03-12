<template>
  <view class="bg-white rounded-2xl overflow-hidden shadow-md transition-smooth hover:shadow-lg hover:-translate-y-1 cursor-pointer border border-transparent group">
    <div class="w-full h-36 overflow-hidden relative">
      <img :src="cover" class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105" alt="Cabin Cover" />
    </div>
    <view class="p-4 flex flex-col gap-1.5">
      <h3 class="text-base font-bold text-text mb-0.5">{{ title }}</h3>
      <p class="text-[13px] text-gray-500 mb-1 leading-relaxed">{{ description }}</p>
      <view class="flex items-center justify-between mt-1">
        <text class="text-[18px] font-bold text-cta">{{ priceLabel }}</text>
        <text class="text-[11px] font-medium px-2.5 py-1 rounded-full border" :class="stockClass">{{ stockLabel }}</text>
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
  if (stockCount.value <= 0) return 'bg-gray-50 text-gray-500 border-gray-200'
  if (stockPercent.value < 20) return 'bg-red-50 text-red-600 border-red-100'
  if (stockPercent.value < 50) return 'bg-orange-50 text-orange-600 border-orange-100'
  return 'bg-green-50 text-green-700 border-green-100'
})
</script>
