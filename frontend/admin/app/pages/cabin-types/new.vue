<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'

const { request } = useApi()

const loading = ref(false)
const error = ref<string | null>(null)
const companies = ref<Record<string, any>[]>([])
const cruises = ref<Record<string, any>[]>([])
const categories = ref<Record<string, any>[]>([])
const selectedCruiseIds = ref<number[]>([])

const form = ref({
  company_id: 0,
  category_id: 0,
  name: '',
  english_name: '',
  code: '',
  area_min: 0,
  area_max: 0,
  occupancy: 2,
  capacity: 2,
  max_capacity: 2,
  intro: '',
  tags: '',
  amenities: '',
  sort_order: 0,
  status: 1,
})

const canSubmit = computed(() => {
  return selectedCruiseIds.value.length > 0 && form.value.category_id > 0 && form.value.name.trim().length > 0
})

async function loadCompanies() {
  try {
    const res = await request('/companies', { query: { page: 1, page_size: 100 } })
    const payload = res?.data ?? res ?? {}
    companies.value = Array.isArray(payload) ? payload : payload?.list ?? []
    if (companies.value.length > 0) {
      const first = Number(companies.value[0].id)
      if (Number.isFinite(first) && first > 0) form.value.company_id = first
    }
  } catch {
    companies.value = []
  }
}

async function loadCruises() {
  try {
    const query: Record<string, any> = { page: 1, page_size: 200 }
    if (form.value.company_id > 0) query.company_id = form.value.company_id
    const res = await request('/cruises', { query })
    const payload = res?.data ?? res ?? {}
    cruises.value = Array.isArray(payload) ? payload : payload?.list ?? []
    const allowed = new Set(cruises.value.map((item) => Number(item.id)))
    selectedCruiseIds.value = selectedCruiseIds.value.filter((id) => allowed.has(id))
  } catch {
    cruises.value = []
    selectedCruiseIds.value = []
  }
}

async function loadCategories() {
  try {
    const res = await request('/cabin-type-categories')
    categories.value = (res?.data ?? res ?? []).filter((item: Record<string, any>) => Number(item.status ?? 1) !== 0)
    if (categories.value.length > 0 && form.value.category_id <= 0) {
      form.value.category_id = Number(categories.value[0].id) || 0
    }
  } catch {
    categories.value = []
  }
}

function toggleCruise(cruiseID: number, checked: boolean) {
  const next = new Set(selectedCruiseIds.value)
  if (checked) next.add(cruiseID)
  else next.delete(cruiseID)
  selectedCruiseIds.value = Array.from(next)
}

async function handleSubmit() {
  if (!canSubmit.value || loading.value) return
  loading.value = true
  error.value = null
  try {
    await request('/cabin-types/batch-create', {
      method: 'POST',
      body: {
        cruise_ids: selectedCruiseIds.value,
        category_id: Number(form.value.category_id),
        name: form.value.name.trim(),
        english_name: form.value.english_name.trim(),
        code: form.value.code.trim(),
        area_min: Number(form.value.area_min),
        area_max: Number(form.value.area_max),
        occupancy: Number(form.value.occupancy),
        capacity: Number(form.value.capacity),
        max_capacity: Number(form.value.max_capacity),
        intro: form.value.intro.trim(),
        description: form.value.intro.trim(),
        tags: form.value.tags.trim(),
        amenities: form.value.amenities.trim(),
        sort_order: Number(form.value.sort_order),
        floor_plan_url: '',
        deck: '',
        bed_type: '',
        status: Number(form.value.status),
      },
    })
    await navigateTo('/cabin-types')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to create cabin type'
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await loadCompanies()
  await Promise.all([loadCruises(), loadCategories()])
})

watch(
  () => form.value.company_id,
  async () => {
    await loadCruises()
  },
)
</script>

