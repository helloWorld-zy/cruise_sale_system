<script setup lang="ts">
const route = useRoute()
const raw = route.params.slug
const list = Array.isArray(raw) ? raw : raw ? [raw] : []
const target = `/${list.join('/')}`
await navigateTo(target === '/' ? '/login' : target, { replace: true })
</script>

<template>
  <div />
</template>