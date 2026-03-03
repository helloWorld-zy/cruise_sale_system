<!-- web/app/layouts/default.vue — 前台默认布局 -->
<!-- 包含顶部导航栏、主内容区域和底部版权信息 -->
<template>
  <div class="min-h-screen">
    <header class="border-b bg-white/85 backdrop-blur">
      <div class="mx-auto flex h-auto max-w-6xl flex-col gap-3 px-4 py-3 md:h-16 md:flex-row md:items-center md:justify-between md:py-0">
        <NuxtLink to="/search" class="flex items-center gap-3">
          <span class="h-9 w-9 rounded-full bg-[#113d5c] text-center text-lg leading-9 text-white">⚓</span>
          <span class="font-['Bodoni_MT','Times_New_Roman',serif] text-xl tracking-wide text-[#12263a]">Azure Deck</span>
          <span class="lux-badge">Cruise Club</span>
        </NuxtLink>
        <nav class="flex flex-wrap items-center gap-2 text-sm md:gap-3">
          <NuxtLink to="/search">搜索</NuxtLink>
          <NuxtLink to="/cruises">邮轮</NuxtLink>
          <NuxtLink to="/cabins">舱位</NuxtLink>
          <NuxtLink to="/booking">预订</NuxtLink>
          <NuxtLink to="/orders">订单</NuxtLink>
          <NuxtLink to="/account/login" class="rounded-full bg-[#113d5c] px-3 py-1.5 text-white">登录</NuxtLink>
        </nav>
      </div>
    </header>

    <main class="mx-auto max-w-6xl px-4 py-8">
      <slot />
    </main>

    <footer class="border-t border-[#d7e2ea] bg-white/70 py-8 text-center text-sm text-[#56708a]">
      © 2026 Azure Deck · Crafted for modern cruise journeys
    </footer>
  </div>
</template>

<script setup lang="ts">
/* 全局组件 NuxtLink 由 Nuxt 框架自动注册 */
</script>
