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
    <view class="header-bg"></view>
    <view class="search-wrap">
      <view class="search-box">
        <text class="search-icon">🔍</text>
        <input v-model="keyword" class="search-input" placeholder="搜索邮轮" placeholder-style="color: rgba(255,255,255,0.7)" @confirm="loadItems" />
      </view>
    </view>

    <view v-if="loading" class="state">Loading...</view>
    <view v-else-if="error" class="state error">{{ error }}</view>
    <view v-else class="list-container">
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
  background: #f5f7fa;
  padding-top: 140rpx;
  position: relative;
}
.header-bg {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 380rpx;
  background: linear-gradient(135deg, #0cebeb 0%, #20e3b2 50%, #29ffc6 100%);
  border-bottom-left-radius: 40rpx;
  border-bottom-right-radius: 40rpx;
  z-index: 0;
}
.search-wrap {
  position: fixed;
  left: 0;
  right: 0;
  top: 0;
  z-index: 10;
  padding: 24rpx;
  background: #0cebeb;
}
.search-box {
  display: flex;
  align-items: center;
  height: 76rpx;
  border-radius: 38rpx;
  background: rgba(255, 255, 255, 0.25);
  padding: 0 32rpx;
}
.search-icon {
  font-size: 28rpx;
  color: #fff;
  margin-right: 12rpx;
}
.search-input {
  flex: 1;
  font-size: 28rpx;
  color: #fff;
}
.list-container {
  position: relative;
  z-index: 1;
  padding-bottom: 40rpx;
}
.card {
  margin: 24rpx;
  border-radius: 32rpx;
  background: #fff;
  box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.04);
  overflow: hidden;
  border: none;
}
.cover {
  width: 100%;
  height: 260rpx;
}
.body {
  padding: 32rpx;
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}
.name {
  font-size: 36rpx;
  font-weight: 700;
  color: #222;
}
.meta {
  font-size: 24rpx;
  color: #888;
}
.detail-link {
  align-self: flex-start;
  margin-top: 16rpx;
  padding: 12rpx 32rpx;
  border-radius: 40rpx;
  background: #ff6b6b;
  color: #fff;
  font-size: 26rpx;
  font-weight: 600;
}
.state {
  position: relative;
  z-index: 1;
  padding: 60rpx 24rpx;
  text-align: center;
  font-size: 28rpx;
  color: #fff;
}
.error {
  color: #ffcccc;
}
</style>
