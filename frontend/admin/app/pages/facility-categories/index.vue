<template>
  <div class="admin-page">
    <AdminPageHeader title="设施分类管理" subtitle="维护分类基础信息、图标和启停状态，支持行内编辑与快速发布。">
      <template #actions>
        <button type="button" class="admin-btn facility-category-new-btn" @click="appendRow">新增分类</button>
      </template>
    </AdminPageHeader>

    <AdminDataCard flush>
      <div class="facility-category-toolbar" aria-live="polite">
        <div class="facility-category-toolbar__item">
          <span class="facility-category-toolbar__label">分类数量</span>
          <strong class="facility-category-toolbar__value">{{ rows.length }}</strong>
        </div>
        <div class="facility-category-toolbar__item">
          <span class="facility-category-toolbar__label">启用数量</span>
          <strong class="facility-category-toolbar__value">{{ rows.filter((item) => Number(item.status) === 1).length }}</strong>
        </div>
        <div class="facility-category-toolbar__hint">建议：直接按图标选择分类样式，排序值越小展示越靠前。</div>
      </div>

      <div class="facility-category-table-wrap overflow-x-auto">
          <table class="facility-category-table w-full min-w-[920px] text-sm">
          <thead class="facility-category-table__head bg-slate-50 text-left text-slate-600">
            <tr>
              <th class="p-3">名称</th>
              <th class="p-3">图标</th>
              <th class="p-3">排序</th>
              <th class="p-3">状态</th>
              <th class="p-3 whitespace-nowrap">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td class="facility-category-table__state p-3" colspan="5">加载中...</td>
            </tr>
            <tr v-else-if="error">
              <td class="facility-category-table__state p-3 text-rose-500" colspan="5">{{ error }}</td>
            </tr>
            <tr v-else-if="rows.length === 0">
              <td class="facility-category-table__state p-3" colspan="5">暂无数据</td>
            </tr>
            <tr
              v-for="(row, idx) in rows"
              v-else
              :key="row.id || idx"
              class="facility-category-table__row"
              :class="idx % 2 === 1 ? 'bg-slate-50/60' : ''"
            >
              <td class="p-3">
                <label class="facility-category-sr-only" :for="`facility-category-name-${row.id || idx}`">分类名称</label>
                <input
                  :id="`facility-category-name-${row.id || idx}`"
                  v-model="row.name"
                  class="facility-category-input h-9 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2"
                  placeholder="请输入分类名称"
                />
              </td>
              <td class="p-3">
                <div class="facility-category-icon-input-group flex items-center gap-2">
                  <button
                    type="button"
                    class="facility-category-picker-toggle"
                    :data-test="`facility-category-icon-trigger-${row.id || idx}`"
                    :aria-expanded="row.iconPickerOpen ? 'true' : 'false'"
                    @click="row.iconPickerOpen = !row.iconPickerOpen"
                  >
                    <span class="facility-category-picker-toggle__preview" :data-test="`facility-category-icon-current-${row.id || idx}`" :data-icon="row.icon || ''">
                      <FacilityCategoryIcon :name="row.icon" />
                    </span>
                    <span class="facility-category-picker-toggle__text">{{ getFacilityCategoryIconLabel(row.icon) }}</span>
                    <span class="facility-category-picker-toggle__hint">选择图标</span>
                  </button>
                </div>
                <div v-if="row.iconPickerOpen" class="facility-category-icon-picker mt-2 grid grid-cols-3 gap-2 rounded-md border border-slate-200 p-2">
                  <button
                    v-for="icon in getIconOptions(row.icon)"
                    :key="icon.value"
                    type="button"
                    class="facility-category-icon-option"
                    :class="row.icon === icon.value ? 'facility-category-icon-option--active' : ''"
                    :data-test="`facility-category-icon-option-${row.id || idx}-${icon.value}`"
                    @click="selectIcon(row, icon.value)"
                  >
                    <FacilityCategoryIcon :name="icon.value" />
                    <span class="facility-category-icon-option__label">{{ icon.label }}</span>
                  </button>
                </div>
              </td>
              <td class="p-3">
                <label class="facility-category-sr-only" :for="`facility-category-sort-${row.id || idx}`">排序值</label>
                <input
                  :id="`facility-category-sort-${row.id || idx}`"
                  v-model.number="row.sort_order"
                  type="number"
                  class="facility-category-input h-9 w-24 rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2"
                />
              </td>
              <td class="p-3">
                <label class="facility-category-sr-only" :for="`facility-category-status-${row.id || idx}`">状态</label>
                <select
                  :id="`facility-category-status-${row.id || idx}`"
                  v-model.number="row.status"
                  class="facility-category-select h-9 rounded-md border border-slate-200 px-2 outline-none ring-indigo-500 focus:ring-2"
                >
                  <option :value="1">启用</option>
                  <option :value="0">停用</option>
                </select>
              </td>
              <td class="p-3 whitespace-nowrap">
                <div class="facility-category-actions flex items-center gap-2">
                  <button type="button" class="admin-btn admin-btn--sm" @click="saveRow(row)">保存</button>
                  <AdminActionLink v-if="row.id" :to="`/facility-categories/${row.id}`" class="admin-btn admin-btn--sm">编辑页</AdminActionLink>
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
      title="确认删除设施分类"
      :message="`确认删除分类「${deleteTarget?.name || `#${deleteTarget?.id ?? ''}`}」吗？删除后不可恢复。`"
      :loading="deleteSubmitting"
      loading-text="删除中..."
      @close="closeDeleteDialog"
      @confirm="confirmDelete"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import AdminConfirmDialog from '../../components/AdminConfirmDialog.vue'
