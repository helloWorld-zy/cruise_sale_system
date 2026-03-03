<template>
  <div class="min-h-screen bg-slate-50 p-4 md:p-6">
    <div class="mx-auto max-w-5xl rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
      <h1 class="mb-6 text-xl font-semibold text-slate-900">编辑设施</h1>
      <p v-if="loading" class="mb-3 text-sm text-slate-600">加载中...</p>
      <p v-else-if="empty" data-test="empty" class="mb-3 text-sm text-slate-600">暂无设施数据</p>
      <form class="space-y-6" @submit.prevent="handleSubmit">
        <section>
          <h2 class="mb-4 border-b border-slate-200 pb-2 text-sm font-semibold text-slate-700">基本信息</h2>
          <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
            <label class="space-y-1 text-sm text-slate-600"><span>所属邮轮</span><select v-model.number="form.cruise_id" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2"><option :value="0">请选择邮轮</option><option v-for="cruise in cruises" :key="cruise.id" :value="Number(cruise.id)">{{ cruise.name || `邮轮 #${cruise.id}` }}</option></select></label>
            <label class="space-y-1 text-sm text-slate-600"><span>设施分类</span><select v-model.number="form.category_id" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2"><option :value="0">请选择分类</option><option v-for="cat in categories" :key="cat.id" :value="Number(cat.id)">{{ cat.name || `分类 #${cat.id}` }}</option></select></label>
            <label class="space-y-1 text-sm text-slate-600"><span>名称</span><input v-model="form.name" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" /></label>
            <label class="space-y-1 text-sm text-slate-600"><span>英文名</span><input v-model="form.english_name" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" /></label>
            <label class="space-y-1 text-sm text-slate-600"><span>位置</span><input v-model="form.location" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" /></label>
            <label class="space-y-1 text-sm text-slate-600"><span>开放时间</span><input v-model="form.open_hours" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" /></label>
          </div>
        </section>

        <section>
          <h2 class="mb-4 border-b border-slate-200 pb-2 text-sm font-semibold text-slate-700">收费与人群</h2>
          <div class="space-y-4">
            <label class="flex items-center gap-2 text-sm text-slate-700"><input v-model="form.extra_charge" type="checkbox" /><span>是否额外收费</span></label>
            <label v-if="form.extra_charge" class="block space-y-1 text-sm text-slate-600"><span>收费说明</span><input v-model="form.charge_price_tip" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" /></label>
            <div>
              <p class="mb-2 text-sm text-slate-600">适合人群</p>
              <div class="flex flex-wrap gap-3">
                <label v-for="aud in audienceOptions" :key="aud" class="flex items-center gap-2 text-sm text-slate-700"><input type="checkbox" :checked="form.target_audience.includes(aud)" @change="toggleAudience(aud, ($event.target as HTMLInputElement).checked)" /><span>{{ aud }}</span></label>
              </div>
            </div>
            <label class="block space-y-1 text-sm text-slate-600"><span>状态</span><select v-model.number="form.status" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2"><option :value="1">开放</option><option :value="0">关闭</option></select></label>
          </div>
        </section>

        <section>
          <h2 class="mb-4 border-b border-slate-200 pb-2 text-sm font-semibold text-slate-700">描述</h2>
          <label class="block space-y-1 text-sm text-slate-600"><span>设施描述</span><textarea v-model="form.description" rows="5" class="w-full rounded-md border border-slate-200 px-3 py-2 outline-none ring-indigo-500 focus:ring-2" /></label>
        </section>

        <div class="flex items-center justify-end gap-3 border-t border-slate-200 pt-4">
          <button type="button" class="rounded-md border border-rose-200 px-4 py-2 text-sm text-rose-600 hover:bg-rose-50" :disabled="loading" @click="handleDelete">删除</button>
          <NuxtLink to="/facilities" class="rounded-md border border-slate-200 px-4 py-2 text-sm text-slate-700 hover:bg-slate-50">取消</NuxtLink>
          <button type="submit" :disabled="loading" class="rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-500">{{ loading ? '提交中...' : '保存' }}</button>
        </div>
        <p v-if="error" class="text-sm text-rose-500">{{ error }}</p>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

