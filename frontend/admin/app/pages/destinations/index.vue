<template>
  <div class="admin-page">
    <AdminPageHeader title="港口城市词典" subtitle="维护航次站点搜索使用的标准港口城市、别名与坐标。">
      <template #actions>
        <button type="button" class="admin-btn admin-btn--secondary dest-toolbar-btn" data-test="dest-export" @click="exportCsv">导出 CSV</button>
        <label class="admin-btn admin-btn--secondary dest-toolbar-btn" for="dest-import-input">导入 CSV</label>
        <button type="button" class="admin-btn dest-new-btn" @click="appendRow">新增城市</button>
      </template>
    </AdminPageHeader>
    <input id="dest-import-input" data-test="dest-import-input" type="file" accept=".csv,text/csv" class="dest-import-input" @change="importCsv" />

    <AdminDataCard flush>
      <div class="dest-toolbar" aria-live="polite">
        <div class="dest-toolbar__item">
          <span class="dest-toolbar__label">目的地数量</span>
          <strong class="dest-toolbar__value">{{ filteredRows.length }}</strong>
        </div>
        <div class="dest-toolbar__item">
          <span class="dest-toolbar__label">启用数量</span>
          <strong class="dest-toolbar__value">{{ filteredRows.filter((r) => Number(r.status) === 1).length }}</strong>
        </div>
        <div class="dest-toolbar__item dest-toolbar__item--search">
          <span class="dest-toolbar__label">搜索词典</span>
          <input
            v-model="searchKeyword"
            data-test="dest-search"
            class="dest-search-input"
            placeholder="搜索名称、国家/地区或关键词"
          />
        </div>
        <div class="dest-toolbar__hint">提示：这里维护的是本地港口城市词典。城市坐标会直接用于行程地图与站点坐标补全，纬度、经度必填；关键词请用逗号分隔，写入常见英文名、繁体名和别名。</div>
        <div v-if="importSummary" class="dest-toolbar__hint">{{ importSummary }}</div>
      </div>

      <div class="dest-table-wrap overflow-x-auto">
        <table class="dest-table w-full min-w-[1140px] text-sm">
          <thead class="dest-table__head bg-slate-50 text-left text-slate-600">
            <tr>
              <th class="p-3">名称</th>
              <th class="p-3">国家/地区</th>
              <th class="p-3">纬度</th>
              <th class="p-3">经度</th>
              <th class="p-3">关键词</th>
              <th class="p-3">排序</th>
              <th class="p-3">状态</th>
              <th class="p-3 whitespace-nowrap">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td class="dest-table__state p-3" colspan="8">加载中...</td>
            </tr>
            <tr v-else-if="error">
              <td class="dest-table__state p-3 text-rose-500" colspan="8">{{ error }}</td>
            </tr>
            <tr v-else-if="filteredRows.length === 0">
              <td class="dest-table__state p-3" colspan="8">暂无数据，点击「新增城市」添加。</td>
            </tr>
            <tr
              v-for="(row, idx) in filteredRows"
              v-else
              :key="row.id || idx"
              class="dest-table__row"
              :class="idx % 2 === 1 ? 'bg-slate-50/60' : ''"
            >
              <td class="p-3">
                <input
                  v-model="row.name"
                  class="dest-input dest-input--name h-9 rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2"
                  placeholder="如 迈阿密"
                />
              </td>
              <td class="p-3">
                <input
                  v-model="row.country"
                  class="dest-input dest-input--country h-9 rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2"
                  placeholder="如 美国"
                />
              </td>
              <td class="p-3">
                <input
                  v-model.number="row.latitude"
                  type="number"
                  step="0.0001"
                  class="dest-input dest-input--coord h-9 rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2"
                  placeholder="25.7617"
                />
              </td>
              <td class="p-3">
                <input
                  v-model.number="row.longitude"
                  type="number"
                  step="0.0001"
                  class="dest-input dest-input--coord h-9 rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2"
                  placeholder="-77.9407"
                />
              </td>
              <td class="p-3">
                <input
                  v-model="row.keywords"
                  class="dest-input dest-input--keywords h-9 rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2"
                  placeholder="如 迈阿密,邁阿密,miami"
                />
              </td>
              <td class="p-3">
                <input
                  v-model.number="row.sort_order"
                  type="number"
                  class="dest-input dest-input--sort h-9 rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2"
                />
              </td>
              <td class="p-3">
                <select
                  v-model.number="row.status"
                  class="dest-select h-9 rounded-md border border-slate-200 px-2 outline-none ring-indigo-500 focus:ring-2"
                >
                  <option :value="1">启用</option>
                  <option :value="0">停用</option>
                </select>
              </td>
              <td class="p-3 whitespace-nowrap">
                <div class="dest-actions flex items-center gap-2">
                  <button type="button" class="admin-btn admin-btn--sm" @click="saveRow(row)">保存</button>
                  <button v-if="row.id" type="button" class="admin-btn admin-btn--danger admin-btn--sm" @click="removeRow(row)">删除</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </AdminDataCard>

    <AdminConfirmDialog
      :visible="deleteDialogVisible"
      title="确认删除目的地"
      :message="`确认删除目的地「${deleteTarget?.name || `#${deleteTarget?.id ?? ''}`}」吗？删除后不可恢复。`"
      :loading="deleteSubmitting"
      loading-text="删除中..."
      @close="closeDeleteDialog"
      @confirm="confirmDelete"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import AdminConfirmDialog from '../../components/AdminConfirmDialog.vue'
