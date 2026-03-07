<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import AdminConfirmDialog from '../../components/AdminConfirmDialog.vue'
import { useAdminDeleteDialog } from '../../composables/useAdminDeleteDialog'

const { request } = useApi()

type BookingRow = {
  id: number
  booking_no?: string
  status: string
  total_cents: number
  phone?: string
  voyage_code?: string
  cruise_name?: string
  voyage_id?: number
  user_id?: number
  created_at?: string
}

const items = ref<BookingRow[]>([])
const total = ref(0)
const loading = ref(false)
const error = ref<string | null>(null)
const {
  visible: deleteDialogVisible,
  submitting: deleteSubmitting,
  target: deleteTargetId,
  open: openDeleteDialog,
  close: closeDeleteDialog,
  run: runDelete,
} = useAdminDeleteDialog<number>()

const statusTabs = [
  { key: '', label: '全部' },
  { key: 'pending_payment', label: '待支付' },
  { key: 'paid', label: '已支付' },
  { key: 'cancelled', label: '已取消' },
  { key: 'refunding', label: '退改中' },
  { key: 'refunded', label: '已退款' },
]

const filters = ref({
  keyword: '',
  bookingNo: '',
  phone: '',
  voyageCode: '',
  cruiseName: '',
  status: '',
  startDate: '',
  endDate: '',
  page: 1,
  pageSize: 20,
})

const activeFilterCount = computed(() => {
  return [
    filters.value.keyword,
    filters.value.bookingNo,
    filters.value.phone,
    filters.value.voyageCode,
    filters.value.cruiseName,
    filters.value.status,
    filters.value.startDate,
    filters.value.endDate,
  ].filter((item) => String(item).trim() !== '').length
})

function buildQuery(includePagination = true) {
  const q: Record<string, string | number> = {}
  if (includePagination) {
    q.page = filters.value.page
    q.page_size = filters.value.pageSize
  }
  if (filters.value.keyword.trim()) q.keyword = filters.value.keyword.trim()
  if (filters.value.bookingNo.trim()) q.booking_no = filters.value.bookingNo.trim()
  if (filters.value.phone.trim()) q.phone = filters.value.phone.trim()
  if (filters.value.voyageCode.trim()) q.voyage_code = filters.value.voyageCode.trim()
  if (filters.value.cruiseName.trim()) q.cruise_name = filters.value.cruiseName.trim()
  if (filters.value.status) q.status = filters.value.status
  if (filters.value.startDate) q.start_date = filters.value.startDate
  if (filters.value.endDate) q.end_date = filters.value.endDate
  return q
}

function applyFilters() {
  filters.value.page = 1
  void loadItems()
}

function submitKeywordSearch() {
  applyFilters()
}

