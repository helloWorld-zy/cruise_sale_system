<!-- admin/app/pages/cabins/index.vue — 舱房管理页面 -->
<script setup lang="ts">
import { onMounted, ref } from 'vue'
// 调用舱位列表接口并处理三态展示。
const { request } = useApi()
const items = ref<any[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

async function loadItems() {
  loading.value = true
  error.value = null
  try {
    const res = await request('/cabins')
    const payload = res?.data ?? res ?? []
    items.value = Array.isArray(payload) ? payload : payload?.list ?? []
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load cabins'
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
  if (!confirm(`确认删除舱房 #${id} 吗？`)) return
  try {
    await request(`/cabins/${id}`, { method: 'DELETE' })
    await loadItems()
  } catch (e: any) {
    error.value = e?.message ?? 'failed to delete cabin'
  }
}

onMounted(loadItems)
</script>

<template>
  <div class="page">
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:12px;">
      <h1>Cabins</h1>
      <NuxtLink to="/cabins/new"><button type="button">新建舱房</button></NuxtLink>
    </div>
    <p v-if="loading">Loading...</p>
    <p v-else-if="error" class="text-red-500">{{ error }}</p>
    <p v-else-if="items.length === 0">No data</p>
    <table v-else>
      <thead>
        <tr>
          <th>ID</th>
          <th>Voyage ID</th>
          <th>Cabin Type ID</th>
          <th>Total</th>
          <th>Available</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in items" :key="item.id">
          <td>{{ item.id }}</td>
          <td>{{ item.voyage_id }}</td>
          <td>{{ item.cabin_type_id }}</td>
          <td>{{ item.total ?? '-' }}</td>
          <td>{{ item.available ?? '-' }}</td>
          <td>
            <NuxtLink :to="`/cabins/${item.id}`">编辑</NuxtLink>
            <button type="button" style="margin-left:8px" @click="handleDelete(item.id)">删除</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

