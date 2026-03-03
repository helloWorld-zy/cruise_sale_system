<script setup lang="ts">
import { onMounted, ref } from 'vue'

const { request } = useApi()
const route = useRoute()
const skuId = Number(route?.query?.skuId ?? 0)

type PriceRow = { date: string; occupancy: number; price_cents: number; price_type: string }

const rows = ref<PriceRow[]>([])
const loading = ref(false)
const submitting = ref(false)
const error = ref<string | null>(null)
const activeType = ref('base')
const showBatchModal = ref(false)
const form = ref({ date: '', price_cents: 0, price_type: 'base' })
const batch = ref({ start_date: '', end_date: '', price_cents: 0 })
const types = [
  { key: 'base', label: '基础' },
  { key: 'child', label: '儿童' },
  { key: 'single_supplement', label: '单人补差' },
  { key: 'holiday', label: '节假日' },
  { key: 'early_bird', label: '早鸟' },
]

async function loadRows() {
  if (!skuId) {
    error.value = '缺少 skuId 参数'
    return
  }
  loading.value = true
  error.value = null
  try {
    const res = await request(`/cabins/${skuId}/prices`)
    const raw = res?.data ?? res ?? []
    rows.value = raw.map((item: any) => ({
      date: item.date,
      occupancy: Number(item.occupancy ?? 2),
      price_cents: Number(item.price_cents ?? item.price ?? 0),
      price_type: item.price_type || 'base',
    }))
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load prices'
  } finally {
    loading.value = false
  }
}

async function submitPrice() {
  if (!skuId || submitting.value || !form.value.date) return
  submitting.value = true
  error.value = null
  try {
    await request(`/cabins/${skuId}/prices`, {
      method: 'POST',
      body: {
        date: form.value.date,
        price_cents: Number(form.value.price_cents),
        price_type: form.value.price_type || activeType.value,
      },
    })
    await loadRows()
  } catch (e: any) {
    error.value = e?.message ?? 'failed to save price'
  } finally {
    submitting.value = false
  }
}

function rowsByType() {
  return rows.value.filter((row) => row.price_type === activeType.value)
}

function findPrice(date: string) {
  const row = rowsByType().find((item) => normalizeDate(item.date) === date)
  return row ? Number(row.price_cents) : 0
}

function normalizeDate(raw: string) {
  if (!raw) return ''
  return raw.slice(0, 10)
}

function monthDays() {
  const baseDate = new Date()
  const year = baseDate.getFullYear()
  const month = baseDate.getMonth()
  const end = new Date(year, month + 1, 0).getDate()
  const cells: string[] = []
  for (let day = 1; day <= end; day++) {
    const d = new Date(year, month, day)
    const yyyy = d.getFullYear()
    const mm = `${d.getMonth() + 1}`.padStart(2, '0')
    const dd = `${d.getDate()}`.padStart(2, '0')
    cells.push(`${yyyy}-${mm}-${dd}`)
  }
  return cells
}

function cellClass(date: string) {
  const value = findPrice(date)
  if (value > 0) return 'bg-emerald-50 text-emerald-700'
  return 'bg-slate-50 text-slate-500'
}

async function submitBatch() {
  if (!skuId || submitting.value) return
  const start = new Date(batch.value.start_date)
  const end = new Date(batch.value.end_date)
  if (Number.isNaN(start.getTime()) || Number.isNaN(end.getTime()) || end < start) {
    error.value = '无效日期范围'
    return
  }
  submitting.value = true
  error.value = null
  try {
    const cursor = new Date(start)
    while (cursor <= end) {
      const yyyy = cursor.getFullYear()
      const mm = `${cursor.getMonth() + 1}`.padStart(2, '0')
      const dd = `${cursor.getDate()}`.padStart(2, '0')
      await request(`/cabins/${skuId}/prices`, {
        method: 'POST',
        body: {
          date: `${yyyy}-${mm}-${dd}`,
          price_cents: Number(batch.value.price_cents),
          price_type: activeType.value,
        },
      })
      cursor.setDate(cursor.getDate() + 1)
    }
    showBatchModal.value = false
    await loadRows()
  } catch (e: any) {
    error.value = e?.message ?? 'failed to batch save prices'
  } finally {
    submitting.value = false
  }
}

