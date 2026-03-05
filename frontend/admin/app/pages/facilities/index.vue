<template>
  <div class="min-h-screen bg-slate-50 p-4 md:p-6">
    <div class="mx-auto max-w-7xl">
      <div class="mb-4 flex items-center justify-between">
        <h1 class="text-xl font-semibold text-slate-900">设施管理</h1>
        <AdminActionLink to="/facilities/new" variant="primary" size="md">新建设施</AdminActionLink>
      </div>

      <div class="mb-4 rounded-lg border border-slate-200 bg-white p-3 shadow-sm">
        <div class="flex flex-wrap items-center gap-3">
          <select v-model.number="filters.cruiseId" data-test="facility-cruise-filter" class="h-10 min-w-56 rounded-md border border-slate-200 px-3 text-sm text-slate-700 outline-none ring-indigo-500 focus:ring-2">
            <option :value="0">选择邮轮</option>
            <option v-for="cruise in cruises" :key="cruise.id" :value="Number(cruise.id)">{{ cruise.name || `邮轮 #${cruise.id}` }}</option>
          </select>
          <select v-model.number="filters.categoryId" class="h-10 min-w-56 rounded-md border border-slate-200 px-3 text-sm text-slate-700 outline-none ring-indigo-500 focus:ring-2">
            <option :value="0">全部分类</option>
            <option v-for="cat in categories" :key="cat.id" :value="Number(cat.id)">{{ cat.name || `分类 #${cat.id}` }}</option>
          </select>
          <button type="button" class="h-10 rounded-md border border-slate-200 px-3 text-sm text-slate-700 hover:bg-slate-50" @click="loadItems">筛选</button>
        </div>
      </div>

      <div class="rounded-lg border border-slate-200 bg-white shadow-sm">
        <table class="w-full text-sm">
          <thead class="bg-slate-50 text-left text-slate-600">
            <tr>
              <th class="p-3">名称</th>
              <th class="p-3">位置</th>
              <th class="p-3">开放时间</th>
              <th class="p-3">收费</th>
              <th class="p-3">状态</th>
              <th class="p-3">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading"><td class="p-3" colspan="6">加载中...</td></tr>
            <tr v-else-if="error"><td class="p-3 text-rose-500" colspan="6">{{ error }}</td></tr>
            <tr v-else-if="items.length === 0"><td class="p-3" colspan="6">暂无数据</td></tr>
            <tr v-for="(item, idx) in items" v-else :key="item.id" :class="idx % 2 === 1 ? 'bg-slate-50' : ''">
              <td class="p-3 font-medium text-slate-900">{{ item.name || '-' }}</td>
              <td class="p-3 text-slate-600">{{ item.location || '-' }}</td>
              <td class="p-3 text-slate-600">{{ item.open_hours || '-' }}</td>
              <td class="p-3">
                <span v-if="item.extra_charge" class="rounded-full bg-amber-50 px-2.5 py-0.5 text-xs font-medium text-amber-700">收费</span>
                <span v-else class="rounded-full bg-emerald-50 px-2.5 py-0.5 text-xs font-medium text-emerald-700">免费</span>
              </td>
              <td class="p-3">
                <span :class="Number(item.status ?? 1) === 1 ? 'rounded-full bg-emerald-50 px-2.5 py-0.5 text-xs font-medium text-emerald-700' : 'rounded-full bg-slate-100 px-2.5 py-0.5 text-xs font-medium text-slate-600'">
                  {{ Number(item.status ?? 1) === 1 ? '开放' : '关闭' }}
                </span>
              </td>
              <td class="p-3">
                <AdminActionLink :to="`/facilities/${item.id}`">编辑</AdminActionLink>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'

const { request } = useApi()
const loading = ref(false)
const error = ref<string | null>(null)
const items = ref<Record<string, any>[]>([])
const cruises = ref<Record<string, any>[]>([])
const categories = ref<Record<string, any>[]>([])

const filters = ref({
  cruiseId: 0,
  categoryId: 0,
})

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
    if (!filters.value.cruiseId && cruises.value.length > 0) {
      filters.value.cruiseId = Number(cruises.value[0].id) || 0
    }
  } catch {
    cruises.value = []
    categories.value = []
  }
}

async function loadItems() {
  if (!filters.value.cruiseId) {
    items.value = []
    return
  }
  loading.value = true
  error.value = null
  try {
    const query: Record<string, any> = { cruise_id: filters.value.cruiseId }
    if (filters.value.categoryId > 0) query.category_id = filters.value.categoryId
    const res = await request('/facilities', { query })
    const payload = res?.data ?? res ?? []
    items.value = Array.isArray(payload) ? payload : payload?.list ?? []
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load facilities'
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await loadOptions()
  await loadItems()
})
</script>
