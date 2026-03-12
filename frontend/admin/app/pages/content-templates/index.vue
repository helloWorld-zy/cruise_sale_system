<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

type TemplateKind = 'fee_note' | 'booking_notice'
type TextItem = { text: string; emphasis?: boolean }
type FeeNoteContent = { included: TextItem[]; excluded: TextItem[] }
type BookingNoticeSection = { key: string; title: string; items: TextItem[] }
type BookingNoticeContent = { sections: BookingNoticeSection[] }
type TemplateItem = {
  id: number
  name: string
  kind: TemplateKind
  status: number
  content?: FeeNoteContent | BookingNoticeContent
}

const { request } = useApi()
const loading = ref(false)
const saving = ref(false)
const error = ref('')
const activeKind = ref<TemplateKind>('fee_note')
const items = ref<TemplateItem[]>([])
const editingId = ref<number | null>(null)
const form = ref({
  name: '',
  feeIncluded: [''],
  feeExcluded: [''],
  bookingSections: [{ key: 'booking_limit', title: '预订限制', items: [''] }],
})

const pageTitle = computed(() => activeKind.value === 'fee_note' ? '费用说明模板' : '预订须知模板')

function normalizeTextList(list?: Array<{ text?: string }> | string[]) {
  const values = (list || []).map((item) => typeof item === 'string' ? item : (item?.text || '')).map((item) => item.trim()).filter(Boolean)
  return values.length > 0 ? values : ['']
}

function normalizeBookingSections(content?: BookingNoticeContent) {
  const sections = (content?.sections || [])
    .map((section, index) => ({
      key: (section?.key || `section_${index + 1}`).trim(),
      title: (section?.title || '').trim(),
      items: normalizeTextList(section?.items),
    }))
    .filter((section) => section.title || section.items.some((item) => item))
  return sections.length > 0 ? sections : [{ key: 'booking_limit', title: '预订限制', items: [''] }]
}

function buildContent() {
  if (activeKind.value === 'fee_note') {
    return {
      included: form.value.feeIncluded.map((text) => ({ text: text.trim() })).filter((item) => item.text),
      excluded: form.value.feeExcluded.map((text) => ({ text: text.trim() })).filter((item) => item.text),
    }
  }
  return {
    sections: form.value.bookingSections
      .map((section, index) => ({
        key: section.key.trim() || `section_${index + 1}`,
        title: section.title.trim(),
        items: section.items.map((text) => ({ text: text.trim() })).filter((item) => item.text),
      }))
      .filter((section) => section.title || section.items.length > 0),
  }
}

function resetForm() {
  form.value = { name: '', feeIncluded: [''], feeExcluded: [''], bookingSections: [{ key: 'booking_limit', title: '预订限制', items: [''] }] }
  editingId.value = null
}

function fillForm(item: TemplateItem) {
  editingId.value = item.id
  activeKind.value = item.kind
  if (item.kind === 'fee_note') {
    const content = item.content as FeeNoteContent | undefined
    form.value = {
      name: item.name,
      feeIncluded: normalizeTextList(content?.included),
      feeExcluded: normalizeTextList(content?.excluded),
      bookingSections: [{ key: 'booking_limit', title: '预订限制', items: [''] }],
    }
    return
  }
  const content = item.content as BookingNoticeContent | undefined
  form.value = {
    name: item.name,
    feeIncluded: [''],
    feeExcluded: [''],
    bookingSections: normalizeBookingSections(content),
  }
}

function addFeeItem(kind: 'feeIncluded' | 'feeExcluded') {
  form.value[kind].push('')
}

function removeFeeItem(kind: 'feeIncluded' | 'feeExcluded', index: number) {
  if (form.value[kind].length === 1) {
    form.value[kind][0] = ''
    return
  }
  form.value[kind].splice(index, 1)
}

function addBookingSection() {
  form.value.bookingSections.push({ key: `section_${form.value.bookingSections.length + 1}`, title: '', items: [''] })
}

function removeBookingSection(index: number) {
  if (form.value.bookingSections.length === 1) {
    form.value.bookingSections[0] = { key: 'booking_limit', title: '预订限制', items: [''] }
    return
  }
  form.value.bookingSections.splice(index, 1)
}

