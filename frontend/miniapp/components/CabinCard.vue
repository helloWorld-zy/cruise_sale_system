<!-- miniapp/components/CabinCard.vue — 小程序端舱房卡片组件 -->
<!-- 展示舱房简要信息（小程序使用 view/text 原生组件） -->
<template>
  <view class="card">
    <text class="name">阳台海景舱</text>
    <text class="desc">含双人早餐与海上观景位</text>
    <text class="price">￥199.00 / 晚</text>
    <!-- TODO: 对接舱房数据展示 -->
  </view>
</template>

<style scoped>
.card {
  display: flex;
  flex-direction: column;
  gap: 10rpx;
  border-radius: 20rpx;
  padding: 26rpx;
  background: linear-gradient(135deg, #ffffff 0%, #eef4ff 100%);
  box-shadow: 0 12rpx 30rpx rgba(29, 74, 126, 0.14);
}

.name {
  font-size: 36rpx;
  font-weight: 700;
  color: #1d4877;
}

.desc {
  color: #5e7490;
}

.price {
  font-size: 34rpx;
  font-weight: 700;
  color: #e05f1d;
}
</style>
