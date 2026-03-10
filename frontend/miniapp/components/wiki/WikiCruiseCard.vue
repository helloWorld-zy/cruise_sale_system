<script setup lang="ts">
type WikiCruiseItem = {
  id: number
  name: string
  english_name?: string
  cover_url?: string
}

const props = defineProps<{ item: WikiCruiseItem }>()

const emit = defineEmits<{
  (e: 'select', id: number): void
}>()

function handleClick() {
  emit('select', props.item.id)
}
</script>

<template>
  <button
    type="button"
    class="wiki-card w-full overflow-hidden rounded-[12px] border border-slate-200 bg-white text-left shadow-[0_4px_14px_rgba(15,23,42,0.08)] transition-smooth hover:border-sky-300"
    @click="handleClick"
  >
    <div class="aspect-[1.2] w-full overflow-hidden bg-slate-100">
      <img
        class="h-full w-full object-cover"
        :src="item.cover_url || `https://picsum.photos/seed/wiki-cruise-${item.id}/800/600`"
        :alt="item.name"
      />
    </div>
    <div class="px-3 py-2.5">
      <div class="text-[16px] font-bold leading-tight text-slate-900">{{ item.name }}</div>
      <div v-if="item.english_name" class="mt-1 text-[12px] italic leading-5 text-slate-500">
        {{ item.english_name }}
      </div>
    </div>
  </button>
</template>