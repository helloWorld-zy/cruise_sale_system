<script setup lang="ts">
defineProps<{
  title?: string
}>()
</script>

<template>
  <section class="admin-form-card">
    <header v-if="title" class="admin-form-card__header">{{ title }}</header>
    <div class="admin-form-card__body">
      <slot />
    </div>
  </section>
</template>
