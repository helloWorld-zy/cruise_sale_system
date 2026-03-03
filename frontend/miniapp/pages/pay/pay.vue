<!-- miniapp/pages/pay/pay.vue — 小程序端支付页面 -->
<!-- 展示支付按钮，引导用户完成支付流程 -->
<script setup lang="ts">
import { onMounted, ref } from 'vue'
// 引入支付按钮组件
import PayButton from '../../components/PayButton.vue'
import { request } from '../../src/utils/request'

// 通过 props 接收 bookingId，便于运行时和测试复用。
const props = defineProps<{ bookingId?: number | string; preview?: boolean }>()
const loading = ref(false)
const error = ref('')
const payUrl = ref('')
const booking = ref<{ id: number; total_cents: number } | null>(null)

function resolveBookingId() {
  return Number(props.bookingId ?? 0)
}

async function loadBooking() {
  if (props.preview) {
    booking.value = {
      id: resolveBookingId() || 2026030501,
      total_cents: 536000,
    }
    return
  }
  const id = resolveBookingId()
  if (!id) {
    error.value = '缺少 bookingId 参数'
    return
  }
  loading.value = true
  error.value = ''
  try {
    const res = await request(`/bookings/${id}`)
    booking.value = res?.data ?? res
  } catch (e: any) {
    error.value = e?.message ?? '加载订单失败'
  } finally {
    loading.value = false
  }
}

function onPaid(url: string) {
  payUrl.value = url
}

function onError(msg: string) {
  error.value = msg
}

onMounted(loadBooking)
</script>

<template>
  <view class="page">
    <text class="title">支付订单</text>
    <text class="subtitle">确认订单后完成支付，舱位将优先为你保留。</text>
    <text v-if="loading" class="hint">Loading...</text>
    <text v-else-if="error" class="error">{{ error }}</text>
    <view v-else-if="booking" class="panel">
      <text class="meta">订单号：{{ booking.id }}</text>
      <text class="amount">金额：¥{{ (booking.total_cents / 100).toFixed(2) }}</text>
      <PayButton :booking-id="booking.id" :amount-cents="booking.total_cents" @paid="onPaid" @error="onError" />
      <text v-if="payUrl" class="hint">支付链接：{{ payUrl }}</text>
    </view>
  </view>
</template>

<style scoped>
.page {
  min-height: 100vh;
  padding: 30rpx;
  background:
    radial-gradient(circle at 8% 0, #dbeaf6 0, transparent 30%),
    linear-gradient(180deg, #f3f8fb 0%, #edf3f7 100%);
}

.title {
  display: block;
  margin-bottom: 8rpx;
  font-size: 46rpx;
  font-weight: 700;
  color: #122b42;
}

.subtitle {
  display: block;
  margin-bottom: 16rpx;
  font-size: 24rpx;
  color: #5b728a;
}

.panel {
  display: flex;
  flex-direction: column;
  gap: 14rpx;
  background: #fff;
  border-radius: 24rpx;
  padding: 26rpx;
  border: 1rpx solid #d4e0ea;
  box-shadow: 0 16rpx 36rpx rgba(16, 47, 72, 0.12);
}

.meta {
  color: #5a7189;
}

.amount {
  font-size: 34rpx;
  font-weight: 700;
  color: #113d5c;
}

.hint {
  color: #5f6288;
}

.error {
  color: #d13e5b;
}
</style>
