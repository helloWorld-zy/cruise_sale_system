<script setup lang="ts">
const route = useRoute()
const payUrl = route.query.url as string
</script>

<template>
  <div class="container mx-auto p-4 text-center">
    <h1 class="text-2xl font-bold mb-4">Payment</h1>
    <p>Order created successfully. Please pay within 15 minutes.</p>
    
    <div class="mt-8">
      <a :href="payUrl" target="_blank" class="bg-blue-500 text-white px-6 py-2 rounded">
        Pay Now
      </a>
    </div>
  </div>
</template>