import FacilityCategoryIcon from '../../components/facility-categories/FacilityCategoryIcon.vue'
import { useAdminDeleteDialog } from '../../composables/useAdminDeleteDialog'
import { ensureFacilityCategoryIconOptions, getFacilityCategoryIconLabel } from '../../constants/facilityCategoryIcons'

const { request } = useApi()
const loading = ref(false)
const error = ref<string | null>(null)

type Row = {
  id?: number
  name: string
  icon: string
  sort_order: number
  status: number
  iconPickerOpen?: boolean
}

const rows = ref<Row[]>([])
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

function getIconOptions(currentIcon: string) {
  return ensureFacilityCategoryIconOptions(currentIcon)
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
      await request(`/facility-categories/${id}`, { method: 'DELETE' })
      await loadItems()
    })
  } catch (e: any) {
    error.value = e?.message ?? '删除设施分类失败，请稍后重试。'
  }
}

onMounted(loadItems)
</script>

<style scoped>
.facility-category-table-wrap {
  scrollbar-gutter: stable;
}

.facility-category-new-btn {
  min-height: 40px;
}

.facility-category-toolbar {
  display: grid;
  gap: 10px;
  grid-template-columns: repeat(1, minmax(0, 1fr));
  padding: 14px;
  border-bottom: 1px solid #edf2f7;
  background: linear-gradient(135deg, #f8fbff 0%, #f2f7ff 100%);
}

.facility-category-toolbar__item {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-height: 40px;
}

.facility-category-toolbar__label {
  color: #64748b;
  font-size: 13px;
}

.facility-category-toolbar__value {
  color: #334155;
  font-size: 18px;
  line-height: 1;
}

.facility-category-toolbar__hint {
  color: #475569;
  font-size: 13px;
  line-height: 1.6;
}

.facility-category-table {
  border-top: 0;
}

.facility-category-table__head th {
  position: sticky;
  top: 0;
  z-index: 1;
}

.facility-category-table__state {
  color: #475569;
  line-height: 1.6;
}

.facility-category-table__row {
  transition: background-color 0.2s ease;
}

.facility-category-input,
.facility-category-select,
.facility-category-icon-option {
  min-height: 38px;
}

.facility-category-icon-input-group {
  position: relative;
}

.facility-category-picker-toggle {
  min-height: 56px;
  width: 100%;
  padding: 10px 12px;
  display: grid;
  grid-template-columns: 36px minmax(0, 1fr) auto;
  align-items: center;
  gap: 10px;
  border: 1px solid #dbe4f0;
  border-radius: 10px;
  background: #fff;
  color: #334155;
  text-align: left;
  transition: border-color 0.2s ease, box-shadow 0.2s ease, background-color 0.2s ease;
}

.facility-category-picker-toggle:hover,
.facility-category-picker-toggle:focus-visible {
  border-color: #93c5fd;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.12);
  outline: none;
}

.facility-category-picker-toggle__preview {
  width: 36px;
  height: 36px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  background: linear-gradient(135deg, #eff6ff 0%, #dbeafe 100%);
  color: #1d4ed8;
}

.facility-category-picker-toggle__text {
  font-size: 13px;
  font-weight: 600;
  color: #0f172a;
}

.facility-category-picker-toggle__hint {
  font-size: 12px;
  color: #64748b;
}

.facility-category-icon-picker {
  background: #fff;
  border: 1px solid #dbe4f0;
  box-shadow: 0 10px 26px rgba(15, 23, 42, 0.1);
}

.facility-category-icon-option {
  border: 1px solid #dbe4f0;
  border-radius: 10px;
  background: #fff;
  color: #334155;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 8px;
  font-size: 12px;
  cursor: pointer;
  transition: border-color 0.2s ease, color 0.2s ease, background-color 0.2s ease, box-shadow 0.2s ease;
}

.facility-category-icon-option:hover,
.facility-category-icon-option:focus-visible {
  border-color: #93c5fd;
  background: #eff6ff;
  color: #1d4ed8;
  outline: none;
}

.facility-category-icon-option--active {
  border-color: #60a5fa;
  background: #dbeafe;
  color: #1e40af;
  box-shadow: inset 0 0 0 1px rgba(59, 130, 246, 0.28);
}

.facility-category-icon-option__label {
  line-height: 1.3;
}

.facility-category-actions {
  display: grid;
  grid-template-columns: repeat(3, 90px);
  gap: 8px;
}

.facility-category-actions :deep(.admin-btn) {
  width: 90px;
  min-width: 90px;
  min-height: 34px;
  padding: 0 10px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
}

.facility-category-sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}

@media (min-width: 768px) {
  .facility-category-toolbar {
    grid-template-columns: auto auto 1fr;
    align-items: center;
  }
}

@media (prefers-reduced-motion: reduce) {
  .facility-category-table__row,
  .facility-category-icon-option {
    transition: none;
  }
}
</style>
