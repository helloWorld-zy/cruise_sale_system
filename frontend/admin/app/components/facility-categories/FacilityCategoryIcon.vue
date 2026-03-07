<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    name?: string
  }>(),
  {
    name: '',
  },
)
</script>

<template>
  <span class="facility-category-icon" :data-icon="props.name || 'unknown'" aria-hidden="true">
    <svg v-if="props.name === 'utensils'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <path d="M4 3v7a2 2 0 0 0 2 2h1v9" />
      <path d="M8 3v7" />
      <path d="M12 3v7" />
      <path d="M18 3v18" />
      <path d="M18 3c2 2 2 5 0 7" />
    </svg>
    <svg v-else-if="props.name === 'music'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <path d="M10 18V6l9-2v12" />
      <circle cx="7" cy="18" r="3" />
      <circle cx="19" cy="16" r="3" />
    </svg>
    <svg v-else-if="props.name === 'spa'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <path d="M12 20c4 0 7-2.4 7-5.5 0-1.5-.8-2.8-2.1-3.8" />
      <path d="M12 20c-4 0-7-2.4-7-5.5 0-1.5.8-2.8 2.1-3.8" />
      <path d="M12 11c2.4-1.3 3.9-3.5 3.9-6" />
      <path d="M12 11C9.6 9.7 8.1 7.5 8.1 5" />
    </svg>
    <svg v-else-if="props.name === 'dumbbell'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <path d="M3 10v4" />
      <path d="M6 8v8" />
      <path d="M10 9v6" />
      <path d="M14 9v6" />
      <path d="M18 8v8" />
      <path d="M21 10v4" />
      <path d="M6 12h12" />
    </svg>
    <svg v-else-if="props.name === 'swimmer'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <circle cx="8" cy="7" r="2.2" />
      <path d="M10 10l3 2 2-1" />
      <path d="M5 15c1.4 1 2.7 1.5 4 1.5s2.6-.5 4-1.5 2.7-1.5 4-1.5 2.6.5 4 1.5" />
      <path d="M4 19c1.3.8 2.7 1.2 4 1.2s2.7-.4 4-1.2 2.7-1.2 4-1.2 2.7.4 4 1.2" />
    </svg>
    <svg v-else-if="props.name === 'gamepad'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <path d="M7 9h10a4 4 0 0 1 3.9 4.8l-.8 3.3a2 2 0 0 1-3.2 1.1l-2.4-2.1H9.5l-2.4 2.1a2 2 0 0 1-3.2-1.1l-.8-3.3A4 4 0 0 1 7 9" />
      <path d="M8 12v4" />
      <path d="M6 14h4" />
      <path d="M16 13h.01" />
      <path d="M18 15h.01" />
    </svg>
    <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <rect x="4" y="4" width="16" height="16" rx="4" />
      <path d="M8 12h8" />
      <path d="M12 8v8" />
    </svg>
  </span>
</template>

<style scoped>
.facility-category-icon {
  width: 22px;
  height: 22px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: currentColor;
  flex: 0 0 auto;
}

.facility-category-icon svg {
  width: 100%;
  height: 100%;
}
</style>