<script setup lang="ts">
import { computed, ref } from 'vue'
import { request } from '../../src/utils/request'

// 表单状态：航次、舱位与出行人数。
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
    <text>Create Booking</text>

    <view class="form">
      <input v-model="voyageId" type="number" placeholder="Voyage ID" :disabled="submitting" />
      <input v-model="cabinSkuId" type="number" placeholder="Cabin SKU ID" :disabled="submitting" />
      <input v-model.number="guests" type="number" min="1" placeholder="Guests" :disabled="submitting" />

      <button :disabled="!canSubmit" @click="handleSubmit">
        {{ submitting ? 'Submitting...' : 'Submit Booking' }}
      </button>

      <text v-if="errorMsg">{{ errorMsg }}</text>
      <text v-if="successMsg">{{ successMsg }}</text>
    </view>
  </view>
</template>
