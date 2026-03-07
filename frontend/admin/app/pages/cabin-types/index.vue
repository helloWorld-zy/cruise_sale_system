<template>
  <div class="admin-page">
    <AdminPageHeader title="舱型管理">
      <template #actions>
        <div class="flex items-center gap-2">
          <AdminActionLink to="/cabin-types/pricing" size="md">价格管理</AdminActionLink>
          <AdminActionLink to="/cabin-types/new" variant="primary" size="md">新建舱型</AdminActionLink>
        </div>
      </template>
    </AdminPageHeader>

    <AdminFilterBar>
      <div class="flex flex-wrap items-center gap-3">
          <select
            v-model.number="filters.companyId"
            class="h-10 min-w-56 rounded-md border border-slate-200 px-3 text-sm text-slate-700 outline-none ring-indigo-500 focus:ring-2"
          >
            <option :value="0">全部公司</option>
            <option v-for="company in companies" :key="company.id" :value="Number(company.id)">{{ company.name || `公司 #${company.id}` }}</option>
          </select>
          <select
            v-model.number="filters.cruiseId"
            data-test="cruise-filter"
            class="h-10 min-w-64 rounded-md border border-slate-200 px-3 text-sm text-slate-700 outline-none ring-indigo-500 focus:ring-2"
          >
            <option :value="0">全部邮轮</option>
            <option v-for="cruise in cruises" :key="cruise.id" :value="Number(cruise.id)">{{ cruise.name || `邮轮 #${cruise.id}` }}</option>
          </select>
          <select
            v-model.number="filters.categoryId"
            class="h-10 min-w-56 rounded-md border border-slate-200 px-3 text-sm text-slate-700 outline-none ring-indigo-500 focus:ring-2"
          >
            <option :value="0">全部舱型大类</option>
            <option v-for="category in categories" :key="category.id" :value="Number(category.id)">{{ category.name || `分类 #${category.id}` }}</option>
          </select>
          <button type="button" class="h-10 rounded-md border border-slate-200 px-3 text-sm text-slate-700 hover:bg-slate-50" @click="loadItems">筛选</button>
      </div>
    </AdminFilterBar>

    <AdminDataCard flush>
      <div class="cabin-type-table-wrap overflow-x-auto">
          <table class="w-full min-w-[1080px] text-sm">
          <thead class="bg-slate-50 text-left text-slate-600">
            <tr>
              <th class="p-3">邮轮</th>
              <th class="p-3">分类</th>
              <th class="p-3">名称</th>
              <th class="p-3">代码</th>
              <th class="p-3">面积范围</th>
              <th class="p-3">容量</th>
              <th class="p-3">简介</th>
              <th class="p-3">状态</th>
              <th class="p-3 whitespace-nowrap">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td class="p-3" colspan="9">加载中...</td>
            </tr>
            <tr v-else-if="error">
              <td class="p-3 text-rose-500" colspan="9">{{ error }}</td>
            </tr>
            <tr v-else-if="items.length === 0">
              <td class="p-3" colspan="9">暂无数据</td>
            </tr>
            <tr v-for="(item, idx) in items" v-else :key="item.id" :class="idx % 2 === 1 ? 'bg-slate-50' : ''">
              <td class="p-3 text-slate-600">{{ cruiseName(item.cruise_id) }}</td>
              <td class="p-3 text-slate-600">{{ categoryName(item.category_id) }}</td>
              <td class="p-3 font-medium text-slate-900">{{ item.name || '-' }}</td>
              <td class="p-3 text-slate-600">{{ item.code || '-' }}</td>
              <td class="p-3 text-slate-600">{{ areaText(item) }}</td>
              <td class="p-3 text-slate-600">{{ item.occupancy || item.max_capacity || item.capacity || '-' }}</td>
              <td class="p-3 text-slate-600">{{ item.intro || item.description || '-' }}</td>
              <td class="p-3">
                <AdminStatusTag :type="statusType(item.status)" :text="statusText(item.status)" />
              </td>
              <td class="p-3 whitespace-nowrap">
                <div class="cabin-type-actions flex items-center gap-2">
                  <AdminActionLink :to="`/cabin-types/${item.id}`">编辑</AdminActionLink>
                </div>
              </td>
            </tr>
          </tbody>
          </table>
      </div>
    </AdminDataCard>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'

