<script setup lang="ts">
withDefaults(
  defineProps<{
    name:
      | 'included-dining'
      | 'specialty-dining'
      | 'bar-lounge'
      | 'leisure-entertainment'
      | 'kids-family'
      | 'suite-privilege'
      | 'sports-fitness'
      | 'other'
  }>(),
  {
    name: 'other',
  },
)
</script>

<template>
  <span class="cruise-facility-icon" :data-icon="name" aria-hidden="true">
    <svg v-if="name === 'included-dining'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
      <path d="M4 3v7a2 2 0 0 0 2 2h1v9" />
      <path d="M8 3v7" />
      <path d="M12 3v7" />
      <path d="M18 3v18" />
      <path d="M18 3c2 2 2 5 0 7" />
    </svg>
    <svg v-else-if="name === 'specialty-dining'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
      <path d="M5 6h14" />
      <path d="M7 6l1.5 10h7L17 6" />
      <path d="M9 10h6" />
      <path d="M9 19h6" />
      <path d="M12 3v3" />
    </svg>
    <svg v-else-if="name === 'bar-lounge'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
      <path d="M5 4h14l-5.5 7v5l-3 2v-7L5 4Z" />
      <path d="M9 20h6" />
    </svg>
    <svg v-else-if="name === 'leisure-entertainment'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
      <rect x="4" y="5" width="16" height="14" rx="3" />
      <path d="M9 10v4" />
      <path d="M7 12h4" />
      <circle cx="15.5" cy="10.5" r="1" />
      <circle cx="17.5" cy="13.5" r="1" />
    </svg>
    <svg v-else-if="name === 'kids-family'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
      <circle cx="9" cy="8" r="2.2" />
      <circle cx="15.5" cy="9" r="1.8" />
      <path d="M6.5 18v-2.5A3.5 3.5 0 0 1 10 12h0a3.5 3.5 0 0 1 3.5 3.5V18" />
      <path d="M14 18v-1.5a2.5 2.5 0 0 1 2.5-2.5h0A2.5 2.5 0 0 1 19 16.5V18" />
    </svg>
    <svg v-else-if="name === 'suite-privilege'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
      <path d="M6 18v-6a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v6" />
      <path d="M4 18h16" />
      <path d="M9 10V7a3 3 0 0 1 6 0v3" />
      <path d="m12 3 1 2 2.2.3-1.6 1.6.4 2.1L12 8l-2 1 .4-2.1-1.6-1.6L11 5l1-2Z" />
    </svg>
    <svg v-else-if="name === 'sports-fitness'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
      <path d="M3 10v4" />
      <path d="M6 8v8" />
      <path d="M10 9v6" />
      <path d="M14 9v6" />
      <path d="M18 8v8" />
      <path d="M21 10v4" />
      <path d="M6 12h12" />
    </svg>
    <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
      <circle cx="12" cy="12" r="8" />
      <path d="M12 8v4" />
      <path d="M12 16h.01" />
    </svg>
  </span>
</template>

<style scoped>
.cruise-facility-icon {
  width: 22px;
  height: 22px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex: 0 0 auto;
}

.cruise-facility-icon svg {
  width: 100%;
  height: 100%;
}
</style>