<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import AdminConfirmDialog from '../../components/AdminConfirmDialog.vue'

type ItineraryFormItem = {
  day_no: number
  stop_index: number
  city: string
  summary: string
  eta_time: string
  etd_time: string
  has_breakfast: boolean
  has_lunch: boolean
  has_dinner: boolean
  has_accommodation: boolean
  accommodation_text: string
}

type VoyageForm = {
  cruise_id: number
  code: string
  brief_info: string
  depart_date: string
  return_date: string
  status: number
  itineraries: ItineraryFormItem[]
}

const route = useRoute()
const { request } = useApi()
const id = Number(route.params.id)

const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const deleteDialogVisible = ref(false)
const error = ref<string | null>(null)
const empty = ref(false)
const cruises = ref<Array<{ id: number; name: string }>>([])
const hasCruiseOptions = computed(() => cruises.value.length > 0)
const form = ref<VoyageForm>({
  cruise_id: 1,
  code: '',
  brief_info: '',
  depart_date: '',
  return_date: '',
  status: 1,
  itineraries: [],
})

const dayNumbers = computed(() => {
  const days = new Set<number>()
  for (const item of form.value.itineraries) days.add(item.day_no)
  return Array.from(days).sort((a, b) => a - b)
})

function itineraryByDay(dayNo: number) {
  return form.value.itineraries
    .filter((item) => item.day_no === dayNo)
    .sort((a, b) => a.stop_index - b.stop_index)
}

function normalizeItineraries(items: ItineraryFormItem[]): ItineraryFormItem[] {
  const grouped = new Map<number, ItineraryFormItem[]>()
  for (const item of items) {
    if (!grouped.has(item.day_no)) grouped.set(item.day_no, [])
    grouped.get(item.day_no)!.push(item)
  }

  const sortedDays = Array.from(grouped.keys()).sort((a, b) => a - b)
  const normalized: ItineraryFormItem[] = []
  sortedDays.forEach((dayNo, dayIndex) => {
    const dayItems = grouped.get(dayNo)!.sort((a, b) => a.stop_index - b.stop_index)
    dayItems.forEach((item, stopIndex) => {
      normalized.push({
        ...item,
        day_no: dayIndex + 1,
        stop_index: stopIndex + 1,
        city: item.city.trim(),
        summary: item.summary.trim(),
        eta_time: item.eta_time.trim(),
        etd_time: item.etd_time.trim(),
        accommodation_text: item.accommodation_text.trim(),
      })
    })
  })
  return normalized
}

function dateOnly(input: string) {
  return input?.slice?.(0, 10) ?? ''
}

function toRFC3339Date(v: string) {
  if (!v) return v
  return `${v}T00:00:00Z`
}

function addDay() {
  const nextDay = Math.max(0, ...form.value.itineraries.map((item) => item.day_no)) + 1
  form.value.itineraries.push({
    day_no: nextDay,
    stop_index: 1,
    city: '',
    summary: '',
    eta_time: '',
    etd_time: '',
    has_breakfast: false,
    has_lunch: false,
    has_dinner: false,
    has_accommodation: false,
    accommodation_text: '',
  })
}

function removeDay(dayNo: number) {
  const remain = form.value.itineraries.filter((item) => item.day_no !== dayNo)
  form.value.itineraries = normalizeItineraries(remain)
}

function addStop(dayNo: number) {
  const maxStop = Math.max(
    0,
    ...form.value.itineraries.filter((item) => item.day_no === dayNo).map((item) => item.stop_index),
  )
  form.value.itineraries.push({
    day_no: dayNo,
    stop_index: maxStop + 1,
    city: '',
    summary: '',
    eta_time: '',
    etd_time: '',
    has_breakfast: false,
    has_lunch: false,
    has_dinner: false,
    has_accommodation: false,
    accommodation_text: '',
  })
}

function removeStop(dayNo: number, stopIndex: number) {
  const remain = form.value.itineraries.filter(
    (item) => !(item.day_no === dayNo && item.stop_index === stopIndex),
  )
  form.value.itineraries = normalizeItineraries(remain)
}

function ensureAtLeastOneItinerary() {
  if (form.value.itineraries.length > 0) return
  form.value.itineraries = [
    {
      day_no: 1,
      stop_index: 1,
      city: '',
      summary: '',
      eta_time: '',
      etd_time: '',
      has_breakfast: false,
      has_lunch: false,
      has_dinner: false,
      has_accommodation: false,
      accommodation_text: '',
    },
  ]
}

