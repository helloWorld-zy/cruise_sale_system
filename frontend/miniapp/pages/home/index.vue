<script setup lang="ts">
import { Headset } from 'lucide-vue-next'

const popularRoutes = [
  { id: 1, name: '日韩', image: 'https://picsum.photos/id/1015/100/100' },
  { id: 2, name: '东南亚', image: 'https://picsum.photos/id/1043/100/100' },
  { id: 3, name: '地中海', image: 'https://picsum.photos/id/1018/100/100' },
  { id: 4, name: '三峡', image: 'https://picsum.photos/id/1036/100/100' },
  { id: 5, name: '南北极', image: 'https://picsum.photos/id/1019/100/100' },
  { id: 6, name: '欧洲河轮', image: 'https://picsum.photos/id/1040/100/100' },
  { id: 7, name: '加勒比', image: 'https://picsum.photos/id/1045/100/100' },
  { id: 8, name: '澳新', image: 'https://picsum.photos/id/1016/100/100' },
  { id: 9, name: '中东', image: 'https://picsum.photos/id/1000/100/100' },
  { id: 10, name: '西沙群岛', image: 'https://picsum.photos/id/1011/100/100' },
]

const recentDeals = [
  {
    id: 1,
    ship: '海洋光谱号',
    route: '上海-福冈-长崎-上海',
    date: '2023-09-15 出发',
    image: 'https://picsum.photos/id/1050/300/180',
  },
  {
    id: 2,
    ship: '荣耀号',
    route: '天津-釜山-济州岛-天津',
    date: '2023-10-05 出发',
    image: 'https://picsum.photos/id/1051/300/180',
  },
  {
    id: 3,
    ship: '量子号',
    route: '香港-冲绳-香港',
    date: '2023-11-20 出发',
    image: 'https://picsum.photos/id/1052/300/180',
  }
]
</script>

<template>
  <div class="pb-[20px] bg-background min-h-screen">
    <!-- Header / Nav -->
    <header class="pt-10 pb-3 px-3 bg-white sticky top-0 z-30 flex items-center justify-between shadow-sm">
      <div class="flex-1"></div>
      <h1 class="text-lg font-heading font-bold text-text flex-1 text-center">首页</h1>
      <div class="flex-1 flex justify-end">
        <div class="flex gap-2 bg-gray-50 rounded-full px-3 py-1.5 items-center shadow-sm border border-gray-200 cursor-pointer hover:bg-gray-100 transition-smooth">
          <div class="flex gap-1">
            <div class="w-1.5 h-1.5 rounded-full bg-text"></div>
            <div class="w-1.5 h-1.5 rounded-full bg-text"></div>
            <div class="w-1.5 h-1.5 rounded-full bg-text"></div>
          </div>
          <div class="w-px h-3 bg-gray-300 mx-1"></div>
          <div class="w-4 h-4 rounded-full border-2 border-text relative">
             <div class="w-1.5 h-1.5 bg-text rounded-full absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2"></div>
          </div>
        </div>
      </div>
    </header>

    <!-- Banner -->
    <section class="px-3 mt-3 relative">
      <div class="relative w-full h-36 rounded-xl overflow-hidden shadow-md group cursor-pointer">
        <img src="https://picsum.photos/id/1029/600/300" alt="Cruise Banner" class="w-full h-full object-cover transition-transform duration-700 group-hover:scale-105" />
        <div class="absolute bottom-2.5 left-0 w-full flex justify-center gap-1.5">
          <div class="w-3.5 h-1 rounded-full bg-white transition-smooth"></div>
          <div class="w-1.5 h-1 rounded-full bg-white/50 transition-smooth hover:bg-white"></div>
          <div class="w-1.5 h-1 rounded-full bg-white/50 transition-smooth hover:bg-white"></div>
        </div>
      </div>
    </section>

    <!-- Popular Routes -->
    <section class="mt-5 px-3">
      <h2 class="text-[15px] font-heading font-bold text-text mb-3">热门航线</h2>
      <div class="grid grid-cols-5 gap-y-3.5 gap-x-1">
        <div v-for="route in popularRoutes" :key="route.id" class="flex flex-col items-center gap-1.5 group cursor-pointer">
          <div class="w-11 h-11 rounded-full overflow-hidden shadow-sm transition-smooth group-hover:shadow-md group-hover:-translate-y-0.5 ring-2 ring-transparent group-hover:ring-primary/20">
            <img :src="route.image" :alt="route.name" class="w-full h-full object-cover" />
          </div>
          <span class="text-[11px] text-text font-medium text-center leading-tight">{{ route.name }}</span>
        </div>
      </div>
    </section>

    <!-- Recent Deals -->
    <section class="mt-5 pl-3 mb-3">
      <h2 class="text-[15px] font-heading font-bold text-text mb-3">近期特价</h2>
      <div class="flex overflow-x-auto hide-scrollbar gap-2.5 pb-4 pr-3">
        <div 
          v-for="deal in recentDeals" 
          :key="deal.id" 
          class="min-w-[180px] max-w-[180px] bg-white rounded-xl overflow-hidden shadow-sm flex-shrink-0 transition-smooth hover:shadow-lg hover:-translate-y-1 cursor-pointer border border-gray-100"
        >
          <div class="h-24 w-full overflow-hidden relative">
            <img :src="deal.image" :alt="deal.ship" class="w-full h-full object-cover transition-transform duration-500 hover:scale-105" />
            <div class="absolute top-1.5 right-1.5 bg-cta text-white text-[9px] font-bold px-1.5 py-0.5 rounded-full shadow-sm">特价</div>
          </div>
          <div class="px-2.5 py-2">
            <h3 class="font-bold text-text text-[13px] mb-0.5 truncate">{{ deal.ship }}</h3>
            <p class="text-[11px] text-gray-500 mb-1 truncate">{{ deal.route }}</p>
            <p class="text-[10px] text-primary font-medium">{{ deal.date }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Floating CS Button -->
    <button class="fixed right-4 bottom-[90px] w-12 h-12 bg-white/90 backdrop-blur rounded-full shadow-lg flex items-center justify-center text-primary transition-smooth hover:scale-110 z-20 hover:text-cta border border-gray-100" aria-label="Customer Service">
      <Headset class="w-6 h-6" />
    </button>
  </div>
</template>
