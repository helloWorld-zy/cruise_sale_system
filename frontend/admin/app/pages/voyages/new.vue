<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useApi } from '../../composables/useApi'

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

const form = ref<VoyageForm>({
  cruise_id: 1,
  code: '',
  brief_info: '',
  depart_date: '',
  return_date: '',
  status: 1,
  itineraries: [
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
  ],
})

const loading = ref(false)
const error = ref<string | null>(null)
const { request } = useApi()
const cruises = ref<Array<{ id: number; name: string }>>([])
const hasCruiseOptions = computed(() => cruises.value.length > 0)

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

function toRFC3339Date(v: string) {
  if (!v) return v
  return `${v}T00:00:00Z`
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

function buildSubmitBody() {
  return {
    cruise_id: Number(form.value.cruise_id),
    code: form.value.code.trim(),
    brief_info: form.value.brief_info.trim(),
    depart_date: toRFC3339Date(form.value.depart_date),
    return_date: toRFC3339Date(form.value.return_date),
    status: Number(form.value.status) || 1,
    itineraries: normalizeItineraries(form.value.itineraries),
  }
}

async function loadDefaultCruise() {
  try {
    const res = await request('/cruises', { query: { page: 1, page_size: 100 } })
    const payload = res?.data ?? res ?? {}
    const rows = Array.isArray(payload) ? payload : payload?.list ?? []
    cruises.value = rows
      .filter((item: any) => Number(item?.id) > 0)
      .map((item: any) => ({ id: Number(item.id), name: String(item.name || '') }))

    const first = cruises.value[0]
    form.value.cruise_id = Number(first?.id || 0)
    if (!form.value.cruise_id) {
      error.value = '暂无可用邮轮，请先创建邮轮后再新建航次。'
    }
  } catch {
    form.value.cruise_id = 0
    cruises.value = []
    error.value = '加载邮轮列表失败，请刷新后重试。'
  }
}

async function handleSubmit() {
  if (loading.value) return
  if (!form.value.cruise_id) {
    error.value = '暂无可用邮轮，请先创建邮轮后再新建航次。'
    return
  }
  loading.value = true
  error.value = null
  try {
    await request('/voyages', {
      method: 'POST',
      body: buildSubmitBody(),
    })
    await navigateTo('/voyages')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to create voyage'
  } finally {
    loading.value = false
  }
}

onMounted(loadDefaultCruise)
</script>

<template>
  <div class="page">
    <h1>新建航次</h1>
    <form style="display:grid;gap:10px;max-width:980px;" @submit.prevent="handleSubmit">
      <input v-model="form.code" placeholder="航次代码（如 TJ-20260701）" :disabled="loading" />
      <input v-model="form.brief_info" placeholder="航次简介（手动输入）" :disabled="loading" />
      <select v-model.number="form.cruise_id" data-test="cruise-select" :disabled="loading || cruises.length === 0">
        <option :value="0">请选择邮轮</option>
        <option v-for="cruise in cruises" :key="cruise.id" :value="cruise.id">{{ cruise.name || `邮轮 #${cruise.id}` }}</option>
      </select>
      <p v-if="!hasCruiseOptions" data-test="cruise-empty-hint" style="margin:0;color:#b45309;font-size:13px;">
        暂无可绑定邮轮，请先创建邮轮后再新建航次。
      </p>
      <div style="display:flex;gap:8px;flex-wrap:wrap;">
        <label style="display:grid;gap:4px;">
          出发日期
          <input v-model="form.depart_date" type="date" :disabled="loading" />
        </label>
        <label style="display:grid;gap:4px;">
          结束日期
          <input v-model="form.return_date" type="date" :disabled="loading" />
        </label>
        <label style="display:grid;gap:4px;">
          状态
          <select v-model.number="form.status" :disabled="loading">
            <option :value="1">开放预订</option>
            <option :value="0">关闭</option>
          </select>
        </label>
      </div>

      <div style="display:flex;justify-content:space-between;align-items:center;">
        <h2 style="margin:0;">按天行程</h2>
        <button type="button" :disabled="loading" @click="addDay">新增一天</button>
      </div>

      <section
        v-for="dayNo in dayNumbers"
        :key="dayNo"
        style="border:1px solid #e5e7eb;border-radius:10px;padding:10px;display:grid;gap:8px;"
      >
        <div style="display:flex;justify-content:space-between;align-items:center;">
          <strong>第 {{ dayNo }} 天</strong>
          <div style="display:flex;gap:8px;">
            <button type="button" :disabled="loading" @click="addStop(dayNo)">新增站点</button>
            <button type="button" :disabled="loading || dayNumbers.length === 1" @click="removeDay(dayNo)">删除当天</button>
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
              :disabled="loading || itineraryByDay(dayNo).length === 1"
              @click="removeStop(dayNo, item.stop_index)"
            >
              删除站点
            </button>
          </div>
          <input v-model="item.city" placeholder="城市（必填）" :disabled="loading" />
          <textarea v-model="item.summary" placeholder="行程简介" :disabled="loading" />
          <div style="display:flex;gap:8px;flex-wrap:wrap;">
            <input v-model="item.eta_time" placeholder="靠港时间（可选，如 08:30）" :disabled="loading" />
            <input v-model="item.etd_time" placeholder="离港时间（可选，如 18:00）" :disabled="loading" />
          </div>
          <div style="display:flex;gap:10px;flex-wrap:wrap;">
            <label><input v-model="item.has_breakfast" type="checkbox" :disabled="loading" /> 早餐</label>
            <label><input v-model="item.has_lunch" type="checkbox" :disabled="loading" /> 午餐</label>
            <label><input v-model="item.has_dinner" type="checkbox" :disabled="loading" /> 晚餐</label>
            <label><input v-model="item.has_accommodation" type="checkbox" :disabled="loading" /> 住宿</label>
          </div>
          <input v-model="item.accommodation_text" placeholder="住宿说明（可选）" :disabled="loading" />
        </article>
      </section>

      <p v-if="error" class="text-red-500">{{ error }}</p>
      <button type="submit" :disabled="loading || !hasCruiseOptions || !form.cruise_id">{{ loading ? '提交中...' : '提交' }}</button>
    </form>
  </div>
</template>

