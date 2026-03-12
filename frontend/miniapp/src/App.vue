<script setup lang="ts">
import { ref } from 'vue'
import { Home, BookOpen, Search, ShoppingCart, User } from 'lucide-vue-next'
import HomePage from '../pages/home/index.vue'
import WikiHomePage from '../pages/wiki/index.vue'
import LoginPage from '../pages/login/login.vue'
import BookingPage from '../pages/booking/create.vue'
import ProductsPage from '../pages/products/index.vue'
import CabinDetailPage from '../pages/cabin/detail.vue'
import CruiseDetailPage from '../pages/cruise/detail.vue'
import VoyageDetailPage from '../pages/voyage/detail.vue'
import PayPage from '../pages/pay/pay.vue'
import OrdersPage from '../pages/orders/list.vue'

type TabKey = 'home' | 'wiki' | 'cabin-list' | 'cart' | 'my' | 'cabin-detail' | 'cruise-detail' | 'voyage-detail' | 'pay' | 'orders' | 'login' | 'booking'

const activeTab = ref<TabKey>('home')
const previousTab = ref<TabKey>('home')
const selectedCruiseId = ref<number | null>(null)
const selectedVoyageId = ref<number | null>(null)
const selectedBookingId = ref<number | null>(null)

/** 导航到子页面，记住来源以便返回 */
function navigateTo(tab: TabKey) {
  previousTab.value = activeTab.value
  activeTab.value = tab
}

/** 返回上一个页面 */
function goBack() {
  activeTab.value = previousTab.value
}

function openCruiseDetail(cruiseId: number) {
  selectedCruiseId.value = cruiseId
  navigateTo('cruise-detail')
}

function openVoyageDetail(voyageId: number) {
  selectedVoyageId.value = voyageId
  navigateTo('voyage-detail')
}

function openPay(bookingId: number) {
  selectedBookingId.value = bookingId
  navigateTo('pay')
}

const tabs: Array<{ key: TabKey; label: string; icon: any }> = [
  { key: 'home', label: '首页', icon: Home },
  { key: 'wiki', label: '邮轮百科', icon: BookOpen },
  { key: 'cabin-list', label: '全部商品', icon: Search },
  { key: 'cart', label: '购物车', icon: ShoppingCart },
  { key: 'my', label: '我的', icon: User },
]
</script>

<template>
  <div class="mobile-app-container relative">
    <main class="min-h-screen bg-background pb-[80px]">
      <HomePage v-if="activeTab === 'home'" />
      <WikiHomePage v-else-if="activeTab === 'wiki'" @open-cruise="openCruiseDetail" />
      <OrdersPage v-else-if="activeTab === 'my'" @open-pay="openPay" />
      <BookingPage v-else-if="activeTab === 'cart'" />

      <ProductsPage v-else-if="activeTab === 'cabin-list'" @open-voyage="openVoyageDetail" />
      <CabinDetailPage v-else-if="activeTab === 'cabin-detail'" :cabin-sku-id="1" :preview="true" @back="goBack" />
      <CruiseDetailPage v-else-if="activeTab === 'cruise-detail'" :cruise-id="selectedCruiseId || 1" @back="goBack" @open-voyage="openVoyageDetail" />
      <VoyageDetailPage v-else-if="activeTab === 'voyage-detail'" :voyage-id="selectedVoyageId || 101" @back="goBack" />
      <PayPage v-else-if="activeTab === 'pay'" :booking-id="selectedBookingId || 1" :preview="!selectedBookingId" @back="goBack" />
      <LoginPage v-else-if="activeTab === 'login'" @back="goBack" />
      <OrdersPage v-else-if="activeTab === 'orders'" @open-pay="openPay" />
    </main>

    <!-- Bottom Tab Bar -->
    <nav class="fixed bottom-0 w-full max-w-[480px] bg-white/95 backdrop-blur-md border-t border-gray-100 flex justify-around items-center px-2 py-2 shadow-[0_-4px_20px_rgba(0,0,0,0.03)] z-40" style="left: 50%; transform: translateX(-50%);">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        class="flex flex-col items-center justify-center gap-1.5 min-w-[64px] h-[50px] transition-smooth group bg-transparent border-0"
        @click="activeTab = tab.key"
      >
        <component 
          :is="tab.icon" 
          class="w-6 h-6 transition-smooth"
          :class="activeTab === tab.key ? 'text-cta scale-110 drop-shadow-sm' : 'text-gray-400 group-hover:text-primary group-hover:-translate-y-0.5'"
        />
        <span 
          class="text-[10px] font-medium transition-smooth"
          :class="activeTab === tab.key ? 'text-cta' : 'text-gray-500 group-hover:text-primary'"
        >
          {{ tab.label }}
        </span>
      </button>
    </nav>
  </div>
</template>
