<!-- admin/app/pages/routes/index.vue — 航线列表页面 -->
<!-- 展示所有航线数据表格 -->
<script setup lang="ts">
import { onMounted, ref } from 'vue'
// 航线列表数据，改为真实 API 请求。
const { request } = useApi()
const items = ref<{ id: number; code: string; name: string }[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

async function loadItems() {
  loading.value = true
  error.value = null
  try {
    const res = await request('/routes')
    const payload = res?.data ?? res ?? []
    items.value = Array.isArray(payload) ? payload : payload?.list ?? []
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load routes'
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
  if (!confirm(`确认删除航线 #${id} 吗？`)) return
  try {
    await request(`/routes/${id}`, { method: 'DELETE' })
    await loadItems()
  } catch (e: any) {
    error.value = e?.message ?? 'failed to delete route'
  }
}

onMounted(loadItems)
</script>

<template>
  <div class="page">
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:12px;">
      <h1>Routes</h1>
      <NuxtLink to="/routes/new"><button type="button">新建航线</button></NuxtLink>
    </div>
    <p v-if="loading">Loading...</p>
    <p v-else-if="error" class="text-red-500">{{ error }}</p>
    <p v-else-if="items.length === 0">No data</p>
    <table v-else>
      <thead>
        <tr>
          <th>ID</th>
          <th>Code</th>
          <th>Name</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="r in items" :key="r.id">
          <td>{{ r.id }}</td>
          <td>{{ r.code }}</td>
          <td>{{ r.name }}</td>
          <td>
            <NuxtLink :to="`/routes/${r.id}`">编辑</NuxtLink>
            <button type="button" style="margin-left:8px" @click="handleDelete(r.id)">删除</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

