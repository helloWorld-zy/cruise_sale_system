<script setup lang="ts">
const columns = [
  { key: 'order_no', label: 'Order No' },
  { key: 'total_amount', label: 'Amount' },
  { key: 'status', label: 'Status' },
  { key: 'created_at', label: 'Date' }
]
const { data: orders } = await useFetch('/api/v1/admin/orders', { baseURL: 'http://localhost:8080' })
</script>

<template>
  <div class="p-4">
    <h1 class="text-2xl font-bold mb-4">Order Management</h1>
    <UTable :rows="orders?.data || []" :columns="columns" />
  </div>
</template>
