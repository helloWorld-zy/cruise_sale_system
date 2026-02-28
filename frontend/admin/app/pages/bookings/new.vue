<!-- admin/app/pages/bookings/new.vue -->
<script setup lang="ts">
import { ref } from 'vue'
import { useApi } from '../../composables/useApi'
const { request } = useApi()
const loading = ref(false)
const error = ref<string | null>(null)
const success = ref<string | null>(null)
const form = ref({
  user_id: 2,
  voyage_id: 0,
  cabin_sku_id: 0,
  guests: 1,
})

async function handleSubmit() {
  if (loading.value) return
  loading.value = true
  error.value = null
  success.value = null
  try {
    await request('/bookings', {
      method: 'POST',
      body: {
        user_id: Number(form.value.user_id),
        voyage_id: Number(form.value.voyage_id),
        cabin_sku_id: Number(form.value.cabin_sku_id),
        guests: Number(form.value.guests),
      },
    })
    success.value = '创建成功'
  } catch (e: any) {
    error.value = e?.message ?? 'failed to create booking'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="page">
    <h1>新建订单</h1>
    <form @submit.prevent="handleSubmit">
      <input v-model.number="form.user_id" type="number" min="1" placeholder="User ID" :disabled="loading" />
      <input v-model.number="form.voyage_id" type="number" min="1" placeholder="Voyage ID" :disabled="loading" />
      <input v-model.number="form.cabin_sku_id" type="number" min="1" placeholder="Cabin SKU ID" :disabled="loading" />
      <input v-model.number="form.guests" type="number" min="1" placeholder="Guests" :disabled="loading" />
      <p v-if="error" class="error">{{ error }}</p>
      <p v-if="success" style="color:#11865d">{{ success }}</p>
      <button type="button" :disabled="loading" @click="handleSubmit">{{ loading ? '提交中...' : '创建' }}</button>
    </form>
  </div>
</template>