import { useAdminDeleteDialog } from '../../composables/useAdminDeleteDialog'

const { request } = useApi()
const loading = ref(false)
const error = ref<string | null>(null)
const importSummary = ref('')
const searchKeyword = ref('')

type Row = {
  id?: number
  name: string
  country: string
  latitude: number | null
  longitude: number | null
  keywords: string
  sort_order: number
  status: number
}

const rows = ref<Row[]>([])

const filteredRows = computed(() => {
  const keyword = searchKeyword.value.trim().toLowerCase()
  if (!keyword) {
    return rows.value
  }
  return rows.value.filter((row) => {
    return [row.name, row.country, row.keywords].some((field) => String(field || '').toLowerCase().includes(keyword))
  })
})

const {
  visible: deleteDialogVisible,
  submitting: deleteSubmitting,
  target: deleteTarget,
  open: openDeleteDialog,
  close: closeDeleteDialog,
  run: runDelete,
} = useAdminDeleteDialog<Row>()

async function loadItems() {
  loading.value = true
  error.value = null
  try {
    const res = await request('/custom-destinations')
    const payload = res?.data ?? res ?? []
    const list = Array.isArray(payload) ? payload : payload?.list ?? []
    rows.value = list.map((item: Record<string, any>) => ({
      id: Number(item.id || 0) || undefined,
      name: item.name || '',
      country: item.country || '',
      latitude: item.latitude != null ? Number(item.latitude) : null,
      longitude: item.longitude != null ? Number(item.longitude) : null,
      keywords: item.keywords || '',
      sort_order: Number(item.sort_order || 0),
      status: Number(item.status ?? 1),
    }))
  } catch (e: any) {
    error.value = e?.message ?? '加载港口城市词典失败'
  } finally {
    loading.value = false
  }
}

