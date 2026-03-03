<script setup lang="ts">
import { ref } from 'vue'
import LoginPage from '../pages/login/login.vue'
import BookingPage from '../pages/booking/create.vue'
import CabinListPage from '../pages/cabin/list.vue'
import CabinDetailPage from '../pages/cabin/detail.vue'
import PayPage from '../pages/pay/pay.vue'
import OrdersPage from '../pages/orders/list.vue'

type TabKey = 'login' | 'booking' | 'cabin-list' | 'cabin-detail' | 'pay' | 'orders'

const activeTab = ref<TabKey>('login')

const tabs: Array<{ key: TabKey; label: string }> = [
  { key: 'login', label: '登录' },
  { key: 'booking', label: '下单' },
  { key: 'cabin-list', label: '舱房列表' },
  { key: 'cabin-detail', label: '舱房详情' },
  { key: 'pay', label: '支付' },
  { key: 'orders', label: '订单' },
]
</script>

<template>
  <div class="preview-shell">
    <header class="preview-header">
      <h1>Azure Deck Mini Program</h1>
      <p>邮轮小程序高保真风格预览：强调重点、流程明确、反馈清晰。</p>
    </header>

    <nav class="preview-tabs">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        class="tab-btn"
        :class="{ active: activeTab === tab.key }"
        @click="activeTab = tab.key"
      >
        {{ tab.label }}
      </button>
    </nav>

    <main class="preview-stage">
      <LoginPage v-if="activeTab === 'login'" />
      <BookingPage v-else-if="activeTab === 'booking'" />
      <CabinListPage v-else-if="activeTab === 'cabin-list'" :preview="true" />
      <CabinDetailPage v-else-if="activeTab === 'cabin-detail'" :cabin-sku-id="1" :preview="true" />
      <PayPage v-else-if="activeTab === 'pay'" :booking-id="1" :preview="true" />
      <OrdersPage v-else-if="activeTab === 'orders'" />
    </main>
  </div>
</template>
