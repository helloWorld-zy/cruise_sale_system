<script setup lang="ts">
import PricingRow from '../../components/PricingRow.vue'
const rows = [{ date: '2026-05-01', occupancy: 2, price: 19900 }]
</script>

<template>
  <div class="page">
    <h1>Pricing Matrix</h1>
    <PricingRow v-for="r in rows" :key="r.date" :row="r" />
  </div>
</template>
