<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import AdminPortCitySelect from '../../components/AdminPortCitySelect.vue'
import AdminConfirmDialog from '../../components/AdminConfirmDialog.vue'
import { useAdminDeleteDialog } from '../../composables/useAdminDeleteDialog'

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
  image_url: string
  brief_info: string
  depart_date: string
  return_date: string
  status: number
  fee_note_template_id: number
  fee_note_mode: 'template' | 'snapshot' | ''
  fee_note_content: { included: Array<{ text: string }>; excluded: Array<{ text: string }> }
  booking_notice_template_id: number
  booking_notice_mode: 'template' | 'snapshot' | ''
  booking_notice_content: { sections: Array<{ key: string; title: string; items: Array<{ text: string }> }> }
  itineraries: ItineraryFormItem[]
}

type TemplateItem = {
  id: number
  name: string
  kind: 'fee_note' | 'booking_notice'
  status?: number
  content?: any
}

function blankFeeNoteContent() {
  return { included: [{ text: '' }], excluded: [{ text: '' }] }
}

function blankBookingNoticeContent() {
  return { sections: [{ key: 'booking_limit', title: '预订限制', items: [{ text: '' }] }] }
}

function normalizeTextItems(items?: Array<{ text?: string }>) {
  const normalized = (items || []).map((item) => ({ text: String(item?.text || '').trim() })).filter((item) => item.text)
  return normalized.length > 0 ? normalized : [{ text: '' }]
}

function normalizeFeeNoteContent(content?: { included?: Array<{ text?: string }>; excluded?: Array<{ text?: string }> }) {
  return {
    included: normalizeTextItems(content?.included),
    excluded: normalizeTextItems(content?.excluded),
  }
}

function normalizeBookingNoticeContent(content?: { sections?: Array<{ key?: string; title?: string; items?: Array<{ text?: string }> }> }) {
  const sections = (content?.sections || []).map((section, index) => ({
    key: String(section?.key || `section_${index + 1}`).trim(),
    title: String(section?.title || '').trim(),
    items: normalizeTextItems(section?.items),
  })).filter((section) => section.title || section.items.some((item) => item.text))

  return { sections: sections.length > 0 ? sections : blankBookingNoticeContent().sections }
}

const route = useRoute()
const { request } = useApi()
const id = Number(route.params.id)

