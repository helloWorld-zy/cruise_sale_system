<script setup lang="ts">
import { ref } from 'vue'
import LoginPage from '../pages/login/login.vue'
import BookingPage from '../pages/booking/create.vue'
import CabinListPage from '../pages/cabin/list.vue'
import CabinDetailPage from '../pages/cabin/detail.vue'
import PayPage from '../pages/pay/pay.vue'

type TabKey = 'login' | 'booking' | 'cabin-list' | 'cabin-detail' | 'pay'

const activeTab = ref<TabKey>('login')

const tabs: Array<{ key: TabKey; label: string }> = [
  { key: 'login', label: '登录' },
  { key: 'booking', label: '下单' },
  { key: 'cabin-list', label: '舱房列表' },
  { key: 'cabin-detail', label: '舱房详情' },
  { key: 'pay', label: '支付' },
]
</script>

<template>
  <div class="preview-shell">
    <header class="preview-header">
      <h1>Miniapp UI Preview</h1>
      <p>用于快速验证小程序端页面样式</p>
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
      <CabinListPage v-else-if="activeTab === 'cabin-list'" />
      <CabinDetailPage v-else-if="activeTab === 'cabin-detail'" :cabin-sku-id="1" />
      <PayPage v-else-if="activeTab === 'pay'" :booking-id="1" />
    </main>
  </div>
</template>