function addBookingItem(sectionIndex: number) {
  form.value.bookingSections[sectionIndex]?.items.push('')
}

function removeBookingItem(sectionIndex: number, itemIndex: number) {
  const section = form.value.bookingSections[sectionIndex]
  if (!section) return
  if (section.items.length === 1) {
    section.items[0] = ''
    return
  }
  section.items.splice(itemIndex, 1)
}

async function loadItems() {
  loading.value = true
  error.value = ''
  try {
    const path = activeKind.value === 'fee_note' ? '/content-templates' : `/content-templates?kind=${activeKind.value}`
    const res = await request(path)
    const payload = res?.data ?? res ?? []
    items.value = Array.isArray(payload) ? payload : payload?.list ?? []
  } catch (e: any) {
    error.value = e?.message ?? '加载模板失败'
  } finally {
    loading.value = false
  }
}

async function handleSubmit() {
  saving.value = true
  error.value = ''
  try {
    await request(editingId.value ? `/content-templates/${editingId.value}` : '/content-templates', {
      method: editingId.value ? 'PUT' : 'POST',
      body: {
        name: form.value.name.trim(),
        kind: activeKind.value,
        status: 1,
        content: buildContent(),
      },
    })
    resetForm()
    await loadItems()
  } catch (e: any) {
    error.value = e?.message ?? '创建模板失败'
  } finally {
    saving.value = false
  }
}

async function switchKind(kind: TemplateKind) {
  activeKind.value = kind
  resetForm()
  await loadItems()
}

onMounted(loadItems)
</script>

