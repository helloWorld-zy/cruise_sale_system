<!-- web/components/PayButton.vue — 支付按钮组件 -->
<!-- 前台 Web 端的支付按钮，点击后触发支付流程 -->
<script setup lang="ts">
import { ref } from 'vue'

declare const useApi: any

// 支付按钮接收订单 ID 和金额，发起支付并向外抛出结果事件。
const props = defineProps<{
  bookingId: number
  amountCents: number
}>()

const emit = defineEmits<{
  (e: 'paid', payUrl: string): void
  (e: 'error', msg: string): void
}>()

const { request } = useApi()
const loading = ref(false)

async function handlePay() {
  if (loading.value) return
  loading.value = true
  try {
    const res = await request('/payments', {
      method: 'POST',
      body: {
        order_id: props.bookingId,
        amount_cents: props.amountCents,
        provider: 'wechat',
      },
    })
    emit('paid', res?.pay_url ?? res?.data?.pay_url ?? '')
  } catch (e: any) {
    emit('error', e?.message ?? 'payment failed')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <button :disabled="loading" @click="handlePay">
    {{ loading ? 'Processing...' : 'Pay Now' }}
  </button>
</template>
