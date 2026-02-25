<script setup lang="ts">
// 登录页当前为轻量入口，复用 PrimaryButton 组件承载后续登录交互。
import PrimaryButton from '../../components/PrimaryButton.vue'
</script>

<template>
  <view class="page">
    <text>Login</text>
    <PrimaryButton>Submit</PrimaryButton>
  </view>
</template>
