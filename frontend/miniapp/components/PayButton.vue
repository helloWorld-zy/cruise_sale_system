<!-- miniapp/components/PayButton.vue — 支付按钮组件 -->
<!-- 小程序端的支付按钮，点击后触发支付流程 -->
<script setup lang="ts">
import { ref } from 'vue'
import { request } from '../src/utils/request'

declare const uni: any

// 小程序支付按钮：调用支付下单接口并触发 requestPayment。
const props = defineProps<{
  bookingId: number
  amountCents: number
}>()

const emit = defineEmits<{
  (e: 'paid', payUrl: string): void
  (e: 'error', msg: string): void
}>()

const loading = ref(false)

async function handlePay() {
  if (loading.value) return
  loading.value = true
  try {
    const raw = await request('/payments', {
      method: 'POST',
      data: {
        order_id: props.bookingId,
        amount_cents: props.amountCents,
        provider: 'wechat',
      },
    })
    const res = raw as Record<string, any>
    const payUrl = res?.pay_url ?? res?.data?.pay_url ?? ''
    const payParams = res?.pay_params ?? res?.data?.pay_params
    if (payParams?.timeStamp) {
      uni.requestPayment({
        ...payParams,
        success: () => emit('paid', payUrl),
        fail: (err: any) => emit('error', err?.errMsg ?? 'payment failed'),
      })
    } else {
      emit('paid', payUrl)
    }
  } catch (e: any) {
    emit('error', e?.message ?? 'payment failed')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <button class="pay-btn" :disabled="loading" @click="handlePay">
    {{ loading ? 'Processing...' : 'Pay Now' }}
  </button>
</template>

<style scoped>
.pay-btn {
  border: 0;
  border-radius: 18rpx;
  height: 86rpx;
  line-height: 86rpx;
  background: linear-gradient(135deg, #0f3d5c, #1f5f86);
  color: #fff;
  font-weight: 700;
  box-shadow: 0 12rpx 26rpx rgba(16, 47, 72, 0.25);
}

.pay-btn[disabled] {
  opacity: 0.6;
}
</style>