const route = useRoute()
const { request } = useApi()
const loading = ref(false)
const error = ref<string | null>(null)
const empty = ref(false)
const cruises = ref<Record<string, any>[]>([])
const categories = ref<Record<string, any>[]>([])
const audienceOptions = ['儿童', '家庭', '情侣', '老年', '商务']

const id = computed(() => {
  const value = Number(route.params.id)
  return Number.isFinite(value) && value > 0 ? value : 0
})

const form = ref({
  category_id: 0,
  cruise_id: 0,
  name: '',
  english_name: '',
  location: '',
  open_hours: '',
  extra_charge: false,
  charge_price_tip: '',
  target_audience: [] as string[],
  description: '',
  status: 1,
  sort_order: 0,
})

function splitCsv(raw: unknown) {
  if (typeof raw !== 'string') return []
  return raw.split(',').map((part) => part.trim()).filter(Boolean)
}

async function loadOptions() {
  try {
    const [cruiseRes, categoryRes] = await Promise.all([
      request('/cruises', { query: { page: 1, page_size: 100 } }),
      request('/facility-categories'),
    ])
    const cruisePayload = cruiseRes?.data ?? cruiseRes ?? {}
    cruises.value = Array.isArray(cruisePayload) ? cruisePayload : cruisePayload?.list ?? []
    const categoryPayload = categoryRes?.data ?? categoryRes ?? []
    categories.value = Array.isArray(categoryPayload) ? categoryPayload : categoryPayload?.list ?? []
  } catch {
    cruises.value = []
    categories.value = []
  }
}

async function loadDetail() {
  if (!id.value) return
  loading.value = true
  error.value = null
  empty.value = false
  try {
    const res = await request(`/facilities/${id.value}`)
    const payload = res?.data ?? res ?? {}
    if (Object.keys(payload).length === 0) {
      empty.value = true
      return
    }
    form.value = {
      category_id: Number(payload.category_id || 0),
      cruise_id: Number(payload.cruise_id || 0),
      name: payload.name || '',
      english_name: payload.english_name || '',
      location: payload.location || '',
      open_hours: payload.open_hours || '',
      extra_charge: Boolean(payload.extra_charge),
      charge_price_tip: payload.charge_price_tip || '',
      target_audience: splitCsv(payload.target_audience),
      description: payload.description || '',
      status: Number(payload.status ?? 1),
      sort_order: Number(payload.sort_order || 0),
    }
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load facility detail'
  } finally {
    loading.value = false
  }
}

function toggleAudience(value: string, checked: boolean) {
  const next = new Set(form.value.target_audience)
  if (checked) next.add(value)
  else next.delete(value)
  form.value.target_audience = Array.from(next)
}

async function handleSubmit() {
  if (!id.value || loading.value) return
  loading.value = true
  error.value = null
  try {
    await request(`/facilities/${id.value}`, {
      method: 'PUT',
      body: {
        category_id: Number(form.value.category_id),
        cruise_id: Number(form.value.cruise_id),
        name: form.value.name,
        english_name: form.value.english_name,
        location: form.value.location,
        open_hours: form.value.open_hours,
        extra_charge: Boolean(form.value.extra_charge),
        charge_price_tip: form.value.extra_charge ? form.value.charge_price_tip : '',
        target_audience: form.value.target_audience.join(','),
        description: form.value.description,
        status: Number(form.value.status),
        sort_order: Number(form.value.sort_order || 0),
      },
    })
    await navigateTo('/facilities')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to update facility'
  } finally {
    loading.value = false
  }
}

async function handleDelete() {
  if (!id.value || loading.value) return
  if (!confirm(`确认删除设施 #${id.value} 吗？`)) return
  loading.value = true
  error.value = null
  try {
    await request(`/facilities/${id.value}`, { method: 'DELETE' })
    await navigateTo('/facilities')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to delete facility'
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await loadOptions()
  await loadDetail()
})
</script>
