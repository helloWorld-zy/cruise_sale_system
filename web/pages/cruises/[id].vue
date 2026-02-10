<script setup lang="ts">
const route = useRoute()
const { data: detail, pending, error } = await useFetch(`/api/v1/cruises/${route.params.id}`, {
  baseURL: 'http://localhost:8080'
})
</script>

<template>
  <div v-if="pending">Loading...</div>
  <div v-else-if="error">Error loading cruise details</div>
  <div v-else class="container mx-auto p-4">
    <div class="mb-8">
      <h1 class="text-3xl font-bold">{{ detail.data.name_en }}</h1>
      <h2 class="text-xl text-gray-500">{{ detail.data.name_cn }}</h2>
      <div class="grid grid-cols-2 gap-4 mt-4">
        <div>
          <img :src="detail.data.gallery?.[0]" class="w-full rounded-lg" />
        </div>
        <div>
          <p>{{ detail.data.description }}</p>
          <div class="mt-4 grid grid-cols-2 gap-2 text-sm">
            <div>Tonnage: {{ detail.data.tonnage }}</div>
            <div>Capacity: {{ detail.data.capacity }}</div>
            <div>Decks: {{ detail.data.decks }}</div>
          </div>
        </div>
      </div>
    </div>

    <h3 class="text-2xl font-semibold mb-4">Cabins</h3>
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <UCard v-for="cabin in detail.data.cabin_types" :key="cabin.id">
        <h4 class="font-bold">{{ cabin.name }}</h4>
        <p>Area: {{ cabin.base_area }} m²</p>
        <p>Capacity: {{ cabin.capacity }}</p>
        <UButton class="mt-2">Check Availability</UButton>
      </UCard>
    </div>
  </div>
</template>
