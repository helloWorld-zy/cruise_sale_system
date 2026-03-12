<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import CruiseFacilityCategoryIcon from '../../components/cruise/CruiseFacilityCategoryIcon.vue'
import NavBar from '../../components/NavBar.vue'
import { cruiseFacilityCategories, groupCruiseFacilities, type CruiseFacilityCategoryId } from '../../src/constants/cruiseFacilityCategories'
import { request } from '../../src/utils/request'

const props = defineProps<{ cruiseId?: number | string }>()
const emit = defineEmits<{ (e: 'back'): void }>()

const loading = ref(false)
const error = ref('')
const detail = ref<Record<string, any> | null>(null)
const facilities = ref<Record<string, any>[]>([])
const activeFacilityCategoryId = ref<CruiseFacilityCategoryId | null>(null)

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

const shipLabel = computed(() => {
  const raw = detail.value?.company_name || detail.value?.company?.name || detail.value?.brand_name
  return String(raw || '邮轮档案')
})

const basicSpecs = computed(() => {
  const source = detail.value ?? {}
  return [
    { label: '总吨位', value: source.tonnage, unit: '吨' },
    { label: '建造年份', value: source.build_year, unit: '年' },
    { label: '载客量', value: source.passenger_capacity, unit: '人' },
    { label: '甲板数', value: source.deck_count, unit: '层' },
    { label: '船长', value: source.length, unit: '米' },
    { label: '船宽', value: source.width, unit: '米' },
    { label: '舱房数', value: source.room_count, unit: '间' },
    { label: '船员数', value: source.crew_count, unit: '人' },
    { label: '翻新年份', value: source.refurbish_year, unit: '年' },
  ]
    .filter((item) => hasDisplayValue(item.value))
    .map((item) => ({ ...item, displayValue: formatSpecValue(item.value) }))
})

const descriptionText = computed(() => {
  const raw = detail.value?.description
  return typeof raw === 'string' ? raw.trim() : ''
})

const facilityGroups = computed(() => groupCruiseFacilities(facilities.value))

const activeFacilityCategory = computed(() =>
  cruiseFacilityCategories.find((category) => category.id === activeFacilityCategoryId.value) ?? null,
)

const activeFacilityItems = computed(() => {
  if (!activeFacilityCategoryId.value) {
    return []
  }
  return facilityGroups.value.get(activeFacilityCategoryId.value) ?? []
})

function resolveCruiseId() {
  return Number(props.cruiseId ?? 0)
}

function hasDisplayValue(value: unknown) {
  if (value === null || value === undefined) return false
  if (typeof value === 'number') return Number.isFinite(value) && value > 0
  if (typeof value === 'string') return value.trim().length > 0 && value.trim() !== '0'
  return Boolean(value)
}

function formatSpecValue(value: unknown) {
  if (typeof value === 'number' && Number.isFinite(value)) {
    return Number.isInteger(value) ? value.toLocaleString('zh-CN') : value.toFixed(1)
  }
  return String(value ?? '')
}

function toggleFacilityCategory(categoryId: CruiseFacilityCategoryId) {
  activeFacilityCategoryId.value = activeFacilityCategoryId.value === categoryId ? null : categoryId
}

function buildFacilityMeta(item: Record<string, any>) {
  return [
    item.location ? `位置：${item.location}` : '',
    item.open_hours ? `开放时间：${item.open_hours}` : '',
    item.target_audience ? `适合人群：${item.target_audience}` : '',
  ].filter(Boolean)
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
    const [detailRes, , facilityRes] = await Promise.all([
      request(`/cruises/${id}`),
      request(`/cabin-types?cruise_id=${id}&page=1&page_size=20`),
      request(`/facilities?cruise_id=${id}`),
    ])

    detail.value = (detailRes as any)?.data ?? detailRes ?? null
    const facilityPayload = (facilityRes as any)?.data ?? facilityRes ?? []
    facilities.value = Array.isArray(facilityPayload) ? facilityPayload : facilityPayload?.list ?? []
    activeFacilityCategoryId.value = null
  } catch (e: any) {
    error.value = e?.message ?? '加载邮轮详情失败'
  } finally {
    loading.value = false
  }
}

onMounted(loadAll)
</script>

