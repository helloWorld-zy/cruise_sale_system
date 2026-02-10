<script setup lang="ts">
const columns = [{ key: 'name_en', label: 'Name (EN)' }, { key: 'status', label: 'Status' }]
const { data: cruises } = await useFetch('/api/v1/cruises', { baseURL: 'http://localhost:8080' })
</script>

<template>
  <div class="p-4">
    <div class="flex justify-between mb-4">
      <h1 class="text-2xl font-bold">Cruise Management</h1>
      <UButton>Add Cruise</UButton>
    </div>
    <UTable :rows="cruises?.data || []" :columns="columns" />
  </div>
</template>
