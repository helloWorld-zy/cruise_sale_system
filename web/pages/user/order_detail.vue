<script setup lang="ts">
const route = useRoute()
const orderId = route.params.id as string
const { data: order } = await useFetch(`/api/v1/orders/${orderId}`, { baseURL: 'http://localhost:8080' })

const showNotice = ref(false)

onMounted(() => {
  if (order.value?.data?.departure_notice_url) {
    showNotice.value = true
  }
})
</script>

<template>
  <div class="container mx-auto p-4">
    <h1 class="text-2xl font-bold mb-4">Order Details</h1>
    <div v-if="order?.data">
      <p>Order No: {{ order.data.order_no }}</p>
      <p>Status: {{ order.data.status }}</p>
      
      <div v-if="order.data.departure_notice_url" class="mt-4 p-4 bg-blue-50 border border-blue-200 rounded">
        <h3 class="font-bold text-blue-800">Departure Notice Available!</h3>
        <UButton :to="order.data.departure_notice_url" target="_blank" icon="i-heroicons-document-arrow-down" class="mt-2">
          Download Departure Notice
        </UButton>
      </div>
    </div>

    <UModal v-model="showNotice">
      <div class="p-6">
        <h3 class="text-lg font-bold mb-2">Notice</h3>
        <p>Your departure notice is ready for download.</p>
        <div class="mt-4 flex justify-end">
          <UButton :to="order?.data?.departure_notice_url" target="_blank" @click="showNotice = false">Download</UButton>
        </div>
      </div>
    </UModal>
  </div>
</template>
