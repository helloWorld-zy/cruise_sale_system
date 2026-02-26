<script setup lang="ts">
import StatCard from '../../components/StatCard.vue'
const summary = { sales: 1000, orders: 12 }
</script>

<template>
  <div class="page">
    <h1>Dashboard</h1>
    <StatCard title="Sales" :value="summary.sales" />
    <StatCard title="Orders" :value="summary.orders" />
  </div>
</template>
