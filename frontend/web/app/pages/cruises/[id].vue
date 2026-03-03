<template>
  <div class="min-h-screen bg-[#f8f4ed] text-slate-900">
    <section class="relative h-[70vh] overflow-hidden">
      <img :src="gallery[currentSlide]" alt="cruise" class="h-full w-full object-cover" />
      <div class="absolute inset-0 bg-gradient-to-b from-slate-900/25 via-transparent to-slate-900/50" />
      <button type="button" class="absolute left-4 top-1/2 -translate-y-1/2 rounded-full bg-white/80 px-3 py-2 text-sm" @click="prevSlide">‹</button>
      <button type="button" class="absolute right-4 top-1/2 -translate-y-1/2 rounded-full bg-white/80 px-3 py-2 text-sm" @click="nextSlide">›</button>
      <div class="absolute bottom-4 left-1/2 flex -translate-x-1/2 gap-2">
        <span v-for="(_, idx) in gallery" :key="'gallery-dot-' + idx" class="h-2 w-2 rounded-full" :class="idx === currentSlide ? 'bg-white' : 'bg-white/40'" />
      </div>
    </section>

    <main class="mx-auto max-w-6xl space-y-6 px-6 py-8">
      <section class="rounded-2xl border border-[#eadfcb] bg-white p-5 shadow-sm">
        <h1 class="font-['Playfair_Display','Georgia',serif] text-3xl text-[#12263a]">{{ detail?.name || '邮轮详情' }}</h1>
        <div class="mt-4 flex flex-wrap items-center gap-4 text-sm text-slate-700">
          <span>⚓ {{ detail?.tonnage || '-' }} 吨</span>
          <span class="text-[#c9a96e]">|</span>
          <span>👤 {{ detail?.passenger_capacity || '-' }} 人</span>
          <span class="text-[#c9a96e]">|</span>
          <span>🧭 {{ detail?.length || '-' }} m</span>
          <span class="text-[#c9a96e]">|</span>
          <span>🛳 {{ detail?.deck_count || '-' }} 层</span>
        </div>
      </section>

      <section class="rounded-2xl border border-[#eadfcb] bg-white p-5 shadow-sm">
        <h2 class="mb-3 text-lg font-semibold text-[#12263a]">舱房类型</h2>
        <div class="space-y-3">
          <article v-for="type in cabinTypes" :key="type.id" class="rounded-xl border border-[#f0e6d5]">
            <button type="button" class="flex w-full items-center justify-between gap-3 p-3 text-left" @click="toggleType(type.id)">
              <div>
                <p class="font-medium text-slate-900">{{ type.name }}</p>
                <p class="text-xs text-slate-500">{{ type.area_min || type.area || '-' }}m2 · {{ type.max_capacity || type.capacity || '-' }}人</p>
              </div>
              <p class="font-['Playfair_Display','Georgia',serif] text-xl text-[#c9a96e]">¥{{ priceFrom(type) }}</p>
            </button>
            <div v-if="expandedTypeIds.has(Number(type.id))" class="border-t border-[#f0e6d5] px-3 py-3 text-sm text-slate-600">
              <p class="mb-2">{{ type.description || '暂无描述' }}</p>
              <div class="flex flex-wrap gap-2">
                <span v-for="tag in splitCsv(type.amenities || '')" :key="tag" class="rounded-full bg-[#f4efe6] px-2.5 py-1 text-xs">{{ tag }}</span>
              </div>
            </div>
          </article>
        </div>
      </section>

      <section class="rounded-2xl border border-[#eadfcb] bg-white p-5 shadow-sm">
        <h2 class="mb-3 text-lg font-semibold text-[#12263a]">设施导览</h2>
        <div class="mb-3 flex flex-wrap gap-2">
          <button type="button" class="rounded-full px-3 py-1 text-sm" :class="activeCategory === 0 ? 'bg-[#12263a] text-white' : 'bg-[#f4efe6] text-slate-700'" @click="activeCategory = 0">全部</button>
          <button v-for="cat in facilityCategories" :key="cat.id" type="button" class="rounded-full px-3 py-1 text-sm" :class="activeCategory === Number(cat.id) ? 'bg-[#12263a] text-white' : 'bg-[#f4efe6] text-slate-700'" @click="activeCategory = Number(cat.id)">{{ cat.name }}</button>
        </div>
        <div class="grid grid-cols-1 gap-3 md:grid-cols-2">
          <article v-for="fac in filteredFacilities" :key="fac.id" class="rounded-xl border border-[#f0e6d5] p-3">
            <p class="font-medium text-slate-900">{{ fac.name }}</p>
            <p class="mt-1 text-xs text-slate-500">{{ fac.open_hours || '时间待定' }}</p>
            <span class="mt-2 inline-block rounded-full px-2.5 py-1 text-xs" :class="fac.extra_charge ? 'bg-amber-50 text-amber-700' : 'bg-emerald-50 text-emerald-700'">{{ fac.extra_charge ? '收费' : '免费' }}</span>
          </article>
        </div>
      </section>

      <section class="rounded-2xl border border-[#eadfcb] bg-white p-5 shadow-sm">
        <h2 class="mb-3 text-lg font-semibold text-[#12263a]">关联航线</h2>
        <div class="space-y-3">
          <article v-for="(route, idx) in relatedRoutes" :key="route.id || idx" class="flex items-center justify-between rounded-xl border border-[#f0e6d5] p-3">
            <div>
              <p class="text-xs text-slate-500">{{ route.date || '-' }}</p>
              <p class="font-medium text-slate-900">{{ route.name || '-' }}</p>
            </div>
            <p class="font-['Playfair_Display','Georgia',serif] text-xl text-[#c9a96e]">¥{{ route.price || '-' }}</p>
          </article>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

