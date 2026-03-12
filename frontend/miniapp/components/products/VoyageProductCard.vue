<script setup lang="ts">
import { ShoppingCart } from 'lucide-vue-next'

type VoyageProductItem = {
  id: number
  brief_info?: string
  image_url?: string
  depart_date?: string
  sold_count?: number
  min_price_cents?: number
}

const props = defineProps<{
  item: VoyageProductItem
}>()

const emit = defineEmits<{
  (e: 'select', voyageId: number): void
}>()

function formatDate(value?: string) {
  if (!value) return '待定'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value.slice(0, 10)
  }
  const year = date.getFullYear()
  const month = `${date.getMonth() + 1}`.padStart(2, '0')
  const day = `${date.getDate()}`.padStart(2, '0')
  return `${year}-${month}-${day}`
}

function formatPrice(priceCents?: number) {
  if (!priceCents || priceCents <= 0) return '待定'
  return `${Math.round(priceCents / 100)}`
}

function handleClick() {
  emit('select', props.item.id)
}
</script>

<template>
  <button
    type="button"
    class="w-full overflow-hidden rounded-[18px] border border-slate-200 bg-white text-left shadow-[0_6px_20px_rgba(15,23,42,0.08)] transition-smooth hover:border-sky-300"
    @click="handleClick"
  >
    <div class="relative h-[124px] overflow-hidden bg-slate-100">
      <img
        class="h-full w-full object-cover"
        :src="item.image_url || `https://picsum.photos/seed/voyage-card-${item.id}/960/560`"
        :alt="item.brief_info || '航次商品图'"
      />
      <span class="absolute left-3 top-3 rounded-md bg-[#ff5b57] px-2 py-1 text-[11px] font-semibold text-white">热门</span>
      <span class="absolute bottom-3 right-3 rounded-full bg-black/55 px-2.5 py-1 text-[11px] text-white">{{ formatDate(item.depart_date) }}</span>
    </div>

    <div class="flex items-end justify-between gap-3 px-3 py-3">
      <div class="min-w-0 flex-1">
        <div class="line-clamp-2 text-[16px] font-semibold leading-6 text-slate-900">{{ item.brief_info || '未命名航次' }}</div>
        <div class="mt-2 text-[12px] text-slate-400">已售{{ item.sold_count || 0 }}人</div>
        <div class="mt-1 text-[28px] font-extrabold leading-none text-[#ff4d2d]">
          <span class="mr-0.5 text-[15px] font-semibold">¥</span>{{ formatPrice(item.min_price_cents) }}
        </div>
      </div>

      <div class="mb-1 flex h-10 w-10 shrink-0 items-center justify-center rounded-full border border-slate-200 text-slate-500">
        <ShoppingCart class="h-5 w-5" />
      </div>
    </div>
  </button>
</template>