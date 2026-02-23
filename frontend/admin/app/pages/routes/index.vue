<script setup lang="ts">
import RouteTable from '../../components/RouteTable.vue'
const items = [{ id: 1, code: 'R1', name: 'Route 1' }]
</script>

<template>
  <div class="page">
    <h1>Routes</h1>
    <RouteTable :items="items" />
  </div>
</template>
