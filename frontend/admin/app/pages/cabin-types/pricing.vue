<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'

const { request } = useApi()

const companies = ref<Record<string, any>[]>([])
const cruises = ref<Record<string, any>[]>([])
const voyages = ref<Record<string, any>[]>([])
const cabinTypes = ref<Record<string, any>[]>([])
const selectedVoyageIds = ref<number[]>([])

const loading = ref(false)
const applying = ref(false)
const historyLoading = ref(false)
const error = ref<string | null>(null)

const filters = ref({
  company_id: 0,
  cruise_id: 0,
  depart_start: '',
  depart_end: '',
})

const form = ref({
  cabin_type_id: 0,
  inventory_total: 0,
  settlement_price_cents: 0,
  sale_price_cents: 0,
  effective_at: '',
})

const history = ref<Record<string, any>[]>([])

const canApply = computed(() => {
  return (
    selectedVoyageIds.value.length > 0 &&
    form.value.cabin_type_id > 0 &&
    Number(form.value.inventory_total) >= 0 &&
    Number(form.value.settlement_price_cents) > 0 &&
    Number(form.value.sale_price_cents) > 0
  )
})

const selectedVoyageRows = computed(() => {
  const set = new Set(selectedVoyageIds.value)
  return voyages.value.filter((item) => set.has(Number(item.id)))
})

function cruiseName(cruiseID: unknown) {
  const id = Number(cruiseID)
  const found = cruises.value.find((item) => Number(item.id) === id)
  return found?.name || (id > 0 ? `邮轮 #${id}` : '-')
}

function voyageLabel(voyage: Record<string, any>) {
  const departDate = String(voyage.depart_date || '').slice(0, 10)
  const name = voyage.name || voyage.voyage_no || `航次 #${voyage.id}`
  return `${name} (${departDate || '未设日期'})`
}

function cabinTypeLabel(item: Record<string, any>) {
  const code = item.code ? ` [${item.code}]` : ''
  const category = item.category_name ? `${item.category_name} / ` : ''
  return `${category}${item.name || `舱型 #${item.id}`}${code}`
}

function toggleVoyage(voyageID: number, checked: boolean) {
  const next = new Set(selectedVoyageIds.value)
  if (checked) next.add(voyageID)
  else next.delete(voyageID)
  selectedVoyageIds.value = Array.from(next)
}

function normalizeEffectiveAt(raw: string) {
  const value = raw.trim()
  if (!value) return undefined
  if (value.includes('T')) {
    const withSeconds = value.length === 16 ? `${value}:00` : value
    return withSeconds.replace('T', ' ')
  }
  return value
}

async function loadCompanies() {
  try {
    const res = await request('/companies', { query: { page: 1, page_size: 100 } })
    const payload = res?.data ?? res ?? {}
    companies.value = Array.isArray(payload) ? payload : payload?.list ?? []
    const firstCompany = companies.value[0]
    if (firstCompany && filters.value.company_id <= 0) {
      const first = Number(firstCompany.id)
      if (Number.isFinite(first) && first > 0) filters.value.company_id = first
    }
  } catch {
    companies.value = []
  }
}

async function loadCruises() {
  try {
    const query: Record<string, any> = { page: 1, page_size: 300 }
    if (filters.value.company_id > 0) query.company_id = filters.value.company_id
    const res = await request('/cruises', { query })
    const payload = res?.data ?? res ?? {}
    cruises.value = Array.isArray(payload) ? payload : payload?.list ?? []
    if (filters.value.cruise_id > 0 && !cruises.value.find((item) => Number(item.id) === filters.value.cruise_id)) {
      filters.value.cruise_id = 0
    }
  } catch {
    cruises.value = []
    filters.value.cruise_id = 0
  }
}

async function loadVoyages() {
  loading.value = true
  error.value = null
  try {
    const query: Record<string, any> = {}
    if (filters.value.company_id > 0) query.company_id = filters.value.company_id
    if (filters.value.cruise_id > 0) query.cruise_id = filters.value.cruise_id
    if (filters.value.depart_start) query.depart_start = filters.value.depart_start
    if (filters.value.depart_end) query.depart_end = filters.value.depart_end

    const res = await request('/cabin-pricing/voyages', { query })
    const payload = res?.data ?? res ?? {}
    voyages.value = Array.isArray(payload) ? payload : payload?.list ?? []

    const allowed = new Set(voyages.value.map((item) => Number(item.id)))
    selectedVoyageIds.value = selectedVoyageIds.value.filter((id) => allowed.has(id))

    await loadCabinTypesByVoyages()
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load voyages'
    voyages.value = []
    selectedVoyageIds.value = []
    cabinTypes.value = []
  } finally {
    loading.value = false
  }
}

