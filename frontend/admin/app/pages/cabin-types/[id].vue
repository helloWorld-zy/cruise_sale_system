<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

const route = useRoute()
const { request } = useApi()

const id = computed(() => {
  const value = Number(route.params.id)
  return Number.isFinite(value) && value > 0 ? value : 0
})

const loading = ref(false)
const error = ref<string | null>(null)
const companies = ref<Record<string, any>[]>([])
const cruises = ref<Record<string, any>[]>([])
const categories = ref<Record<string, any>[]>([])

const mediaItems = ref<Record<string, any>[]>([])
const mediaUploading = ref(false)
const mediaType = ref<'image' | 'floor_plan'>('image')
const mediaTitle = ref('')
const mediaSortOrder = ref(0)
const mediaPrimary = ref(false)
const mediaFile = ref<File | null>(null)

const form = ref({
  company_id: 0,
  cruise_id: 0,
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

function normalizedMedia(type: 'image' | 'floor_plan') {
  return mediaItems.value.filter((item) => item.media_type === type)
}

async function loadCompanies() {
  try {
    const res = await request('/companies', { query: { page: 1, page_size: 100 } })
    const payload = res?.data ?? res ?? {}
    companies.value = Array.isArray(payload) ? payload : payload?.list ?? []
  } catch {
    companies.value = []
  }
}

async function loadCruises(companyID = 0) {
  try {
    const query: Record<string, any> = { page: 1, page_size: 200 }
    if (companyID > 0) query.company_id = companyID
    const res = await request('/cruises', { query })
    const payload = res?.data ?? res ?? {}
    cruises.value = Array.isArray(payload) ? payload : payload?.list ?? []
  } catch {
    cruises.value = []
  }
}

async function loadCategories() {
  try {
    const res = await request('/cabin-type-categories')
    categories.value = res?.data ?? res ?? []
  } catch {
    categories.value = []
  }
}

async function loadDetail() {
  if (!id.value) {
    error.value = '无效舱型 ID'
    return
  }

  loading.value = true
  error.value = null
  try {
    let detail: Record<string, any> | null = null

    const companyCandidates = companies.value.length > 0 ? companies.value : [{ id: 0 }]
    for (const company of companyCandidates) {
      const companyID = Number(company.id || 0)
      await loadCruises(companyID)
      for (const cruise of cruises.value) {
        const cruiseID = Number(cruise.id)
        if (!Number.isFinite(cruiseID) || cruiseID <= 0) continue
        const res = await request('/cabin-types', {
          query: { cruise_id: cruiseID, page: 1, page_size: 200 },
        })
        const payload = res?.data ?? res ?? {}
        const list = Array.isArray(payload) ? payload : payload?.list ?? []
        const found = list.find((item: Record<string, any>) => Number(item.id) === id.value)
        if (found) {
          detail = found
          break
        }
      }
      if (detail) break
    }

    if (!detail) {
      error.value = '未找到舱型详情'
      return
    }

    const companyID = Number(detail.company_id || 0)
    await loadCruises(companyID)

    form.value = {
      company_id: companyID,
      cruise_id: Number(detail.cruise_id || 0),
      category_id: Number(detail.category_id || 0),
      name: detail.name || '',
      english_name: detail.english_name || '',
      code: detail.code || '',
      area_min: Number(detail.area_min || 0),
      area_max: Number(detail.area_max || 0),
      occupancy: Number(detail.occupancy || detail.capacity || 2),
      capacity: Number(detail.capacity || 2),
      max_capacity: Number(detail.max_capacity || 2),
      intro: detail.intro || detail.description || '',
      tags: detail.tags || '',
      amenities: detail.amenities || '',
      sort_order: Number(detail.sort_order || 0),
      status: Number(detail.status ?? 1),
    }
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load cabin type'
  } finally {
    loading.value = false
  }
}

async function loadMedia() {
  if (!id.value) return
  try {
    const res = await request(`/cabin-types/${id.value}/media`)
    const payload = res?.data ?? res ?? []
    mediaItems.value = Array.isArray(payload) ? payload : payload?.list ?? []
  } catch {
    mediaItems.value = []
  }
}

async function handleSubmit() {
  if (!id.value || loading.value) return
  loading.value = true
  error.value = null
  try {
    await request(`/cabin-types/${id.value}`, {
      method: 'PUT',
      body: {
        cruise_id: Number(form.value.cruise_id),
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
    error.value = e?.message ?? 'failed to update cabin type'
  } finally {
    loading.value = false
  }
}

async function handleDelete() {
  if (!id.value || loading.value) return
  if (!confirm('确认删除该舱型吗？')) return
  loading.value = true
  error.value = null
  try {
    await request(`/cabin-types/${id.value}`, { method: 'DELETE' })
    await navigateTo('/cabin-types')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to delete cabin type'
  } finally {
    loading.value = false
  }
}

async function uploadMedia() {
  if (!id.value || mediaUploading.value || !mediaFile.value) return
  mediaUploading.value = true
  error.value = null
  try {
    const body = new FormData()
    body.append('file', mediaFile.value)
    body.append('media_type', mediaType.value)
    body.append('title', mediaTitle.value)
    body.append('sort_order', String(Number(mediaSortOrder.value || 0)))
    body.append('is_primary', mediaPrimary.value ? 'true' : 'false')

    await request(`/cabin-types/${id.value}/media/upload`, {
      method: 'POST',
      body,
    })

    mediaFile.value = null
    mediaTitle.value = ''
    mediaSortOrder.value = 0
    mediaPrimary.value = false
    await loadMedia()
  } catch (e: any) {
    error.value = e?.message ?? 'failed to upload media'
  } finally {
    mediaUploading.value = false
  }
}

async function saveMedia(item: Record<string, any>) {
  try {
    await request(`/cabin-types/${id.value}/media/${item.id}`, {
      method: 'PUT',
      body: {
        media_type: item.media_type,
        url: item.url,
        title: item.title || '',
        sort_order: Number(item.sort_order || 0),
        is_primary: Boolean(item.is_primary),
      },
    })
    await loadMedia()
  } catch (e: any) {
    error.value = e?.message ?? 'failed to save media metadata'
  }
}

async function removeMedia(item: Record<string, any>) {
  try {
    await request(`/cabin-types/${id.value}/media/${item.id}`, { method: 'DELETE' })
    await loadMedia()
  } catch (e: any) {
    error.value = e?.message ?? 'failed to delete media'
  }
}

onMounted(async () => {
  await Promise.all([loadCompanies(), loadCategories()])
  await loadDetail()
  await loadMedia()
})
</script>

<template>
  <div class="min-h-screen bg-slate-50 p-4 md:p-6">
    <div class="mx-auto max-w-6xl rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
      <div class="mb-4 flex items-center justify-between">
        <h1 class="text-xl font-semibold text-slate-900">编辑舱型</h1>
        <div class="flex items-center gap-2">
          <button type="button" class="rounded-md border border-rose-200 px-4 py-2 text-sm text-rose-600 hover:bg-rose-50" :disabled="loading" @click="handleDelete">删除</button>
          <AdminActionLink to="/cabin-types">返回列表</AdminActionLink>
        </div>
      </div>

      <form class="space-y-6" @submit.prevent="handleSubmit">
        <section class="grid grid-cols-1 gap-4 md:grid-cols-3">
          <label class="space-y-1 text-sm text-slate-600">
            <span>公司</span>
            <select v-model.number="form.company_id" disabled class="h-10 w-full rounded-md border border-slate-200 bg-slate-50 px-3">
              <option :value="0">未知公司</option>
              <option v-for="company in companies" :key="company.id" :value="Number(company.id)">{{ company.name || `公司 #${company.id}` }}</option>
            </select>
          </label>
          <label class="space-y-1 text-sm text-slate-600">
            <span>邮轮</span>
            <select v-model.number="form.cruise_id" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2">
              <option :value="0">请选择邮轮</option>
              <option v-for="cruise in cruises" :key="cruise.id" :value="Number(cruise.id)">{{ cruise.name || `邮轮 #${cruise.id}` }}</option>
            </select>
          </label>
          <label class="space-y-1 text-sm text-slate-600">
            <span>舱型大类</span>
            <select v-model.number="form.category_id" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2">
              <option :value="0">请选择大类</option>
              <option v-for="category in categories" :key="category.id" :value="Number(category.id)">{{ category.name || `分类 #${category.id}` }}</option>
            </select>
          </label>
        </section>

        <section class="grid grid-cols-1 gap-4 md:grid-cols-3">
          <label class="space-y-1 text-sm text-slate-600">
            <span>舱型名称</span>
            <input v-model="form.name" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
          </label>
          <label class="space-y-1 text-sm text-slate-600">
            <span>英文名</span>
            <input v-model="form.english_name" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
          </label>
          <label class="space-y-1 text-sm text-slate-600">
            <span>代码</span>
            <input v-model="form.code" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
          </label>
          <label class="space-y-1 text-sm text-slate-600">
            <span>面积最小值</span>
            <input v-model.number="form.area_min" type="number" min="0" step="0.1" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
          </label>
          <label class="space-y-1 text-sm text-slate-600">
            <span>面积最大值</span>
            <input v-model.number="form.area_max" type="number" min="0" step="0.1" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
          </label>
          <label class="space-y-1 text-sm text-slate-600">
            <span>排序权重</span>
            <input v-model.number="form.sort_order" type="number" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
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

        <section class="rounded-md border border-slate-200 p-4">
          <h2 class="mb-3 text-sm font-semibold text-slate-700">舱型媒体</h2>
          <div class="grid grid-cols-1 gap-3 md:grid-cols-5">
            <select v-model="mediaType" class="h-10 rounded-md border border-slate-200 px-3 text-sm">
              <option value="image">图片</option>
              <option value="floor_plan">平面图</option>
            </select>
            <input v-model="mediaTitle" placeholder="标题" class="h-10 rounded-md border border-slate-200 px-3 text-sm" />
            <input v-model.number="mediaSortOrder" type="number" placeholder="排序" class="h-10 rounded-md border border-slate-200 px-3 text-sm" />
            <label class="flex items-center gap-2 text-sm text-slate-600"><input v-model="mediaPrimary" type="checkbox" />设为主图</label>
            <input type="file" accept="image/*" class="text-sm" @change="mediaFile = ($event.target as HTMLInputElement).files?.[0] || null" />
          </div>
          <div class="mt-3">
            <button type="button" class="rounded-md border border-slate-200 px-3 py-2 text-sm text-slate-700 hover:bg-slate-50" :disabled="mediaUploading || !mediaFile" @click="uploadMedia">
              {{ mediaUploading ? '上传中...' : '上传媒体' }}
            </button>
          </div>

          <div class="mt-4 grid grid-cols-1 gap-4 md:grid-cols-2">
            <div>
              <h3 class="mb-2 text-sm font-medium text-slate-700">图片</h3>
              <div class="space-y-2">
                <div v-for="item in normalizedMedia('image')" :key="item.id" class="rounded border border-slate-200 p-2">
                  <a :href="item.url" target="_blank" class="text-xs text-indigo-600 hover:text-indigo-500">{{ item.url }}</a>
                  <div class="mt-2 grid grid-cols-1 gap-2 md:grid-cols-3">
                    <input v-model="item.title" class="h-9 rounded border border-slate-200 px-2 text-xs" />
                    <input v-model.number="item.sort_order" type="number" class="h-9 rounded border border-slate-200 px-2 text-xs" />
                    <label class="flex items-center gap-1 text-xs"><input v-model="item.is_primary" type="checkbox" />主图</label>
                  </div>
                  <div class="mt-2 flex items-center gap-2">
                    <button type="button" class="rounded border border-slate-200 px-2 py-1 text-xs" @click="saveMedia(item)">保存</button>
                    <button type="button" class="rounded border border-rose-200 px-2 py-1 text-xs text-rose-600" @click="removeMedia(item)">删除</button>
                  </div>
                </div>
                <p v-if="normalizedMedia('image').length === 0" class="text-xs text-slate-500">暂无图片</p>
              </div>
            </div>
            <div>
              <h3 class="mb-2 text-sm font-medium text-slate-700">平面图</h3>
              <div class="space-y-2">
                <div v-for="item in normalizedMedia('floor_plan')" :key="item.id" class="rounded border border-slate-200 p-2">
                  <a :href="item.url" target="_blank" class="text-xs text-indigo-600 hover:text-indigo-500">{{ item.url }}</a>
                  <div class="mt-2 grid grid-cols-1 gap-2 md:grid-cols-3">
                    <input v-model="item.title" class="h-9 rounded border border-slate-200 px-2 text-xs" />
                    <input v-model.number="item.sort_order" type="number" class="h-9 rounded border border-slate-200 px-2 text-xs" />
                    <label class="flex items-center gap-1 text-xs"><input v-model="item.is_primary" type="checkbox" />主图</label>
                  </div>
                  <div class="mt-2 flex items-center gap-2">
                    <button type="button" class="rounded border border-slate-200 px-2 py-1 text-xs" @click="saveMedia(item)">保存</button>
                    <button type="button" class="rounded border border-rose-200 px-2 py-1 text-xs text-rose-600" @click="removeMedia(item)">删除</button>
                  </div>
                </div>
                <p v-if="normalizedMedia('floor_plan').length === 0" class="text-xs text-slate-500">暂无平面图</p>
              </div>
            </div>
          </div>
        </section>

        <div class="flex items-center justify-end gap-3 border-t border-slate-200 pt-4">
          <AdminActionLink to="/cabin-types">取消</AdminActionLink>
          <button type="submit" :disabled="loading" class="rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-500 disabled:cursor-not-allowed disabled:opacity-60">
            {{ loading ? '提交中...' : '保存' }}
          </button>
        </div>
        <p v-if="error" class="text-sm text-rose-500">{{ error }}</p>
      </form>
    </div>
  </div>
</template>
