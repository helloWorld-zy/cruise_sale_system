<script setup lang="ts">
import CabinCard from '../../components/CabinCard.vue'
const results = [{ id: 1, name: 'Balcony Cabin', price: 19900 }]
</script>

<template>
  <div class="page">
    <h1>Search</h1>
    <div>Filter</div>
    <CabinCard v-for="r in results" :key="r.id" :cabin="r" />
  </div>
</template>
