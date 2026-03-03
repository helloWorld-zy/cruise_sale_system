<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

type OrderRow = {
  id: number
  status: string
  total_cents: number
  created_at?: string
}

const { request } = useApi()
const loading = ref(false)
const error = ref('')
const items = ref<OrderRow[]>([])
const activeStatus = ref('')

const tabs = [
  { key: '', label: '全部' },
  { key: 'pending_payment', label: '待支付' },
  { key: 'paid', label: '已支付' },
  { key: 'cancelled', label: '已取消' },
  { key: 'refunding', label: '退改中' },
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
    const token = globalThis?.sessionStorage?.getItem('auth_token') || ''
    if (!token) {
      items.value = []
      return
    }

    const res = await request('/admin/bookings', { query: { page: 1, page_size: 50 } })
    const payload = res?.data ?? res ?? {}
    items.value = Array.isArray(payload) ? payload : payload?.list ?? []
  } catch (e: any) {
    error.value = e?.message ?? '加载订单失败'
  } finally {
    loading.value = false
  }
}

onMounted(loadOrders)
</script>

<template>
  <div class="min-h-screen bg-[#f8f4ed] pb-24 text-slate-900">
    <main class="mx-auto max-w-5xl px-6 py-8">
      <h1 class="mb-5 font-['Playfair_Display','Georgia',serif] text-3xl text-[#12263a]">我的订单</h1>

      <div class="mb-4 flex flex-wrap gap-2">
        <button
          v-for="tab in tabs"
          :key="tab.key || 'all'"
          type="button"
          class="rounded-full px-4 py-2 text-sm"
          :class="activeStatus === tab.key ? 'bg-[#12263a] text-white' : 'bg-white text-slate-700 border border-[#eadfcb]'"
          @click="activeStatus = tab.key"
        >
          {{ tab.label }}
        </button>
      </div>

      <p v-if="loading" data-test="loading" class="py-10 text-center text-slate-600">加载中...</p>
      <p v-else-if="error" data-test="error" class="py-10 text-center text-rose-600">{{ error }}</p>
      <p v-else-if="filtered.length === 0" data-test="empty" class="py-10 text-center text-slate-600">暂无订单</p>

      <div v-else class="space-y-3" data-test="list">
        <div v-for="item in filtered" :key="item.id" class="rounded-xl border border-[#eadfcb] bg-white p-4 shadow-sm">
          <div class="flex items-start justify-between gap-4">
            <div>
              <p class="text-sm text-slate-500">订单 #{{ item.id }}</p>
              <p class="mt-1 text-sm text-slate-600">状态：{{ statusLabel(item.status) }}</p>
              <p class="mt-1 font-semibold text-[#12263a]">¥{{ (item.total_cents / 100).toFixed(2) }}</p>
            </div>
            <div class="flex flex-wrap justify-end gap-2">
              <NuxtLink
                v-if="item.status === 'pending_payment'"
                :to="`/pay/${item.id}`"
                class="rounded-md bg-[#12263a] px-3 py-1.5 text-xs text-white"
              >
                去支付
              </NuxtLink>
              <NuxtLink :to="`/orders/${item.id}`" class="rounded-md border border-[#d6cab4] px-3 py-1.5 text-xs text-[#12263a]">
                查看详情
              </NuxtLink>
              <NuxtLink :to="`/orders/${item.id}?action=refund`" class="rounded-md border border-[#d6cab4] px-3 py-1.5 text-xs text-[#12263a]">
                申请退改
              </NuxtLink>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>
