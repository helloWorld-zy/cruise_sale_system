<template>
  <div class="admin-page">
    <AdminPageHeader title="设施管理">
      <template #actions>
        <AdminActionLink to="/facilities/new" variant="primary" size="md">新建设施</AdminActionLink>
      </template>
    </AdminPageHeader>

    <AdminFilterBar>
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
    </AdminFilterBar>

    <AdminDataCard flush>
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
                <AdminStatusTag v-if="item.extra_charge" type="warning" text="收费" />
                <AdminStatusTag v-else type="success" text="免费" />
              </td>
              <td class="p-3">
                <AdminStatusTag :type="Number(item.status ?? 1) === 1 ? 'success' : 'info'" :text="Number(item.status ?? 1) === 1 ? '开放' : '关闭'" />
              </td>
              <td class="p-3">
                <AdminActionLink :to="`/facilities/${item.id}`">编辑</AdminActionLink>
              </td>
            </tr>
          </tbody>
        </table>
    </AdminDataCard>
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
