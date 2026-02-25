<!-- miniapp/pages/login/login.vue — 小程序端登录页面 -->
<!-- 当前为轻量入口，复用 PrimaryButton 组件承载后续登录交互 -->
<script setup lang="ts">
import PrimaryButton from '../../components/PrimaryButton.vue'
</script>

<template>
  <view class="page">
    <text>Login</text>
    <PrimaryButton>Submit</PrimaryButton>
  </view>
</template>
