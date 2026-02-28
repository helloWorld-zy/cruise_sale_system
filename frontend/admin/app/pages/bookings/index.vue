<!-- admin/app/pages/bookings/index.vue — 订单列表页面 -->
<!-- H-02 修复：使用 useApi 获取真实订单数据，去掉硬编码模拟数据 -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
// 使用 useApi 请求后端订单列表数据
const { request } = useApi()
const items = ref<{ id: number; status: string; total_cents: number }[]>([])
const total = ref(0)
const loading = ref(false)
const error = ref<string | null>(null)

async function loadItems() {
  loading.value = true
  try {
    const res = await request('/bookings')
    const payload = res?.data ?? res ?? {}
    items.value = payload?.list ?? []
    total.value = Number(payload?.total ?? 0)
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load bookings'
  } finally {
    loading.value = false
  }
}

function resolveId(raw: unknown) {
  const id = Number(raw)
  return Number.isFinite(id) && id > 0 ? id : 0
}

async function handleDelete(rawId: unknown) {
  const id = resolveId(rawId)
  if (!id) {
    error.value = '无效记录 ID，无法删除'
    return
  }
  if (!confirm(`确认删除订单 #${id} 吗？`)) return
  try {
    await request(`/bookings/${id}`, { method: 'DELETE' })
    await loadItems()
  } catch (e: any) {
    error.value = e?.message ?? 'failed to delete booking'
  }
}

onMounted(loadItems)
</script>

<template>
  <div class="page">
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:12px;">
      <h1>Bookings</h1>
      <NuxtLink to="/bookings/new"><button type="button">新建订单</button></NuxtLink>
    </div>
    <p>总数：{{ total }}</p>
    <p v-if="loading">Loading...</p>
    <p v-else-if="error" class="error">{{ error }}</p>
    <table v-else>
      <thead>
        <tr>
          <th>ID</th>
          <th>状态</th>
          <th>总额</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="b in items" :key="b.id">
          <td>{{ b.id }}</td>
          <td>{{ b.status }}</td>
          <td>{{ b.total_cents }}</td>
          <td>
            <NuxtLink :to="`/bookings/${b.id}`">编辑</NuxtLink>
            <button type="button" style="margin-left:8px" @click="handleDelete(b.id)">删除</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

