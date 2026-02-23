<!-- miniapp/pages/cabin/list.vue — 小程序端舱房列表页面 -->
<!-- 展示舱房卡片列表 -->
<template>
  <view class="page">
    <text>Cabins</text>
    <!-- 舱房卡片组件 -->
    <CabinCard />
  </view>
</template>

<script setup lang="ts">
// 引入舱房卡片组件
import CabinCard from '../../components/CabinCard.vue'
</script>
