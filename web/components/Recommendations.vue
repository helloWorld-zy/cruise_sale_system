<script setup lang="ts">
const { data: recommendations } = await useFetch('/api/v1/recommendations', { baseURL: 'http://localhost:8080' })
</script>

<template>
  <div class="p-4 bg-gray-50 rounded-lg">
    <h3 class="font-bold text-lg mb-4">Recommended for You</h3>
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <UCard v-for="cruise in recommendations?.data || []" :key="cruise.id">
        <template #header>
          <img :src="cruise.gallery?.[0] || 'https://via.placeholder.com/300'" class="w-full h-32 object-cover" />
        </template>
        <div class="text-sm font-semibold">{{ cruise.name_en }}</div>
        <UButton :to="`/cruises/${cruise.id}`" size="xs" class="mt-2">View</UButton>
      </UCard>
    </div>
  </div>
</template>