const loading = ref(false)
const saving = ref(false)
const uploadingImage = ref(false)
const error = ref<string | null>(null)
const empty = ref(false)
const cruises = ref<Array<{ id: number; name: string }>>([])
const templates = ref<TemplateItem[]>([])
const hasCruiseOptions = computed(() => cruises.value.length > 0)
const feeNoteTemplates = computed(() => templates.value.filter((item) => item.kind === 'fee_note'))
const bookingNoticeTemplates = computed(() => templates.value.filter((item) => item.kind === 'booking_notice'))
const {
  visible: deleteDialogVisible,
  submitting: deleting,
  open: openDeleteDialog,
  close: closeDeleteDialog,
  run: runDelete,
} = useAdminDeleteDialog()
const form = ref<VoyageForm>({
  cruise_id: 0,
  code: '',
  image_url: '',
  brief_info: '',
  depart_date: '',
  return_date: '',
  status: 1,
  fee_note_template_id: 0,
  fee_note_mode: '',
  fee_note_content: blankFeeNoteContent(),
  booking_notice_template_id: 0,
  booking_notice_mode: '',
  booking_notice_content: blankBookingNoticeContent(),
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

function cloneTemplateContent<T>(value: T, fallback: T): T {
  return JSON.parse(JSON.stringify(value ?? fallback))
}

function findTemplate(id: number) {
  return templates.value.find((item) => item.id === id)
}

function copyBookingNoticeToSnapshot() {
  const template = findTemplate(form.value.booking_notice_template_id)
  form.value.booking_notice_mode = 'snapshot'
  form.value.booking_notice_content = normalizeBookingNoticeContent(cloneTemplateContent(template?.content, blankBookingNoticeContent()))
}

function copyFeeNoteToSnapshot() {
  const template = findTemplate(form.value.fee_note_template_id)
  form.value.fee_note_mode = 'snapshot'
  form.value.fee_note_content = normalizeFeeNoteContent(cloneTemplateContent(template?.content, blankFeeNoteContent()))
}

function onFeeNoteTemplateChange() {
  form.value.fee_note_mode = form.value.fee_note_template_id ? 'template' : ''
}

function onBookingNoticeTemplateChange() {
  form.value.booking_notice_mode = form.value.booking_notice_template_id ? 'template' : ''
}

function addFeeNoteLine(kind: 'included' | 'excluded') {
  form.value.fee_note_content[kind].push({ text: '' })
}

function removeFeeNoteLine(kind: 'included' | 'excluded', index: number) {
  const list = form.value.fee_note_content[kind]
  if (list.length === 1) {
    list[0] = { text: '' }
    return
  }
  list.splice(index, 1)
}

function addBookingNoticeSection() {
  form.value.booking_notice_content.sections.push({ key: `section_${form.value.booking_notice_content.sections.length + 1}`, title: '', items: [{ text: '' }] })
}

function removeBookingNoticeSection(index: number) {
  if (form.value.booking_notice_content.sections.length === 1) {
    form.value.booking_notice_content.sections = blankBookingNoticeContent().sections
    return
  }
  form.value.booking_notice_content.sections.splice(index, 1)
}

function addBookingNoticeItem(sectionIndex: number) {
  form.value.booking_notice_content.sections[sectionIndex]?.items.push({ text: '' })
}

function removeBookingNoticeItem(sectionIndex: number, itemIndex: number) {
  const section = form.value.booking_notice_content.sections[sectionIndex]
  if (!section) return
  if (section.items.length === 1) {
    section.items[0] = { text: '' }
    return
  }
  section.items.splice(itemIndex, 1)
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
      cruise_id: Number(data.cruise_id || 0),
      code: data.code ?? '',
      image_url: String(data.image_url || ''),
      brief_info: data.brief_info ?? '',
      depart_date: dateOnly(data.depart_date),
      return_date: dateOnly(data.return_date),
      status: Number(data.status ?? 1),
      fee_note_template_id: Number(data.fee_note_template_id || 0),
      fee_note_mode: data.fee_note_mode || '',
      fee_note_content: normalizeFeeNoteContent(cloneTemplateContent(data.fee_note_content, blankFeeNoteContent())),
      booking_notice_template_id: Number(data.booking_notice_template_id || 0),
      booking_notice_mode: data.booking_notice_mode || '',
      booking_notice_content: normalizeBookingNoticeContent(cloneTemplateContent(data.booking_notice_content, blankBookingNoticeContent())),
      itineraries: normalizeItineraries(itineraries),
    }
    ensureAtLeastOneItinerary()
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load voyage detail'
  } finally {
    loading.value = false
  }
}

