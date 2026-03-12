<script setup lang="ts">
type CompanyItem = {
  id: number
  name: string
}

type CruiseItem = {
  id: number
  name: string
}

type CompanyTreeItem = {
  company: CompanyItem
  cruises: CruiseItem[]
}

defineProps<{
  items: CompanyTreeItem[]
  expandedCompanyIds: number[]
  selectedCruiseId: number | null
  loading?: boolean
}>()

const emit = defineEmits<{
  (e: 'toggle-company', companyId: number): void
  (e: 'select-cruise', cruiseId: number): void
}>()

function isExpanded(companyId: number, expandedCompanyIds: number[]) {
  return expandedCompanyIds.includes(companyId)
}
</script>

<template>
  <aside class="flex h-full w-[108px] shrink-0 flex-col border-r border-slate-200 bg-[#f5f5f5]">
    <div v-if="loading" class="px-3 py-4 text-center text-[12px] text-slate-500">加载中</div>
    <div v-for="item in items" :key="item.company.id" class="border-b border-slate-200">
      <button
        type="button"
        class="flex min-h-[62px] w-full items-center justify-between border-0 bg-transparent px-3 py-3 text-left text-[13px] text-slate-700 transition-smooth hover:bg-slate-100"
        @click="emit('toggle-company', item.company.id)"
      >
        <span class="line-clamp-2 flex-1 font-medium leading-5">{{ item.company.name }}</span>
        <span class="ml-2 text-[12px] text-slate-400">{{ isExpanded(item.company.id, expandedCompanyIds) ? '−' : '+' }}</span>
      </button>

      <div v-if="isExpanded(item.company.id, expandedCompanyIds)" class="bg-[#efefef] px-2 py-1.5">
        <button
          v-for="cruise in item.cruises"
          :key="cruise.id"
          type="button"
          class="mb-1 flex min-h-[40px] w-full items-center rounded-xl border-0 px-3 py-2 text-left text-[12px] leading-4 transition-smooth"
          :class="selectedCruiseId === cruise.id ? 'bg-white font-semibold text-[#0f5ba9] shadow-sm' : 'bg-transparent text-slate-500 hover:bg-white/70'"
          @click="emit('select-cruise', cruise.id)"
        >
          {{ cruise.name }}
        </button>
      </div>
    </div>
  </aside>
</template>