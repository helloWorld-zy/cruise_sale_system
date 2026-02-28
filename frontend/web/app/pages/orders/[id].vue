<!-- web/pages/orders/[id].vue — 订单详情页面 -->
<!-- 根据路由参数加载并展示订单状态和金额 -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'

// 获取路由参数中的订单 ID
const route = useRoute()
const orderId = route.params.id as string

// 订单详情接口类型
interface OrderDetail {
  id: string     // 订单 ID
  status: string // 订单状态
  amount: number // 订单金额（单位：分）
}

// 响应式状态
const order = ref<OrderDetail | null>(null) // 订单详情数据
const loading = ref(true)                   // 加载状态
const error = ref<string | null>(null)      // 错误信息

// 页面挂载时加载订单详情数据
onMounted(async () => {
  try {
    const res = await fetch(`/api/v1/bookings/${orderId}`)
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    order.value = await res.json()
  } catch (e) {
    error.value = e instanceof Error ? e.message : '加载订单失败'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="page">
    <h1>订单状态</h1>
    <div v-if="loading">加载中…</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <template v-else-if="order">
      <p>订单号 #{{ order.id }}</p>
      <p>状态: {{ order.status }}</p>
      <p>金额: ¥{{ (order.amount / 100).toFixed(2) }}</p>
    </template>
    <div v-else>订单未找到。</div>
  </div>
</template>