<template>
  <div class="admin-page">
    <AdminPageHeader title="文案模板" subtitle="集中管理费用说明与预订须知模板" />

    <AdminDataCard flush>
      <div class="flex flex-wrap gap-2 border-b border-slate-100 px-4 py-3">
        <button data-test="kind-fee-note" type="button" class="rounded-full px-4 py-2 text-sm" :class="activeKind === 'fee_note' ? 'bg-slate-900 text-white' : 'bg-slate-100 text-slate-600'" @click="switchKind('fee_note')">费用说明模板</button>
        <button data-test="kind-booking-notice" type="button" class="rounded-full px-4 py-2 text-sm" :class="activeKind === 'booking_notice' ? 'bg-slate-900 text-white' : 'bg-slate-100 text-slate-600'" @click="switchKind('booking_notice')">预订须知模板</button>
      </div>
      <div class="overflow-x-auto px-4 py-4">
        <p v-if="loading" class="text-sm text-slate-500">加载中...</p>
        <p v-else-if="error" class="text-sm text-rose-500">{{ error }}</p>
        <p v-else-if="items.length === 0" class="text-sm text-slate-500">暂无模板</p>
        <table v-else class="w-full text-sm">
          <thead>
            <tr class="text-left text-slate-500">
              <th class="px-3 py-2">模板名称</th>
              <th class="px-3 py-2">类型</th>
              <th class="px-3 py-2">状态</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in items" :key="item.id" class="border-t border-slate-100">
              <td class="px-3 py-3 font-medium text-slate-900">{{ item.name }}</td>
              <td class="px-3 py-3 text-slate-600">{{ item.kind }}</td>
              <td class="px-3 py-3 text-slate-600">
                <div class="flex items-center gap-3">
                  <span>{{ item.status === 1 ? '启用' : '停用' }}</span>
                  <button :data-test="`edit-template-${item.id}`" type="button" class="text-indigo-600" @click="fillForm(item)">编辑</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </AdminDataCard>

    <AdminFormCard :title="editingId ? `编辑${pageTitle}` : `新增${pageTitle}`">
      <form class="grid gap-4" @submit.prevent="handleSubmit">
        <input v-model="form.name" data-test="template-name" type="text" placeholder="模板名称" class="h-10 rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
        <template v-if="activeKind === 'fee_note'">
          <div class="grid gap-4 md:grid-cols-2">
            <div class="rounded-xl border border-slate-200 p-3">
              <div class="flex items-center justify-between">
                <strong class="text-sm text-slate-900">费用包含</strong>
                <button data-test="add-fee-included" type="button" class="text-sm text-indigo-600" @click="addFeeItem('feeIncluded')">新增条目</button>
              </div>
              <div class="mt-3 grid gap-2">
                <div v-for="(item, index) in form.feeIncluded" :key="`fee-in-${index}`" class="flex gap-2">
                  <input v-model="form.feeIncluded[index]" :data-test="`fee-included-${index}`" type="text" placeholder="费用包含条目" class="h-10 flex-1 rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
                  <button type="button" class="rounded-md bg-slate-100 px-3 text-sm text-slate-600" @click="removeFeeItem('feeIncluded', index)">删除</button>
                </div>
              </div>
            </div>
            <div class="rounded-xl border border-slate-200 p-3">
              <div class="flex items-center justify-between">
                <strong class="text-sm text-slate-900">费用不包含</strong>
                <button data-test="add-fee-excluded" type="button" class="text-sm text-indigo-600" @click="addFeeItem('feeExcluded')">新增条目</button>
              </div>
              <div class="mt-3 grid gap-2">
                <div v-for="(item, index) in form.feeExcluded" :key="`fee-out-${index}`" class="flex gap-2">
                  <input v-model="form.feeExcluded[index]" :data-test="`fee-excluded-${index}`" type="text" placeholder="费用不包含条目" class="h-10 flex-1 rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
                  <button type="button" class="rounded-md bg-slate-100 px-3 text-sm text-slate-600" @click="removeFeeItem('feeExcluded', index)">删除</button>
                </div>
              </div>
            </div>
          </div>
        </template>
        <template v-else>
          <div class="grid gap-3">
            <div v-for="(section, sectionIndex) in form.bookingSections" :key="`section-${sectionIndex}`" class="rounded-xl border border-slate-200 p-3">
              <div class="flex items-center justify-between gap-3">
                <strong class="text-sm text-slate-900">分组 {{ sectionIndex + 1 }}</strong>
                <div class="flex gap-2">
                  <button :data-test="`add-booking-item-${sectionIndex}`" type="button" class="text-sm text-indigo-600" @click="addBookingItem(sectionIndex)">新增条目</button>
                  <button type="button" class="text-sm text-slate-500" @click="removeBookingSection(sectionIndex)">删除分组</button>
                </div>
              </div>
              <div class="mt-3 grid gap-2 md:grid-cols-2">
                <input v-model="section.title" :data-test="`booking-section-title-${sectionIndex}`" type="text" placeholder="分组标题" class="h-10 rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
                <input v-model="section.key" :data-test="`booking-section-key-${sectionIndex}`" type="text" placeholder="分组 key" class="h-10 rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
              </div>
              <div class="mt-3 grid gap-2">
                <div v-for="(item, itemIndex) in section.items" :key="`booking-${sectionIndex}-${itemIndex}`" class="flex gap-2">
                  <textarea v-model="section.items[itemIndex]" :data-test="`booking-item-${sectionIndex}-${itemIndex}`" rows="3" placeholder="预订须知正文" class="flex-1 rounded-md border border-slate-200 px-3 py-3 outline-none ring-indigo-500 focus:ring-2" />
                  <button type="button" class="rounded-md bg-slate-100 px-3 text-sm text-slate-600" @click="removeBookingItem(sectionIndex, itemIndex)">删除</button>
                </div>
              </div>
            </div>
            <button data-test="add-booking-section" type="button" class="rounded-md border border-dashed border-slate-300 px-4 py-2 text-sm text-slate-600" @click="addBookingSection">新增须知分组</button>
          </div>
        </template>
        <div>
          <button type="submit" class="rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white" :disabled="saving">{{ saving ? '提交中...' : '保存模板' }}</button>
          <button v-if="editingId" type="button" class="ml-2 rounded-md bg-slate-100 px-4 py-2 text-sm font-medium text-slate-600" @click="resetForm">取消编辑</button>
        </div>
      </form>
    </AdminFormCard>
  </div>
</template>