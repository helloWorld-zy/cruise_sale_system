<!-- miniapp/pages/booking/create.vue — 小程序端创建预订页面 -->
<!-- 填写航次、舱房和出行人数，提交后创建预订 -->
<script setup lang="ts">
import { computed, ref } from 'vue'
import { request } from '../../src/utils/request'

// 表单状态：航次、舱位与出行人数
const voyageId = ref('')
const cabinSkuId = ref('')
const guests = ref(1)
const submitting = ref(false)
const errorMsg = ref('')
const successMsg = ref('')

// canSubmit 控制按钮可提交状态，避免非法参数提交。
const canSubmit = computed(() => {
  return !submitting.value && Number(voyageId.value) > 0 && Number(cabinSkuId.value) > 0 && guests.value > 0
})

// handleSubmit 调用后端创建预订接口并处理成功/失败反馈。
async function handleSubmit() {
  if (!canSubmit.value) return

  submitting.value = true
  errorMsg.value = ''
  successMsg.value = ''

  try {
    await request('/bookings', {
      method: 'POST',
      data: {
        voyage_id: Number(voyageId.value),
        cabin_sku_id: Number(cabinSkuId.value),
        guests: guests.value,
      },
    })
    successMsg.value = 'Booking Created'
  } catch (err: any) {
    errorMsg.value = err?.message ?? 'Booking failed'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <view class="page">
    <text class="title">创建预订</text>
    <text class="subtitle">三步完成舱房锁定：航次、舱位、人数。</text>

    <view class="form">
      <text class="label">航次 ID</text>
      <input v-model="voyageId" type="number" placeholder="输入航次 ID" :disabled="submitting" />

      <text class="label">舱位 SKU ID</text>
      <input v-model="cabinSkuId" type="number" placeholder="输入舱位 SKU ID" :disabled="submitting" />

      <text class="label">乘客人数</text>
      <input v-model.number="guests" type="number" min="1" placeholder="输入乘客人数" :disabled="submitting" />

      <button class="submit-btn" :disabled="!canSubmit" @click="handleSubmit">
        {{ submitting ? 'Submitting...' : 'Submit Booking' }}
      </button>

      <text v-if="errorMsg" class="error">{{ errorMsg }}</text>
      <text v-if="successMsg" class="success">{{ successMsg }}</text>
    </view>
  </view>
</template>

<style scoped>
.page {
  min-height: 100vh;
  padding: 40rpx;
  background: #f5f7fa;
  position: relative;
}
.page::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 360rpx;
  background: linear-gradient(135deg, #0cebeb 0%, #20e3b2 50%, #29ffc6 100%);
  border-bottom-left-radius: 40rpx;
  border-bottom-right-radius: 40rpx;
  z-index: 0;
}
.title {
  position: relative;
  z-index: 1;
  display: block;
  margin-bottom: 12rpx;
  font-size: 48rpx;
  font-weight: 800;
  color: #fff;
}
.subtitle {
  position: relative;
  z-index: 1;
  display: block;
  margin-bottom: 32rpx;
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.9);
}
.form {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  gap: 16rpx;
  padding: 32rpx;
  border-radius: 32rpx;
  background: #fff;
  border: none;
  box-shadow: 0 12rpx 32rpx rgba(0, 0, 0, 0.05);
}
.label {
  font-size: 26rpx;
  font-weight: 700;
  color: #333;
  margin-top: 8rpx;
}
input {
  border: none;
  border-radius: 20rpx;
  padding: 24rpx;
  background: #f5f7fa;
  font-size: 28rpx;
}
.submit-btn {
  margin-top: 24rpx;
  border-radius: 999rpx;
  background: linear-gradient(135deg, #ff8e53 0%, #ff6b6b 100%);
  color: #fff;
  font-size: 32rpx;
  font-weight: 700;
  border: none;
  padding: 2rpx 0;
}
.error { color: #ffcccc; margin-top: 16rpx; text-align: center; }
.success { color: #0cebeb; margin-top: 16rpx; text-align: center; font-weight: bold; }
</style>