async function loadCabinTypesByVoyages() {
  const cruiseIDSet = new Set<number>()
  for (const voyage of voyages.value) {
    const cruiseID = Number(voyage.cruise_id)
    if (Number.isFinite(cruiseID) && cruiseID > 0) cruiseIDSet.add(cruiseID)
  }

  const merged = new Map<number, Record<string, any>>()
  for (const cruiseID of cruiseIDSet) {
    const res = await request('/cabin-types', {
      query: {
        cruise_id: cruiseID,
        page: 1,
        page_size: 500,
      },
    })
    const payload = res?.data ?? res ?? {}
    const list = Array.isArray(payload) ? payload : payload?.list ?? []
    for (const item of list) {
      const id = Number(item.id)
      if (!Number.isFinite(id) || id <= 0) continue
      if (!merged.has(id)) merged.set(id, item)
    }
  }

  cabinTypes.value = Array.from(merged.values())
  if (form.value.cabin_type_id > 0 && !cabinTypes.value.find((item) => Number(item.id) === form.value.cabin_type_id)) {
    form.value.cabin_type_id = 0
  }
}

async function applyPrices() {
  if (!canApply.value || applying.value) return
  applying.value = true
  error.value = null
  try {
    await request('/cabin-pricing/batch-apply', {
      method: 'POST',
      body: {
        voyage_ids: selectedVoyageIds.value,
        cabin_type_id: Number(form.value.cabin_type_id),
        inventory_total: Number(form.value.inventory_total),
        settlement_price_cents: Number(form.value.settlement_price_cents),
        sale_price_cents: Number(form.value.sale_price_cents),
        effective_at: normalizeEffectiveAt(form.value.effective_at),
      },
    })

    if (selectedVoyageIds.value.length > 0 && form.value.cabin_type_id > 0) {
      await loadHistoryFor(Number(selectedVoyageIds.value[0]), Number(form.value.cabin_type_id))
    }
  } catch (e: any) {
    error.value = e?.message ?? 'failed to apply prices'
  } finally {
    applying.value = false
  }
}

async function loadHistoryFor(voyageID: number, cabinTypeID: number) {
  if (voyageID <= 0 || cabinTypeID <= 0) {
    history.value = []
    return
  }
  historyLoading.value = true
  error.value = null
  try {
    const res = await request('/cabin-pricing/history', {
      query: {
        voyage_id: voyageID,
        cabin_type_id: cabinTypeID,
        page: 1,
        page_size: 30,
      },
    })
    const payload = res?.data ?? res ?? {}
    history.value = Array.isArray(payload) ? payload : payload?.list ?? []
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load price history'
    history.value = []
  } finally {
    historyLoading.value = false
  }
}

onMounted(async () => {
  await loadCompanies()
  await loadCruises()
  await loadVoyages()
})

watch(
  () => filters.value.company_id,
  async () => {
    await loadCruises()
    await loadVoyages()
  },
)

watch(
  () => [form.value.cabin_type_id, selectedVoyageIds.value[0]],
  async ([cabinTypeID, voyageID]) => {
    await loadHistoryFor(Number(voyageID || 0), Number(cabinTypeID || 0))
  },
)
</script>

