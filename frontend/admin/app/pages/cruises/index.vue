<template>
  <div class="min-h-screen bg-slate-50 p-4 md:p-6">
    <div class="mx-auto max-w-7xl">
      <div class="mb-4 flex items-center justify-between">
        <h1 class="text-xl font-semibold text-slate-900">邮轮管理</h1>
        <NuxtLink to="/cruises/create" class="rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-500">新建邮轮</NuxtLink>
      </div>

      <div class="mb-4 rounded-lg border border-slate-200 bg-white p-3 shadow-sm">
        <div class="flex flex-wrap items-center gap-3">
          <input
            v-model="filters.keyword"
            data-test="filter-keyword"
            class="h-10 w-64 rounded-md border border-slate-200 px-3 text-sm text-slate-700 outline-none ring-indigo-500 focus:ring-2"
            placeholder="搜索名称/代码"
            @keyup.enter="loadItems"
          />
          <select
            v-model.number="filters.status"
            data-test="filter-status"
            class="h-10 rounded-md border border-slate-200 px-3 text-sm text-slate-700 outline-none ring-indigo-500 focus:ring-2"
          >
            <option :value="0">全部状态</option>
            <option :value="1">上架</option>
            <option :value="2">维护中</option>
            <option :value="-1">下架</option>
          </select>
          <select
            v-model="filters.sortBy"
            class="h-10 rounded-md border border-slate-200 px-3 text-sm text-slate-700 outline-none ring-indigo-500 focus:ring-2"
          >
            <option value="">默认排序</option>
            <option value="tonnage_desc">吨位降序</option>
            <option value="tonnage_asc">吨位升序</option>
            <option value="name_asc">名称 A-Z</option>
            <option value="name_desc">名称 Z-A</option>
          </select>
          <button type="button" class="h-10 rounded-md border border-slate-200 px-3 text-sm text-slate-700 hover:bg-slate-50" @click="loadItems">筛选</button>
        </div>
      </div>

      <div class="rounded-lg border border-slate-200 bg-white shadow-sm">
        <table class="w-full text-sm">
          <thead class="bg-slate-50 text-left text-slate-600">
            <tr>
              <th class="w-10 p-3">
                <input type="checkbox" :checked="allChecked" @change="toggleCheckAll(($event.target as HTMLInputElement).checked)" />
              </th>
              <th class="p-3">名称</th>
              <th class="p-3">代码</th>
              <th class="p-3">公司ID</th>
              <th class="p-3">吨位</th>
              <th class="p-3">载客量</th>
              <th class="p-3">状态</th>
              <th class="p-3">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td class="p-3" colspan="8">加载中...</td>
            </tr>
            <tr v-else-if="error">
              <td class="p-3 text-rose-500" colspan="8">{{ error }}</td>
            </tr>
            <tr v-else-if="items.length === 0">
              <td class="p-3" colspan="8">暂无数据</td>
            </tr>
            <tr
              v-for="(item, idx) in items"
              v-else
              :key="item.id"
              :class="idx % 2 === 1 ? 'bg-slate-50' : ''"
            >
              <td class="p-3">
                <input type="checkbox" :checked="selectedIds.has(item.id)" @change="toggleSingle(item.id, ($event.target as HTMLInputElement).checked)" />
              </td>
              <td class="p-3 font-medium text-slate-900">{{ item.name || '-' }}</td>
              <td class="p-3 text-slate-600">{{ item.code || '-' }}</td>
              <td class="p-3 text-slate-600">{{ item.company_id || '-' }}</td>
              <td class="p-3 text-slate-600">{{ item.tonnage || '-' }}</td>
              <td class="p-3 text-slate-600">{{ item.passenger_capacity || '-' }}</td>
              <td class="p-3">
                <span :class="statusClass(item.status)">{{ statusText(item.status) }}</span>
              </td>
              <td class="p-3">
                <NuxtLink :to="`/cruises/${item.id}`" class="text-indigo-600 hover:text-indigo-500">编辑</NuxtLink>
                <button type="button" class="ml-2 text-rose-500 hover:text-rose-400" @click="handleDelete(item.id)">删除</button>
              </td>
            </tr>
          </tbody>
        </table>

        <div class="flex items-center justify-end gap-3 border-t border-slate-200 p-3 text-sm text-slate-600">
          <button type="button" class="rounded border border-slate-200 px-3 py-1.5 hover:bg-slate-50" :disabled="filters.page <= 1" @click="changePage(filters.page - 1)">上一页</button>
          <span>第 {{ filters.page }} 页</span>
          <button type="button" class="rounded border border-slate-200 px-3 py-1.5 hover:bg-slate-50" :disabled="items.length < filters.pageSize" @click="changePage(filters.page + 1)">下一页</button>
        </div>
      </div>

      <div
        v-if="selectedIds.size > 0"
        data-test="batch-action"
        class="fixed bottom-0 left-0 right-0 flex items-center justify-center gap-3 bg-indigo-600 px-4 py-3 text-sm text-white"
      >
        <span>已选 {{ selectedIds.size }} 项</span>
        <button type="button" class="rounded bg-white/20 px-3 py-1.5 hover:bg-white/30" @click="batchUpdateStatus(1)">批量上架</button>
        <button type="button" class="rounded bg-white/20 px-3 py-1.5 hover:bg-white/30" @click="batchUpdateStatus(-1)">批量下架</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

