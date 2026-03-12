<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { Headset, Search } from 'lucide-vue-next'
import NavBar from '../../components/NavBar.vue'
import ProductCompanyTree from '../../components/products/ProductCompanyTree.vue'
import VoyageProductCard from '../../components/products/VoyageProductCard.vue'
import { request } from '../../src/utils/request'

type CompanyItem = {
  id: number
  name: string
  english_name?: string
}

type CruiseItem = {
  id: number
  company_id: number
  name: string
  english_name?: string
}

type VoyageItem = {
  id: number
  cruise_id: number
  cruise?: CruiseItem & { company?: CompanyItem }
  code?: string
  brief_info?: string
  image_url?: string
  depart_date?: string
  sold_count?: number
  min_price_cents?: number
}

const emit = defineEmits<{
  (e: 'open-voyage', voyageId: number): void
}>()

const companies = ref<CompanyItem[]>([])
const cruises = ref<CruiseItem[]>([])
const voyages = ref<VoyageItem[]>([])
const expandedCompanyIds = ref<number[]>([])
const selectedCruiseId = ref<number | null>(null)
const keyword = ref('')
const loading = ref(false)
const error = ref('')

function buildVoyagesPath() {
  const params = new URLSearchParams()
  if (selectedCruiseId.value !== null) {
    params.set('cruise_id', String(selectedCruiseId.value))
  }
  const trimmedKeyword = keyword.value.trim()
  if (trimmedKeyword) {
    params.set('keyword', trimmedKeyword)
  }
  params.set('page', '1')
  params.set('page_size', '100')
  return `/voyages?${params.toString()}`
}

const companyTree = computed(() => {
  return companies.value.map((company) => ({
    company,
    cruises: cruises.value.filter((cruise) => cruise.company_id === company.id),
  }))
})

const cruiseNameMap = computed(() => {
  return new Map(cruises.value.map((item) => [item.id, item.name]))
})

const visibleVoyages = computed(() => {
  let rows = [...voyages.value]

  rows.sort((left, right) => {
    const leftCruiseName = cruiseNameMap.value.get(left.cruise_id) || ''
    const rightCruiseName = cruiseNameMap.value.get(right.cruise_id) || ''
    if (leftCruiseName !== rightCruiseName) {
      return leftCruiseName.localeCompare(rightCruiseName, 'zh-CN')
    }
    return String(left.depart_date || '').localeCompare(String(right.depart_date || ''))
  })

  return rows
})

async function loadCatalogData() {
  loading.value = true
  error.value = ''
  try {
    const [companyRes, cruiseRes] = await Promise.all([
      request('/companies?page=1&page_size=50'),
      request('/cruises?page=1&page_size=100'),
    ])
    const companyPayload = (companyRes as any)?.data ?? companyRes ?? {}
    const cruisePayload = (cruiseRes as any)?.data ?? cruiseRes ?? {}

    companies.value = Array.isArray(companyPayload) ? companyPayload : companyPayload?.list ?? []
    cruises.value = Array.isArray(cruisePayload) ? cruisePayload : cruisePayload?.list ?? []
  } catch (requestError: any) {
    error.value = requestError?.message ?? '全部商品加载失败'
    companies.value = []
    cruises.value = []
  } finally {
    loading.value = false
  }
}

async function loadVoyages() {
  loading.value = true
  error.value = ''
  try {
    const voyageRes = await request(buildVoyagesPath())
    const voyagePayload = (voyageRes as any)?.data ?? voyageRes ?? {}
    voyages.value = Array.isArray(voyagePayload) ? voyagePayload : voyagePayload?.list ?? []
  } catch (requestError: any) {
    error.value = requestError?.message ?? '航次商品加载失败'
    voyages.value = []
  } finally {
    loading.value = false
  }
}

function toggleCompany(companyId: number) {
  if (expandedCompanyIds.value.includes(companyId)) {
    expandedCompanyIds.value = expandedCompanyIds.value.filter((id) => id !== companyId)
    return
  }
  expandedCompanyIds.value = [...expandedCompanyIds.value, companyId]
}

function selectCruise(cruiseId: number) {
  selectedCruiseId.value = cruiseId
}

function openVoyage(voyageId: number) {
  emit('open-voyage', voyageId)
}

watch([selectedCruiseId, keyword], async () => {
  await loadVoyages()
})

onMounted(async () => {
  await Promise.all([loadCatalogData(), loadVoyages()])
})
</script>

<template>
  <div class="min-h-screen bg-[#f2f2f2] pb-[92px]">
    <NavBar title="全部商品" />
    <div class="bg-white px-4 py-3 shadow-sm">
      <label class="flex h-10 items-center gap-2 rounded-full bg-[#f4f4f4] px-4 text-slate-400">
        <Search class="h-4 w-4" />
        <input v-model="keyword" class="w-full border-0 bg-transparent text-[14px] text-slate-700 outline-none" placeholder="请输入搜索的商品" />
      </label>
    </div>

    <section class="flex min-h-[calc(100vh-162px)] bg-[#f2f2f2]">
      <ProductCompanyTree
        :items="companyTree"
        :expanded-company-ids="expandedCompanyIds"
        :selected-cruise-id="selectedCruiseId"
        :loading="loading"
        @toggle-company="toggleCompany"
        @select-cruise="selectCruise"
      />

      <div class="min-w-0 flex-1 px-2.5 py-2.5">
        <div v-if="error" class="rounded-xl border border-rose-200 bg-rose-50 px-3 py-2 text-[13px] text-rose-600">
          {{ error }}
        </div>

        <div v-else-if="loading" class="space-y-3">
          <div v-for="item in 3" :key="item" class="overflow-hidden rounded-[16px] border border-slate-200 bg-white shadow-sm">
            <div class="h-[124px] animate-pulse bg-slate-200"></div>
            <div class="space-y-2 px-3 py-3">
              <div class="h-5 w-4/5 animate-pulse rounded bg-slate-200"></div>
              <div class="h-4 w-2/5 animate-pulse rounded bg-slate-100"></div>
              <div class="h-7 w-1/3 animate-pulse rounded bg-slate-100"></div>
            </div>
          </div>
        </div>

        <div v-else-if="visibleVoyages.length === 0" class="rounded-[16px] border border-dashed border-slate-300 bg-white px-4 py-10 text-center text-[14px] text-slate-500">
          当前条件下暂无航次商品
        </div>

        <div v-else class="space-y-3">
          <VoyageProductCard
            v-for="item in visibleVoyages"
            :key="item.id"
            :item="item"
            @select="openVoyage"
          />
        </div>
      </div>
    </section>

    <button
      type="button"
      class="fixed bottom-[108px] right-5 flex h-12 w-12 items-center justify-center rounded-full border border-slate-200 bg-white text-slate-500 shadow-[0_6px_18px_rgba(15,23,42,0.16)] transition-smooth hover:text-[#0f5ba9]"
      aria-label="联系客服"
    >
      <Headset class="h-5 w-5" />
    </button>
  </div>
</template>