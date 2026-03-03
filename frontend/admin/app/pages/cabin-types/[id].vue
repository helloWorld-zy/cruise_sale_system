<template>
  <div class="min-h-screen bg-slate-50 p-4 md:p-6">
    <div class="mx-auto max-w-5xl rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
      <h1 class="mb-6 text-xl font-semibold text-slate-900">编辑舱房类型</h1>
      <p v-if="loading" class="mb-3 text-sm text-slate-600">加载中...</p>
      <p v-else-if="empty" data-test="empty" class="mb-3 text-sm text-slate-600">暂无舱房类型数据</p>
      <form class="space-y-6" @submit.prevent="handleSubmit">
        <section>
          <h2 class="mb-4 border-b border-slate-200 pb-2 text-sm font-semibold text-slate-700">基本信息</h2>
          <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
            <label class="space-y-1 text-sm text-slate-600">
              <span>所属邮轮</span>
              <select v-model.number="form.cruise_id" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2">
                <option :value="0">请选择邮轮</option>
                <option v-for="cruise in cruises" :key="cruise.id" :value="Number(cruise.id)">{{ cruise.name || `邮轮 #${cruise.id}` }}</option>
              </select>
            </label>
            <label class="space-y-1 text-sm text-slate-600">
              <span>名称</span>
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
            <label class="space-y-1 text-sm text-slate-600 md:col-span-2">
              <span>描述</span>
              <textarea v-model="form.description" rows="4" class="w-full rounded-md border border-slate-200 px-3 py-2 outline-none ring-indigo-500 focus:ring-2" />
            </label>
          </div>
        </section>

        <section>
          <h2 class="mb-4 border-b border-slate-200 pb-2 text-sm font-semibold text-slate-700">面积与容量</h2>
          <div class="grid grid-cols-1 gap-4 md:grid-cols-4">
            <label class="space-y-1 text-sm text-slate-600">
              <span>面积最小值</span>
              <input v-model.number="form.area_min" type="number" min="0" step="0.1" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
            </label>
            <label class="space-y-1 text-sm text-slate-600">
              <span>面积最大值</span>
              <input v-model.number="form.area_max" type="number" min="0" step="0.1" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
            </label>
            <label class="space-y-1 text-sm text-slate-600">
              <span>基础容量</span>
              <input v-model.number="form.capacity" type="number" min="1" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
            </label>
            <label class="space-y-1 text-sm text-slate-600">
              <span>最大容量</span>
              <input v-model.number="form.max_capacity" type="number" min="1" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
            </label>
          </div>
        </section>

        <section>
          <h2 class="mb-4 border-b border-slate-200 pb-2 text-sm font-semibold text-slate-700">床型与设施</h2>
          <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
            <label class="space-y-1 text-sm text-slate-600">
              <span>床型</span>
              <input v-model="form.bed_type" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
            </label>
            <label class="space-y-1 text-sm text-slate-600">
              <span>楼层</span>
              <input v-model="form.deck" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
            </label>
          </div>

          <div class="mt-4">
            <p class="mb-2 text-sm text-slate-600">特色标签</p>
            <div class="grid grid-cols-1 gap-2 md:grid-cols-3">
              <label v-for="tag in tagOptions" :key="tag" class="flex items-center gap-2 rounded border border-slate-200 px-3 py-2 text-sm text-slate-700">
                <input type="checkbox" :checked="form.tags.includes(tag)" @change="toggleTags(tag, ($event.target as HTMLInputElement).checked)" />
                <span>{{ tag }}</span>
              </label>
            </div>
          </div>

          <div class="mt-4">
            <p class="mb-2 text-sm text-slate-600">设施清单</p>
            <div class="grid grid-cols-1 gap-2 md:grid-cols-2">
              <label v-for="amenity in amenityOptions" :key="amenity" class="flex items-center gap-2 rounded border border-slate-200 px-3 py-2 text-sm text-slate-700">
                <input type="checkbox" :checked="form.amenities.includes(amenity)" @change="toggleAmenities(amenity, ($event.target as HTMLInputElement).checked)" />
                <span>{{ amenity }}</span>
              </label>
            </div>
          </div>
        </section>

        <section>
          <h2 class="mb-4 border-b border-slate-200 pb-2 text-sm font-semibold text-slate-700">图片与介绍</h2>
          <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
            <label class="space-y-1 text-sm text-slate-600 md:col-span-2">
              <span>图片画廊（URL，多行）</span>
              <textarea v-model="form.gallery_text" rows="3" placeholder="每行一个图片 URL" class="w-full rounded-md border border-slate-200 px-3 py-2 outline-none ring-indigo-500 focus:ring-2" />
            </label>
            <div class="rounded-lg border-2 border-dashed border-slate-300 p-4 text-sm text-slate-500">图片上传区（占位）</div>
            <div class="rounded-lg border-2 border-dashed border-slate-300 p-4 text-sm text-slate-500">平面图上传区（占位）</div>
            <label class="space-y-1 text-sm text-slate-600 md:col-span-2">
              <span>平面图 URL</span>
              <input v-model="form.floor_plan_url" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
            </label>
          </div>
        </section>

        <div class="flex items-center justify-end gap-3 border-t border-slate-200 pt-4">
          <button type="button" :disabled="loading" class="rounded-md border border-rose-200 px-4 py-2 text-sm text-rose-600 hover:bg-rose-50" @click="handleDelete">删除</button>
          <NuxtLink to="/cabin-types" class="rounded-md border border-slate-200 px-4 py-2 text-sm text-slate-700 hover:bg-slate-50">取消</NuxtLink>
          <button type="submit" :disabled="loading" class="rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-500 disabled:cursor-not-allowed disabled:opacity-60">{{ loading ? '提交中...' : '保存' }}</button>
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
const tagOptions = ['亲子优选', '高性价比', '景观优先', '静音', '无障碍']
const amenityOptions = ['独立卫浴', '迷你吧', '阳台桌椅', '智能电视', '胶囊咖啡机', '浴缸']

