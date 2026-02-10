<script setup lang="ts">
const route = useRoute()
const { data: cruises, pending, error } = await useFetch('/api/v1/cruises', {
  baseURL: 'http://localhost:8080', // In real app, use runtime config
  query: {
    destination: route.query.destination,
    date: route.query.date
  }
})
</script>

<template>
  <div class="container mx-auto p-4">
    <h1 class="text-2xl font-bold mb-4">Find Your Cruise</h1>
    
    <!-- Filters (Simplified) -->
    <div class="flex gap-4 mb-6">
      <UInput placeholder="Destination" />
      <UInput type="date" />
      <UButton>Search</UButton>
    </div>

    <div v-if="pending">Loading...</div>
    <div v-else-if="error">Error loading cruises</div>
    <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <UCard v-for="cruise in cruises.data" :key="cruise.id">
        <template #header>
          <img :src="cruise.gallery?.[0] || 'https://via.placeholder.com/300'" class="w-full h-48 object-cover" />
        </template>
        <h2 class="text-xl font-semibold">{{ cruise.name_en }} / {{ cruise.name_cn }}</h2>
        <p class="text-gray-600">{{ cruise.description?.substring(0, 100) }}...</p>
        <template #footer>
          <UButton :to="`/cruises/${cruise.id}`" block>View Details</UButton>
        </template>
      </UCard>
    </div>
  </div>
</template>
