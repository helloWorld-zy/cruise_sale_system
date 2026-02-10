<script setup lang="ts">
const { data: orders, refresh } = await useFetch('/api/v1/orders/mine', {
  baseURL: 'http://localhost:8080'
})

const cancelOrder = async (id: string) => {
  if (!confirm('Are you sure?')) return
  await useFetch(`/api/v1/orders/${id}/cancel`, {
    method: 'POST',
    baseURL: 'http://localhost:8080'
  })
  refresh()
}
</script>

<template>
  <div class="container mx-auto p-4">
    <h1 class="text-2xl font-bold mb-4">My Orders</h1>
    <div class="space-y-4">
      <UCard v-for="order in orders.data" :key="order.id">
        <div class="flex justify-between items-center">
          <div>
            <h3 class="font-bold">Order {{ order.order_no }}</h3>
            <p>Status: {{ order.status }}</p>
            <p>Total: {{ order.currency }} {{ order.total_amount }}</p>
          </div>
          <div>
            <UButton v-if="['pending', 'confirmed'].includes(order.status)" 
                     color="red" variant="soft" @click="cancelOrder(order.id)">
              Cancel
            </UButton>
            <UButton v-if="order.status === 'pending'" :to="`/booking/pay?orderId=${order.id}`" class="ml-2">
              Pay Now
            </UButton>
          </div>
        </div>
      </UCard>
    </div>
  </div>
</template>
