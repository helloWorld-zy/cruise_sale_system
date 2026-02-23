<script setup lang="ts">
import { useRoute } from 'vue-router'
const route = useRoute()
</script>

<template>
  <div class="page">
    <h1>Booking Detail {{ route.params.id }}</h1>
  </div>
</template>
