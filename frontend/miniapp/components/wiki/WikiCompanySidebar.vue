<script setup lang="ts">
type SidebarItem = {
  id: number | 'all'
  name: string
}

defineProps<{
  items: SidebarItem[]
  selectedCompanyId: number | 'all'
  loading?: boolean
}>()

const emit = defineEmits<{
  (e: 'select', id: number | 'all'): void
}>()
</script>

<template>
  <aside class="wiki-sidebar flex h-full w-[104px] shrink-0 flex-col border-r border-slate-200 bg-[#f4f4f4]">
    <div v-if="loading" class="px-3 py-4 text-center text-[12px] text-slate-500">加载中</div>
    <button
      v-for="item in items"
      :key="item.id"
      type="button"
      class="relative min-h-[66px] border-0 border-b border-slate-200 bg-transparent px-2 py-4 text-center text-[13px] leading-5 text-slate-600 transition-smooth"
      :class="selectedCompanyId === item.id ? 'bg-white font-semibold text-[#0f5ba9]' : 'hover:bg-slate-100'"
      @click="emit('select', item.id)"
    >
      <span
        v-if="selectedCompanyId === item.id"
        class="absolute left-0 top-1/2 h-9 w-1 -translate-y-1/2 rounded-r-full bg-[#0f5ba9]"
      ></span>
      {{ item.name }}
    </button>
  </aside>
</template>