const id = computed(() => {
  const value = Number(route.params.id)
  return Number.isFinite(value) && value > 0 ? value : 0
})

const form = ref({
  cruise_id: 0,
  name: '',
  english_name: '',
  code: '',
  area_min: 0,
  area_max: 0,
  capacity: 2,
  max_capacity: 2,
  bed_type: '',
  tags: [] as string[],
  amenities: [] as string[],
  floor_plan_url: '',
  area: 0,
  deck: '',
  description: '',
  sort_order: 0,
  gallery_text: '',
})

async function loadCruises() {
  try {
    const res = await request('/cruises', { query: { page: 1, page_size: 100 } })
    const payload = res?.data ?? res ?? {}
    cruises.value = Array.isArray(payload) ? payload : payload?.list ?? []
  } catch {
    cruises.value = []
  }
}

async function loadDetail() {
  if (!id.value) {
    error.value = '无效的舱房类型 ID'
    return
  }
  loading.value = true
  error.value = null
  empty.value = false
  try {
    const res = await request('/cabin-types', {
      query: {
        cruise_id: 1,
        page: 1,
        page_size: 200,
      },
    })
    const payload = res?.data ?? res ?? {}
    const list = Array.isArray(payload) ? payload : payload?.list ?? []
    const detail = list.find((item: Record<string, any>) => Number(item.id) === id.value)
    if (!detail) {
      empty.value = true
      return
    }
    form.value = {
      cruise_id: Number(detail.cruise_id || 0),
      name: detail.name || '',
      english_name: detail.english_name || '',
      code: detail.code || '',
      area_min: Number(detail.area_min || 0),
      area_max: Number(detail.area_max || 0),
      capacity: Number(detail.capacity || 2),
      max_capacity: Number(detail.max_capacity || 2),
      bed_type: detail.bed_type || '',
      tags: splitToArray(detail.tags),
      amenities: splitToArray(detail.amenities),
      floor_plan_url: detail.floor_plan_url || '',
      area: Number(detail.area || 0),
      deck: detail.deck || '',
      description: detail.description || '',
      sort_order: Number(detail.sort_order || 0),
      gallery_text: '',
    }
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load cabin type detail'
  } finally {
    loading.value = false
  }
}

function splitToArray(raw: unknown) {
  if (typeof raw !== 'string') return []
  return raw.split(',').map((part) => part.trim()).filter(Boolean)
}

function toggleTags(tag: string, checked: boolean) {
  const next = new Set(form.value.tags)
  if (checked) next.add(tag)
  else next.delete(tag)
  form.value.tags = Array.from(next)
}

function toggleAmenities(amenity: string, checked: boolean) {
  const next = new Set(form.value.amenities)
  if (checked) next.add(amenity)
  else next.delete(amenity)
  form.value.amenities = Array.from(next)
}

async function handleSubmit() {
  if (loading.value || !id.value) return
  loading.value = true
  error.value = null
  try {
    await request(`/cabin-types/${id.value}`, {
      method: 'PUT',
      body: {
        cruise_id: Number(form.value.cruise_id),
        name: form.value.name,
        english_name: form.value.english_name,
        code: form.value.code,
        area_min: Number(form.value.area_min),
        area_max: Number(form.value.area_max),
        capacity: Number(form.value.capacity),
        max_capacity: Number(form.value.max_capacity),
        bed_type: form.value.bed_type,
        tags: form.value.tags.join(','),
        amenities: form.value.amenities.join(','),
        floor_plan_url: form.value.floor_plan_url,
        area: Number(form.value.area),
        deck: form.value.deck,
        description: form.value.description,
        sort_order: Number(form.value.sort_order),
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
  if (loading.value || !id.value) return
  if (!confirm(`确认删除舱房类型 #${id.value} 吗？`)) return
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

onMounted(async () => {
  await loadCruises()
  await loadDetail()
})
</script>