const { request } = useApi()
const items = ref<Record<string, any>[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const selectedIds = ref<Set<number>>(new Set())
const total = ref(0)

const filters = ref({
  keyword: '',
  status: 0,
  sortBy: '',
  page: 1,
  pageSize: 10,
})

const allChecked = computed(() => items.value.length > 0 && items.value.every((it) => selectedIds.value.has(Number(it.id))))

async function loadItems() {
  loading.value = true
  error.value = null
  try {
    const params: Record<string, any> = {
      page: filters.value.page,
      page_size: filters.value.pageSize,
    }
    if (filters.value.keyword.trim()) params.keyword = filters.value.keyword.trim()
    if (filters.value.status !== 0) params.status = filters.value.status === -1 ? 0 : filters.value.status
    if (filters.value.sortBy) params.sort_by = filters.value.sortBy
    const res = await request('/cruises', { query: params })
    const payload = res?.data ?? res ?? {}
    items.value = Array.isArray(payload) ? payload : payload?.list ?? []
    total.value = Number(payload?.total ?? items.value.length)
    selectedIds.value = new Set()
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load cruises'
  } finally {
    loading.value = false
  }
}

function toggleSingle(idRaw: unknown, checked: boolean) {
  const id = Number(idRaw)
  if (!Number.isFinite(id) || id <= 0) return
  const next = new Set(selectedIds.value)
  if (checked) next.add(id)
  else next.delete(id)
  selectedIds.value = next
}

function toggleCheckAll(checked: boolean) {
  if (!checked) {
    selectedIds.value = new Set()
    return
  }
  selectedIds.value = new Set(items.value.map((it) => Number(it.id)).filter((it) => Number.isFinite(it) && it > 0))
}

function statusText(statusRaw: unknown) {
  const status = Number(statusRaw)
  if (status === 1) return '上架'
  if (status === 2) return '维护中'
  return '下架'
}

function statusClass(statusRaw: unknown) {
  const status = Number(statusRaw)
  if (status === 1) return 'rounded-full bg-emerald-50 px-2.5 py-0.5 text-xs font-medium text-emerald-700'
  if (status === 2) return 'rounded-full bg-amber-50 px-2.5 py-0.5 text-xs font-medium text-amber-700'
  return 'rounded-full bg-rose-50 px-2.5 py-0.5 text-xs font-medium text-rose-700'
}

async function batchUpdateStatus(status: number) {
  if (selectedIds.value.size === 0) return
  try {
    await request('/cruises/batch-status', {
      method: 'PUT',
      body: {
        ids: Array.from(selectedIds.value),
        status,
      },
    })
    await loadItems()
  } catch (e: any) {
    error.value = e?.message ?? 'failed to batch update status'
  }
}

function changePage(page: number) {
  if (page <= 0) return
  filters.value.page = page
  void loadItems()
}

function resolveId(raw: unknown) {
  const id = Number(raw)
  return Number.isFinite(id) && id > 0 ? id : 0
}

async function handleDelete(rawId: unknown) {
  const id = resolveId(rawId)
  if (!id) {
    error.value = '无效记录 ID，无法删除'
    return
  }
  if (!confirm(`确认删除邮轮 #${id} 吗？`)) return
  try {
    await request(`/cruises/${id}`, { method: 'DELETE' })
    await loadItems()
  } catch (e: any) {
    error.value = e?.message ?? 'failed to delete cruise'
  }
}

onMounted(loadItems)
</script>

