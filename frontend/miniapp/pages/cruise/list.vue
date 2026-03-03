<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { request } from '../../src/utils/request'

const keyword = ref('')
const loading = ref(false)
const error = ref('')
const items = ref<Record<string, any>[]>([])

async function loadItems() {
  loading.value = true
  error.value = ''
  try {
    const query = keyword.value.trim() ? `?keyword=${encodeURIComponent(keyword.value.trim())}&page=1&page_size=30` : '?page=1&page_size=30'
    const res = await request<any>(`/cruises${query}`)
    const payload = (res as any)?.data ?? res ?? {}
    items.value = Array.isArray(payload) ? payload : payload?.list ?? []
  } catch (e: any) {
    error.value = e?.message ?? '加载邮轮失败'
  } finally {
    loading.value = false
  }
}

onMounted(loadItems)
</script>

<template>
  <view class="page">
    <view class="search-wrap">
      <input v-model="keyword" class="search-input" placeholder="搜索邮轮" @confirm="loadItems" />
    </view>

    <view v-if="loading" class="state">Loading...</view>
    <view v-else-if="error" class="state error">{{ error }}</view>
    <view v-else>
      <view v-for="item in items" :key="item.id" class="card">
        <image class="cover" mode="aspectFill" :src="item.cover_url || `https://picsum.photos/seed/mini-cruise-${item.id}/800/600`" />
        <view class="body">
          <text class="name">{{ item.name || '-' }}</text>
          <text class="meta">{{ item.tonnage || '-' }}吨 · {{ item.passenger_capacity || '-' }}人 · {{ item.deck_count || '-' }}层甲板</text>
          <text class="detail-link">查看详情 ></text>
        </view>
      </view>
    </view>
  </view>
</template>

<style scoped>
.page {
  min-height: 100vh;
  background: #f6f7fb;
  padding-top: 108rpx;
}
.search-wrap {
  position: fixed;
  left: 0;
  right: 0;
  top: 0;
  z-index: 10;
  background: #fff;
  padding: 18rpx 24rpx;
}
.search-input {
  height: 68rpx;
  border-radius: 32rpx;
  background: #f5f5f5;
  padding: 0 24rpx;
  font-size: 26rpx;
}
.card {
  margin: 16rpx 24rpx;
  border-radius: 16rpx;
  background: #fff;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.06);
  overflow: hidden;
}
.cover {
  width: 100%;
  height: 320rpx;
}
.body {
  padding: 22rpx;
  display: flex;
  flex-direction: column;
  gap: 10rpx;
}
.name {
  font-size: 32rpx;
  font-weight: 700;
  color: #172533;
}
.meta {
  font-size: 24rpx;
  color: #7a8796;
}
.detail-link {
  font-size: 24rpx;
  color: #ff6b6b;
}
.state {
  padding: 36rpx 24rpx;
  font-size: 26rpx;
  color: #6b7c8f;
}
.error {
  color: #d7415f;
}
</style>
