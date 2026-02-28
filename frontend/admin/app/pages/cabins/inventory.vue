<!-- admin/app/pages/cabins/inventory.vue — 库存管理页面 -->
<script setup lang="ts">
import { onMounted, ref } from 'vue'
// 从路由 query 读取 SKU ID，并对接库存查询和调整接口。
const { request } = useApi()
const route = useRoute()
const skuId = Number(route?.query?.skuId ?? 0)
const loading = ref(false)
const submitting = ref(false)
const error = ref<string | null>(null)
const inventory = ref<{ total: number; available: number } | null>(null)
const delta = ref(0)

async function loadInventory() {
  if (!skuId) {
    error.value = '缺少 skuId 参数'
    return
  }
  loading.value = true
  error.value = null
  try {
    const res = await request(`/cabins/${skuId}/inventory`)
    inventory.value = res?.data ?? res ?? null
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load inventory'
  } finally {
    loading.value = false
  }
}

async function adjustInventory() {
  if (!skuId || submitting.value) return
  submitting.value = true
  error.value = null
  try {
    await request(`/cabins/${skuId}/inventory/adjust`, {
      method: 'POST',
      body: { delta: Number(delta.value) },
    })
    await loadInventory()
  } catch (e: any) {
    error.value = e?.message ?? 'failed to adjust inventory'
  } finally {
    submitting.value = false
  }
}

onMounted(loadInventory)
</script>

<template>
  <div class="page">
    <h1>Inventory</h1>
    <p v-if="loading">Loading...</p>
    <p v-else-if="error" class="text-red-500">{{ error }}</p>
    <div v-else-if="inventory">
      <p>Total: {{ inventory.total }}</p>
      <p>Available: {{ inventory.available }}</p>
      <div>
        <input v-model.number="delta" type="number" placeholder="delta" />
        <button :disabled="submitting" @click="adjustInventory">
          {{ submitting ? 'Submitting...' : 'Adjust' }}
        </button>
      </div>
    </div>
  </div>
</template>

