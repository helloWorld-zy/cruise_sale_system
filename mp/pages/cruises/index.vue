<template>
  <view class="container">
    <view class="cruise-list">
      <view class="cruise-item" v-for="item in cruises" :key="item.id" @click="goToDetail(item.id)">
        <text>{{ item.name_cn }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'

const cruises = ref([])

onLoad(() => {
  uni.request({
    url: 'http://localhost:8080/api/v1/cruises',
    success: (res) => {
      cruises.value = res.data.data
    }
  })
})

const goToDetail = (id: string) => {
  uni.navigateTo({ url: `/pages/cruises/detail?id=${id}` })
}
</script>

<style>
.container { padding: 20px; }
.cruise-item { padding: 15px; border-bottom: 1px solid #eee; }
</style>
