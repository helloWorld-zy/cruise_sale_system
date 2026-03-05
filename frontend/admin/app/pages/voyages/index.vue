<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import AdminConfirmDialog from '../../components/AdminConfirmDialog.vue'

type VoyageRow = {
  id: number
  code: string
  brief_info?: string
  depart_date?: string
  return_date?: string
  itinerary_days?: number
  first_stop_city?: string
}

const { request } = useApi()
const items = ref<VoyageRow[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const deleteDialogVisible = ref(false)
const deleteSubmitting = ref(false)
const deleteTarget = ref<{ id: number; name: string } | null>(null)
const filters = ref({
  code: '',
  departDate: '',
  returnDate: '',
  port: '',
})

async function loadItems() {
  loading.value = true
  error.value = null
  try {
    const res = await request('/voyages')
    const payload = res?.data ?? res ?? []
    items.value = Array.isArray(payload) ? payload : payload?.list ?? []
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load voyages'
  } finally {
    loading.value = false
  }
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
  const item = items.value.find((it) => resolveId(it?.id) === id)
  deleteTarget.value = { id, name: String(item?.name ?? '') }
  deleteDialogVisible.value = true
}

function closeDeleteDialog() {
  if (deleteSubmitting.value) return
  deleteDialogVisible.value = false
  deleteTarget.value = null
}

async function confirmDelete() {
  const id = resolveId(deleteTarget.value?.id)
  if (!id) {
    error.value = '无效记录 ID，无法删除'
    closeDeleteDialog()
    return
  }
  error.value = null
  deleteSubmitting.value = true
  try {
    await request(`/voyages/${id}`, { method: 'DELETE' })
    closeDeleteDialog()
    await loadItems()
  } catch (e: any) {
    error.value = e?.message ?? '删除航次失败，请稍后重试。'
  } finally {
    deleteSubmitting.value = false
  }
}

function dateOnly(input?: string) {
  if (!input) return '-'
  return input.slice(0, 10)
}

const filteredItems = computed(() => {
  const codeQ = filters.value.code.trim().toLowerCase()
  const departQ = filters.value.departDate
  const returnQ = filters.value.returnDate
  const portQ = filters.value.port.trim().toLowerCase()
  return items.value.filter((item) => {
    const code = String(item.code || '').toLowerCase()
    const depart = dateOnly(item.depart_date)
    const ret = dateOnly(item.return_date)
    const port = String(item.first_stop_city || item.brief_info || '').toLowerCase()

    if (codeQ && !code.includes(codeQ)) return false
    if (departQ && depart !== departQ) return false
    if (returnQ && ret !== returnQ) return false
    if (portQ && !port.includes(portQ)) return false
    return true
  })
})

onMounted(loadItems)
</script>

<template>
  <div class="page">
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:12px;">
      <h1>航次管理</h1>
      <AdminActionLink to="/voyages/new" variant="primary" size="md">新建航次</AdminActionLink>
    </div>
    <div style="display:flex;gap:8px;flex-wrap:wrap;margin-bottom:10px;">
      <input v-model="filters.code" data-test="filter-code" placeholder="按航次代码搜索" />
      <input v-model="filters.departDate" data-test="filter-depart" type="date" />
      <input v-model="filters.returnDate" data-test="filter-return" type="date" />
      <input v-model="filters.port" data-test="filter-port" placeholder="按港口关键字搜索" />
    </div>

    <p v-if="loading">Loading...</p>
    <p v-else-if="error" class="text-red-500">{{ error }}</p>
    <p v-else-if="filteredItems.length === 0">No data</p>
    <table v-else>
      <thead>
        <tr>
          <th>ID</th>
          <th>Code</th>
          <th>简介</th>
          <th>首站港口</th>
          <th>行程天数</th>
          <th>出发</th>
          <th>结束</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="r in filteredItems" :key="r.id">
          <td>{{ r.id }}</td>
          <td>{{ r.code }}</td>
          <td>{{ r.brief_info || '-' }}</td>
          <td>{{ r.first_stop_city || '-' }}</td>
          <td>{{ r.itinerary_days || '-' }}</td>
          <td>{{ dateOnly(r.depart_date) }}</td>
          <td>{{ dateOnly(r.return_date) }}</td>
          <td>
            <AdminActionLink :to="`/voyages/${r.id}`">编辑</AdminActionLink>
            <button type="button" style="margin-left:8px" @click="handleDelete(r.id)">删除</button>
          </td>
        </tr>
      </tbody>
    </table>

    <AdminConfirmDialog
      :visible="deleteDialogVisible"
      title="确认删除航次"
      :message="`确认删除航次「${deleteTarget?.name || `#${deleteTarget?.id ?? ''}`}」吗？删除后不可恢复。`"
      :loading="deleteSubmitting"
      loading-text="删除中..."
      @close="closeDeleteDialog"
      @confirm="confirmDelete"
    />
  </div>
</template>

