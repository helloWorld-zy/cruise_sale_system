<!-- web/pages/booking/index.vue — 预订第一步页面 -->
<!-- H-01 修复：航次和舱房选择，填写后跳转到确认页 -->
<script setup lang="ts">
import { ref, computed } from 'vue'

// 航次和舱房选择表单，验证通过后跳转到预订确认页
const voyageId = ref('')
const cabinSkuId = ref('')
const guests = ref(1)
const loading = ref(false)
const errorMsg = ref('')

// Nuxt 自动导入 useRouter
const router = useRouter()

const canNext = computed(
    () => voyageId.value !== '' && cabinSkuId.value !== '' && guests.value > 0 && !loading.value
)

async function handleNext() {
    if (!canNext.value) return
    loading.value = true
    errorMsg.value = ''
    try {
        await router.push({
            path: '/booking/confirm',
            query: {
                voyage_id: voyageId.value,
                cabin_sku_id: cabinSkuId.value,
                guests: String(guests.value),
            },
        })
    } catch (e: any) {
        errorMsg.value = e?.message ?? '跳转失败'
    } finally {
        loading.value = false
    }
}
</script>

<template>
  <div class="page">
    <h1>Booking Step 1</h1>

    <form data-testid="booking-step1-form" @submit.prevent="handleNext">
      <div class="field">
        <label for="voyage-id">航次 ID</label>
        <input
          id="voyage-id"
          v-model="voyageId"
          type="number"
          min="1"
          placeholder="请输入航次 ID"
          :disabled="loading"
        />
      </div>

      <div class="field">
        <label for="cabin-sku-id">舱房 SKU ID</label>
        <input
          id="cabin-sku-id"
          v-model="cabinSkuId"
          type="number"
          min="1"
          placeholder="请输入舱房 SKU ID"
          :disabled="loading"
        />
      </div>

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

      <p v-if="errorMsg" class="error" role="alert">{{ errorMsg }}</p>

      <button type="submit" :disabled="!canNext">
        {{ loading ? '处理中…' : '下一步' }}
      </button>
    </form>
  </div>
</template>
