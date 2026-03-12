<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { request } from '../../src/utils/request'
import NavBar from '../../components/NavBar.vue'
import { buildScheduleTableStops } from '../../src/utils/route-map'

type CompanyItem = {
  id: number
  name: string
}

type CruiseItem = {
  id: number
  name: string
  company?: CompanyItem
}

type VoyageItinerary = {
  id: number
  day_no: number
  stop_index: number
  city: string
  summary?: string
  eta_time?: string
  etd_time?: string
  has_breakfast?: boolean
  has_lunch?: boolean
  has_dinner?: boolean
  has_accommodation?: boolean
  accommodation_text?: string
}

type ContentTextItem = {
  text: string
  emphasis?: boolean
}

type FeeNoteContent = {
  included?: ContentTextItem[]
  excluded?: ContentTextItem[]
}

type BookingNoticeSection = {
  key: string
  title: string
  items?: ContentTextItem[]
}

type BookingNoticeContent = {
  sections?: BookingNoticeSection[]
}

type VoyageDetail = {
  id: number
  cruise_id: number
  code?: string
  brief_info?: string
  image_url?: string
  depart_date?: string
  return_date?: string
  min_price_cents?: number
  sold_count?: number
  cruise?: CruiseItem
  itineraries?: VoyageItinerary[]
  route_map?: {
    provider?: string
    geometry_type?: string
    coordinates?: number[][][]
    distance_km?: number
    resolution_km?: number
  }
  fee_note?: FeeNoteContent
  booking_notice?: BookingNoticeContent
}

type CabinTypeItem = {
  id: number
  name?: string
  min_price_cents?: number
}

type FacilityItem = {
  id: number
  category_id?: number
  name?: string
  extra_charge?: boolean
}

const props = defineProps<{ voyageId?: number | string }>()
const emit = defineEmits<{ (e: 'back'): void }>()

const loading = ref(false)
const error = ref('')
const detail = ref<VoyageDetail | null>(null)
const cabinTypes = ref<CabinTypeItem[]>([])
const facilities = ref<FacilityItem[]>([])
const activeNoticeSectionKey = ref('')

const resolvedVoyageId = computed(() => Number(props.voyageId ?? 0))

const heroImage = computed(() => {
  if (detail.value?.image_url) return detail.value.image_url
  return `https://picsum.photos/seed/voyage-detail-${resolvedVoyageId.value || 101}/1200/720`
})

const companyCruiseLabel = computed(() => {
  const companyName = detail.value?.cruise?.company?.name || '未知公司'
  const cruiseName = detail.value?.cruise?.name || '未知邮轮'
  return `${companyName} · ${cruiseName}`
})

const itineraryRows = computed(() => {
  return [...(detail.value?.itineraries || [])].sort((left, right) => {
    if (left.day_no !== right.day_no) {
      return left.day_no - right.day_no
    }
    return left.stop_index - right.stop_index
  })
})

const scheduleTableRows = computed(() => {
  return buildScheduleTableStops(
    itineraryRows.value.map((item) => ({
      city: item.city,
      dayNo: item.day_no,
      summary: item.summary,
      etaTime: item.eta_time,
      etdTime: item.etd_time,
    })),
  )
})

const feeIncluded = computed(() => detail.value?.fee_note?.included || [])
const feeExcluded = computed(() => detail.value?.fee_note?.excluded || [])
const bookingNoticeSections = computed(() => detail.value?.booking_notice?.sections || [])
const activeBookingNoticeSection = computed(() => {
  const sections = bookingNoticeSections.value
  if (sections.length === 0) return null
  return sections.find((item) => item.key === activeNoticeSectionKey.value) || sections[0]
})

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

function formatPrice(value?: number) {
  if (!value || value <= 0) return '待定'
  return `${Math.round(value / 100)}`
}

function formatTime(value?: string) {
  if (!value) return '--'
  return value
}