async function uploadVoyageImage(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  uploadingImage.value = true
  error.value = null
  try {
    const body = new FormData()
    body.append('file', file)
    const res = await request('/upload/image', { method: 'POST', body })
    const payload = res?.data ?? res ?? {}
    const url = String(payload?.url || '')
    if (!url) throw new Error('上传成功但未返回图片地址')
    form.value.image_url = url
  } catch (e: any) {
    error.value = e?.message ?? '上传航次图片失败'
  } finally {
    uploadingImage.value = false
    input.value = ''
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

async function loadContentTemplates() {
  try {
    const res = await request('/content-templates')
    const payload = res?.data ?? res ?? []
    templates.value = Array.isArray(payload) ? payload : payload?.list ?? []
  } catch {
    templates.value = []
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
        image_url: form.value.image_url.trim(),
        brief_info: form.value.brief_info.trim(),
        depart_date: toRFC3339Date(form.value.depart_date),
        return_date: toRFC3339Date(form.value.return_date),
        status: Number(form.value.status) || 1,
        fee_note_template_id: Number(form.value.fee_note_template_id || 0),
        fee_note_mode: form.value.fee_note_mode || (form.value.fee_note_template_id ? 'template' : ''),
        fee_note_content: form.value.fee_note_mode === 'snapshot' ? form.value.fee_note_content : undefined,
        booking_notice_template_id: Number(form.value.booking_notice_template_id || 0),
        booking_notice_mode: form.value.booking_notice_mode || (form.value.booking_notice_template_id ? 'template' : ''),
        booking_notice_content: form.value.booking_notice_mode === 'snapshot' ? form.value.booking_notice_content : undefined,
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
  openDeleteDialog()
}

async function confirmDelete() {
  if (deleting.value) return
  error.value = null
  try {
    await runDelete(async () => {
      await request(`/voyages/${id}`, { method: 'DELETE' })
      await navigateTo('/voyages')
    })
  } catch (e: any) {
    error.value = e?.message ?? '删除航次失败，请稍后重试。'
  }
}

onMounted(async () => {
  await Promise.all([loadCruises(), loadContentTemplates(), loadDetail()])
})
</script>

<template>
  <div class="admin-page">
    <AdminPageHeader :title="`编辑航次 #${id}`" />
    <AdminFormCard>
      <p v-if="loading" class="text-sm text-slate-600">Loading...</p>
      <p v-else-if="empty" data-test="empty" class="text-sm text-slate-600">暂无航次数据</p>
      <form v-else style="display:grid;gap:10px;max-width:980px;" @submit.prevent="handleSave">
      <input v-model="form.code" placeholder="航次代码" :disabled="saving || deleting" />
      <div style="display:grid;gap:6px;">
        <label style="font-size:13px;color:#475569;">航次图片</label>
        <div style="display:flex;align-items:center;gap:10px;flex-wrap:wrap;">
          <input type="file" accept="image/*" :disabled="saving || deleting || uploadingImage" @change="uploadVoyageImage" />
          <input v-model="form.image_url" placeholder="或直接填写图片 URL" :disabled="saving || deleting || uploadingImage" />
          <span v-if="uploadingImage" style="font-size:12px;color:#64748b;">上传中...</span>
        </div>
        <img v-if="form.image_url" :src="form.image_url" alt="航次图片预览" style="width:96px;height:96px;object-fit:cover;border-radius:8px;border:1px solid #e2e8f0;" />
      </div>
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

      <section style="border:1px solid #e5e7eb;border-radius:12px;padding:12px;display:grid;gap:10px;">
        <strong>费用说明</strong>
        <select v-model.number="form.fee_note_template_id" data-test="fee-note-template-select" :disabled="saving || deleting" @change="onFeeNoteTemplateChange">
          <option :value="0">请选择费用说明模板</option>
          <option v-for="item in feeNoteTemplates" :key="item.id" :value="item.id">{{ item.name }}</option>
        </select>
        <div style="display:flex;align-items:center;gap:10px;flex-wrap:wrap;">
          <span style="font-size:13px;color:#64748b;">当前模式：{{ form.fee_note_mode || '未选择模板' }}</span>
          <button type="button" data-test="copy-fee-note-to-snapshot" :disabled="saving || deleting || !form.fee_note_template_id" @click="copyFeeNoteToSnapshot">复制为航次专属内容</button>
        </div>
        <div v-if="form.fee_note_mode === 'snapshot'" style="display:grid;gap:12px;">
          <div style="border:1px solid #e5e7eb;border-radius:10px;padding:10px;display:grid;gap:8px;">
            <div style="display:flex;justify-content:space-between;align-items:center;">
              <strong>费用包含</strong>
              <button type="button" :disabled="saving || deleting" @click="addFeeNoteLine('included')">新增条目</button>
            </div>
            <div v-for="(item, index) in form.fee_note_content.included" :key="`fee-included-${index}`" style="display:flex;gap:8px;align-items:flex-start;">
              <input v-model="item.text" :data-test="`fee-note-included-${index}`" placeholder="费用包含条目" :disabled="saving || deleting" style="flex:1;" />
              <button type="button" :disabled="saving || deleting" @click="removeFeeNoteLine('included', index)">删除</button>
            </div>
          </div>
          <div style="border:1px solid #e5e7eb;border-radius:10px;padding:10px;display:grid;gap:8px;">
            <div style="display:flex;justify-content:space-between;align-items:center;">
              <strong>费用不包含</strong>
              <button type="button" :disabled="saving || deleting" @click="addFeeNoteLine('excluded')">新增条目</button>
            </div>
            <div v-for="(item, index) in form.fee_note_content.excluded" :key="`fee-excluded-${index}`" style="display:flex;gap:8px;align-items:flex-start;">
              <input v-model="item.text" :data-test="`fee-note-excluded-${index}`" placeholder="费用不包含条目" :disabled="saving || deleting" style="flex:1;" />
              <button type="button" :disabled="saving || deleting" @click="removeFeeNoteLine('excluded', index)">删除</button>
            </div>
          </div>
        </div>
      </section>

      <section style="border:1px solid #e5e7eb;border-radius:12px;padding:12px;display:grid;gap:10px;">
        <strong>预订须知</strong>
        <select v-model.number="form.booking_notice_template_id" data-test="booking-notice-template-select" :disabled="saving || deleting" @change="onBookingNoticeTemplateChange">
          <option :value="0">请选择预订须知模板</option>
          <option v-for="item in bookingNoticeTemplates" :key="item.id" :value="item.id">{{ item.name }}</option>
        </select>
        <div style="display:flex;align-items:center;gap:10px;flex-wrap:wrap;">
          <span style="font-size:13px;color:#64748b;">当前模式：{{ form.booking_notice_mode || '未选择模板' }}</span>
          <button type="button" data-test="copy-booking-notice-to-snapshot" :disabled="saving || deleting || !form.booking_notice_template_id" @click="copyBookingNoticeToSnapshot">复制为航次专属内容</button>
        </div>
        <div v-if="form.booking_notice_mode === 'snapshot'" style="display:grid;gap:10px;">
          <section v-for="(section, sectionIndex) in form.booking_notice_content.sections" :key="`booking-section-${sectionIndex}`" style="border:1px solid #e5e7eb;border-radius:10px;padding:10px;display:grid;gap:8px;">
            <div style="display:flex;justify-content:space-between;align-items:center;gap:8px;">
              <strong>须知分组 {{ sectionIndex + 1 }}</strong>
              <div style="display:flex;gap:8px;">
                <button type="button" :disabled="saving || deleting" @click="addBookingNoticeItem(sectionIndex)">新增条目</button>
                <button type="button" :disabled="saving || deleting" @click="removeBookingNoticeSection(sectionIndex)">删除分组</button>
              </div>
            </div>
            <div style="display:flex;gap:8px;flex-wrap:wrap;">
              <input v-model="section.title" :data-test="`booking-notice-section-title-${sectionIndex}`" placeholder="分组标题" :disabled="saving || deleting" style="flex:1;min-width:180px;" />
              <input v-model="section.key" :data-test="`booking-notice-section-key-${sectionIndex}`" placeholder="分组 key" :disabled="saving || deleting" style="flex:1;min-width:180px;" />
            </div>
            <div v-for="(item, itemIndex) in section.items" :key="`booking-item-${sectionIndex}-${itemIndex}`" style="display:flex;gap:8px;align-items:flex-start;">
              <textarea v-model="item.text" :data-test="sectionIndex === 0 && itemIndex === 0 ? 'booking-notice-text-0' : `booking-notice-text-${sectionIndex}-${itemIndex}`" rows="4" placeholder="预订须知正文" :disabled="saving || deleting" style="flex:1;" />
              <button type="button" :disabled="saving || deleting" @click="removeBookingNoticeItem(sectionIndex, itemIndex)">删除</button>
            </div>
          </section>
          <button type="button" :disabled="saving || deleting" @click="addBookingNoticeSection">新增须知分组</button>
        </div>
      </section>

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
          <AdminPortCitySelect
            v-model="item.city"
            :disabled="saving || deleting"
            :test-id-base="`itinerary-city-${dayNo}-${item.stop_index}`"
          />
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

      <p v-if="error" class="text-sm text-rose-500">{{ error }}</p>
      <div>
        <button type="submit" :disabled="saving || deleting || !hasCruiseOptions || !form.cruise_id">{{ saving ? '保存中...' : '保存' }}</button>
        <button type="button" data-test="delete-voyage" style="margin-left:8px" :disabled="saving || deleting" @click="handleDelete">{{ deleting ? '删除中...' : '删除' }}</button>
      </div>
      </form>
    </AdminFormCard>

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
