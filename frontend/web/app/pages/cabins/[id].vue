<template>
  <div class="min-h-screen bg-[#f8f4ed] pb-24 text-slate-900">
    <main class="mx-auto max-w-6xl px-6 py-8">
      <div class="grid grid-cols-1 gap-6 lg:grid-cols-[1.2fr_1fr]">
        <section class="overflow-hidden rounded-2xl border border-[#eadfcb] bg-white shadow-sm">
          <div class="relative aspect-[16/10]">
            <img :src="gallery[currentSlide]" alt="cabin" class="h-full w-full object-cover" />
            <button type="button" class="absolute left-3 top-1/2 -translate-y-1/2 rounded-full bg-white/80 px-3 py-2 text-sm" @click="prevSlide">‹</button>
            <button type="button" class="absolute right-3 top-1/2 -translate-y-1/2 rounded-full bg-white/80 px-3 py-2 text-sm" @click="nextSlide">›</button>
            <div class="absolute bottom-3 left-1/2 flex -translate-x-1/2 gap-2">
              <span v-for="(_, idx) in gallery" :key="'gallery-dot-' + idx" class="h-2 w-2 rounded-full" :class="idx === currentSlide ? 'bg-white' : 'bg-white/40'" />
            </div>
            <span class="absolute right-3 top-3 rounded-full bg-slate-900/60 px-2 py-1 text-xs text-white">{{ currentSlide + 1 }}/{{ gallery.length }}</span>
          </div>
        </section>

        <section class="space-y-4 rounded-2xl border border-[#eadfcb] bg-white p-5 shadow-sm">
          <h1 class="font-['Playfair_Display','Georgia',serif] text-3xl text-[#12263a]">{{ detail?.code || '舱位详情' }}</h1>
          <p class="text-sm text-slate-600">{{ detail?.deck || '-' }}甲板 · {{ detail?.area || '-' }}m2 · 最多{{ detail?.max_guests || '-' }}人</p>

          <div class="flex flex-wrap gap-2" data-test="attr-tags">
            <span class="rounded-full bg-blue-50 px-3 py-1 text-xs text-blue-700">床型: {{ detail?.bed_type || '-' }}</span>
            <span class="rounded-full bg-blue-50 px-3 py-1 text-xs text-blue-700">位置: {{ detail?.position || '-' }}</span>
            <span class="rounded-full bg-blue-50 px-3 py-1 text-xs text-blue-700">朝向: {{ detail?.orientation || '-' }}</span>
            <span class="rounded-full bg-blue-50 px-3 py-1 text-xs text-blue-700">窗户: {{ detail?.has_window ? '有' : '无' }}</span>
            <span class="rounded-full bg-blue-50 px-3 py-1 text-xs text-blue-700">阳台: {{ detail?.has_balcony ? '有' : '无' }}</span>
          </div>

          <div>
            <p class="mb-2 text-sm text-slate-600">库存状态</p>
            <div class="h-2 rounded bg-slate-200">
              <div class="h-2 rounded" :class="inventoryBarClass" :style="{ width: `${inventoryPercent}%` }" />
            </div>
            <p class="mt-2 text-sm" :class="inventoryLabelClass">{{ inventoryLabel }}</p>
          </div>
        </section>
      </div>

      <section class="mt-6 rounded-2xl border border-[#eadfcb] bg-white p-5 shadow-sm">
        <h2 class="mb-3 text-lg font-semibold text-[#12263a]">价格日历</h2>
        <div class="mb-3 flex flex-wrap gap-2">
          <button v-for="type in priceTypes" :key="type.key" type="button" class="rounded-full px-3 py-1 text-sm" :class="activePriceType === type.key ? 'bg-[#12263a] text-white' : 'bg-[#f4efe6] text-slate-700'" @click="activePriceType = type.key">{{ type.label }}</button>
        </div>
        <div class="grid grid-cols-7 gap-2" data-test="price-calendar">
          <button
            v-for="date in calendarDays"
            :key="date"
            type="button"
            class="rounded-lg border px-2 py-2 text-center text-xs"
            :class="calendarCellClass(date)"
            @click="selectedDate = date"
          >
            <p>{{ date.slice(-2) }}</p>
            <p>{{ priceText(date) }}</p>
          </button>
        </div>
      </section>
    </main>

    <div class="fixed bottom-0 left-0 right-0 border-t border-[#eadfcb] bg-[#102e45] px-4 py-3">
      <div class="mx-auto flex max-w-6xl items-center justify-between">
        <div>
          <p class="text-xs text-[#c8d7e3]">当前选中价格</p>
          <p class="font-['Playfair_Display','Georgia',serif] text-2xl text-white">¥{{ selectedPriceDisplay }}</p>
        </div>
        <button type="button" class="rounded-xl bg-white px-5 py-3 text-sm font-medium text-[#102e45]">立即预订</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