async function loadAll() {
  if (!resolvedVoyageId.value) {
    error.value = '缺少 voyageId 参数'
    return
  }
  loading.value = true
  error.value = ''
  try {
    const detailRes = await request(`/voyages/${resolvedVoyageId.value}`)
    const detailPayload = (detailRes as any)?.data ?? detailRes ?? null
    detail.value = detailPayload
    activeNoticeSectionKey.value = detailPayload?.booking_notice?.sections?.[0]?.key || ''

    const cruiseId = Number(detailPayload?.cruise_id || 0)
    if (!cruiseId) {
      throw new Error('航次缺少关联邮轮信息')
    }

    const [cabinTypeRes, facilityRes] = await Promise.all([
      request(`/cabin-types?cruise_id=${cruiseId}&page=1&page_size=20`),
      request(`/facilities?cruise_id=${cruiseId}`),
    ])
    const cabinTypePayload = (cabinTypeRes as any)?.data ?? cabinTypeRes ?? {}
    const facilityPayload = (facilityRes as any)?.data ?? facilityRes ?? []

    cabinTypes.value = Array.isArray(cabinTypePayload) ? cabinTypePayload : cabinTypePayload?.list ?? []
    facilities.value = Array.isArray(facilityPayload) ? facilityPayload : facilityPayload?.list ?? []
  } catch (loadError: any) {
    error.value = loadError?.message ?? '加载航次详情失败'
    detail.value = null
    cabinTypes.value = []
    facilities.value = []
  } finally {
    loading.value = false
  }
}

onMounted(loadAll)
</script>

