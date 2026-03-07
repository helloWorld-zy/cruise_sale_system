<!-- web/pages/orders/[id].vue — 订单详情页面 -->
<!-- 根据路由参数加载并展示订单状态和金额 -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'

declare const useApi: any

// 获取路由参数中的订单 ID
const route = useRoute()
const orderId = route.params.id as string
const { request } = useApi()

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
    let token = ''
    try {
      token = globalThis?.sessionStorage?.getItem('auth_token') || ''
    } catch {
      token = ''
    }

    // C-end currently has no public booking detail API; skip remote call when unauthenticated.
    if (!token) {
      order.value = {
        id: orderId,
        status: '待登录后查看',
        amount: 0,
      }
      return
    }

    const payload = await request(`/admin/bookings/${orderId}`)
    const data = payload?.data ?? payload
    if (!data) {
      order.value = null
    } else {
      order.value = {
        id: String(data.id ?? orderId),
        status: String(data.status ?? '未知'),
        amount: Number(data.amount ?? data.amount_cents ?? 0),
      }
    }
  } catch (e) {
    error.value = e instanceof Error ? e.message : '加载订单失败'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="page max-w-xl mx-auto mt-10 text-center">
    <h1 class="text-[var(--color-primary)]">订单状态</h1>
    
    <div class="mt-8 p-6 bg-slate-50 rounded-xl border border-slate-100 flex flex-col items-center gap-4">
      <div v-if="loading" class="text-[var(--web-muted)] italic">加载中...</div>
      <div v-else-if="error" class="error bg-red-50 p-3 rounded-lg w-full">{{ error }}</div>
      <template v-else-if="order">
        <div class="text-xs font-bold tracking-widest text-[var(--web-muted)] uppercase mb-2">Order #{{ order.id }}</div>
        
        <div class="inline-flex items-center px-4 py-2 rounded-full font-bold text-sm" 
             :class="order.status.includes('成功') || order.status.includes('已') ? 'bg-emerald-100 text-emerald-800' : 'bg-amber-100 text-amber-800'">
          {{ order.status }}
        </div>
        
        <div class="mt-4 heading-font text-4xl text-[var(--color-cta)]">
          ¥{{ (order.amount / 100).toFixed(2) }}
        </div>
      </template>
      <div v-else class="text-[var(--web-muted)]">订单未找到。</div>
    </div>
    
    <div class="mt-8">
      <NuxtLink to="/orders" class="text-sm font-semibold text-[var(--color-primary)] hover:text-[var(--color-cta)] transition-colors underline underline-offset-4">
        返回订单列表
      </NuxtLink>
    </div>
  </div>
</template>
