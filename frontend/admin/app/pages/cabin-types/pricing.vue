<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'

const { request } = useApi()

const companies = ref<Record<string, any>[]>([])
const cruises = ref<Record<string, any>[]>([])
const voyages = ref<Record<string, any>[]>([])
const cabinTypes = ref<Record<string, any>[]>([])
const selectedVoyageIds = ref<number[]>([])
const applyingField = ref<'inventory_total' | 'settlement_price_cents' | 'sale_price_cents' | 'effective_at' | null>(null)

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

const batchForm = ref({
  inventory_total: '',
  settlement_price_cents: '',
  sale_price_cents: '',
  effective_at: '',
})

const rowForms = ref<Record<number, {
  inventory_total: string
  settlement_price_cents: string
  sale_price_cents: string
  effective_at: string
}>>({})

type RowForm = {
  inventory_total: string
  settlement_price_cents: string
  sale_price_cents: string
  effective_at: string
}

const history = ref<Record<string, any>[]>([])

const selectedVoyageRows = computed(() => {
  const set = new Set(selectedVoyageIds.value)
  return voyages.value.filter((item) => set.has(Number(item.id)))
})

const cruiseMap = computed(() => {
  const map = new Map<number, Record<string, any>>()
  for (const cruise of cruises.value) {
    const id = Number(cruise.id)
    if (Number.isFinite(id) && id > 0) map.set(id, cruise)
  }
  return map
})

const historyCache = new Map<string, Record<string, any> | null>()

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

function voyageCode(voyage: Record<string, any>) {
  return voyage.code || voyage.voyage_no || `航次 #${voyage.id}`
}

function voyageIntro(voyage: Record<string, any>) {
  return voyage.brief_info || voyage.name || '暂无简介'
}

function formatDate(value: unknown) {
  const text = String(value || '')
  if (!text) return '--'
  return text.slice(0, 10)
}

function voyageImage(voyage: Record<string, any>) {
  const cruise = cruiseMap.value.get(Number(voyage.cruise_id))
  return (
    cruise?.cover_url ||
    cruise?.image_url ||
    cruise?.cover_image ||
    voyage.cover_url ||
    voyage.image_url ||
    'data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" width="320" height="180"><defs><linearGradient id="g" x1="0" x2="1" y1="0" y2="1"><stop stop-color="%23dbeafe" offset="0"/><stop stop-color="%23bfdbfe" offset="1"/></linearGradient></defs><rect width="100%" height="100%" fill="url(%23g)"/><text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle" font-family="Arial" font-size="18" fill="%23334155">Voyage</text></svg>'
  )
}

function ensureRowForm(cabinTypeID: number) {
  const id = Number(cabinTypeID)
  if (!Number.isFinite(id) || id <= 0) return
  if (rowForms.value[id]) return
  rowForms.value[id] = {
    inventory_total: '',
    settlement_price_cents: '',
    sale_price_cents: '',
    effective_at: '',
  }
}

function getRowForm(cabinTypeID: number): RowForm {
  ensureRowForm(cabinTypeID)
  return rowForms.value[cabinTypeID] as RowForm
}