async function loadItems() {
  loading.value = true
  error.value = null
  try {
    const res = await request('/bookings', { query: buildQuery() })
    const payload = res?.data ?? res ?? {}
    items.value = payload?.list ?? []
    total.value = Number(payload?.total ?? 0)
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load bookings'
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
  openDeleteDialog(id)
}

async function confirmDelete() {
  const id = resolveId(deleteTargetId.value)
  if (!id) {
    error.value = '无效记录 ID，无法删除'
    closeDeleteDialog()
    return
  }
  error.value = null
  try {
    await runDelete(async () => {
      await request(`/bookings/${id}`, { method: 'DELETE' })
      await loadItems()
    })
  } catch (e: any) {
    error.value = e?.message ?? '删除订单失败，请稍后重试。'
  }
}

function setStatusTab(status: string) {
  filters.value.status = status
  applyFilters()
}

async function exportCurrentRows() {
  error.value = null
  try {
    const response = await request('/bookings/export', {
      query: buildQuery(false),
      responseType: 'blob',
    })
    const blob = response instanceof Blob
      ? response
      : new Blob([typeof response === 'string' ? response : JSON.stringify(response)], { type: 'text/csv;charset=utf-8' })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = `bookings-${Date.now()}.csv`
    link.click()
    URL.revokeObjectURL(link.href)
  } catch (e: any) {
    error.value = e?.message ?? '导出订单失败，请稍后重试。'
  }
}

function statusText(status: string) {
  if (status === 'pending_payment') return '待支付'
  if (status === 'paid') return '已支付'
  if (status === 'cancelled') return '已取消'
  if (status === 'refunding') return '退改中'
  if (status === 'refunded') return '已退款'
  return status || '-'
}

function payActionVisible(status: string) {
  return status === 'pending_payment'
}

function bookingNoText(booking: BookingRow) {
  return booking.booking_no || `#${booking.id}`
}

function bookingMetaText(booking: BookingRow) {
  const parts = [booking.phone, booking.voyage_code, booking.cruise_name].filter(Boolean)
  return parts.length > 0 ? parts.join(' / ') : '暂无附加信息'
}

function formatCreatedAt(value?: string) {
  if (!value) return '-'
  const normalized = String(value)
  if (normalized.includes('T')) {
    return normalized.replace('T', ' ').slice(0, 16)
  }
  return normalized.slice(0, 16)
}

onMounted(loadItems)
</script>

<template>
  <div class="admin-page">
    <AdminPageHeader title="订单管理" subtitle="Booking Orders">
      <template #actions>
        <div data-test="header-actions" class="bookings-header__actions">
          <button type="button" data-test="export" class="admin-btn admin-btn--secondary bookings-header__action" @click="exportCurrentRows">导出 CSV</button>
          <AdminActionLink to="/bookings/new" data-test="create-booking" class="admin-btn bookings-header__action" variant="primary" size="md">新建订单</AdminActionLink>
        </div>
      </template>
    </AdminPageHeader>

    <AdminFilterBar>
      <section class="bookings-search-panel">
        <div class="bookings-search-panel__intro">
          <div>
            <p class="bookings-search-panel__eyebrow">订单搜索台</p>
            <h2 class="bookings-search-panel__title">输入任意订单信息，快速定位目标订单</h2>
          </div>
          <p class="bookings-search-panel__summary">当前已启用 {{ activeFilterCount }} 项筛选条件</p>
        </div>

        <div class="bookings-search-bar">
          <label class="bookings-search-bar__field" for="booking-search-input">
            <span class="bookings-search-bar__icon" aria-hidden="true">检</span>
            <input
              id="booking-search-input"
              v-model="filters.keyword"
              data-test="booking-search-input"
              class="bookings-search-bar__input"
              placeholder="搜索订单号、手机号、航次代码、邮轮名称、状态或金额"
              @keyup.enter="submitKeywordSearch"
            />
          </label>
          <button type="button" data-test="booking-search-submit" class="admin-btn bookings-search-bar__submit" @click="submitKeywordSearch">搜索订单</button>
        </div>

        <div class="bookings-status-tabs">
          <button
            v-for="tab in statusTabs"
            :key="tab.key || 'all'"
            type="button"
            :data-test="`tab-${tab.key || 'all'}`"
            class="bookings-status-tabs__item"
            :class="filters.status === tab.key ? 'bookings-status-tabs__item--active' : ''"
            @click="setStatusTab(tab.key)"
          >
            {{ tab.label }}
          </button>
        </div>

        <form data-test="filter-form" class="bookings-filter-grid" @submit.prevent="applyFilters">
          <input v-model="filters.bookingNo" data-test="filter-booking-no" placeholder="订单号" class="bookings-filter-grid__control" />
          <input v-model="filters.phone" data-test="filter-phone" placeholder="手机号" class="bookings-filter-grid__control" />
          <input v-model="filters.voyageCode" data-test="filter-voyage-code" placeholder="航次代码" class="bookings-filter-grid__control" />
          <input v-model="filters.cruiseName" data-test="filter-cruise-name" placeholder="邮轮名称" class="bookings-filter-grid__control" />
          <input v-model="filters.startDate" data-test="filter-start" type="date" class="bookings-filter-grid__control" />
          <input v-model="filters.endDate" data-test="filter-end" type="date" class="bookings-filter-grid__control" />
          <button type="submit" data-test="filter-submit" class="admin-btn bookings-filter-grid__submit">应用筛选</button>
        </form>
      </section>
    </AdminFilterBar>

    <div class="bookings-overview">
      <p data-test="total" class="bookings-overview__total">总数：{{ total }}</p>
      <p class="bookings-overview__hint">结果可按订单号、联系人、航次和邮轮信息快速识别。</p>
    </div>

    <AdminDataCard flush>
      <p v-if="loading" data-test="loading" class="p-3">Loading...</p>
      <p v-else-if="error" data-test="error" class="error p-3">{{ error }}</p>
      <p v-else-if="items.length === 0" data-test="empty" class="p-3">No data</p>
      <div v-else class="overflow-x-auto">
        <table data-test="table" class="bookings-table">
          <thead class="bookings-table__head">
            <tr>
              <th class="p-3">订单号</th>
              <th class="p-3">联系人</th>
              <th class="p-3">航次 / 邮轮</th>
              <th class="p-3">状态</th>
              <th class="p-3">总额（分）</th>
              <th class="p-3">创建时间</th>
              <th class="p-3">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="b in items" :key="b.id" class="bookings-table__row">
              <td class="p-3 text-slate-900">
                <div class="bookings-table__primary">{{ bookingNoText(b) }}</div>
                <div class="bookings-table__secondary">ID: {{ b.id }}</div>
              </td>
              <td class="p-3 text-slate-700">
                <div class="bookings-table__primary">{{ b.phone || '未绑定手机号' }}</div>
                <div class="bookings-table__secondary">用户 ID: {{ b.user_id ?? '-' }}</div>
              </td>
              <td class="p-3 text-slate-700">
                <div class="bookings-table__primary">{{ b.voyage_code || '未关联航次' }}</div>
                <div class="bookings-table__secondary">{{ b.cruise_name || '未关联邮轮' }}</div>
              </td>
              <td class="p-3">
                <AdminStatusTag :text="statusText(b.status)" :type="b.status === 'pending_payment' ? 'warning' : b.status === 'paid' ? 'success' : b.status === 'cancelled' ? 'info' : 'danger'" />
              </td>
              <td class="p-3 text-slate-700">
                <div class="bookings-table__primary">{{ b.total_cents }}</div>
                <div class="bookings-table__secondary">{{ bookingMetaText(b) }}</div>
              </td>
              <td class="p-3 text-slate-700">{{ formatCreatedAt(b.created_at) }}</td>
              <td class="p-3">
                <div class="bookings-row-actions">
                  <AdminActionLink :to="`/bookings/${b.id}`" class="bookings-row-actions__item">查看详情</AdminActionLink>
                  <AdminActionLink v-if="payActionVisible(b.status)" :to="`/bookings/${b.id}`" class="bookings-row-actions__item">处理支付</AdminActionLink>
                  <AdminActionLink :to="`/bookings/${b.id}`" class="bookings-row-actions__item">处理退改</AdminActionLink>
                  <AdminActionLink :to="`/bookings/${b.id}`" class="bookings-row-actions__item">编辑</AdminActionLink>
                  <button type="button" class="admin-btn admin-btn--danger admin-btn--sm bookings-row-actions__item" @click="handleDelete(b.id)">删除</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </AdminDataCard>

    <AdminConfirmDialog
      :visible="deleteDialogVisible"
      title="确认删除订单"
      :message="`确认删除订单 #${deleteTargetId} 吗？删除后不可恢复。`"
      :loading="deleteSubmitting"
      loading-text="删除中..."
      @close="closeDeleteDialog"
      @confirm="confirmDelete"
    />
  </div>
</template>

<style scoped>
.bookings-header__actions {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 12px;
}

.bookings-header__action {
  min-width: 124px;
  min-height: 44px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.bookings-search-panel {
  display: grid;
  gap: 16px;
  padding: 18px;
  border: 1px solid #dbe7f5;
  border-radius: 16px;
  background: linear-gradient(180deg, #f8fbff 0%, #ffffff 100%);
}

.bookings-search-panel__intro {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
}

.bookings-search-panel__eyebrow {
  margin: 0 0 6px;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #2563eb;
}

.bookings-search-panel__title {
  margin: 0;
  font-size: 20px;
  line-height: 1.35;
  color: #0f172a;
}

.bookings-search-panel__summary {
  margin: 0;
  padding: 8px 12px;
  border-radius: 999px;
  background: #dbeafe;
  color: #1d4ed8;
  font-size: 13px;
  font-weight: 600;
}

.bookings-search-bar {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 12px;
}

.bookings-search-bar__field {
  min-height: 52px;
  display: grid;
  grid-template-columns: 36px minmax(0, 1fr);
  align-items: center;
  gap: 8px;
  padding: 0 14px;
  border: 1px solid #bfdbfe;
  border-radius: 14px;
  background: #ffffff;
  box-shadow: 0 10px 30px rgba(37, 99, 235, 0.08);
}

.bookings-search-bar__field:focus-within {
  border-color: #60a5fa;
  box-shadow: 0 0 0 4px rgba(59, 130, 246, 0.14);
}

.bookings-search-bar__icon {
  width: 32px;
  height: 32px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  background: #dbeafe;
  color: #1d4ed8;
  font-size: 13px;
  font-weight: 700;
}

.bookings-search-bar__input {
  width: 100%;
  min-width: 0;
  border: 0;
  outline: none;
  color: #0f172a;
  font-size: 14px;
  background: transparent;
}

.bookings-search-bar__submit {
  min-width: 124px;
  min-height: 52px;
}

.bookings-status-tabs {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.bookings-status-tabs__item {
  min-height: 42px;
  padding: 0 16px;
  border: 1px solid #dbe4f0;
  border-radius: 999px;
  background: #f8fafc;
  color: #475569;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s ease, border-color 0.2s ease, color 0.2s ease;
}

.bookings-status-tabs__item:hover,
.bookings-status-tabs__item:focus-visible {
  border-color: #93c5fd;
  background: #eff6ff;
  color: #1d4ed8;
  outline: none;
}

.bookings-status-tabs__item--active {
  border-color: #1d4ed8;
  background: #1d4ed8;
  color: #ffffff;
}

.bookings-filter-grid {
  display: grid;
  grid-template-columns: repeat(7, minmax(0, 1fr));
  gap: 12px;
}

.bookings-filter-grid__control {
  min-height: 44px;
  width: 100%;
  border: 1px solid #dbe4f0;
  border-radius: 12px;
  padding: 0 14px;
  outline: none;
  color: #0f172a;
  background: #ffffff;
}

.bookings-filter-grid__control:focus {
  border-color: #60a5fa;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.12);
}

.bookings-filter-grid__submit {
  min-height: 44px;
}

.bookings-overview {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  gap: 10px;
  margin: 12px 0;
}

.bookings-overview__total,
.bookings-overview__hint {
  margin: 0;
  font-size: 13px;
}

.bookings-overview__total {
  color: #334155;
  font-weight: 600;
}

.bookings-overview__hint {
  color: #64748b;
}

.bookings-table {
  width: 100%;
  min-width: 1220px;
  text-align: left;
  font-size: 13px;
}

.bookings-table__head {
  background: #f8fafc;
  color: #64748b;
}

.bookings-table__row {
  border-top: 1px solid #eef2f7;
}

.bookings-table__row:hover {
  background: #f8fbff;
}

.bookings-table__primary {
  color: #0f172a;
  font-weight: 600;
}

.bookings-table__secondary {
  margin-top: 4px;
  color: #64748b;
  line-height: 1.45;
}

.bookings-row-actions {
  display: grid;
  grid-template-columns: repeat(3, 96px);
  gap: 8px;
}

.bookings-row-actions__item {
  min-width: 96px;
  min-height: 36px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
}

@media (max-width: 1024px) {
  .bookings-search-bar,
  .bookings-filter-grid {
    grid-template-columns: 1fr;
  }
}

@media (prefers-reduced-motion: reduce) {
  .bookings-status-tabs__item,
  .bookings-table__row {
    transition: none;
  }
}
</style>
