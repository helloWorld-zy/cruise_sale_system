<!-- miniapp/pages/cabin/detail.vue — 小程序端舱房详情页面 -->
<!-- 展示舱房 SKU 的详细信息和价格日历 -->
<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { request } from '../../src/utils/request'

// 通过 props 接收页面参数，便于运行时与测试场景复用。
const props = defineProps<{ cabinSkuId?: number | string }>()
const loading = ref(false)
const error = ref('')
const detail = ref<any>(null)
const prices = ref<any[]>([])

function resolveSkuId() {
  return Number(props.cabinSkuId ?? 0)
}

async function loadDetail() {
  const skuId = resolveSkuId()
  if (!skuId) {
    error.value = '缺少 cabinSkuId 参数'
    return
  }
  loading.value = true
  error.value = ''
  try {
    const detailRes = await request(`/cabins/${skuId}`)
    detail.value = detailRes?.data ?? detailRes
    const pricesRes = await request(`/cabins/${skuId}/prices`)
    prices.value = pricesRes?.data ?? pricesRes ?? []
  } catch (e: any) {
    error.value = e?.message ?? '加载舱房失败'
  } finally {
    loading.value = false
  }
}

onMounted(loadDetail)
</script>

<template>
  <view class="page">
    <text class="title">舱房详情</text>
    <text v-if="loading" class="hint">Loading...</text>
    <text v-else-if="error" class="error">{{ error }}</text>
    <view v-else-if="detail" class="panel">
      <text>名称：{{ detail.name ?? '-' }}</text>
      <text>描述：{{ detail.description ?? '-' }}</text>
      <text>价格：{{ detail.price_cents ?? detail.price ?? '-' }}</text>
      <text>库存：{{ detail.available ?? '-' }}</text>
      <view v-if="prices.length === 0">
        <text class="hint">No price data</text>
      </view>
      <view v-else class="price-list">
        <text v-for="item in prices" :key="`${item.date}-${item.occupancy ?? 0}`">
          {{ item.date }} / {{ item.price_cents ?? item.price }}
        </text>
      </view>
    </view>
  </view>
</template>

<style scoped>
.page {
  min-height: 100vh;
  padding: 36rpx;
  background: #f5f8ff;
}

.title {
  display: block;
  margin-bottom: 18rpx;
  font-size: 42rpx;
  font-weight: 700;
  color: #20476f;
}

.panel {
  display: flex;
  flex-direction: column;
  gap: 12rpx;
  padding: 26rpx;
  border-radius: 18rpx;
  background: #fff;
  box-shadow: 0 12rpx 30rpx rgba(42, 80, 130, 0.12);
}

.price-list {
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.hint {
  color: #5a7190;
}

.error {
  color: #d13e5b;
}
</style>
