<script setup lang="ts">
const form = reactive({
  passengers: [
    { name_cn: '', name_en: '', doc_type: 'passport', doc_number: '' }
  ]
})

const submitOrder = async () => {
  const { data, error } = await useFetch('/api/v1/orders', {
    method: 'POST',
    baseURL: 'http://localhost:8080',
    body: {
      voyage_id: '...', // From route query or store
      cabin_type_id: '...', 
      passengers: form.passengers
    }
  })
  
  if (data.value) {
    navigateTo(`/booking/pay?orderId=${data.value.data.order_id}&url=${encodeURIComponent(data.value.data.payment_url)}`)
  }
}
</script>

<template>
  <div class="container mx-auto p-4">
    <h1 class="text-2xl font-bold mb-4">Passenger Information</h1>
    <form @submit.prevent="submitOrder">
      <div v-for="(p, index) in form.passengers" :key="index" class="mb-4 p-4 border rounded">
        <h3>Passenger {{ index + 1 }}</h3>
        <UInput v-model="p.name_cn" label="Chinese Name" />
        <UInput v-model="p.name_en" label="English Name" />
        <UInput v-model="p.doc_number" label="Document Number" />
      </div>
      <UButton type="submit">Submit Order</UButton>
    </form>
  </div>
</template>
