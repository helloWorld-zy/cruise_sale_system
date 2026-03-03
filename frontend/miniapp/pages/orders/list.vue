<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { request } from '../../src/utils/request'

type OrderItem = {
  id: number
  status: string
  total_cents: number
  created_at?: string
}

const loading = ref(false)
const error = ref('')
const items = ref<OrderItem[]>([])
const activeStatus = ref('')

const tabs = [
  { key: '', label: '全部' },
  { key: 'pending_payment', label: '待支付' },
  { key: 'paid', label: '已支付' },
  { key: 'cancelled', label: '已取消' },
]

const filtered = computed(() => {
  if (!activeStatus.value) return items.value
  return items.value.filter((item) => item.status === activeStatus.value)
})

function statusLabel(status: string) {
  if (status === 'pending_payment') return '待支付'
  if (status === 'paid') return '已支付'
  if (status === 'cancelled') return '已取消'
  if (status === 'refunding') return '退改中'
  if (status === 'refunded') return '已退款'
  return status
}

async function loadOrders() {
  loading.value = true
  error.value = ''
  try {
    const res: any = await request('/bookings', { data: { page: 1, page_size: 50 } })
    const payload = res?.data ?? res ?? {}
    items.value = Array.isArray(payload) ? payload : payload?.list ?? []
  } catch (e: any) {
    error.value = e?.message || '订单加载失败'
  } finally {
    loading.value = false
  }
}

onMounted(loadOrders)
</script>

<template>
  <view class="page">
    <text class="title">我的订单</text>
    <text class="subtitle">按状态筛选订单，快速完成支付或退改申请。</text>

    <view class="tabs">
      <text
        v-for="tab in tabs"
        :key="tab.key || 'all'"
        class="tab"
        :class="activeStatus === tab.key ? 'tab-active' : ''"
        @click="activeStatus = tab.key"
      >
        {{ tab.label }}
      </text>
    </view>

    <text v-if="loading" class="hint">Loading...</text>
    <text v-else-if="error" class="error">{{ error }}</text>
    <text v-else-if="filtered.length === 0" class="hint">暂无订单</text>

    <view v-else class="list">
      <view v-for="item in filtered" :key="item.id" class="card">
        <text class="order-no">订单 #{{ item.id }}</text>
        <text class="status">状态：{{ statusLabel(item.status) }}</text>
        <text class="amount">¥{{ (item.total_cents / 100).toFixed(2) }}</text>
        <view class="actions">
          <navigator v-if="item.status === 'pending_payment'" :url="`/pages/pay/index?id=${item.id}`" class="btn primary">去支付</navigator>
          <navigator :url="`/pages/orders/detail?id=${item.id}`" class="btn">查看详情</navigator>
          <navigator :url="`/pages/orders/detail?id=${item.id}&action=refund`" class="btn">申请退改</navigator>
        </view>
      </view>
    </view>
  </view>
</template>

<style scoped>
.page {
  min-height: 100vh;
  background:
    radial-gradient(circle at 8% 0, #dceaf6 0, transparent 30%),
    linear-gradient(180deg, #f3f8fb 0%, #edf3f7 100%);
  padding: 28rpx;
}

.title {
  font-size: 46rpx;
  font-weight: 700;
  color: #122b42;
}

.subtitle {
  display: block;
  margin-top: 8rpx;
  font-size: 24rpx;
  color: #5b728a;
}

.tabs {
  display: flex;
  gap: 12rpx;
  margin: 18rpx 0 20rpx;
  flex-wrap: wrap;
}

.tab {
  border-radius: 999rpx;
  border: 1rpx solid #c7d1dd;
  padding: 8rpx 18rpx;
  color: #566a80;
  font-size: 24rpx;
}

.tab-active {
  border-color: #113d5c;
  background: #113d5c;
  color: #fff;
}

.list {
  display: flex;
  flex-direction: column;
  gap: 14rpx;
}

.card {
  border-radius: 20rpx;
  background: #fff;
  border: 1rpx solid #d4e0ea;
  box-shadow: 0 12rpx 30rpx rgba(16, 47, 72, 0.1);
  padding: 18rpx;
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}

.order-no,
.status,
.amount {
  color: #334a62;
  font-size: 24rpx;
}

.amount {
  color: #113d5c;
  font-weight: 700;
}

.actions {
  margin-top: 8rpx;
  display: flex;
  gap: 10rpx;
  flex-wrap: wrap;
}

.btn {
  border-radius: 999rpx;
  border: 1rpx solid #c7d1dd;
  padding: 6rpx 14rpx;
  font-size: 22rpx;
  color: #2d4560;
}

.primary {
  border-color: #113d5c;
  background: #113d5c;
  color: #fff;
}

.hint {
  color: #5a7190;
}

.error {
  color: #d13e5b;
}
</style>
