<script setup lang="ts">
import { onMounted, ref } from 'vue'

const { request } = useApi()

type BookingRow = {
  id: number
  status: string
  total_cents: number
  voyage_id?: number
  user_id?: number
  created_at?: string
}

const items = ref<BookingRow[]>([])
const total = ref(0)
const loading = ref(false)
const error = ref<string | null>(null)

const statusTabs = [
  { key: '', label: '全部' },
  { key: 'pending_payment', label: '待支付' },
  { key: 'paid', label: '已支付' },
  { key: 'cancelled', label: '已取消' },
  { key: 'refunding', label: '退改中' },
  { key: 'refunded', label: '已退款' },
]

const filters = ref({
  bookingNo: '',
  phone: '',
  routeId: '',
  status: '',
  startDate: '',
  endDate: '',
  page: 1,
  pageSize: 20,
})

function buildQuery() {
  const q: Record<string, string | number> = {
    page: filters.value.page,
    page_size: filters.value.pageSize,
  }
  if (filters.value.bookingNo.trim()) q.booking_no = filters.value.bookingNo.trim()
  if (filters.value.phone.trim()) q.phone = filters.value.phone.trim()
  if (filters.value.routeId.trim()) q.route_id = filters.value.routeId.trim()
  if (filters.value.status) q.status = filters.value.status
  if (filters.value.startDate) q.start_date = filters.value.startDate
  if (filters.value.endDate) q.end_date = filters.value.endDate
  return q
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
  if (!confirm(`确认删除订单 #${id} 吗？`)) return
  try {
    await request(`/bookings/${id}`, { method: 'DELETE' })
    await loadItems()
  } catch (e: any) {
    error.value = e?.message ?? 'failed to delete booking'
  }
}

function setStatusTab(status: string) {
  filters.value.status = status
  filters.value.page = 1
  void loadItems()
}

function exportCurrentRows() {
  const head = ['id', 'status', 'total_cents', 'voyage_id', 'user_id', 'created_at']
  const lines = [head.join(',')]
  for (const item of items.value) {
    lines.push([
      item.id,
      item.status || '',
      item.total_cents ?? 0,
      item.voyage_id ?? '',
      item.user_id ?? '',
      item.created_at ?? '',
    ].join(','))
  }

  const blob = new Blob([lines.join('\n')], { type: 'text/csv;charset=utf-8' })
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = `bookings-${Date.now()}.csv`
  link.click()
  URL.revokeObjectURL(link.href)
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

onMounted(loadItems)
</script>

<template>
  <div class="page">
    <div style="display:flex;justify-content:space-between;align-items:center;gap:8px;margin-bottom:12px;flex-wrap:wrap;">
      <h1>Bookings</h1>
      <div style="display:flex;gap:8px;align-items:center;">
        <button type="button" data-test="export" @click="exportCurrentRows">导出 CSV</button>
        <NuxtLink to="/bookings/new"><button type="button">新建订单</button></NuxtLink>
      </div>
    </div>

    <div style="display:flex;gap:8px;flex-wrap:wrap;margin-bottom:8px;">
      <button
        v-for="tab in statusTabs"
        :key="tab.key || 'all'"
        type="button"
        :data-test="`tab-${tab.key || 'all'}`"
        :style="{ background: filters.status === tab.key ? '#111827' : '#f3f4f6', color: filters.status === tab.key ? '#fff' : '#111827' }"
        @click="setStatusTab(tab.key)"
      >
        {{ tab.label }}
      </button>
    </div>

    <form style="display:grid;grid-template-columns:repeat(auto-fit,minmax(160px,1fr));gap:8px;margin-bottom:10px;" @submit.prevent="loadItems">
      <input v-model="filters.bookingNo" data-test="filter-booking-no" placeholder="订单号" />
      <input v-model="filters.phone" data-test="filter-phone" placeholder="手机号" />
      <input v-model="filters.routeId" data-test="filter-route" placeholder="航线ID" />
      <input v-model="filters.startDate" data-test="filter-start" type="date" />
      <input v-model="filters.endDate" data-test="filter-end" type="date" />
      <button type="submit" data-test="filter-submit">筛选</button>
    </form>

    <p data-test="total">总数：{{ total }}</p>

    <p v-if="loading" data-test="loading">Loading...</p>
    <p v-else-if="error" data-test="error" class="error">{{ error }}</p>
    <p v-else-if="items.length === 0" data-test="empty">No data</p>
    <table v-else data-test="table">
      <thead>
        <tr>
          <th>ID</th>
          <th>状态</th>
          <th>总额</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="b in items" :key="b.id">
          <td>{{ b.id }}</td>
          <td>{{ statusText(b.status) }}</td>
          <td>{{ b.total_cents }}</td>
          <td>
            <NuxtLink :to="`/bookings/${b.id}`">查看详情</NuxtLink>
            <NuxtLink v-if="payActionVisible(b.status)" :to="`/bookings/${b.id}`" style="margin-left:8px">处理支付</NuxtLink>
            <NuxtLink :to="`/bookings/${b.id}`" style="margin-left:8px">处理退改</NuxtLink>
            <NuxtLink :to="`/bookings/${b.id}`" style="margin-left:8px">编辑</NuxtLink>
            <button type="button" style="margin-left:8px" @click="handleDelete(b.id)">删除</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
