<!-- web/pages/pay/[id].vue — 前台支付页面 -->
<!-- 根据订单 ID 展示支付按钮，引导用户完成支付 -->
<script setup lang="ts">
import { onMounted, ref } from 'vue'

// 引入支付按钮组件
import PayButton from '../../../components/PayButton.vue'

declare const useApi: any
declare const useRoute: any

// 加载订单详情并把支付参数传递给支付按钮。
const { request } = useApi()
const route = useRoute()
const loading = ref(false)
const error = ref<string | null>(null)
const payUrl = ref('')
const booking = ref<{ id: number; total_cents: number; status?: string } | null>(null)

async function loadBooking() {
  loading.value = true
  error.value = null
  try {
    const id = Number(route?.params?.id ?? 0)
    const res = await request(`/bookings/${id}`)
    booking.value = res?.data ?? res
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load booking'
  } finally {
    loading.value = false
  }
}

function onPaid(url: string) {
  payUrl.value = url
}

function onPayError(msg: string) {
  error.value = msg
}

onMounted(loadBooking)
</script>

<template>
  <div class="page">
    <h1>支付</h1>
    <p v-if="loading">Loading...</p>
    <p v-else-if="error" class="error">{{ error }}</p>
    <div v-else-if="booking">
      <p>订单号：{{ booking.id }}</p>
      <p>订单金额：{{ booking.total_cents }}</p>
      <p v-if="booking.status">订单状态：{{ booking.status }}</p>
      <!-- 支付按钮组件 -->
      <PayButton
        :booking-id="booking.id"
        :amount-cents="booking.total_cents"
        @paid="onPaid"
        @error="onPayError"
      />
      <p v-if="payUrl">支付链接：{{ payUrl }}</p>
    </div>
  </div>
</template>
