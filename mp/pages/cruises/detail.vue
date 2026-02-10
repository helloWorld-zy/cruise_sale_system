<template>
  <view class="container" v-if="detail">
    <view class="header">
      <text class="title">{{ detail.name_cn }}</text>
    </view>
    <view class="cabins">
      <view v-for="cabin in detail.cabin_types" :key="cabin.id">
        <text>{{ cabin.name }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'

const detail = ref(null)

onLoad((option) => {
  if (option.id) {
    uni.request({
      url: `http://localhost:8080/api/v1/cruises/${option.id}`,
      success: (res) => {
        detail.value = res.data.data
      }
    })
  }
})
</script>

<style>
.container { padding: 20px; }
.title { font-size: 20px; font-weight: bold; }
</style>
