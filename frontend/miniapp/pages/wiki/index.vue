<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { Headset } from 'lucide-vue-next'
import { request } from '../../src/utils/request'
import WikiCompanySidebar from '../../components/wiki/WikiCompanySidebar.vue'
import WikiCruiseCard from '../../components/wiki/WikiCruiseCard.vue'

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
  cover_url?: string
}

const emit = defineEmits<{
  (e: 'open-cruise', id: number): void
}>()

const companies = ref<CompanyItem[]>([])
const cruises = ref<CruiseItem[]>([])
const selectedCompanyId = ref<number | 'all'>('all')
const loadingCompanies = ref(false)
const loadingCruises = ref(false)
const companiesError = ref('')
const cruisesError = ref('')

const sidebarItems = computed(() => [
  { id: 'all' as const, name: '全部邮轮' },
  ...companies.value.map((item) => ({ id: item.id, name: item.name })),
])

function buildCruisesPath(companyId: number | 'all') {
  if (companyId === 'all') {
    return '/cruises?page=1&page_size=30'
  }
  return `/cruises?company_id=${companyId}&page=1&page_size=30`
}

async function loadCompanies() {
  loadingCompanies.value = true
  companiesError.value = ''
  try {
    const res = await request<any>('/companies?page=1&page_size=50')
    const payload = (res as any)?.data ?? res ?? {}
    companies.value = Array.isArray(payload) ? payload : payload?.list ?? []
  } catch (error: any) {
    companiesError.value = error?.message ?? '公司加载失败'
    companies.value = []
  } finally {
    loadingCompanies.value = false
  }
}

async function loadCruises(companyId: number | 'all' = selectedCompanyId.value) {
  loadingCruises.value = true
  cruisesError.value = ''
  try {
    const res = await request<any>(buildCruisesPath(companyId))
    const payload = (res as any)?.data ?? res ?? {}
    cruises.value = Array.isArray(payload) ? payload : payload?.list ?? []
  } catch (error: any) {
    cruisesError.value = error?.message ?? '邮轮加载失败'
    cruises.value = []
  } finally {
    loadingCruises.value = false
  }
}

async function handleSelectCompany(companyId: number | 'all') {
  if (selectedCompanyId.value === companyId) {
    return
  }
  selectedCompanyId.value = companyId
  await loadCruises(companyId)
}

function handleOpenCruise(cruiseId: number) {
  emit('open-cruise', cruiseId)
}

onMounted(async () => {
  await Promise.all([loadCompanies(), loadCruises('all')])
})
</script>

<template>
  <div class="min-h-screen bg-[#f1f1f1] pb-[92px]">
    <header class="bg-white pt-7 shadow-sm">
      <div class="flex items-center justify-between px-3 pb-3">
        <div class="w-16 text-[12px] text-slate-500">WeChat</div>
        <div class="text-[13px] font-semibold tracking-[0.5px] text-slate-700">邮轮百科</div>
        <div class="flex w-16 justify-end">
          <div class="rounded-full border border-slate-300 px-2 py-1 text-[11px] text-slate-700">•••</div>
        </div>
      </div>
      <div class="bg-[#0f5ba9] px-4 py-3 text-center text-[21px] font-bold tracking-[1px] text-white">邮轮百科</div>
    </header>

    <section class="flex min-h-[calc(100vh-152px)] bg-[#f1f1f1]">
      <WikiCompanySidebar
        :items="sidebarItems"
        :selected-company-id="selectedCompanyId"
        :loading="loadingCompanies"
        @select="handleSelectCompany"
      />

      <div class="min-w-0 flex-1 px-2.5 py-2.5">
        <div v-if="companiesError" class="mb-3 rounded-xl border border-rose-200 bg-rose-50 px-3 py-2 text-[13px] text-rose-600">
          {{ companiesError }}
        </div>
        <div v-if="cruisesError" class="rounded-xl border border-rose-200 bg-rose-50 px-3 py-2 text-[13px] text-rose-600">
          {{ cruisesError }}
        </div>
        <div v-else-if="loadingCruises" class="space-y-3">
          <div v-for="item in 3" :key="item" class="overflow-hidden rounded-[12px] border border-slate-200 bg-white shadow-sm">
            <div class="aspect-[4/3] animate-pulse bg-slate-200"></div>
            <div class="space-y-2 px-3 py-3">
              <div class="h-5 w-1/2 animate-pulse rounded bg-slate-200"></div>
              <div class="h-4 w-2/3 animate-pulse rounded bg-slate-100"></div>
            </div>
          </div>
        </div>
        <div v-else-if="cruises.length === 0" class="rounded-[16px] border border-dashed border-slate-300 bg-white px-4 py-10 text-center text-[14px] text-slate-500">
          当前分类下暂无邮轮
        </div>
        <div v-else class="space-y-3">
          <WikiCruiseCard
            v-for="item in cruises"
            :key="item.id"
            :item="item"
            @select="handleOpenCruise"
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