async function loadDetail() {
  loading.value = true
  error.value = null
  empty.value = false
  try {
    const res = await request(`/voyages/${id}`)
    const data = res?.data ?? res ?? {}
    if (Object.keys(data).length === 0) {
      empty.value = true
      return
    }

    const itineraries = Array.isArray(data.itineraries)
      ? data.itineraries.map((item: any) => ({
        day_no: Number(item.day_no || 1),
        stop_index: Number(item.stop_index || 1),
        city: String(item.city || ''),
        summary: String(item.summary || ''),
        eta_time: String(item.eta_time || ''),
        etd_time: String(item.etd_time || ''),
        has_breakfast: Boolean(item.has_breakfast),
        has_lunch: Boolean(item.has_lunch),
        has_dinner: Boolean(item.has_dinner),
        has_accommodation: Boolean(item.has_accommodation),
        accommodation_text: String(item.accommodation_text || ''),
      }))
      : []

    form.value = {
      cruise_id: Number(data.cruise_id || 1),
      code: data.code ?? '',
      brief_info: data.brief_info ?? '',
      depart_date: dateOnly(data.depart_date),
      return_date: dateOnly(data.return_date),
      status: Number(data.status ?? 1),
      itineraries: normalizeItineraries(itineraries),
    }
    ensureAtLeastOneItinerary()
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load voyage detail'
  } finally {
    loading.value = false
  }
}

async function loadCruises() {
  try {
    const res = await request('/cruises', { query: { page: 1, page_size: 100 } })
    const payload = res?.data ?? res ?? {}
    const rows = Array.isArray(payload) ? payload : payload?.list ?? []
    cruises.value = rows
      .filter((item: any) => Number(item?.id) > 0)
      .map((item: any) => ({ id: Number(item.id), name: String(item.name || '') }))
  } catch {
    cruises.value = []
  }
}

