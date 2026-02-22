<template>
  <div class="min-h-screen bg-[var(--admin-bg)]">
    <header class="h-14 bg-white border-b px-4 flex items-center">
      <div class="font-semibold">CruiseBooking Admin</div>
    </header>
    <div class="flex">
      <aside class="w-64 bg-white border-r min-h-[calc(100vh-56px)]">
        <nav class="p-3 space-y-1">
          <NuxtLink class="block px-3 py-2 rounded hover:bg-gray-50" to="/cruises">邮轮管理</NuxtLink>
          <NuxtLink class="block px-3 py-2 rounded hover:bg-gray-50" to="/cabin-types">舱房类型</NuxtLink>
          <NuxtLink class="block px-3 py-2 rounded hover:bg-gray-50" to="/facilities">设施管理</NuxtLink>
        </nav>
      </aside>
      <main class="flex-1 p-6">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
/* global NuxtLink */
</script>
