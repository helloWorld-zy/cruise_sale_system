<!-- miniapp/components/PrimaryButton.vue — 主要操作按钮组件 -->
<!-- 小程序端通用的主按钮样式封装，通过 slot 传入按钮文字 -->
<template>
  <button class="primary-button"><slot/></button>
</template>

<style scoped>
.primary-button {
  border: 0;
  border-radius: 18rpx;
  background: linear-gradient(135deg, #0f3d5c, #1f5f86);
  color: #fff;
  font-weight: 700;
  height: 86rpx;
  line-height: 86rpx;
  box-shadow: 0 12rpx 26rpx rgba(16, 47, 72, 0.28);
}

.primary-button:active {
  transform: translateY(1rpx);
  opacity: 0.94;
}
</style>
