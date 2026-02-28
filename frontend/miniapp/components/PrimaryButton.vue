<!-- miniapp/components/PrimaryButton.vue — 主要操作按钮组件 -->
<!-- 小程序端通用的主按钮样式封装，通过 slot 传入按钮文字 -->
<template>
  <button class="primary-button"><slot/></button>
</template>

<style scoped>
.primary-button {
  border-radius: 16rpx;
  background: linear-gradient(135deg, #1f7bd5, #2d9af9);
  color: #fff;
  font-weight: 700;
}
</style>
