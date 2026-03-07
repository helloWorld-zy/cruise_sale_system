<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import AdminConfirmDialog from '../../components/AdminConfirmDialog.vue'
import { useAdminDeleteDialog } from '../../composables/useAdminDeleteDialog'

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
const {
  visible: deleteDialogVisible,
  submitting: deleteSubmitting,
  target: deleteTarget,
  open: openDeleteDialog,
  close: closeDeleteDialog,
  run: runDelete,
} = useAdminDeleteDialog<{ id: number; name: string }>()
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
  openDeleteDialog({ id, name: String(item?.name ?? '') })
}

async function confirmDelete() {
  const id = resolveId(deleteTarget.value?.id)
  if (!id) {
    error.value = '无效记录 ID，无法删除'
    closeDeleteDialog()
    return
  }
  error.value = null
  try {
    await runDelete(async () => {
      await request(`/voyages/${id}`, { method: 'DELETE' })
      await loadItems()
    })
  } catch (e: any) {
    error.value = e?.message ?? '删除航次失败，请稍后重试。'
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
  <div class="admin-page">
    <AdminPageHeader title="航次管理">
      <template #actions>
      <AdminActionLink to="/voyages/new" variant="primary" size="md">新建航次</AdminActionLink>
      </template>
    </AdminPageHeader>

    <AdminFilterBar>
      <div class="flex flex-wrap items-center gap-2">
      <input v-model="filters.code" data-test="filter-code" placeholder="按航次代码搜索" />
      <input v-model="filters.departDate" data-test="filter-depart" type="date" />
      <input v-model="filters.returnDate" data-test="filter-return" type="date" />
      <input v-model="filters.port" data-test="filter-port" placeholder="按港口关键字搜索" />
      </div>
    </AdminFilterBar>

    <AdminDataCard flush>
      <p v-if="loading" class="p-3">Loading...</p>
      <p v-else-if="error" class="p-3 text-red-500">{{ error }}</p>
      <p v-else-if="filteredItems.length === 0" class="p-3">No data</p>
      <div v-else class="voyage-table-wrap">
        <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Code</th>
            <th>简介</th>
            <th>首站港口</th>
            <th>行程天数</th>
            <th class="voyage-date-col">出发</th>
            <th class="voyage-date-col">结束</th>
            <th class="voyage-action-col">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="r in filteredItems" :key="r.id">
            <td>{{ r.id }}</td>
            <td>{{ r.code }}</td>
            <td>{{ r.brief_info || '-' }}</td>
            <td>{{ r.first_stop_city || '-' }}</td>
            <td>{{ r.itinerary_days || '-' }}</td>
            <td class="voyage-date-cell">{{ dateOnly(r.depart_date) }}</td>
            <td class="voyage-date-cell">{{ dateOnly(r.return_date) }}</td>
            <td class="voyage-actions-cell">
              <div class="voyage-actions">
                <AdminActionLink :to="`/voyages/${r.id}`">编辑</AdminActionLink>
                <button type="button" @click="handleDelete(r.id)">删除</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      </div>
    </AdminDataCard>

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

<style scoped>
.voyage-table-wrap {
  overflow-x: auto;
}

.voyage-table-wrap table {
  min-width: 860px;
}

.voyage-date-col,
.voyage-date-cell {
  white-space: nowrap;
}

.voyage-action-col,
.voyage-actions-cell {
  white-space: nowrap;
  width: 1%;
}

.voyage-actions {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}
</style>

