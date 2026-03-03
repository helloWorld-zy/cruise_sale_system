<!-- admin/app/layouts/default.vue — 管理后台默认布局 -->
<!-- 包含顶部导航栏、左侧菜单和主内容区域 -->
<template>
  <div class="admin-shell">
    <!-- 顶部导航栏 -->
    <header class="admin-header">
      <div class="admin-brand">CruiseBooking Admin</div>
    </header>
    <div class="admin-body">
      <!-- 左侧菜单 -->
      <aside class="admin-sidebar">
        <nav class="admin-nav">
          <NuxtLink class="admin-link" to="/dashboard">数据看板</NuxtLink>
          <NuxtLink class="admin-link" to="/cruises">邮轮管理</NuxtLink>
          <NuxtLink class="admin-link" to="/routes">航线管理</NuxtLink>
          <NuxtLink class="admin-link" to="/voyages">航次管理</NuxtLink>
          <NuxtLink class="admin-link" to="/cabin-types">舱型管理</NuxtLink>
          <NuxtLink class="admin-link" to="/cabins">舱房管理</NuxtLink>
          <NuxtLink class="admin-link" to="/cabins/alerts">库存预警</NuxtLink>
          <NuxtLink class="admin-link" to="/facility-categories">设施分类</NuxtLink>
          <NuxtLink class="admin-link" to="/facilities">设施管理</NuxtLink>
          <NuxtLink class="admin-link" to="/bookings">订单管理</NuxtLink>
          <NuxtLink class="admin-link" to="/finance">财务统计</NuxtLink>
          <NuxtLink class="admin-link" to="/staff">员工管理</NuxtLink>
          <NuxtLink class="admin-link" to="/settings/shop">店铺设置</NuxtLink>
          <NuxtLink class="admin-link" to="/notifications/templates">通知模板</NuxtLink>
        </nav>
      </aside>
      <!-- 主内容区域：通过 slot 渲染当前页面 -->
      <main class="admin-main">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
/* 全局组件 NuxtLink 由 Nuxt 框架自动注册 */
</script>
