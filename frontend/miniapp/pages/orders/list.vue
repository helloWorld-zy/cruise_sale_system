<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { request } from '../../src/utils/request'
import NavBar from '../../components/NavBar.vue'

declare const uni: any

type OrderItem = {
  id: number
  status: string
  total_cents: number
  created_at?: string
}

const emit = defineEmits<{
  (e: 'open-pay', bookingId: number): void
  (e: 'open-order-detail', bookingId: number): void
}>()

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

/** 导航到支付页，优先 uni.navigateTo，降级 emit */
function goToPay(id: number) {
  try {
    if (typeof uni !== 'undefined' && typeof uni.navigateTo === 'function') {
      uni.navigateTo({ url: `/pages/pay/pay?id=${id}` })
      return
    }
  } catch { /* ignored */ }
  emit('open-pay', id)
}

/** 导航到订单详情，优先 uni.navigateTo，降级 emit */
function goToDetail(id: number, action?: string) {
  const query = action ? `?id=${id}&action=${action}` : `?id=${id}`
  try {
    if (typeof uni !== 'undefined' && typeof uni.navigateTo === 'function') {
      uni.navigateTo({ url: `/pages/orders/detail${query}` })
      return
    }
  } catch { /* ignored */ }
  emit('open-order-detail', id)
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
  <view class="min-h-screen bg-background pb-8 overflow-x-hidden">
    <NavBar title="我的订单" />

    <view class="px-4 mt-4 relative z-10 flex flex-col gap-4">
      <!-- Tabs -->
      <view class="flex gap-2 flex-wrap mb-2 pl-2">
        <button
          v-for="tab in tabs"
          :key="tab.key || 'all'"
          class="px-4 py-1.5 border-0 rounded-full text-[13px] font-medium transition-smooth shadow-sm cursor-pointer"
          :class="activeStatus === tab.key ? 'bg-white text-primary hover:-translate-y-0.5 mt-0' : 'bg-white/30 text-white hover:bg-white/40 backdrop-blur-sm'"
          @click="activeStatus = tab.key"
        >
          {{ tab.label }}
        </button>
      </view>

      <text v-if="loading" class="text-center text-text mt-12 block font-medium">Loading...</text>
      <text v-else-if="error" class="text-center text-red-500 mt-12 block">{{ error }}</text>
      <text v-else-if="filtered.length === 0" class="text-center text-gray-400 mt-12 block">暂无订单</text>

      <view v-else class="flex flex-col gap-4 pb-6">
        <view v-for="item in filtered" :key="item.id" class="bg-white rounded-2xl p-5 shadow-sm border border-gray-50 flex flex-col gap-3 group transition-smooth hover:shadow-md hover:-translate-y-1">
          <div class="flex justify-between items-start">
            <text class="text-base font-bold text-text">订单 #{{ item.id }}</text>
            <text class="text-[13px] font-medium px-2 py-0.5 rounded-full bg-gray-50 text-gray-600 border border-gray-100">{{ statusLabel(item.status) }}</text>
          </div>
          
          <text class="text-2xl font-bold text-cta my-1">¥{{ (item.total_cents / 100).toFixed(2) }}</text>
          
          <view class="flex gap-2 justify-end mt-2 pt-3 border-t border-gray-50">
            <button v-if="item.status === 'pending_payment'" class="px-5 py-2 rounded-full text-[13px] font-bold bg-cta text-white shadow-sm hover:opacity-90 transition-smooth border-0 cursor-pointer" @click="goToPay(item.id)">去支付</button>
            <button class="px-4 py-2 rounded-full text-[13px] font-medium bg-gray-50 text-text hover:bg-gray-100 transition-smooth border border-gray-100 cursor-pointer" @click="goToDetail(item.id)">查看详情</button>
            <button class="px-4 py-2 rounded-full text-[13px] font-medium bg-gray-50 text-text hover:bg-gray-100 transition-smooth border border-gray-100 cursor-pointer" @click="goToDetail(item.id, 'refund')">申请退改</button>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>
