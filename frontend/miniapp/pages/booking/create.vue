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

    <view class="form">
      <input v-model="voyageId" type="number" placeholder="Voyage ID" :disabled="submitting" />
      <input v-model="cabinSkuId" type="number" placeholder="Cabin SKU ID" :disabled="submitting" />
      <input v-model.number="guests" type="number" min="1" placeholder="Guests" :disabled="submitting" />

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
  background: linear-gradient(180deg, #fff9f2 0%, #fff 100%);
}

.title {
  display: block;
  margin-bottom: 24rpx;
  font-size: 42rpx;
  font-weight: 700;
  color: #7f3f1c;
}

.form {
  display: flex;
  flex-direction: column;
  gap: 18rpx;
  padding: 28rpx;
  border-radius: 20rpx;
  background: #fff;
  box-shadow: 0 14rpx 32rpx rgba(160, 90, 47, 0.14);
}

input {
  border: 2rpx solid #f2d7c6;
  border-radius: 14rpx;
  padding: 18rpx;
  background: #fffdfb;
}

.submit-btn {
  border-radius: 14rpx;
  background: #f06a29;
  color: #fff;
}

.error {
  color: #d13e5b;
}

.success {
  color: #0f8b62;
}
</style>
