<!-- web/pages/booking/confirm.vue — 预订确认页面 -->
<!-- H-01 修复：乘客信息表单 + 输入验证 + 防重复提交 + 调用预订 API -->
<script setup lang="ts">
import { ref, computed } from 'vue'

declare const useApi: any

// 预订确认表单：校验参数、填写乘客人数、提交预订请求
const voyageId = ref(0)
const cabinSkuId = ref(0)
const guests = ref(1)
const loading = ref(false)
const errorMsg = ref('')
const { request } = useApi()

// 从 URL query 参数初始化
// Nuxt 自动导入 useRoute
const route = useRoute()
voyageId.value = Number(route.query.voyage_id) || 0
cabinSkuId.value = Number(route.query.cabin_sku_id) || 0
const guestsQuery = Number(route.query.guests)
if (guestsQuery > 0) guests.value = guestsQuery

const canSubmit = computed(
    () => voyageId.value > 0 && cabinSkuId.value > 0 && guests.value > 0 && !loading.value
)

async function handleSubmit() {
    if (!canSubmit.value) return
    errorMsg.value = ''
    loading.value = true
    try {
      // 统一通过 useApi 发起请求并由 composable 注入 token。
      const res = await request<any>('/bookings', {
            method: 'POST',
            body: {
                voyage_id: voyageId.value,
                cabin_sku_id: cabinSkuId.value,
                guests: guests.value,
            },
        })

      const orderId = Number(res?.data?.id ?? res?.id)
      if (!Number.isFinite(orderId) || orderId <= 0) {
        throw new Error('预订成功，但订单号缺失')
      }

      await navigateTo({
        path: '/booking/success',
        query: { order_id: String(orderId) },
      })
    } catch (e: any) {
        errorMsg.value = e?.data?.message ?? e?.message ?? '预订失败，请重试'
    } finally {
        loading.value = false
    }
}
</script>

<template>
  <div class="page">
    <h1>Confirm Booking</h1>

    <form data-testid="booking-form" @submit.prevent="handleSubmit">
      <div class="field">
        <label for="guests">乘客人数</label>
        <input
          id="guests"
          v-model.number="guests"
          type="number"
          min="1"
          max="9"
          :disabled="loading"
        />
      </div>

      <p v-if="voyageId <= 0" class="hint">缺少航次信息</p>
      <p v-if="cabinSkuId <= 0" class="hint">缺少舱房信息</p>

      <p v-if="errorMsg" class="error" role="alert">{{ errorMsg }}</p>

      <button type="submit" class="submit-btn" :disabled="!canSubmit">
        {{ loading ? '提交中…' : '确认预订' }}
      </button>
    </form>
  </div>
</template>
