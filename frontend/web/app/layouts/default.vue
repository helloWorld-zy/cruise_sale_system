<!-- web/app/layouts/default.vue — 前台默认布局 -->
<!-- 包含顶部导航栏、主内容区域和底部版权信息 -->
<template>
  <div class="min-h-screen bg-[var(--color-background)] font-sans">
    <header class="fixed top-0 w-full z-50 transition-colors duration-300" :class="isScrolled ? 'bg-white shadow-md' : 'bg-white/95 backdrop-blur-md border-b border-[var(--web-border)]'">
      <div class="mx-auto flex h-auto w-full px-6 py-3 md:h-16 md:flex-row md:items-center md:justify-between md:py-0">
        <NuxtLink to="/search" class="flex items-center gap-3 group">
          <span class="h-9 w-9 flex items-center justify-center rounded bg-[var(--color-primary)] text-white shadow-sm transition group-hover:bg-[var(--color-cta)]">⚓</span>
          <span class="heading-font text-2xl tracking-wide text-[var(--color-primary)] font-bold">Azure Deck</span>
          <span class="text-[0.65rem] tracking-widest uppercase font-bold text-[var(--web-muted)] ml-2 border-l border-slate-300 pl-3">Cruise Line</span>
        </NuxtLink>
        <nav class="flex flex-wrap items-center gap-6 text-sm text-[var(--color-primary)] font-bold tracking-wide uppercase">
          <NuxtLink to="/search" class="hover:text-[var(--color-action)] transition-colors">搜索邮轮</NuxtLink>
          <NuxtLink to="/cruises" class="hover:text-[var(--color-action)] transition-colors">旗下舰队</NuxtLink>
          <NuxtLink to="/cabins" class="hover:text-[var(--color-action)] transition-colors">舱位概览</NuxtLink>
          <NuxtLink to="/booking" class="hover:text-[var(--color-action)] transition-colors">在线预订</NuxtLink>
          <NuxtLink to="/orders" class="hover:text-[var(--color-action)] transition-colors">我的订单</NuxtLink>
          <div class="flex items-center gap-3 ml-4">
            <NuxtLink to="/account/login" class="text-xs hover:text-[var(--color-action)] transition-colors flex items-center gap-1">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" /></svg>
              登录
            </NuxtLink>
            <a href="tel:1-800-432-1234" class="text-xs font-bold flex items-center gap-1 hover:text-[var(--color-action)] transition-colors">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" /></svg>
              1-800-AZURE-DECK
            </a>
          </div>
        </nav>
      </div>
    </header>

    <main class="w-full pt-16">
      <slot />
    </main>

    <footer class="border-t border-[var(--web-border)] bg-[#F1F5F9] py-12 px-6">
      <div class="mx-auto max-w-6xl flex flex-col md:flex-row justify-between items-center gap-6">
        <div class="text-[var(--color-primary)] font-bold heading-font text-2xl flex items-center gap-2">
          ⚓ Azure Deck
        </div>
        <div class="text-xs text-[var(--web-muted)] uppercase tracking-wider">
          © 2026 Azure Deck · 专注于现代邮轮旅行体验
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const isScrolled = ref(false)

function handleScroll() {
  isScrolled.value = window.scrollY > 20
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>