function toNumberOrNull(raw: string | number | null | undefined) {
  const value = String(raw ?? '').trim()
  if (!value) return null
  const num = Number(value)
  if (!Number.isFinite(num)) return null
  return num
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

function cacheKey(voyageID: number, cabinTypeID: number) {
  return `${voyageID}:${cabinTypeID}`
}

async function loadLatestVersion(voyageID: number, cabinTypeID: number) {
  const key = cacheKey(voyageID, cabinTypeID)
  if (historyCache.has(key)) return historyCache.get(key)

  const res = await request('/cabin-pricing/history', {
    query: {
      voyage_id: voyageID,
      cabin_type_id: cabinTypeID,
      page: 1,
      page_size: 1,
    },
  })
  const payload = res?.data ?? res ?? {}
  const list = Array.isArray(payload) ? payload : payload?.list ?? []
  const latest = list.length > 0 ? list[0] : null
  historyCache.set(key, latest)
  return latest
}

function clearHistoryCache() {
  historyCache.clear()
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
    clearHistoryCache()

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
  const selectedSet = new Set(selectedVoyageIds.value)
  for (const voyage of voyages.value) {
    if (!selectedSet.has(Number(voyage.id))) continue
    const cruiseID = Number(voyage.cruise_id)
    if (Number.isFinite(cruiseID) && cruiseID > 0) cruiseIDSet.add(cruiseID)
  }

  if (cruiseIDSet.size === 0) {
    cabinTypes.value = []
    return
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
  for (const item of cabinTypes.value) {
    ensureRowForm(Number(item.id))
  }
}

async function applyForCabinType(
  cabinTypeID: number,
  updates: {
    inventory_total?: number
    settlement_price_cents?: number
    sale_price_cents?: number
    effective_at?: string
  },
) {
  const voyageIDs = selectedVoyageIds.value.slice()
  if (voyageIDs.length === 0) return

  const failures: string[] = []
  const normalizedEffectiveAt = updates.effective_at ? normalizeEffectiveAt(updates.effective_at) : undefined

  for (const voyageID of voyageIDs) {
    const needsBase =
      updates.inventory_total === undefined ||
      updates.settlement_price_cents === undefined ||
      updates.sale_price_cents === undefined

    let base: Record<string, any> | null = null
    if (needsBase) {
      base = await loadLatestVersion(voyageID, cabinTypeID)
      if (!base) {
        failures.push(`航次 #${voyageID} 的舱型 #${cabinTypeID} 无历史价格，无法仅更新单项`) 
        continue
      }
    }

    const inventory = updates.inventory_total ?? Number(base?.inventory_total)
    const settlement = updates.settlement_price_cents ?? Number(base?.settlement_price_cents)
    const sale = updates.sale_price_cents ?? Number(base?.sale_price_cents)

    if (!Number.isFinite(inventory) || inventory < 0 || !Number.isFinite(settlement) || settlement < 0 || !Number.isFinite(sale) || sale < 0) {
      failures.push(`航次 #${voyageID} 的舱型 #${cabinTypeID} 价格参数无效`)
      continue
    }

    await request('/cabin-pricing/batch-apply', {
      method: 'POST',
      body: {
        voyage_ids: [voyageID],
        cabin_type_id: cabinTypeID,
        inventory_total: inventory,
        settlement_price_cents: settlement,
        sale_price_cents: sale,
        effective_at: normalizedEffectiveAt,
      },
    })
  }

  if (failures.length > 0) {
    error.value = failures[0] || '部分航次应用失败'
  }
}

async function applyRow(cabinTypeID: number) {
  if (applying.value) return
  const row = rowForms.value[cabinTypeID]
  if (!row) return
  if (selectedVoyageIds.value.length === 0) {
    error.value = '请先选择至少一个航次'
    return
  }

  const inventory = toNumberOrNull(row.inventory_total)
  const settlement = toNumberOrNull(row.settlement_price_cents)
  const sale = toNumberOrNull(row.sale_price_cents)
  const effectiveAt = row.effective_at.trim()

  const updates: {
    inventory_total?: number
    settlement_price_cents?: number
    sale_price_cents?: number
    effective_at?: string
  } = {}

  if (inventory !== null) {
    if (inventory < 0) {
      error.value = '库存总量不能小于 0'
      return
    }
    updates.inventory_total = inventory
  }
  if (settlement !== null) {
    if (settlement < 0) {
      error.value = '结算价不能小于 0'
      return
    }
    updates.settlement_price_cents = settlement
  }
  if (sale !== null) {
    if (sale < 0) {
      error.value = '售价不能小于 0'
      return
    }
    updates.sale_price_cents = sale
  }
  if (effectiveAt) {
    updates.effective_at = effectiveAt
  }

  if (Object.keys(updates).length === 0) {
    error.value = '请至少填写一项后再保存'
    return
  }

  applying.value = true
  error.value = null
  try {
    await applyForCabinType(cabinTypeID, updates)
    await loadHistoryFor(Number(selectedVoyageIds.value[0] || 0), cabinTypeID)
  } catch (e: any) {
    error.value = e?.message ?? 'failed to apply prices'
  } finally {
    applying.value = false
  }
}

async function applyBatchField(field: 'inventory_total' | 'settlement_price_cents' | 'sale_price_cents' | 'effective_at') {
  if (applying.value || applyingField.value) return
  if (selectedVoyageIds.value.length === 0) {
    error.value = '请先选择至少一个航次'
    return
  }
  if (cabinTypes.value.length === 0) {
    error.value = '当前没有可设置舱型'
    return
  }

  const updates: {
    inventory_total?: number
    settlement_price_cents?: number
    sale_price_cents?: number
    effective_at?: string
  } = {}

  if (field === 'inventory_total') {
    const v = toNumberOrNull(batchForm.value.inventory_total)
    if (v === null || v < 0) {
      error.value = '请填写有效库存总量'
      return
    }
    updates.inventory_total = v
  }
  if (field === 'settlement_price_cents') {
    const v = toNumberOrNull(batchForm.value.settlement_price_cents)
    if (v === null || v < 0) {
      error.value = '请填写有效结算价'
      return
    }
    updates.settlement_price_cents = v
  }
  if (field === 'sale_price_cents') {
    const v = toNumberOrNull(batchForm.value.sale_price_cents)
    if (v === null || v < 0) {
      error.value = '请填写有效售价'
      return
    }
    updates.sale_price_cents = v
  }
  if (field === 'effective_at') {
    if (!batchForm.value.effective_at.trim()) {
      error.value = '请填写生效时间'
      return
    }
    updates.effective_at = batchForm.value.effective_at
  }

  applying.value = true
  applyingField.value = field
  error.value = null
  try {
    for (const item of cabinTypes.value) {
      await applyForCabinType(Number(item.id), updates)
    }
    if (cabinTypes.value.length > 0) {
      await loadHistoryFor(Number(selectedVoyageIds.value[0] || 0), Number(cabinTypes.value[0]?.id || 0))
    }
  } catch (e: any) {
    error.value = e?.message ?? 'failed to apply batch field'
  } finally {
    applyingField.value = null
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
  () => selectedVoyageIds.value.join(','),
  async () => {
    await loadCabinTypesByVoyages()
  },
)

watch(
  () => [cabinTypes.value[0]?.id, selectedVoyageIds.value[0]],
  async ([cabinTypeID, voyageID]) => {
    await loadHistoryFor(Number(voyageID || 0), Number(cabinTypeID || 0))
  },
)
</script>

<template>
  <div class="admin-page">
    <AdminPageHeader title="舱型价格管理" subtitle="按公司与航次筛选后，可批量或按舱型逐项设置库存与价格。">
      <template #actions>
        <AdminActionLink to="/cabin-types">返回舱型管理</AdminActionLink>
      </template>
    </AdminPageHeader>

    <AdminFormCard title="价格工作台">
      <div class="admin-cruise-form">
        <section class="admin-cruise-form__intro">
          <h2 class="admin-cruise-form__intro-title">统一价格录入流程</h2>
          <p class="admin-cruise-form__intro-desc">先筛选航次，再选择目标航次与舱型，支持批量字段应用和逐行精细调整。</p>
        </section>

        <section class="admin-cruise-form__section">
          <h3 class="admin-cruise-form__section-title">筛选航次</h3>
          <p class="admin-cruise-form__section-subtitle">按公司、邮轮和出发日期范围快速锁定目标航次。</p>
          <div class="admin-cruise-form__grid pricing-filter-grid">
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">公司</span>
              <select v-model.number="filters.company_id" class="admin-cruise-form__control">
              <option :value="0">全部公司</option>
              <option v-for="company in companies" :key="company.id" :value="Number(company.id)">{{ company.name || `公司 #${company.id}` }}</option>
            </select>
          </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">邮轮</span>
              <select v-model.number="filters.cruise_id" class="admin-cruise-form__control">
              <option :value="0">全部邮轮</option>
              <option v-for="cruise in cruises" :key="cruise.id" :value="Number(cruise.id)">{{ cruise.name || `邮轮 #${cruise.id}` }}</option>
            </select>
          </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">出发开始</span>
              <input v-model="filters.depart_start" type="date" class="admin-cruise-form__control" />
          </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">出发结束</span>
              <input v-model="filters.depart_end" type="date" class="admin-cruise-form__control" />
          </label>
          </div>
          <div class="admin-cruise-form__actions">
            <button type="button" class="admin-btn admin-btn--secondary" @click="loadVoyages">刷新航次</button>
          </div>
        </section>

        <section class="admin-cruise-form__section">
          <h3 class="admin-cruise-form__section-title">可选航次（可多选）</h3>
          <p class="admin-cruise-form__section-subtitle">已选 {{ selectedVoyageIds.length }} 条，选中后自动加载对应舱型。</p>
          <p v-if="loading" class="text-sm text-slate-500">加载中...</p>
          <p v-else-if="voyages.length === 0" class="text-sm text-slate-500">暂无可用航次</p>
          <div v-else class="voyage-cards-grid max-h-[460px] overflow-auto pr-1">
            <label
              v-for="voyage in voyages"
              :key="voyage.id"
              class="voyage-card cursor-pointer"
              :class="{ 'voyage-card--selected': selectedVoyageIds.includes(Number(voyage.id)) }"
            >
              <div class="voyage-card__cover-wrap">
                <img class="voyage-card__cover" :src="voyageImage(voyage)" alt="航次图片" />
              </div>
              <div class="voyage-card__body">
                <div class="voyage-card__top">
                  <span class="voyage-card__cruise">{{ cruiseName(voyage.cruise_id) }}</span>
                </div>
                <p class="voyage-card__code">{{ voyageCode(voyage) }}</p>
                <p class="voyage-card__intro">{{ voyageIntro(voyage) }}</p>
                <p class="voyage-card__dates">{{ formatDate(voyage.depart_date) }} - {{ formatDate(voyage.return_date) }}</p>
              </div>
              <div class="voyage-card__check">
                <input
                  type="checkbox"
                  :checked="selectedVoyageIds.includes(Number(voyage.id))"
                  @change="toggleVoyage(Number(voyage.id), ($event.target as HTMLInputElement).checked)"
                />
              </div>
            </label>
          </div>
        </section>

        <section class="admin-cruise-form__section">
          <h3 class="admin-cruise-form__section-title">批量设置</h3>
          <p class="admin-cruise-form__section-subtitle">以下为所有舱型批量生效，四项可分别单独应用。</p>
          <div class="pricing-batch-grid">
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">库存总量</span>
              <div class="pricing-inline-inputs">
                <input v-model="batchForm.inventory_total" type="number" min="0" class="admin-cruise-form__control" />
                <button type="button" class="admin-btn pricing-apply-btn" :disabled="applying" @click="applyBatchField('inventory_total')">应用</button>
              </div>
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">结算价（分）</span>
              <div class="pricing-inline-inputs">
                <input v-model="batchForm.settlement_price_cents" type="number" min="0" class="admin-cruise-form__control" />
                <button type="button" class="admin-btn pricing-apply-btn" :disabled="applying" @click="applyBatchField('settlement_price_cents')">应用</button>
              </div>
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">售价（分）</span>
              <div class="pricing-inline-inputs">
                <input v-model="batchForm.sale_price_cents" type="number" min="0" class="admin-cruise-form__control" />
                <button type="button" class="admin-btn pricing-apply-btn" :disabled="applying" @click="applyBatchField('sale_price_cents')">应用</button>
              </div>
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">生效时间（可选）</span>
              <div class="pricing-inline-inputs">
                <input v-model="batchForm.effective_at" type="datetime-local" class="admin-cruise-form__control" />
                <button type="button" class="admin-btn pricing-apply-btn" :disabled="applying" @click="applyBatchField('effective_at')">应用</button>
              </div>
            </label>
          </div>
          <p class="mt-2 text-xs text-slate-500">只填某一项并点击该项“应用”，仅更新该项；其他留空不会被清空。</p>
        </section>

        <section class="admin-cruise-form__section">
          <h3 class="admin-cruise-form__section-title">价格设置（基于已选航次）</h3>
          <p class="admin-cruise-form__section-subtitle">支持按舱型逐行录入并立即保存。</p>
        <p v-if="selectedVoyageIds.length === 0" class="text-sm text-slate-500">请先在上方选择至少一个航次</p>
        <p v-else-if="cabinTypes.length === 0" class="text-sm text-slate-500">当前选中航次暂无可设置舱型</p>
        <div v-else class="overflow-x-auto">
            <table class="pricing-table w-full min-w-[1100px] text-sm">
            <thead class="bg-slate-50 text-left text-slate-600">
              <tr>
                <th class="p-3">舱型</th>
                <th class="p-3">库存总量</th>
                <th class="p-3">结算价（分）</th>
                <th class="p-3">售价（分）</th>
                <th class="p-3">生效时间</th>
                <th class="p-3">操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in cabinTypes" :key="item.id" class="border-t border-slate-100">
                <td class="p-3 text-slate-700">{{ cabinTypeLabel(item) }}</td>
                <td class="p-3"><input v-model="getRowForm(Number(item.id)).inventory_total" type="number" min="0" class="admin-cruise-form__control pricing-row-input" /></td>
                <td class="p-3"><input v-model="getRowForm(Number(item.id)).settlement_price_cents" type="number" min="0" class="admin-cruise-form__control pricing-row-input" /></td>
                <td class="p-3"><input v-model="getRowForm(Number(item.id)).sale_price_cents" type="number" min="0" class="admin-cruise-form__control pricing-row-input" /></td>
                <td class="p-3"><input v-model="getRowForm(Number(item.id)).effective_at" type="datetime-local" class="admin-cruise-form__control pricing-row-input" /></td>
                <td class="p-3 whitespace-nowrap">
                  <button type="button" class="admin-btn pricing-save-btn" :disabled="applying" @click="applyRow(Number(item.id))">保存</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        </section>

        <section class="admin-cruise-form__section">
          <h3 class="admin-cruise-form__section-title">价格历史（当前选中第一条航次 + 第一个舱型）</h3>
        <p v-if="historyLoading" class="text-sm text-slate-500">加载中...</p>
        <p v-else-if="history.length === 0" class="text-sm text-slate-500">暂无历史记录</p>
        <div v-else class="overflow-x-auto">
            <table class="pricing-table w-full text-sm">
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

        <section v-if="selectedVoyageRows.length > 0" class="admin-cruise-form__section">
          <h3 class="admin-cruise-form__section-title">已选航次预览</h3>
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
    </AdminFormCard>
  </div>
</template>

<style scoped>
.pricing-filter-grid {
  grid-template-columns: repeat(1, minmax(0, 1fr));
}

.pricing-batch-grid {
  display: grid;
  gap: 14px;
  grid-template-columns: repeat(1, minmax(0, 1fr));
}

.pricing-inline-inputs {
  display: flex;
  gap: 8px;
  align-items: stretch;
}

.pricing-apply-btn,
.pricing-save-btn {
  min-width: 76px;
  min-height: 40px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.pricing-save-btn {
  min-width: 84px;
}

.pricing-row-input {
  min-height: 36px;
  padding: 0 10px;
}

.pricing-table {
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  overflow: hidden;
}

.voyage-cards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 10px;
}

.voyage-card {
  display: grid;
  grid-template-columns: 72px 1fr auto;
  gap: 10px;
  align-items: center;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  padding: 8px;
  background: #fff;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

.voyage-card:hover {
  border-color: #93c5fd;
  box-shadow: 0 3px 10px rgba(59, 130, 246, 0.1);
}

.voyage-card--selected {
  border-color: #3b82f6;
  background: #f8fbff;
}

.voyage-card__cover-wrap {
  width: 72px;
  height: 72px;
  border-radius: 8px;
  overflow: hidden;
  background: #e2e8f0;
}

.voyage-card__cover {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.voyage-card__body {
  min-width: 0;
}

.voyage-card__top {
  margin-bottom: 2px;
}

.voyage-card__code {
  margin: 0;
  font-size: 13px;
  font-weight: 600;
  color: #0f172a;
  line-height: 1.25;
  word-break: break-word;
}

.voyage-card__cruise {
  font-size: 12px;
  color: #475569;
  line-height: 1.25;
  word-break: break-word;
}

.voyage-card__intro,
.voyage-card__dates {
  margin: 2px 0 0;
  font-size: 12px;
  color: #475569;
  line-height: 1.3;
  word-break: break-word;
}

.voyage-card__intro {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

@media (min-width: 768px) {
  .pricing-filter-grid,
  .pricing-batch-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
</style>
