<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted } from 'vue'

const route = useRoute()
let redirectTimer: ReturnType<typeof setTimeout> | null = null

const orderId = computed(() => {
  const raw = Number(route.query.order_id)
  return Number.isFinite(raw) && raw > 0 ? String(Math.trunc(raw)) : ''
})

function goToOrder() {
  if (!orderId.value) return
  navigateTo(`/orders/${orderId.value}`)
}

onMounted(() => {
  if (!orderId.value) return
  redirectTimer = setTimeout(() => {
    goToOrder()
  }, 5000)
})

onBeforeUnmount(() => {
  if (redirectTimer) {
    clearTimeout(redirectTimer)
  }
})
</script>

<template>
  <div class="page">
    <h1>预订成功</h1>
    <p>订单已创建，将在 5 秒后自动跳转到订单详情页。</p>
    <p v-if="orderId">订单号 #{{ orderId }}</p>
    <p v-else class="error">未获取到订单号，请返回订单中心查看。</p>
    <button type="button" :disabled="!orderId" @click="goToOrder">立即查看订单</button>
  </div>
</template>
