<!-- miniapp/pages/pay/pay.vue — 小程序端支付页面 -->
<!-- 展示支付按钮，引导用户完成支付流程 -->
<script setup lang="ts">
import { onMounted, ref } from 'vue'
// 引入支付按钮组件
import PayButton from '../../components/PayButton.vue'
import { request } from '../../src/utils/request'

// 通过 props 接收 bookingId，便于运行时和测试复用。
const props = defineProps<{ bookingId?: number | string }>()
const loading = ref(false)
const error = ref('')
const payUrl = ref('')
const booking = ref<{ id: number; total_cents: number } | null>(null)

function resolveBookingId() {
  return Number(props.bookingId ?? 0)
}

async function loadBooking() {
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
    <text v-if="loading" class="hint">Loading...</text>
    <text v-else-if="error" class="error">{{ error }}</text>
    <view v-else-if="booking" class="panel">
      <text>订单号：{{ booking.id }}</text>
      <text>金额：{{ booking.total_cents }}</text>
      <PayButton :booking-id="booking.id" :amount-cents="booking.total_cents" @paid="onPaid" @error="onError" />
      <text v-if="payUrl" class="hint">支付链接：{{ payUrl }}</text>
    </view>
  </view>
</template>

<style scoped>
.page {
  min-height: 100vh;
  padding: 36rpx;
  background: #f8f6ff;
}

.title {
  display: block;
  margin-bottom: 20rpx;
  font-size: 42rpx;
  font-weight: 700;
  color: #50348b;
}

.panel {
  display: flex;
  flex-direction: column;
  gap: 14rpx;
  background: #fff;
  border-radius: 18rpx;
  padding: 26rpx;
  box-shadow: 0 12rpx 28rpx rgba(73, 53, 133, 0.12);
}

.hint {
  color: #5f6288;
}

.error {
  color: #d13e5b;
}
</style>