const route = useRoute()
const { request } = useApi()
const id = Number(route.params.id)

const detail = ref<Record<string, any> | null>(null)
const prices = ref<Array<Record<string, any>>>([])
const inventory = ref<Record<string, any> | null>(null)
const currentSlide = ref(0)
const selectedDate = ref('')
const activePriceType = ref('base')

const priceTypes = [
  { key: 'base', label: '基础' },
  { key: 'child', label: '儿童' },
  { key: 'single_supplement', label: '单人补差' },
  { key: 'holiday', label: '节假日' },
  { key: 'early_bird', label: '早鸟' },
]

const gallery = computed(() => {
  const images = detail.value && Array.isArray(detail.value.images) ? detail.value.images : []
  if (images.length > 0) return images.map((item: any) => item.url || item)
  return [
    `https://picsum.photos/seed/web-cabin-${id}-1/1400/900`,
    `https://picsum.photos/seed/web-cabin-${id}-2/1400/900`,
    `https://picsum.photos/seed/web-cabin-${id}-3/1400/900`,
  ]
})

const available = computed(() => {
  if (!inventory.value) return 0
  if (inventory.value.available !== undefined) return Number(inventory.value.available)
  return Number(inventory.value.total || 0) - Number(inventory.value.locked || 0) - Number(inventory.value.sold || 0)
})

const inventoryPercent = computed(() => {
  const total = Number(inventory.value?.total || 0)
  if (total <= 0) return 0
  return Math.max(0, Math.min(100, Math.round((available.value / total) * 100)))
})

const inventoryBarClass = computed(() => {
  if (inventoryPercent.value > 50) return 'bg-emerald-500'
  if (inventoryPercent.value >= 20) return 'bg-amber-500'
  return 'bg-rose-500'
})

const inventoryLabelClass = computed(() => {
  if (inventoryPercent.value > 50) return 'text-emerald-600'
  if (inventoryPercent.value >= 20) return 'text-amber-600'
  return 'text-rose-600'
})

const inventoryLabel = computed(() => {
  if (inventoryPercent.value < 20) return `仅剩 ${available.value} 间`
  if (inventoryPercent.value < 50) return '库存紧张'
  return '库存充足'
})

const calendarDays = computed(() => {
  const base = new Date()
  const y = base.getFullYear()
  const m = base.getMonth()
  const end = new Date(y, m + 1, 0).getDate()
  const days: string[] = []
  for (let d = 1; d <= end; d++) {
    const mm = `${m + 1}`.padStart(2, '0')
    const dd = `${d}`.padStart(2, '0')
    days.push(`${y}-${mm}-${dd}`)
  }
  if (!selectedDate.value && days.length > 0) selectedDate.value = days[0]
  return days
})

function normalizeDate(raw: unknown) {
  if (typeof raw !== 'string') return ''
  return raw.slice(0, 10)
}

function priceForDate(date: string) {
  const row = prices.value.find((item) => normalizeDate(item.date) === date && (item.price_type || 'base') === activePriceType.value)
  return row ? Number(row.price_cents || 0) : 0
}

function priceText(date: string) {
  const value = priceForDate(date)
  return value > 0 ? `${Math.round(value / 100)}` : '-'
}

function calendarCellClass(date: string) {
  const value = priceForDate(date)
  if (selectedDate.value === date) return 'ring-2 ring-indigo-500 border-indigo-200 bg-white text-slate-800'
  if (value <= 0) return 'border-slate-200 bg-slate-50 text-slate-400'
  return 'border-emerald-200 bg-emerald-50 text-emerald-700'
}

const selectedPriceDisplay = computed(() => {
  const value = priceForDate(selectedDate.value)
  return value > 0 ? Math.round(value / 100) : '-'
})

function prevSlide() {
  currentSlide.value = (currentSlide.value - 1 + gallery.value.length) % gallery.value.length
}

function nextSlide() {
  currentSlide.value = (currentSlide.value + 1) % gallery.value.length
}

async function loadData() {
  const [detailRes, priceRes, inventoryRes] = await Promise.all([
    request(`/cabins/${id}`),
    request(`/cabins/${id}/prices`),
    request(`/cabins/${id}/inventory`),
  ])

  detail.value = detailRes?.data ?? detailRes ?? null
  const pricePayload = priceRes?.data ?? priceRes ?? []
  prices.value = Array.isArray(pricePayload) ? pricePayload : pricePayload?.list ?? []
  inventory.value = inventoryRes?.data ?? inventoryRes ?? null
}

onMounted(loadData)
</script>
