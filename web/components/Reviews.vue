<script setup lang="ts">
const props = defineProps<{ cruiseId: string }>()
const { data: reviews } = await useFetch(`/api/v1/cruises/${props.cruiseId}/reviews`, { baseURL: 'http://localhost:8080' })
</script>

<template>
  <div class="mt-8">
    <h3 class="text-xl font-bold mb-4">Reviews</h3>
    <div v-if="reviews?.data?.length === 0">No reviews yet.</div>
    <div v-else class="space-y-4">
      <div v-for="review in reviews.data" :key="review.id" class="p-4 border rounded">
        <div class="flex items-center gap-2 mb-2">
          <div class="font-bold">User</div>
          <div class="text-yellow-500">★ {{ review.rating }}</div>
        </div>
        <p>{{ review.comment }}</p>
      </div>
    </div>
  </div>
</template>