async function handleSave() {
  if (saving.value) return
  if (!hasCruiseOptions.value || !form.value.cruise_id) {
    error.value = '暂无可绑定邮轮，请先创建邮轮后再编辑航次。'
    return
  }
  saving.value = true
  error.value = null
  try {
    await request(`/voyages/${id}`, {
      method: 'PUT',
      body: {
        cruise_id: Number(form.value.cruise_id),
        code: form.value.code.trim(),
        brief_info: form.value.brief_info.trim(),
        depart_date: toRFC3339Date(form.value.depart_date),
        return_date: toRFC3339Date(form.value.return_date),
        status: Number(form.value.status) || 1,
        itineraries: normalizeItineraries(form.value.itineraries),
      },
    })
    await navigateTo('/voyages')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to update voyage'
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  if (deleting.value) return
  deleteDialogVisible.value = true
}

function closeDeleteDialog() {
  if (deleting.value) return
  deleteDialogVisible.value = false
}

async function confirmDelete() {
  if (deleting.value) return
  deleting.value = true
  error.value = null
  try {
    await request(`/voyages/${id}`, { method: 'DELETE' })
    closeDeleteDialog()
    await navigateTo('/voyages')
  } catch (e: any) {
    error.value = e?.message ?? '删除航次失败，请稍后重试。'
  } finally {
    deleting.value = false
  }
}

onMounted(async () => {
  await Promise.all([loadCruises(), loadDetail()])
})
</script>

<template>
  <div class="page">
    <h1>编辑航次 #{{ id }}</h1>
    <p v-if="loading">Loading...</p>
    <p v-else-if="empty" data-test="empty">暂无航次数据</p>
    <form v-else style="display:grid;gap:10px;max-width:980px;" @submit.prevent="handleSave">
      <input v-model="form.code" placeholder="航次代码" :disabled="saving || deleting" />
      <input v-model="form.brief_info" placeholder="航次简介（手动输入）" :disabled="saving || deleting" />
      <select v-model.number="form.cruise_id" data-test="cruise-select" :disabled="saving || deleting || cruises.length === 0">
        <option :value="0">请选择邮轮</option>
        <option v-for="cruise in cruises" :key="cruise.id" :value="cruise.id">{{ cruise.name || `邮轮 #${cruise.id}` }}</option>
      </select>
      <p v-if="!hasCruiseOptions" data-test="cruise-empty-hint" style="margin:0;color:#b45309;font-size:13px;">
        暂无可绑定邮轮，请先创建邮轮后再编辑航次。
      </p>
      <div style="display:flex;gap:8px;flex-wrap:wrap;">
        <label style="display:grid;gap:4px;">
          出发日期
          <input v-model="form.depart_date" type="date" :disabled="saving || deleting" />
        </label>
        <label style="display:grid;gap:4px;">
          结束日期
          <input v-model="form.return_date" type="date" :disabled="saving || deleting" />
        </label>
        <label style="display:grid;gap:4px;">
          状态
          <select v-model.number="form.status" :disabled="saving || deleting">
            <option :value="1">开放预订</option>
            <option :value="0">关闭</option>
          </select>
        </label>
      </div>

      <div style="display:flex;justify-content:space-between;align-items:center;">
        <h2 style="margin:0;">按天行程</h2>
        <button type="button" :disabled="saving || deleting" @click="addDay">新增一天</button>
      </div>

      <section
        v-for="dayNo in dayNumbers"
        :key="dayNo"
        style="border:1px solid #e5e7eb;border-radius:10px;padding:10px;display:grid;gap:8px;"
      >
        <div style="display:flex;justify-content:space-between;align-items:center;">
          <strong>第 {{ dayNo }} 天</strong>
          <div style="display:flex;gap:8px;">
            <button type="button" :disabled="saving || deleting" @click="addStop(dayNo)">新增站点</button>
            <button type="button" :disabled="saving || deleting || dayNumbers.length === 1" @click="removeDay(dayNo)">删除当天</button>
          </div>
        </div>

        <article
          v-for="item in itineraryByDay(dayNo)"
          :key="`${item.day_no}-${item.stop_index}`"
          style="border:1px solid #f3f4f6;border-radius:8px;padding:8px;display:grid;gap:8px;"
        >
          <div style="display:flex;justify-content:space-between;align-items:center;">
            <span>站点 {{ item.stop_index }}</span>
            <button
              type="button"
              :disabled="saving || deleting || itineraryByDay(dayNo).length === 1"
              @click="removeStop(dayNo, item.stop_index)"
            >
              删除站点
            </button>
          </div>
          <input v-model="item.city" placeholder="城市（必填）" :disabled="saving || deleting" />
          <textarea v-model="item.summary" placeholder="行程简介" :disabled="saving || deleting" />
          <div style="display:flex;gap:8px;flex-wrap:wrap;">
            <input v-model="item.eta_time" placeholder="靠港时间（可选，如 08:30）" :disabled="saving || deleting" />
            <input v-model="item.etd_time" placeholder="离港时间（可选，如 18:00）" :disabled="saving || deleting" />
          </div>
          <div style="display:flex;gap:10px;flex-wrap:wrap;">
            <label><input v-model="item.has_breakfast" type="checkbox" :disabled="saving || deleting" /> 早餐</label>
            <label><input v-model="item.has_lunch" type="checkbox" :disabled="saving || deleting" /> 午餐</label>
            <label><input v-model="item.has_dinner" type="checkbox" :disabled="saving || deleting" /> 晚餐</label>
            <label><input v-model="item.has_accommodation" type="checkbox" :disabled="saving || deleting" /> 住宿</label>
          </div>
          <input v-model="item.accommodation_text" placeholder="住宿说明（可选）" :disabled="saving || deleting" />
        </article>
      </section>

      <p v-if="error" class="error">{{ error }}</p>
      <div>
        <button type="submit" :disabled="saving || deleting || !hasCruiseOptions || !form.cruise_id">{{ saving ? '保存中...' : '保存' }}</button>
        <button type="button" style="margin-left:8px" :disabled="saving || deleting" @click="handleDelete">{{ deleting ? '删除中...' : '删除' }}</button>
      </div>
    </form>

    <AdminConfirmDialog
      :visible="deleteDialogVisible"
      title="确认删除航次"
      :message="`确认删除航次 #${id} 吗？删除后不可恢复。`"
      :loading="deleting"
      loading-text="删除中..."
      @close="closeDeleteDialog"
      @confirm="confirmDelete"
    />
  </div>
</template>
