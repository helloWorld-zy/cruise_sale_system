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

    <div class="search-card mx-auto -mt-8 max-w-5xl rounded-2xl p-6 relative z-10" data-test="search-card">
      <div class="flex flex-wrap gap-4">
        <input v-model="keyword" placeholder="搜索邮轮名称/代码" class="h-12 min-w-60 flex-1 rounded-sm border border-[#eadfcb] bg-[#fffdf9] px-4 text-sm outline-none ring-[#c9a96e] focus:ring-1 transition-all" @keyup.enter="loadItems" />
        <button type="button" class="h-12 rounded-sm bg-[#0f3d5c] px-8 text-sm tracking-widest text-[#fcfbf9] font-serif hover:bg-[#12496d] transition-colors" @click="loadItems">搜索</button>
      </div>
    </div>

    <main class="mx-auto max-w-6xl px-6 pb-16 pt-12">
      <p v-if="loading" class="text-sm font-serif italic text-slate-500 text-center py-10">加载中...</p>
      <p v-else-if="error" class="text-sm text-rose-800 text-center py-10">{{ error }}</p>
      <div v-else class="grid grid-cols-1 gap-10 md:grid-cols-2 lg:grid-cols-3" data-test="cruise-grid">
        <article
          v-for="(item, idx) in items"
          :key="item.id"
          class="group overflow-hidden bg-white shadow-sm transition-all duration-500 hover:-translate-y-2 hover:shadow-2xl border border-transparent hover:border-[#eadfcb]"
          :style="{ animationDelay: `${idx * 120}ms` }"
        >
          <div class="relative aspect-[4/3] overflow-hidden bg-slate-100">
            <img :src="coverFor(item)" :alt="item.name || 'cruise cover'" class="h-full w-full object-cover transition-transform duration-700 group-hover:scale-110" />
            <div class="absolute inset-0 bg-gradient-to-t from-black/40 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
          </div>
          <div class="flex flex-col space-y-4 p-6 bg-white relative">
            <h2 class="font-['Playfair_Display','Georgia',serif] text-2xl text-[#0f3d5c]">{{ item.name || '-' }}</h2>
            <div class="flex flex-wrap gap-3 text-[11px] font-semibold tracking-wider text-slate-500">
              <span class="border-b border-[#eadfcb] pb-1">⚓ {{ item.tonnage || '-' }} 吨</span>
              <span class="border-b border-[#eadfcb] pb-1">👤 {{ item.passenger_capacity || '-' }} 人</span>
              <span class="border-b border-[#eadfcb] pb-1">🧭 {{ item.length || '-' }} m</span>
            </div>
            <div class="pt-4 mt-2 border-t border-slate-100">
                <NuxtLink :to="`/cruises/${item.id}`" class="inline-flex items-center text-xs font-bold tracking-widest text-[#c9a96e] hover:text-[#b48f4f] transition-colors">
                  查看详情 <span class="ml-2 text-lg leading-none">></span>
                </NuxtLink>
            </div>
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
.magazine-page {
  font-family: 'Inter', system-ui, sans-serif;
  background-color: #fcfbf9;
}
.hero {
  background: url('https://picsum.photos/seed/luxury-cruise-hero/1920/1080') center/cover no-repeat;
  background-attachment: fixed;
}
.hero-overlay {
  background: linear-gradient(to bottom, rgba(15, 61, 92, 0.4), rgba(7, 24, 36, 0.8));
  backdrop-filter: sepia(10%) hue-rotate(190deg);
}
article {
  animation: riseIn 500ms cubic-bezier(0.2, 0.8, 0.2, 1) both;
}
@keyframes riseIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
.search-card {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(201, 169, 110, 0.3);
  box-shadow: 0 20px 40px rgba(15, 61, 92, 0.08);
}
</style>