<template>
  <div class="min-h-screen bg-slate-50 p-4 md:p-6">
    <div class="mx-auto max-w-5xl rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
      <div class="mb-4 flex items-center justify-between">
        <h1 class="text-xl font-semibold text-slate-900">新建舱型</h1>
        <AdminActionLink to="/cabin-types">返回列表</AdminActionLink>
      </div>

      <form class="space-y-6" @submit.prevent="handleSubmit">
        <section class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <label class="space-y-1 text-sm text-slate-600">
            <span>所属公司</span>
            <select v-model.number="form.company_id" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2">
              <option :value="0">请选择公司</option>
              <option v-for="company in companies" :key="company.id" :value="Number(company.id)">{{ company.name || `公司 #${company.id}` }}</option>
            </select>
          </label>
          <label class="space-y-1 text-sm text-slate-600">
            <span>舱型大类</span>
            <select v-model.number="form.category_id" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2">
              <option :value="0">请选择大类</option>
              <option v-for="category in categories" :key="category.id" :value="Number(category.id)">{{ category.name || `分类 #${category.id}` }}</option>
            </select>
          </label>
          <label class="space-y-1 text-sm text-slate-600">
            <span>舱型名称</span>
            <input v-model="form.name" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
          </label>
          <label class="space-y-1 text-sm text-slate-600">
            <span>代码</span>
            <input v-model="form.code" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
          </label>
          <label class="space-y-1 text-sm text-slate-600">
            <span>英文名</span>
            <input v-model="form.english_name" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
          </label>
          <label class="space-y-1 text-sm text-slate-600">
            <span>排序权重</span>
            <input v-model.number="form.sort_order" type="number" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
          </label>
        </section>

        <section class="grid grid-cols-1 gap-4 md:grid-cols-3">
          <label class="space-y-1 text-sm text-slate-600">
            <span>面积最小值</span>
            <input v-model.number="form.area_min" type="number" min="0" step="0.1" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
          </label>
          <label class="space-y-1 text-sm text-slate-600">
            <span>面积最大值</span>
            <input v-model.number="form.area_max" type="number" min="0" step="0.1" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
          </label>
          <label class="space-y-1 text-sm text-slate-600">
            <span>默认入住人数</span>
            <input v-model.number="form.occupancy" type="number" min="1" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
          </label>
          <label class="space-y-1 text-sm text-slate-600">
            <span>标准容量</span>
            <input v-model.number="form.capacity" type="number" min="1" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
          </label>
          <label class="space-y-1 text-sm text-slate-600">
            <span>最大容量</span>
            <input v-model.number="form.max_capacity" type="number" min="1" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
          </label>
        </section>

        <section>
          <label class="space-y-1 text-sm text-slate-600">
            <span>简介</span>
            <textarea v-model="form.intro" rows="4" class="w-full rounded-md border border-slate-200 px-3 py-2 outline-none ring-indigo-500 focus:ring-2" />
          </label>
          <div class="mt-4 grid grid-cols-1 gap-4 md:grid-cols-2">
            <label class="space-y-1 text-sm text-slate-600">
              <span>标签（逗号分隔）</span>
              <input v-model="form.tags" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
            </label>
            <label class="space-y-1 text-sm text-slate-600">
              <span>设施（逗号分隔）</span>
              <input v-model="form.amenities" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
            </label>
          </div>
        </section>

        <section>
          <p class="mb-2 text-sm font-medium text-slate-700">适用邮轮（可多选，按每个邮轮独立创建）</p>
          <div class="grid grid-cols-1 gap-2 rounded-md border border-slate-200 p-3 md:grid-cols-2">
            <label v-for="cruise in cruises" :key="cruise.id" class="flex items-center gap-2 text-sm text-slate-700">
              <input
                type="checkbox"
                :checked="selectedCruiseIds.includes(Number(cruise.id))"
                @change="toggleCruise(Number(cruise.id), ($event.target as HTMLInputElement).checked)"
              />
              <span>{{ cruise.name || `邮轮 #${cruise.id}` }}</span>
            </label>
            <p v-if="cruises.length === 0" class="text-sm text-slate-500">当前公司下暂无邮轮</p>
          </div>
        </section>

        <div class="flex items-center justify-end gap-3 border-t border-slate-200 pt-4">
          <AdminActionLink to="/cabin-types">取消</AdminActionLink>
          <button
            type="submit"
            :disabled="loading || !canSubmit"
            class="rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-500 disabled:cursor-not-allowed disabled:opacity-60"
          >
            {{ loading ? '提交中...' : '保存' }}
          </button>
        </div>
        <p v-if="error" class="text-sm text-rose-500">{{ error }}</p>
      </form>
    </div>
  </div>
</template>