const { request } = useApi()
const companies = ref<Record<string, any>[]>([])
const cruises = ref<Record<string, any>[]>([])
const categories = ref<Record<string, any>[]>([])
const items = ref<Record<string, any>[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

const filters = ref({
  companyId: 0,
  cruiseId: 0,
  categoryId: 0,
})

const cruiseMap = computed(() => {
  const map = new Map<number, string>()
  for (const cruise of cruises.value) {
    const id = Number(cruise.id)
    if (Number.isFinite(id) && id > 0) map.set(id, cruise.name || `邮轮 #${id}`)
  }
  return map
})

const categoryMap = computed(() => {
  const map = new Map<number, string>()
  for (const category of categories.value) {
    const id = Number(category.id)
    if (Number.isFinite(id) && id > 0) map.set(id, category.name || `分类 #${id}`)
  }
  return map
})

async function loadCompanies() {
  try {
    const res = await request('/companies', { query: { page: 1, page_size: 100 } })
    const payload = res?.data ?? res ?? {}
    companies.value = Array.isArray(payload) ? payload : payload?.list ?? []
  } catch {
    companies.value = []
  }
}

async function loadCruises() {
  try {
    const query: Record<string, any> = { page: 1, page_size: 100 }
    if (filters.value.companyId > 0) query.company_id = filters.value.companyId
    const res = await request('/cruises', { query })
    const payload = res?.data ?? res ?? {}
    cruises.value = Array.isArray(payload) ? payload : payload?.list ?? []
    if (filters.value.cruiseId > 0 && !cruises.value.find((item) => Number(item.id) === filters.value.cruiseId)) {
      filters.value.cruiseId = 0
    }
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

async function loadItems() {
  loading.value = true
  error.value = null
  try {
    const all: Record<string, any>[] = []

    if (filters.value.cruiseId > 0) {
      const res = await request('/cabin-types', {
        query: {
          cruise_id: filters.value.cruiseId,
          page: 1,
          page_size: 50,
        },
      })
      const payload = res?.data ?? res ?? {}
      all.push(...(Array.isArray(payload) ? payload : payload?.list ?? []))
    } else {
      for (const cruise of cruises.value) {
        const cruiseID = Number(cruise.id)
        if (!Number.isFinite(cruiseID) || cruiseID <= 0) continue
        const res = await request('/cabin-types', {
          query: {
            cruise_id: cruiseID,
            page: 1,
            page_size: 100,
          },
        })
        const payload = res?.data ?? res ?? {}
        all.push(...(Array.isArray(payload) ? payload : payload?.list ?? []))
      }
    }

    items.value = all.filter((item) => {
      if (filters.value.categoryId > 0 && Number(item.category_id) !== filters.value.categoryId) return false
      return true
    })
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load cabin types'
  } finally {
    loading.value = false
  }
}

function cruiseName(raw: unknown) {
  const id = Number(raw)
  return cruiseMap.value.get(id) || (id > 0 ? `邮轮 #${id}` : '-')
}

function categoryName(raw: unknown) {
  const id = Number(raw)
  return categoryMap.value.get(id) || (id > 0 ? `分类 #${id}` : '-')
}

function statusText(statusRaw: unknown) {
  const status = Number(statusRaw)
  return status === 1 ? '上架' : '下架'
}

function statusClass(statusRaw: unknown) {
  const status = Number(statusRaw)
  if (status === 1) return 'rounded-full bg-emerald-50 px-2.5 py-0.5 text-xs font-medium text-emerald-700'
  return 'rounded-full bg-slate-100 px-2.5 py-0.5 text-xs font-medium text-slate-600'
}

function statusType(statusRaw: unknown): 'success' | 'info' {
  const status = Number(statusRaw)
  return status === 1 ? 'success' : 'info'
}

function areaText(item: Record<string, any>) {
  const min = Number(item.area_min || 0)
  const max = Number(item.area_max || 0)
  if (min > 0 && max > 0) return `${min}-${max} m2`
  const area = Number(item.area || 0)
  return area > 0 ? `${area} m2` : '-'
}

onMounted(async () => {
  await loadCompanies()
  await loadCruises()
  await loadCategories()
  await loadItems()
})

watch(
  () => filters.value.companyId,
  async () => {
    await loadCruises()
    await loadItems()
  },
)
</script>

<style scoped>
.cabin-type-table-wrap {
  scrollbar-gutter: stable;
}
</style>
