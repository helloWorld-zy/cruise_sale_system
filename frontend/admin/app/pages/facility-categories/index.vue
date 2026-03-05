<template>
  <div class="min-h-screen bg-slate-50 p-4 md:p-6">
    <div class="mx-auto max-w-6xl">
      <div class="mb-4 flex items-center justify-between">
        <h1 class="text-xl font-semibold text-slate-900">设施分类管理</h1>
        <button type="button" class="rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-500" @click="appendRow">新增分类</button>
      </div>

      <div class="rounded-lg border border-slate-200 bg-white shadow-sm">
        <table class="w-full text-sm">
          <thead class="bg-slate-50 text-left text-slate-600">
            <tr>
              <th class="p-3">名称</th>
              <th class="p-3">图标</th>
              <th class="p-3">排序</th>
              <th class="p-3">状态</th>
              <th class="p-3">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td class="p-3" colspan="5">加载中...</td>
            </tr>
            <tr v-else-if="error">
              <td class="p-3 text-rose-500" colspan="5">{{ error }}</td>
            </tr>
            <tr v-else-if="rows.length === 0">
              <td class="p-3" colspan="5">暂无数据</td>
            </tr>
            <tr v-for="(row, idx) in rows" v-else :key="row.id || idx" :class="idx % 2 === 1 ? 'bg-slate-50' : ''">
              <td class="p-3">
                <input v-model="row.name" class="h-9 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
              </td>
              <td class="p-3">
                <div class="flex items-center gap-2">
                  <input v-model="row.icon" class="h-9 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" placeholder="icon 名称" />
                  <button type="button" class="rounded border border-slate-200 px-2 py-1 text-xs hover:bg-slate-50" @click="row.iconPickerOpen = !row.iconPickerOpen">选择</button>
                </div>
                <div v-if="row.iconPickerOpen" class="mt-2 grid grid-cols-6 gap-2 rounded-md border border-slate-200 p-2">
                  <button
                    v-for="icon in iconOptions"
                    :key="icon"
                    type="button"
                    class="rounded border border-slate-200 px-2 py-1 text-xs hover:bg-slate-50"
                    @click="selectIcon(row, icon)"
                  >
                    {{ icon }}
                  </button>
                </div>
              </td>
              <td class="p-3">
                <input v-model.number="row.sort_order" type="number" class="h-9 w-24 rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
              </td>
              <td class="p-3">
                <select v-model.number="row.status" class="h-9 rounded-md border border-slate-200 px-2 outline-none ring-indigo-500 focus:ring-2">
                  <option :value="1">启用</option>
                  <option :value="0">停用</option>
                </select>
              </td>
              <td class="p-3">
                <div class="flex items-center gap-2">
                  <button type="button" class="rounded border border-emerald-200 px-2 py-1 text-xs text-emerald-700 hover:bg-emerald-50" @click="saveRow(row)">保存</button>
                  <NuxtLink v-if="row.id" :to="`/facility-categories/${row.id}`" class="rounded border border-slate-200 px-2 py-1 text-xs text-slate-600 hover:bg-slate-50">编辑页</NuxtLink>
                  <button v-if="row.id" type="button" class="rounded border border-rose-200 px-2 py-1 text-xs text-rose-600 hover:bg-rose-50" @click="removeRow(row)">删除</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <AdminConfirmDialog
        :visible="deleteDialogVisible"
        title="确认删除设施分类"
        :message="`确认删除分类「${deleteTarget?.name || `#${deleteTarget?.id ?? ''}`}」吗？删除后不可恢复。`"
        :loading="deleteSubmitting"
        loading-text="删除中..."
        @close="closeDeleteDialog"
        @confirm="confirmDelete"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import AdminConfirmDialog from '../../components/AdminConfirmDialog.vue'

const { request } = useApi()
const loading = ref(false)
const error = ref<string | null>(null)
const iconOptions = ['utensils', 'music', 'spa', 'dumbbell', 'swimmer', 'gamepad']

type Row = {
  id?: number
  name: string
  icon: string
  sort_order: number
  status: number
  iconPickerOpen?: boolean
}

const rows = ref<Row[]>([])
const deleteDialogVisible = ref(false)
const deleteSubmitting = ref(false)
const deleteTarget = ref<Row | null>(null)

async function loadItems() {
  loading.value = true
  error.value = null
  try {
    const res = await request('/facility-categories')
    const payload = res?.data ?? res ?? []
    const list = Array.isArray(payload) ? payload : payload?.list ?? []
    rows.value = list.map((item: Record<string, any>) => ({
      id: Number(item.id || 0) || undefined,
      name: item.name || '',
      icon: item.icon || '',
      sort_order: Number(item.sort_order || 0),
      status: Number(item.status ?? 1),
      iconPickerOpen: false,
    }))
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load facility categories'
  } finally {
    loading.value = false
  }
}

function appendRow() {
  rows.value.unshift({ name: '', icon: '', sort_order: 0, status: 1, iconPickerOpen: false })
}

function selectIcon(row: Row, icon: string) {
  row.icon = icon
  row.iconPickerOpen = false
}

async function saveRow(row: Row) {
  const payload = {
    name: row.name,
    icon: row.icon,
    sort_order: Number(row.sort_order || 0),
    status: Number(row.status ?? 1),
  }
  try {
    if (row.id) {
      await request(`/facility-categories/${row.id}`, { method: 'PUT', body: payload })
    } else {
      await request('/facility-categories', { method: 'POST', body: payload })
    }
    await loadItems()
  } catch (e: any) {
    error.value = e?.message ?? 'failed to save facility category'
  }
}

async function removeRow(row: Row) {
  if (!row.id) return
  deleteTarget.value = row
  deleteDialogVisible.value = true
}

function closeDeleteDialog() {
  if (deleteSubmitting.value) return
  deleteDialogVisible.value = false
  deleteTarget.value = null
}

async function confirmDelete() {
  const id = Number(deleteTarget.value?.id ?? 0)
  if (!Number.isFinite(id) || id <= 0) {
    error.value = '无效记录 ID，无法删除'
    closeDeleteDialog()
    return
  }
  deleteSubmitting.value = true
  error.value = null
  try {
    await request(`/facility-categories/${id}`, { method: 'DELETE' })
    closeDeleteDialog()
    await loadItems()
  } catch (e: any) {
    error.value = e?.message ?? '删除设施分类失败，请稍后重试。'
  } finally {
    deleteSubmitting.value = false
  }
}

onMounted(loadItems)
</script>