<template>
  <div class="min-h-screen bg-[#f5f6fa] pb-[92px]">
    <div v-if="loading" class="flex min-h-screen items-center justify-center text-[14px] text-slate-500">加载中...</div>

    <div v-else-if="error" class="flex min-h-screen items-center justify-center px-6 text-center text-[14px] text-rose-500">{{ error }}</div>

    <template v-else-if="detail">
    <NavBar title="航次详情" show-back transparent @back="emit('back')" />
    <div class="relative h-[220px] overflow-hidden bg-slate-900">
      <img
        class="h-full w-full object-cover opacity-85"
        :src="heroImage"
        :alt="detail.brief_info || '航次详情'"
      />
      <div class="absolute inset-0 bg-gradient-to-t from-black/70 via-black/15 to-transparent"></div>
      <div class="absolute bottom-0 left-0 right-0 px-4 pb-5 text-white">
        <div class="mt-2 text-[20px] font-bold leading-7">{{ detail.brief_info || '未命名航次' }}</div>
      </div>
    </div>

    <div class="space-y-3 px-4 py-4">
      <section class="rounded-[20px] bg-white p-4 shadow-[0_8px_24px_rgba(15,23,42,0.06)]">
        <div class="text-[12px] text-slate-400">航次编号</div>
        <div class="mt-1 text-[20px] font-bold text-slate-900">{{ detail.code || '-' }}</div>
        <div class="mt-2 text-[14px] text-slate-500">{{ companyCruiseLabel }}</div>
        <div class="mt-4 flex flex-wrap gap-3 text-[13px] text-slate-600">
          <div class="rounded-full bg-slate-100 px-3 py-1.5">出发 {{ formatDate(detail.depart_date) }}</div>
          <div class="rounded-full bg-slate-100 px-3 py-1.5">返程 {{ formatDate(detail.return_date) }}</div>
          <div class="rounded-full bg-slate-100 px-3 py-1.5">已售 {{ detail.sold_count || 0 }} 人</div>
          <div class="rounded-full bg-[#fff3ea] px-3 py-1.5 text-[#ff5b20]">¥{{ formatPrice(detail.min_price_cents) }} 起</div>
        </div>
      </section>

      <section class="rounded-[20px] bg-white p-4 shadow-[0_8px_24px_rgba(15,23,42,0.06)]">
        <div class="text-[16px] font-semibold text-slate-900">航次说明</div>
        <p class="mt-3 text-[14px] leading-7 text-slate-600">{{ detail.brief_info || '暂无航次说明' }}</p>
      </section>

      <section class="rounded-[20px] bg-white p-4 shadow-[0_8px_24px_rgba(15,23,42,0.06)]">
        <div class="text-[16px] font-semibold text-slate-900">行程安排</div>
        <div v-if="itineraryRows.length === 0" class="mt-3 text-[13px] text-slate-400">暂无行程安排</div>
        <div v-else class="mt-3 space-y-3">
          <div v-for="item in itineraryRows" :key="item.id" class="rounded-2xl bg-slate-50 px-3 py-3">
            <div>
              <div>
                <div class="text-[12px] font-semibold text-[#0f5ba9]">第{{ item.day_no }}天 · {{ item.city }}</div>
                <div class="mt-1 text-[13px] leading-6 text-slate-600">{{ item.summary || '暂无说明' }}</div>
              </div>
            </div>
            <div class="mt-3 flex flex-wrap gap-2">
              <span v-if="item.has_breakfast" class="rounded-full bg-[#e8f7ee] px-3 py-1 text-[12px] font-medium text-[#16803c]">早餐</span>
              <span v-if="item.has_lunch" class="rounded-full bg-[#eef5ff] px-3 py-1 text-[12px] font-medium text-[#0f5ba9]">午餐</span>
              <span v-if="item.has_dinner" class="rounded-full bg-[#fff3ea] px-3 py-1 text-[12px] font-medium text-[#ff5b20]">晚餐</span>
              <span v-if="item.has_accommodation" class="rounded-full bg-slate-900 px-3 py-1 text-[12px] font-medium text-white">住宿</span>
            </div>
            <div v-if="item.has_accommodation && item.accommodation_text" class="mt-3 rounded-2xl bg-white px-3 py-2 text-[12px] leading-5 text-slate-500">
              {{ item.accommodation_text }}
            </div>
          </div>
        </div>
      </section>
      <section class="rounded-[20px] bg-white p-4 shadow-[0_8px_24px_rgba(15,23,42,0.06)]" data-test="schedule-table">
        <div class="text-[16px] font-semibold text-slate-900">抵离港时间</div>
        <div class="mt-3 overflow-hidden rounded-2xl border border-slate-200 bg-white">
          <div class="grid grid-cols-[68px_minmax(0,1fr)_82px_82px] bg-slate-50 px-3 py-2 text-[12px] font-semibold text-slate-500">
            <div class="whitespace-nowrap">天数</div>
            <div class="whitespace-nowrap">港口</div>
            <div class="whitespace-nowrap">抵港</div>
            <div class="whitespace-nowrap">离港</div>
          </div>
          <div v-for="item in scheduleTableRows" :key="`${item.dayNo}-${item.city}-row`" class="grid grid-cols-[68px_minmax(0,1fr)_82px_82px] items-center border-t border-slate-100 px-3 py-2 text-[12px] text-slate-600">
            <div class="whitespace-nowrap">第{{ item.dayNo }}天</div>
            <div class="truncate pr-2 font-medium text-slate-700">{{ item.city }}</div>
            <div class="whitespace-nowrap font-mono text-[12px]">{{ formatTime(item.etaTime) }}</div>
            <div class="whitespace-nowrap font-mono text-[12px]">{{ formatTime(item.etdTime) }}</div>
          </div>
        </div>
      </section>

      <section class="rounded-[20px] bg-white p-4 shadow-[0_8px_24px_rgba(15,23,42,0.06)]">
        <div class="text-[16px] font-semibold text-slate-900">可选舱型</div>
        <div v-if="cabinTypes.length === 0" class="mt-3 text-[13px] text-slate-400">暂无舱型数据</div>
        <div v-else class="mt-3 space-y-3">
          <div v-for="item in cabinTypes" :key="item.id" class="flex items-center justify-between rounded-2xl bg-slate-50 px-3 py-3">
            <div class="text-[14px] text-slate-700">{{ item.name || '未命名舱型' }}</div>
            <div class="text-[16px] font-semibold text-[#ff5b20]">¥{{ formatPrice(item.min_price_cents) }}</div>
          </div>
        </div>
      </section>

      <section class="rounded-[20px] bg-white p-4 shadow-[0_8px_24px_rgba(15,23,42,0.06)]">
        <div class="text-[16px] font-semibold text-slate-900">船上设施</div>
        <div v-if="facilities.length === 0" class="mt-3 text-[13px] text-slate-400">暂无设施数据</div>
        <div v-else class="mt-3 flex flex-wrap gap-2.5">
          <div v-for="item in facilities" :key="item.id" class="rounded-full bg-slate-100 px-3 py-2 text-[13px] text-slate-600">
            {{ item.name || '未命名设施' }}
          </div>
        </div>
      </section>

      <section class="rounded-[20px] bg-white p-4 shadow-[0_8px_24px_rgba(15,23,42,0.06)]">
        <div class="text-[16px] font-semibold text-slate-900">费用说明</div>
        <div class="mt-4 grid gap-4 md:grid-cols-2">
          <div class="rounded-[18px] bg-[#f8fbff] p-4">
            <div class="text-[15px] font-semibold text-slate-900">费用包含</div>
            <ul v-if="feeIncluded.length > 0" class="mt-3 space-y-2 text-[13px] leading-6 text-slate-600">
              <li v-for="(item, index) in feeIncluded" :key="`in-${index}`" class="flex gap-2"><span>•</span><span>{{ item.text }}</span></li>
            </ul>
            <div v-else class="mt-3 text-[13px] text-slate-400">暂无费用包含说明</div>
          </div>
          <div class="rounded-[18px] bg-[#fff7f3] p-4">
            <div class="text-[15px] font-semibold text-slate-900">费用不包含</div>
            <ul v-if="feeExcluded.length > 0" class="mt-3 space-y-2 text-[13px] leading-6 text-slate-600">
              <li v-for="(item, index) in feeExcluded" :key="`out-${index}`" class="flex gap-2"><span>•</span><span :class="item.emphasis ? 'font-semibold text-[#ff5b20]' : ''">{{ item.text }}</span></li>
            </ul>
            <div v-else class="mt-3 text-[13px] text-slate-400">暂无费用不包含说明</div>
          </div>
        </div>
      </section>

      <section class="rounded-[20px] bg-white p-4 shadow-[0_8px_24px_rgba(15,23,42,0.06)]">
        <div class="text-[16px] font-semibold text-slate-900">预订须知</div>
        <div v-if="bookingNoticeSections.length === 0" class="mt-3 text-[13px] text-slate-400">暂无预订须知</div>
        <template v-else>
          <div class="mt-4 flex gap-2 overflow-x-auto pb-1">
            <button v-for="section in bookingNoticeSections" :key="section.key" type="button" class="shrink-0 rounded-full px-3 py-1.5 text-[12px] font-medium" :class="activeBookingNoticeSection?.key === section.key ? 'bg-[#ffefe7] text-[#ff5b20]' : 'bg-slate-100 text-slate-500'" @click="activeNoticeSectionKey = section.key">
              {{ section.title }}
            </button>
          </div>
          <div v-if="activeBookingNoticeSection" class="mt-4 rounded-[18px] border border-[#fde7da] bg-[linear-gradient(180deg,#fffaf7_0%,#ffffff_100%)] p-4">
            <div class="inline-flex rounded-full bg-[#ffefe7] px-3 py-1 text-[11px] font-semibold tracking-[0.12em] text-[#ff5b20]">NOTICE</div>
            <div class="mt-3 text-[15px] font-semibold text-slate-900">{{ activeBookingNoticeSection.title }}</div>
            <div class="mt-3 space-y-3">
              <div v-for="(item, index) in activeBookingNoticeSection.items || []" :key="`notice-${index}`" class="rounded-2xl border px-3 py-3 text-[13px] leading-6" :class="item.emphasis ? 'border-[#ffd7c2] bg-[#fff3ea] font-semibold text-[#ff5b20]' : 'border-slate-100 bg-white text-slate-600'">
                {{ item.text }}
              </div>
            </div>
          </div>
        </template>
      </section>
    </div>
    </template>
  </div>
</template>