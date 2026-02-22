<template>
  <div class="min-h-screen bg-white">
    <header class="border-b bg-white/90 backdrop-blur">
      <div class="max-w-6xl mx-auto px-4 h-16 flex items-center justify-between">
        <NuxtLink to="/" class="text-lg font-semibold">CruiseBooking</NuxtLink>
        <nav class="space-x-4 text-sm">
          <NuxtLink to="/cruises">邮轮</NuxtLink>
        </nav>
      </div>
    </header>
    <main class="max-w-6xl mx-auto px-4 py-8">
      <slot />
    </main>
    <footer class="border-t py-6 text-center text-sm text-gray-500">
      © 2026 CruiseBooking
    </footer>
  </div>
</template>

<script setup lang="ts">
/* global NuxtLink */
</script>
