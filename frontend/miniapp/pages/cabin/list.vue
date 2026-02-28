<!-- miniapp/pages/cabin/list.vue — 小程序端舱房列表页面 -->
<!-- 展示舱房卡片列表 -->
<template>
  <view class="page">
    <text class="title">精选舱房</text>
    <!-- 舱房卡片组件 -->
    <CabinCard />
  </view>
</template>

<script setup lang="ts">
// 引入舱房卡片组件
import CabinCard from '../../components/CabinCard.vue'
</script>

<style scoped>
.page {
  min-height: 100vh;
  padding: 36rpx;
  background: #f3f8ff;
}

.title {
  display: block;
  margin-bottom: 20rpx;
  font-size: 42rpx;
  font-weight: 700;
  color: #1f4871;
}
</style>