onMounted(loadRows)
</script>

<template>
  <div class="min-h-screen bg-slate-50 p-4 md:p-6">
    <div class="mx-auto max-w-6xl space-y-4">
      <h1 class="text-xl font-semibold text-slate-900">价格矩阵</h1>
      <p v-if="loading" class="text-sm text-slate-600">加载中...</p>
      <p v-else-if="error" class="text-sm text-rose-500">{{ error }}</p>
      <div v-else class="space-y-4">
        <div class="flex flex-wrap items-end justify-between gap-3 rounded-lg border border-slate-200 bg-white p-3 shadow-sm">
          <div class="flex flex-wrap items-center gap-2">
            <button
              v-for="type in types"
              :key="type.key"
              type="button"
              class="border-b-2 px-2 py-1 text-sm"
              :class="activeType === type.key ? 'border-indigo-600 text-indigo-600' : 'border-transparent text-slate-500 hover:text-slate-700'"
              @click="activeType = type.key; form.price_type = type.key"
            >
              {{ type.label }}
            </button>
          </div>
          <button type="button" class="rounded-md bg-indigo-600 px-3 py-2 text-sm text-white hover:bg-indigo-500" @click="showBatchModal = true">批量设价</button>
        </div>

        <div class="rounded-lg border border-slate-200 bg-white p-4 shadow-sm">
          <div class="mb-3 grid grid-cols-7 gap-2 text-center text-xs text-slate-500">
            <span>一</span><span>二</span><span>三</span><span>四</span><span>五</span><span>六</span><span>日</span>
          </div>
          <div class="grid grid-cols-7 gap-2">
            <div v-for="date in monthDays()" :key="date" class="rounded-md px-2 py-2 text-center text-xs ring-1 ring-transparent" :class="cellClass(date)">
              <p>{{ date.slice(-2) }}</p>
              <p>{{ findPrice(date) > 0 ? Math.round(findPrice(date) / 100) : '-' }}</p>
            </div>
          </div>
        </div>

        <div class="rounded-lg border border-slate-200 bg-white p-4 shadow-sm">
          <div class="flex flex-wrap items-center gap-2">
            <input v-model="form.date" type="date" class="h-10 rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
            <input v-model.number="form.price_cents" type="number" min="0" placeholder="price_cents" class="h-10 rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
            <button type="button" class="rounded-md border border-slate-200 px-3 py-2 text-sm text-slate-700 hover:bg-slate-50" :disabled="submitting" @click="submitPrice">{{ submitting ? '提交中...' : '保存单日价格' }}</button>
          </div>
        </div>
      </div>

      <div v-if="showBatchModal" class="fixed inset-0 z-30 flex items-center justify-center bg-slate-900/30">
        <div class="w-full max-w-md rounded-lg border border-slate-200 bg-white p-4 shadow-lg">
          <h2 class="mb-3 text-base font-semibold text-slate-900">批量设价</h2>
          <div class="space-y-3">
            <label class="block space-y-1 text-sm text-slate-600"><span>开始日期</span><input v-model="batch.start_date" type="date" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" /></label>
            <label class="block space-y-1 text-sm text-slate-600"><span>结束日期</span><input v-model="batch.end_date" type="date" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" /></label>
            <label class="block space-y-1 text-sm text-slate-600"><span>价格（分）</span><input v-model.number="batch.price_cents" type="number" min="0" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" /></label>
          </div>
          <div class="mt-4 flex justify-end gap-2">
            <button type="button" class="rounded-md border border-slate-200 px-3 py-2 text-sm text-slate-700 hover:bg-slate-50" @click="showBatchModal = false">取消</button>
            <button type="button" class="rounded-md bg-indigo-600 px-3 py-2 text-sm text-white hover:bg-indigo-500" :disabled="submitting" @click="submitBatch">提交</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

