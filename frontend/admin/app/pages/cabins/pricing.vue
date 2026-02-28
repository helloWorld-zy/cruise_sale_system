<!-- admin/app/pages/cabins/pricing.vue — 定价管理页面 -->
<!-- 展示价格日历矩阵，每行对应一个日期/入住人数的价格 -->
<script setup lang="ts">
import { onMounted, ref } from 'vue'
import PricingRow from '../../components/PricingRow.vue'
// 对接价格列表与新增价格接口。
const { request } = useApi()
const route = useRoute()
const skuId = Number(route?.query?.skuId ?? 0)
const rows = ref<{ date: string; occupancy: number; price: number }[]>([])
const loading = ref(false)
const submitting = ref(false)
const error = ref<string | null>(null)
const form = ref({ date: '', price_cents: 0 })

async function loadRows() {
  if (!skuId) {
    error.value = '缺少 skuId 参数'
    return
  }
  loading.value = true
  error.value = null
  try {
    const res = await request(`/cabins/${skuId}/prices`)
    const raw = res?.data ?? res ?? []
    rows.value = raw.map((item: any) => ({
      date: item.date,
      occupancy: item.occupancy ?? 2,
      price: item.price ?? item.price_cents,
    }))
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load prices'
  } finally {
    loading.value = false
  }
}

async function submitPrice() {
  if (!skuId || submitting.value || !form.value.date) return
  submitting.value = true
  error.value = null
  try {
    await request(`/cabins/${skuId}/prices`, {
      method: 'POST',
      body: {
        date: form.value.date,
        price_cents: Number(form.value.price_cents),
      },
    })
    await loadRows()
  } catch (e: any) {
    error.value = e?.message ?? 'failed to save price'
  } finally {
    submitting.value = false
  }
}

onMounted(loadRows)
</script>

<template>
  <div class="page">
    <h1>Pricing Matrix</h1>
    <p v-if="loading">Loading...</p>
    <p v-else-if="error" class="text-red-500">{{ error }}</p>
    <p v-else-if="rows.length === 0">No data</p>
    <!-- 渲染每行价格数据 -->
    <PricingRow v-for="r in rows" v-else :key="`${r.date}-${r.occupancy}`" :row="r" />

    <div class="mt-4">
      <input v-model="form.date" type="date" />
      <input v-model.number="form.price_cents" type="number" min="0" placeholder="price_cents" />
      <button :disabled="submitting" @click="submitPrice">
        {{ submitting ? 'Submitting...' : 'Save Price' }}
      </button>
    </div>
  </div>
</template>

