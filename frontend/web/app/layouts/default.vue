<!-- web/app/layouts/default.vue — 前台默认布局 -->
<!-- 包含顶部导航栏、主内容区域和底部版权信息 -->
<template>
  <div class="min-h-screen bg-[var(--color-background)]">
    <header class="border-b border-[var(--web-border)] bg-white/85 backdrop-blur sticky top-0 z-50 shadow-sm">
      <div class="mx-auto flex h-auto max-w-6xl flex-col gap-3 px-4 py-3 md:h-16 md:flex-row md:items-center md:justify-between md:py-0">
        <NuxtLink to="/search" class="flex items-center gap-3 group">
          <span class="h-9 w-9 flex items-center justify-center rounded-full bg-[var(--color-primary)] text-white shadow-md transition group-hover:bg-[var(--color-cta)]">⚓</span>
          <span class="heading-font text-xl tracking-wide text-[var(--color-primary)] font-bold">Azure Deck</span>
          <span class="lux-badge">Cruise Club</span>
        </NuxtLink>
        <nav class="flex flex-wrap items-center gap-2 text-sm md:gap-3 text-[var(--color-primary)] font-semibold">
          <NuxtLink to="/search" class="hover:text-[var(--color-cta)] transition-colors">搜索</NuxtLink>
          <NuxtLink to="/cruises" class="hover:text-[var(--color-cta)] transition-colors">邮轮</NuxtLink>
          <NuxtLink to="/cabins" class="hover:text-[var(--color-cta)] transition-colors">舱位</NuxtLink>
          <NuxtLink to="/booking" class="hover:text-[var(--color-cta)] transition-colors">预订</NuxtLink>
          <NuxtLink to="/orders" class="hover:text-[var(--color-cta)] transition-colors">订单</NuxtLink>
          <NuxtLink to="/account/login" class="btn-primary ml-2 py-1.5 px-4 rounded-full shadow-md text-white text-sm">登录</NuxtLink>
        </nav>
      </div>
    </header>

    <main class="mx-auto max-w-6xl px-4 py-8">
      <slot />
    </main>

    <footer class="border-t border-[var(--web-border)] bg-white/70 py-8 text-center text-sm text-[var(--web-muted)]">
      © 2026 Azure Deck · Crafted for modern cruise journeys
    </footer>
  </div>
</template>

<script setup lang="ts">
/* 全局组件 NuxtLink 由 Nuxt 框架自动注册 */
</script>