const route = useRoute()
const { request } = useApi()
const id = Number(route.params.id)

const detail = ref<Record<string, any> | null>(null)
const cabinTypes = ref<Record<string, any>[]>([])
const facilities = ref<Record<string, any>[]>([])
const facilityCategories = ref<Record<string, any>[]>([])
const relatedRoutes = ref<Record<string, any>[]>([])
const expandedTypeIds = ref<Set<number>>(new Set())
const activeCategory = ref(0)
const currentSlide = ref(0)

const gallery = computed(() => {
  if (!detail.value) return ['https://picsum.photos/seed/cruise-hero/1600/900']
  const images = Array.isArray(detail.value.images) ? detail.value.images : []
  if (images.length > 0) return images.map((img: any) => img.url || img)
  return [
    `https://picsum.photos/seed/cruise-${id}-1/1600/900`,
    `https://picsum.photos/seed/cruise-${id}-2/1600/900`,
    `https://picsum.photos/seed/cruise-${id}-3/1600/900`,
  ]
})

const filteredFacilities = computed(() => {
  if (activeCategory.value === 0) return facilities.value
  return facilities.value.filter((item) => Number(item.category_id) === activeCategory.value)
})

function splitCsv(raw: unknown) {
  if (typeof raw !== 'string') return []
  return raw.split(',').map((part) => part.trim()).filter(Boolean)
}

function toggleType(typeId: number) {
  const next = new Set(expandedTypeIds.value)
  if (next.has(typeId)) next.delete(typeId)
  else next.add(typeId)
  expandedTypeIds.value = next
}

function priceFrom(type: Record<string, any>) {
  return Math.round(Number(type.min_price_cents || type.price_cents || 0) / 100) || '-'
}

function prevSlide() {
  if (gallery.value.length === 0) return
  currentSlide.value = (currentSlide.value - 1 + gallery.value.length) % gallery.value.length
}

function nextSlide() {
  if (gallery.value.length === 0) return
  currentSlide.value = (currentSlide.value + 1) % gallery.value.length
}

async function loadAll() {
  const [detailRes, typeRes, facilityRes, categoryRes] = await Promise.allSettled([
    request(`/cruises/${id}`),
    request('/cabin-types', { query: { cruise_id: id, page: 1, page_size: 50 } }),
    request('/facilities', { query: { cruise_id: id } }),
    request('/facility-categories'),
  ])

  const unwrap = (result: PromiseSettledResult<any>, fallback: any) => {
    if (result.status === 'fulfilled') return result.value
    return fallback
  }

  const detailPayload = unwrap(detailRes, null)
  detail.value = detailPayload?.data ?? detailPayload ?? null

  const typePayload = unwrap(typeRes, [])
  const typeData = typePayload?.data ?? typePayload ?? {}
  cabinTypes.value = Array.isArray(typeData) ? typeData : typeData?.list ?? []

  const facilityPayload = unwrap(facilityRes, [])
  const facilityData = facilityPayload?.data ?? facilityPayload ?? []
  facilities.value = Array.isArray(facilityData) ? facilityData : facilityData?.list ?? []

  const categoryPayload = unwrap(categoryRes, [])
  const categoryData = categoryPayload?.data ?? categoryPayload ?? []
  facilityCategories.value = Array.isArray(categoryData) ? categoryData : categoryData?.list ?? []

  // Public API currently has no dedicated route list endpoint for C-end, so keep this optional.
  relatedRoutes.value = []
}

onMounted(loadAll)
</script>