async function exportCsv() {
  error.value = null
  const content = await request('/custom-destinations/export', { responseType: 'text' })
  const blob = new Blob([content], { type: 'text/csv;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = 'port_city_dictionary.csv'
  link.click()
  URL.revokeObjectURL(url)
}

async function importCsv(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) {
    return
  }
  const formData = new FormData()
  formData.append('file', file)
  try {
    const res = await request('/custom-destinations/import', { method: 'POST', body: formData })
    const payload = res?.data ?? res ?? {}
    importSummary.value = `已导入 ${Number(payload.imported || 0)} 条记录`
    error.value = null
    await loadItems()
  } catch (e: any) {
    error.value = e?.message ?? '导入港口城市词典失败'
  } finally {
    target.value = ''
  }
}

function appendRow() {
  rows.value.unshift({
    name: '',
    country: '',
    latitude: null,
    longitude: null,
    keywords: '',
    sort_order: 0,
    status: 1,
  })
}

async function saveRow(row: Row) {
  if (!row.name?.trim()) {
    error.value = '名称不能为空'
    return
  }
  if (!row.country?.trim()) {
    error.value = '国家/地区不能为空'
    return
  }
  if (row.latitude == null || row.longitude == null || Number.isNaN(Number(row.latitude)) || Number.isNaN(Number(row.longitude))) {
    error.value = '纬度和经度不能为空'
    return
  }
  const payload = {
    name: row.name.trim(),
    country: row.country.trim(),
    latitude: row.latitude,
    longitude: row.longitude,
    keywords: row.keywords || '',
    sort_order: Number(row.sort_order || 0),
    status: Number(row.status ?? 1),
  }
  try {
    if (row.id) {
      await request(`/custom-destinations/${row.id}`, { method: 'PUT', body: payload })
    } else {
      await request('/custom-destinations', { method: 'POST', body: payload })
    }
    error.value = null
    await loadItems()
  } catch (e: any) {
    error.value = e?.message ?? '保存港口城市词典失败'
  }
}

function removeRow(row: Row) {
  if (!row.id) return
  openDeleteDialog(row)
}

async function confirmDelete() {
  const id = Number(deleteTarget.value?.id ?? 0)
  if (!Number.isFinite(id) || id <= 0) {
    error.value = '无效记录 ID，无法删除'
    closeDeleteDialog()
    return
  }
  error.value = null
  try {
    await runDelete(async () => {
      await request(`/custom-destinations/${id}`, { method: 'DELETE' })
      await loadItems()
    })
  } catch (e: any) {
    error.value = e?.message ?? '删除港口城市词典记录失败，请稍后重试。'
  }
}

onMounted(loadItems)
</script>

<style scoped>
.dest-import-input {
	display: none;
}

.dest-table-wrap {
  scrollbar-gutter: stable;
  padding: 0 10px 10px;
}

.dest-new-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 40px;
  line-height: 1;
}

.dest-toolbar-btn {
	display: inline-flex;
	align-items: center;
	justify-content: center;
	min-height: 40px;
	line-height: 1;
}

.dest-toolbar {
  display: grid;
  gap: 10px;
  grid-template-columns: repeat(1, minmax(0, 1fr));
  padding: 14px;
  border-bottom: 1px solid #edf2f7;
  background: linear-gradient(135deg, #f8fbff 0%, #f2f7ff 100%);
}

.dest-toolbar__item {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-height: 40px;
}

.dest-toolbar__item--search {
  flex-wrap: wrap;
}

.dest-toolbar__label {
  color: #64748b;
  font-size: 13px;
}

.dest-toolbar__value {
  color: #334155;
  font-size: 18px;
  line-height: 1;
}

.dest-toolbar__hint {
  color: #475569;
  font-size: 13px;
  line-height: 1.6;
}

.dest-search-input {
  width: min(320px, 100%);
  min-height: 38px;
  border: 1px solid #cbd5e1;
  border-radius: 4px;
  padding: 0 12px;
  outline: none;
}

.dest-search-input:focus {
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.15);
}

.dest-table {
  border-top: 0;
}

.dest-table__head th {
  position: sticky;
  top: 0;
  z-index: 1;
}

.dest-table__state {
  color: #475569;
  line-height: 1.6;
}

.dest-table__row {
  transition: background-color 0.2s ease;
}

.dest-input,
.dest-select {
  min-height: 38px;
}

.dest-input--name {
  width: 188px;
  min-width: 188px;
}

.dest-input--country {
  width: 200px;
  min-width: 200px;
}

.dest-input--coord {
  width: 112px;
  min-width: 112px;
}

.dest-input--keywords {
  width: 176px;
  min-width: 176px;
}

.dest-input--sort {
  width: 88px;
  min-width: 88px;
}

.dest-actions {
  display: grid;
  grid-template-columns: repeat(2, 72px);
  gap: 8px;
}

.dest-actions :deep(.admin-btn) {
  width: 72px;
  min-width: 72px;
  min-height: 34px;
  padding: 0 8px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
}
</style>
