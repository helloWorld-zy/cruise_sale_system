<!-- web/app/layouts/default.vue — 前台默认布局 -->
<!-- 包含顶部导航栏、主内容区域和底部版权信息 -->
<template>
  <div class="min-h-screen bg-white">
    <!-- 顶部导航栏：磨砂半透明效果 -->
    <header class="border-b bg-white/90 backdrop-blur">
      <div class="max-w-6xl mx-auto px-4 h-16 flex items-center justify-between">
        <!-- 网站 Logo / 首页链接 -->
        <NuxtLink to="/search" class="text-lg font-semibold">CruiseBooking</NuxtLink>
        <!-- 导航菜单 -->
        <nav class="space-x-4 text-sm">
          <NuxtLink to="/search">搜索</NuxtLink>
          <NuxtLink to="/booking">预订</NuxtLink>
          <NuxtLink to="/account/login">登录</NuxtLink>
        </nav>
      </div>
    </header>
    <!-- 主内容区域 -->
    <main class="max-w-6xl mx-auto px-4 py-8">
      <slot />
    </main>
    <!-- 底部版权 -->
    <footer class="border-t py-6 text-center text-sm text-gray-500">
      © 2026 CruiseBooking
    </footer>
  </div>
</template>

<script setup lang="ts">
/* 全局组件 NuxtLink 由 Nuxt 框架自动注册 */
</script>