<template>
  <div class="min-h-screen bg-[linear-gradient(180deg,#f5f6f8_0%,#eef1f5_100%)] pb-10">
    <div v-if="loading" class="flex min-h-screen items-center justify-center px-8 text-center text-[14px] text-slate-500">Loading...</div>
    <div v-else-if="error" class="flex min-h-screen items-center justify-center px-8 text-center text-[14px] text-rose-500">{{ error }}</div>
    <div v-else-if="detail" class="pb-10">
      <NavBar title="邮轮详情" show-back @back="emit('back')" />

      <div class="mx-5 overflow-x-auto rounded-[28px] shadow-[0_14px_34px_rgba(15,23,42,0.10)] hide-scrollbar">
        <div class="flex snap-x snap-mandatory overflow-x-auto hide-scrollbar">
          <div v-for="img in gallery" :key="img" class="h-[240px] min-w-full snap-start overflow-hidden bg-slate-200">
            <img class="h-full w-full object-cover" :src="img" alt="邮轮详情图" />
          </div>
        </div>
      </div>

      <div class="relative z-10 mx-5 -mt-7 flex items-start justify-between gap-4 rounded-[28px] bg-white/92 px-6 py-5 shadow-[0_10px_28px_rgba(15,23,42,0.06)] backdrop-blur-sm">
        <div class="min-w-0 flex-1">
          <span class="block text-[20px] font-bold leading-tight text-slate-900">{{ detail.name || '-' }}</span>
          <span v-if="detail.english_name" class="mt-2 block text-[15px] italic text-slate-500">{{ detail.english_name }}</span>
        </div>
        <span class="shrink-0 rounded-full bg-[linear-gradient(135deg,#0f5ba9_0%,#2f7ed8_100%)] px-4 py-2 text-[15px] font-semibold text-white">{{ shipLabel }}</span>
      </div>

      <section v-if="basicSpecs.length" class="mx-5 mt-4 rounded-[24px] bg-white/92 px-6 py-6 shadow-[0_10px_28px_rgba(15,23,42,0.06)] backdrop-blur-sm">
        <div class="mb-4 flex items-center gap-3">
          <span class="h-7 w-2 rounded-full bg-[linear-gradient(180deg,#0f5ba9_0%,#5a9ff0_100%)]"></span>
          <span class="text-[16px] font-bold text-slate-900">基本参数</span>
        </div>
        <div class="grid grid-cols-3 gap-3 max-[360px]:grid-cols-2">
          <div v-for="item in basicSpecs" :key="item.label" class="min-h-[90px] rounded-[18px] bg-[linear-gradient(180deg,#f8fbff_0%,#f3f6fb_100%)] px-4 py-4">
            <span class="block text-[13px] text-slate-500">{{ item.label }}</span>
            <div class="mt-3 flex items-baseline gap-1.5">
              <span class="text-[18px] font-bold text-slate-900">{{ item.displayValue }}</span>
              <span class="text-[13px] text-slate-600">{{ item.unit }}</span>
            </div>
          </div>
        </div>
      </section>

      <section v-if="descriptionText" class="mx-5 mt-4 rounded-[24px] bg-white/92 px-6 py-6 shadow-[0_10px_28px_rgba(15,23,42,0.06)] backdrop-blur-sm">
        <div class="mb-4 flex items-center gap-3">
          <span class="h-7 w-2 rounded-full bg-[linear-gradient(180deg,#0f5ba9_0%,#5a9ff0_100%)]"></span>
          <span class="text-[16px] font-bold text-slate-900">邮轮介绍</span>
        </div>
        <p class="m-0 whitespace-pre-wrap text-[14px] leading-7 text-slate-600">{{ descriptionText }}</p>
      </section>

      <section class="mx-5 mt-4 rounded-[24px] bg-white/92 px-6 py-6 shadow-[0_10px_28px_rgba(15,23,42,0.06)] backdrop-blur-sm">
        <div class="mb-4 flex items-center gap-3">
          <span class="h-7 w-2 rounded-full bg-[linear-gradient(180deg,#0f5ba9_0%,#5a9ff0_100%)]"></span>
          <span class="text-[16px] font-bold text-slate-900">邮轮设施</span>
        </div>
        <div class="grid grid-cols-4 gap-x-2 gap-y-5 max-[360px]:grid-cols-2">
          <button
            v-for="category in cruiseFacilityCategories"
            :key="category.id"
            type="button"
            :data-test="`facility-category-${category.id}`"
            class="flex flex-col items-center gap-2 border-0 bg-transparent p-0 text-slate-500"
            :class="activeFacilityCategoryId === category.id ? 'text-[#0f5ba9]' : ''"
            @click="toggleFacilityCategory(category.id)"
          >
            <span class="flex h-[72px] w-[72px] items-center justify-center rounded-full bg-slate-100 shadow-[inset_0_0_0_1px_rgba(148,163,184,0.12)]" :class="activeFacilityCategoryId === category.id ? 'bg-[linear-gradient(180deg,#e7f1ff_0%,#dcebff_100%)] text-[#0f5ba9] shadow-[inset_0_0_0_1px_rgba(15,91,169,0.15)]' : ''">
              <CruiseFacilityCategoryIcon :name="category.icon" />
            </span>
            <span class="text-center text-[15px] leading-5">{{ category.label }}</span>
          </button>
        </div>

        <div v-if="activeFacilityCategory" class="relative mt-6 overflow-hidden rounded-[22px] bg-[linear-gradient(180deg,#f8fbff_0%,#eef5ff_100%)] px-6 py-6 shadow-[inset_0_0_0_1px_rgba(148,163,184,0.08),0_12px_30px_rgba(15,91,169,0.06)]">
          <div class="pointer-events-none absolute -right-10 -top-14 h-[180px] w-[180px] rounded-full bg-[radial-gradient(circle,rgba(118,175,255,0.22)_0%,rgba(118,175,255,0)_72%)]"></div>
          <div class="relative z-10 mb-4">
            <span class="inline-block rounded-full bg-[rgba(15,91,169,0.08)] px-3.5 py-2 text-[13px] font-semibold tracking-[1px] text-[#0f5ba9]">服务亮点</span>
            <div class="mt-3 flex items-center gap-4">
              <span class="flex h-[64px] w-[64px] items-center justify-center rounded-[20px] bg-white text-[#0f5ba9]">
                <CruiseFacilityCategoryIcon :name="activeFacilityCategory.icon" />
              </span>
              <div>
                <span class="block text-[18px] font-bold text-slate-900">{{ activeFacilityCategory.label }}</span>
                <span class="mt-1 block text-[13px] text-slate-500">{{ activeFacilityItems.length }} 项设施</span>
              </div>
            </div>
          </div>

          <div v-if="activeFacilityItems.length === 0" class="relative z-10 rounded-[18px] bg-white px-6 py-6 text-center text-[15px] text-slate-500">
            当前类目暂未配置设施。
          </div>
          <div v-else class="relative z-10 space-y-4">
            <div v-for="(item, index) in activeFacilityItems" :key="item.id" class="rounded-[18px] border border-sky-100 bg-[linear-gradient(180deg,#ffffff_0%,#fbfdff_100%)] px-5 py-5 shadow-[0_8px_20px_rgba(15,23,42,0.05)]">
              <div class="grid grid-cols-[56px_minmax(0,1fr)_auto] items-center gap-3">
                <div class="flex h-[56px] w-[56px] items-center justify-center rounded-[18px] bg-[linear-gradient(135deg,#0f5ba9_0%,#65a3ee_100%)] text-[18px] font-bold tracking-[1px] text-white">{{ String(index + 1).padStart(2, '0') }}</div>
                <div class="min-w-0">
                  <span class="block text-[12px] tracking-[1px] text-slate-500">{{ item.category_name || activeFacilityCategory.label }}</span>
                  <span class="block text-[16px] font-bold text-slate-900">{{ item.name }}</span>
                </div>
                <span class="rounded-full px-3.5 py-2 text-[13px] font-medium" :class="item.extra_charge ? 'bg-orange-50 text-orange-700' : 'bg-emerald-50 text-emerald-700'">
                  {{ item.extra_charge ? '收费' : '免费' }}
                </span>
              </div>

              <div class="mt-4 border-t border-dashed border-slate-200 pt-4">
                <span class="block text-[15px] leading-7 text-slate-600">{{ item.description || '推荐在航海日与靠港日晚间安排体验，避开高峰时段会更从容。' }}</span>
              </div>
              <div v-if="buildFacilityMeta(item).length" class="mt-4 flex flex-wrap gap-2">
                <span v-for="meta in buildFacilityMeta(item)" :key="meta" class="rounded-full bg-[linear-gradient(180deg,#f2f7ff_0%,#edf3fb_100%)] px-3.5 py-2 text-[13px] text-[#35506d]">{{ meta }}</span>
              </div>
              <div v-if="item.charge_price_tip" class="mt-4 rounded-[16px] bg-[linear-gradient(180deg,#fff9ef_0%,#fff6e8_100%)] px-4 py-4">
                <span class="mb-1.5 block text-[12px] font-semibold text-amber-700">贴心提示</span>
                <span class="block text-[15px] leading-7 text-slate-600">{{ item.charge_price_tip }}</span>
              </div>
              <p v-if="item.description" class="mb-0 mt-4 text-[15px] leading-7 text-slate-600">{{ item.description }}</p>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>
