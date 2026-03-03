<template>
  <div class="magazine-page min-h-screen bg-[#f8f4ed] text-slate-900">
    <section class="hero relative overflow-hidden">
      <div class="hero-overlay absolute inset-0" />
      <div class="relative mx-auto flex h-[50vh] max-w-6xl items-center justify-center px-6 text-center">
        <div>
          <p class="mb-3 text-xs uppercase tracking-[0.35em] text-[#f7e7c8]">Blue Horizon</p>
          <h1 class="font-['Playfair_Display','Georgia',serif] text-4xl font-semibold text-white md:text-5xl">全球邮轮精选航季</h1>
        </div>
      </div>
    </section>

    <div class="mx-auto -mt-8 max-w-5xl rounded-2xl border border-[#e8ddca] bg-white/95 p-4 shadow-xl backdrop-blur" data-test="search-card">
      <div class="flex flex-wrap gap-3">
        <input v-model="keyword" placeholder="搜索邮轮名称/代码" class="h-11 min-w-60 flex-1 rounded-xl border border-[#eadfcb] bg-[#fffdf9] px-4 text-sm outline-none ring-[#c9a96e] focus:ring-2" @keyup.enter="loadItems" />
        <button type="button" class="h-11 rounded-xl bg-[#0f3d5c] px-5 text-sm font-medium text-white hover:bg-[#12496d]" @click="loadItems">搜索</button>
      </div>
    </div>

    <main class="mx-auto max-w-6xl px-6 pb-12 pt-10">
      <p v-if="loading" class="text-sm text-slate-500">加载中...</p>
      <p v-else-if="error" class="text-sm text-rose-500">{{ error }}</p>
      <div v-else class="grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-3" data-test="cruise-grid">
        <article
          v-for="(item, idx) in items"
          :key="item.id"
          class="group overflow-hidden rounded-2xl border border-[#eadfcb] bg-white shadow-sm transition duration-300 hover:-translate-y-1 hover:shadow-xl"
          :style="{ animationDelay: `${idx * 90}ms` }"
        >
          <div class="relative aspect-[4/3] overflow-hidden bg-slate-200">
            <img :src="coverFor(item)" :alt="item.name || 'cruise cover'" class="h-full w-full object-cover transition duration-500 group-hover:scale-105" />
          </div>
          <div class="space-y-3 p-4">
            <h2 class="font-['Playfair_Display','Georgia',serif] text-2xl text-[#12263a]">{{ item.name || '-' }}</h2>
            <div class="flex flex-wrap gap-2 text-xs text-slate-600">
              <span class="rounded-full bg-[#f5efe4] px-2.5 py-1">⚓ {{ item.tonnage || '-' }} 吨</span>
              <span class="rounded-full bg-[#f5efe4] px-2.5 py-1">👤 {{ item.passenger_capacity || '-' }} 人</span>
              <span class="rounded-full bg-[#f5efe4] px-2.5 py-1">🧭 {{ item.length || '-' }} m</span>
            </div>
            <NuxtLink :to="`/cruises/${item.id}`" class="inline-block text-sm font-medium text-[#c9a96e] hover:text-[#b48f4f]">查看详情 ></NuxtLink>
          </div>
        </article>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'

const { request } = useApi()
const loading = ref(false)
const error = ref<string | null>(null)
const keyword = ref('')
const items = ref<Record<string, any>[]>([])

async function loadItems() {
  loading.value = true
  error.value = null
  try {
    const query: Record<string, any> = { page: 1, page_size: 30 }
    if (keyword.value.trim()) query.keyword = keyword.value.trim()
    const res = await request('/cruises', { query })
    const payload = res?.data ?? res ?? {}
    items.value = Array.isArray(payload) ? payload : payload?.list ?? []
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load cruises'
  } finally {
    loading.value = false
  }
}

function coverFor(item: Record<string, any>) {
  if (item.cover_url) return item.cover_url
  return `https://picsum.photos/seed/cruise-${item.id || 0}/960/720`
}

onMounted(loadItems)
</script>

<style scoped>
.hero {
  background: radial-gradient(circle at 20% 20%, rgba(48, 91, 123, 0.65), transparent 40%),
    radial-gradient(circle at 80% 30%, rgba(201, 169, 110, 0.4), transparent 42%),
    linear-gradient(120deg, #0f3d5c 0%, #155178 46%, #1f6a8f 100%);
}
.hero-overlay {
  background: linear-gradient(180deg, rgba(7, 24, 36, 0.3), rgba(7, 24, 36, 0.55));
}
article {
  animation: riseIn 420ms ease both;
}
@keyframes riseIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
