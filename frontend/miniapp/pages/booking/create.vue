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
      <input v-model="voyageId" type="number" placeholder="例如 1" :disabled="submitting" />

      <text class="label">舱位 SKU ID</text>
      <input v-model="cabinSkuId" type="number" placeholder="例如 1" :disabled="submitting" />

      <text class="label">乘客人数</text>
      <input v-model.number="guests" type="number" min="1" placeholder="例如 2" :disabled="submitting" />

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
  padding: 30rpx;
  background:
    radial-gradient(circle at 8% 0, #dfeaf4 0, transparent 28%),
    linear-gradient(180deg, #f4f8fb 0%, #edf3f7 100%);
}

.title {
  display: block;
  margin-bottom: 8rpx;
  font-size: 46rpx;
  font-weight: 700;
  color: #122b42;
}

.subtitle {
  display: block;
  margin-bottom: 18rpx;
  font-size: 24rpx;
  color: #5a7189;
}

.form {
  display: flex;
  flex-direction: column;
  gap: 10rpx;
  padding: 24rpx;
  border-radius: 24rpx;
  background: #fff;
  border: 1rpx solid #d4e0ea;
  box-shadow: 0 16rpx 36rpx rgba(16, 47, 72, 0.12);
}

.label {
  font-size: 23rpx;
  font-weight: 600;
  color: #24435f;
}

input {
  border: 2rpx solid #cfdde7;
  border-radius: 14rpx;
  padding: 18rpx;
  background: #f9fcff;
}

.submit-btn {
  margin-top: 8rpx;
  border-radius: 16rpx;
  background: linear-gradient(135deg, #0f3d5c, #1f5f86);
  color: #fff;
}

.error {
  color: #d13e5b;
}

.success {
  color: #0f8b62;
}
</style>
