<template>
  <div class="admin-shell" :class="{ 'admin-shell--collapsed': collapsed }">
    <aside class="admin-sidebar">
      <div class="admin-sidebar__brand" :title="collapsed ? 'Cruise Admin' : ''">
        <span class="admin-sidebar__logo">CB</span>
        <span v-if="!collapsed" class="admin-sidebar__title">Cruise Admin</span>
      </div>
      <nav class="admin-nav">
        <NuxtLink
          v-for="item in menuItems"
          :key="item.to"
          class="admin-link"
          :to="item.to"
          :title="collapsed ? item.label : ''"
        >
          <span class="admin-link__dot" />
          <span v-if="!collapsed">{{ item.label }}</span>
        </NuxtLink>
      </nav>
    </aside>

    <div class="admin-body">
      <header class="admin-header">
        <button type="button" class="admin-collapse-btn" @click="toggleSidebar">
          {{ collapsed ? '>>' : '<<' }}
        </button>
        <div class="admin-breadcrumb" aria-label="breadcrumb">
          <span class="admin-breadcrumb__root">Admin</span>
          <span class="admin-breadcrumb__sep">/</span>
          <span class="admin-breadcrumb__current">{{ currentTitle }}</span>
        </div>
        <div class="admin-header__actions">
          <button type="button" class="admin-header__action" @click="refreshPage">刷新</button>
        </div>
      </header>

      <main class="admin-main">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

const route = useRoute()
const collapsed = ref(false)

const menuItems = [
  { to: '/dashboard', label: '数据看板' },
  { to: '/companies', label: '邮轮公司管理' },
  { to: '/cruises', label: '邮轮管理' },
  { to: '/voyages', label: '航次管理' },
  { to: '/cabin-types', label: '舱型管理' },
  { to: '/cabin-types/pricing', label: '舱型价格管理' },
  { to: '/facility-categories', label: '设施分类' },
  { to: '/facilities', label: '设施管理' },
  { to: '/destinations', label: '港口城市词典' },
  { to: '/bookings', label: '订单管理' },
  { to: '/content-templates', label: '文案模板' },
  { to: '/finance', label: '财务统计' },
  { to: '/staff', label: '员工管理' },
  { to: '/settings/shop', label: '店铺设置' },
  { to: '/notifications/templates', label: '通知模板' },
]

const titleMap = new Map(menuItems.map((item) => [item.to, item.label]))

const currentTitle = computed(() => {
  const path = route.path
  if (titleMap.has(path)) return titleMap.get(path)
  for (const item of menuItems) {
    if (path.startsWith(`${item.to}/`)) return item.label
  }
  return '控制台'
})

function toggleSidebar() {
  collapsed.value = !collapsed.value
  const ls = storage()
  if (ls) {
    ls.setItem('admin_sidebar_collapsed', collapsed.value ? '1' : '0')
  }
}

function refreshPage() {
  if (typeof window !== 'undefined') {
    window.location.reload()
  }
}

if (typeof window !== 'undefined') {
  const ls = storage()
  collapsed.value = !!ls && ls.getItem('admin_sidebar_collapsed') === '1'
}

function storage() {
  if (typeof window === 'undefined') return null
  const ls = (window as any).localStorage
  if (!ls) return null
  if (typeof ls.getItem !== 'function' || typeof ls.setItem !== 'function') return null
  return ls as Storage
}
</script>
