<template>
  <view class="container">
    <view class="order-list">
      <view class="order-item" v-for="order in orders" :key="order.id">
        <view class="info">
          <text>Order: {{ order.order_no }}</text>
          <text>Status: {{ order.status }}</text>
        </view>
        <view class="actions">
          <button v-if="order.status === 'pending'" @click="cancelOrder(order.id)">Cancel</button>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'

const orders = ref([])

const loadOrders = () => {
  uni.request({
    url: 'http://localhost:8080/api/v1/orders/mine',
    success: (res) => {
      orders.value = res.data.data
    }
  })
}

const cancelOrder = (id: string) => {
  uni.request({
    url: `http://localhost:8080/api/v1/orders/${id}/cancel`,
    method: 'POST',
    success: () => {
      loadOrders()
    }
  })
}

onLoad(() => {
  loadOrders()
})
</script>

<style>
.container { padding: 20px; }
.order-item { padding: 15px; border-bottom: 1px solid #eee; display: flex; justify-content: space-between; }
</style>
