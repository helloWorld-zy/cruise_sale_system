<!-- admin/app/pages/cruises/index.vue — 邮轮管理列表页面 -->
<!-- 展示邮轮列表，提供新建邮轮按钮入口 -->
<template>
  <div class="page">
    <!-- 页面标题栏 -->
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:12px;">
      <h1>邮轮管理</h1>
      <NuxtLink to="/cruises/create">
        <button type="button">新建邮轮</button>
      </NuxtLink>
    </div>
    <!-- 邮轮数据表格 -->
    <table class="w-full border">
      <thead>
        <tr class="bg-gray-50">
          <th class="text-left p-2">ID</th>
          <th class="text-left p-2">名称</th>
          <th class="text-left p-2">公司ID</th>
          <th class="text-left p-2">状态</th>
          <th class="text-left p-2">操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="loading">
          <td class="p-2" colspan="5">Loading...</td>
        </tr>
        <tr v-else-if="error">
          <td class="p-2 text-red-500" colspan="5">{{ error }}</td>
        </tr>
        <tr v-else-if="items.length === 0">
          <td class="p-2" colspan="5">No data</td>
        </tr>
        <tr v-for="item in items" v-else :key="item.id">
          <td class="p-2">{{ item.id }}</td>
          <td class="p-2">{{ item.name }}</td>
          <td class="p-2">{{ item.company_id }}</td>
          <td class="p-2">{{ item.status ?? '-' }}</td>
          <td class="p-2">
            <NuxtLink :to="`/cruises/${item.id}`">编辑</NuxtLink>
            <button type="button" style="margin-left: 8px" @click="handleDelete(item.id)">删除</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
// 邮轮列表页面：调用后端接口加载真实数据。
const { request } = useApi()
const items = ref<any[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

async function loadItems() {
  loading.value = true
  error.value = null
  try {
    const res = await request('/cruises')
    const payload = res?.data ?? res ?? []
    items.value = Array.isArray(payload) ? payload : payload?.list ?? []
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load cruises'
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
  if (!confirm(`确认删除邮轮 #${id} 吗？`)) return
  try {
    await request(`/cruises/${id}`, { method: 'DELETE' })
    await loadItems()
  } catch (e: any) {
    error.value = e?.message ?? 'failed to delete cruise'
  }
}

onMounted(loadItems)
</script>