<template>
  <div class="min-h-screen bg-slate-50 p-4 md:p-6">
    <div class="mx-auto max-w-7xl space-y-4">
      <div class="flex items-center justify-between">
        <h1 class="text-xl font-semibold text-slate-900">舱型价格管理</h1>
        <AdminActionLink to="/cabin-types">返回舱型管理</AdminActionLink>
      </div>

      <div class="rounded-lg border border-slate-200 bg-white p-4 shadow-sm">
        <p class="mb-3 text-sm font-medium text-slate-800">筛选航次</p>
        <div class="grid grid-cols-1 gap-3 md:grid-cols-4">
          <label class="space-y-1 text-sm text-slate-600">
            <span>公司</span>
            <select v-model.number="filters.company_id" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2">
              <option :value="0">全部公司</option>
              <option v-for="company in companies" :key="company.id" :value="Number(company.id)">{{ company.name || `公司 #${company.id}` }}</option>
            </select>
          </label>
          <label class="space-y-1 text-sm text-slate-600">
            <span>邮轮</span>
            <select v-model.number="filters.cruise_id" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2">
              <option :value="0">全部邮轮</option>
              <option v-for="cruise in cruises" :key="cruise.id" :value="Number(cruise.id)">{{ cruise.name || `邮轮 #${cruise.id}` }}</option>
            </select>
          </label>
          <label class="space-y-1 text-sm text-slate-600">
            <span>出发开始</span>
            <input v-model="filters.depart_start" type="date" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
          </label>
          <label class="space-y-1 text-sm text-slate-600">
            <span>出发结束</span>
            <input v-model="filters.depart_end" type="date" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
          </label>
        </div>
        <div class="mt-3 flex justify-end">
          <button type="button" class="rounded-md border border-slate-200 px-3 py-2 text-sm text-slate-700 hover:bg-slate-50" @click="loadVoyages">刷新航次</button>
        </div>
      </div>

      <div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
        <section class="rounded-lg border border-slate-200 bg-white p-4 shadow-sm">
          <div class="mb-2 flex items-center justify-between">
            <h2 class="text-sm font-semibold text-slate-800">可选航次（可多选）</h2>
            <span class="text-xs text-slate-500">已选 {{ selectedVoyageIds.length }} 条</span>
          </div>
          <p v-if="loading" class="text-sm text-slate-500">加载中...</p>
          <p v-else-if="voyages.length === 0" class="text-sm text-slate-500">暂无可用航次</p>
          <div v-else class="max-h-80 space-y-2 overflow-auto pr-1">
            <label
              v-for="voyage in voyages"
              :key="voyage.id"
              class="flex items-center justify-between rounded-md border border-slate-200 px-3 py-2 text-sm text-slate-700"
            >
              <div class="flex items-center gap-2">
                <input
                  type="checkbox"
                  :checked="selectedVoyageIds.includes(Number(voyage.id))"
                  @change="toggleVoyage(Number(voyage.id), ($event.target as HTMLInputElement).checked)"
                />
                <span>{{ voyageLabel(voyage) }}</span>
              </div>
              <span class="text-xs text-slate-500">{{ cruiseName(voyage.cruise_id) }}</span>
            </label>
          </div>
        </section>

        <section class="rounded-lg border border-slate-200 bg-white p-4 shadow-sm">
          <h2 class="mb-2 text-sm font-semibold text-slate-800">批量生效设置</h2>
          <div class="grid grid-cols-1 gap-3 md:grid-cols-2">
            <label class="space-y-1 text-sm text-slate-600 md:col-span-2">
              <span>舱型</span>
              <select v-model.number="form.cabin_type_id" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2">
                <option :value="0">请选择舱型</option>
                <option v-for="item in cabinTypes" :key="item.id" :value="Number(item.id)">{{ cabinTypeLabel(item) }}</option>
              </select>
            </label>
            <label class="space-y-1 text-sm text-slate-600">
              <span>库存总量</span>
              <input v-model.number="form.inventory_total" type="number" min="0" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
            </label>
            <label class="space-y-1 text-sm text-slate-600">
              <span>结算价（分）</span>
              <input v-model.number="form.settlement_price_cents" type="number" min="0" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
            </label>
            <label class="space-y-1 text-sm text-slate-600">
              <span>售价（分）</span>
              <input v-model.number="form.sale_price_cents" type="number" min="0" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
            </label>
            <label class="space-y-1 text-sm text-slate-600">
              <span>生效时间（可选）</span>
              <input v-model="form.effective_at" type="datetime-local" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
            </label>
          </div>
          <div class="mt-3 flex items-center justify-between">
            <span class="text-xs text-slate-500">不填生效时间时，按上海当前时间生效</span>
            <button
              type="button"
              :disabled="!canApply || applying"
              class="rounded-md bg-indigo-600 px-3 py-2 text-sm text-white hover:bg-indigo-500 disabled:cursor-not-allowed disabled:opacity-60"
              @click="applyPrices"
            >
              {{ applying ? '提交中...' : '批量应用价格' }}
            </button>
          </div>
        </section>
      </div>

      <section class="rounded-lg border border-slate-200 bg-white p-4 shadow-sm">
        <h2 class="mb-2 text-sm font-semibold text-slate-800">价格历史（当前选中第一条航次）</h2>
        <p v-if="historyLoading" class="text-sm text-slate-500">加载中...</p>
        <p v-else-if="history.length === 0" class="text-sm text-slate-500">暂无历史记录</p>
        <div v-else class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead class="bg-slate-50 text-left text-slate-600">
              <tr>
                <th class="p-3">生效时间</th>
                <th class="p-3">库存</th>
                <th class="p-3">结算价（分）</th>
                <th class="p-3">售价（分）</th>
                <th class="p-3">创建时间</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in history" :key="item.id" class="border-t border-slate-100">
                <td class="p-3 text-slate-700">{{ String(item.effective_at || '').replace('T', ' ').slice(0, 19) }}</td>
                <td class="p-3 text-slate-700">{{ item.inventory_total ?? '-' }}</td>
                <td class="p-3 text-slate-700">{{ item.settlement_price_cents ?? '-' }}</td>
                <td class="p-3 text-slate-700">{{ item.sale_price_cents ?? '-' }}</td>
                <td class="p-3 text-slate-500">{{ String(item.created_at || '').replace('T', ' ').slice(0, 19) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <p v-if="error" class="text-sm text-rose-500">{{ error }}</p>

      <section v-if="selectedVoyageRows.length > 0" class="rounded-lg border border-slate-200 bg-white p-4 shadow-sm">
        <h2 class="mb-2 text-sm font-semibold text-slate-800">已选航次预览</h2>
        <div class="flex flex-wrap gap-2">
          <span
            v-for="voyage in selectedVoyageRows"
            :key="voyage.id"
            class="rounded-full bg-indigo-50 px-2.5 py-1 text-xs font-medium text-indigo-700"
          >
            {{ voyageLabel(voyage) }}
          </span>
        </div>
      </section>
    </div>
  </div>
</template>
