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
  background: #f5f7fa;
  padding: 40rpx;
  position: relative;
}

.page::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 380rpx;
  background: linear-gradient(135deg, #0cebeb 0%, #20e3b2 50%, #29ffc6 100%);
  border-bottom-left-radius: 40rpx;
  border-bottom-right-radius: 40rpx;
  z-index: 0;
}

.title {
  position: relative;
  z-index: 1;
  font-size: 48rpx;
  font-weight: 800;
  color: #fff;
  display: block;
  margin-bottom: 8rpx;
}

.subtitle {
  position: relative;
  z-index: 1;
  display: block;
  margin-top: 8rpx;
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.9);
}

.tabs {
  position: relative;
  z-index: 1;
  display: flex;
  gap: 16rpx;
  margin: 32rpx 0 24rpx;
  flex-wrap: wrap;
}

.tab {
  border-radius: 999rpx;
  border: none;
  background: rgba(255, 255, 255, 0.25);
  padding: 10rpx 28rpx;
  color: #fff;
  font-size: 26rpx;
  font-weight: 600;
}

.tab-active {
  background: #fff;
  color: #0cebeb;
}

.list {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  gap: 24rpx;
  padding-bottom: 40rpx;
}

.card {
  border-radius: 32rpx;
  background: #fff;
  border: none;
  box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.05);
  padding: 32rpx;
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}

.order-no,
.status,
.amount {
  color: #666;
  font-size: 26rpx;
  font-weight: 500;
}

.order-no {
  font-size: 28rpx;
  font-weight: 700;
  color: #222;
}

.amount {
  color: #ff6b6b;
  font-weight: 800;
  font-size: 36rpx;
  margin: 8rpx 0;
}

.actions {
  margin-top: 16rpx;
  display: flex;
  gap: 16rpx;
  flex-wrap: wrap;
  justify-content: flex-end;
}

.btn {
  border-radius: 999rpx;
  border: none;
  background: #f5f7fa;
  padding: 12rpx 32rpx;
  font-size: 24rpx;
  color: #333;
  font-weight: 600;
}

.primary {
  background: linear-gradient(135deg, #ff8e53 0%, #ff6b6b 100%);
  color: #fff;
}

.hint {
  position: relative;
  z-index: 1;
  color: #fff;
  text-align: center;
  display: block;
  margin-top: 60rpx;
}

.error {
  position: relative;
  z-index: 1;
  color: #ffcccc;
  text-align: center;
  display: block;
  margin-top: 60rpx;
}
</style>
