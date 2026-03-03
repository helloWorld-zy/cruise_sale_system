<script setup lang="ts">
import { onMounted, ref } from 'vue'

const { request } = useApi()

type StaffItem = {
  id: number
  real_name?: string
  email?: string
  role?: string
  status?: number
}

const loading = ref(false)
const error = ref('')
const items = ref<StaffItem[]>([])

async function loadItems() {
  loading.value = true
  error.value = ''
  try {
    const res = await request('/staffs')
    const payload = res?.data ?? res
    items.value = Array.isArray(payload) ? payload : (payload?.list ?? [])
  } catch (e: any) {
    error.value = e?.message ?? '加载员工列表失败'
  } finally {
    loading.value = false
  }
}

onMounted(loadItems)
</script>

<template>
  <div class="page">
    <h1>员工管理</h1>
    <p v-if="loading" data-test="loading">加载中...</p>
    <p v-else-if="error" data-test="error" class="error">{{ error }}</p>
    <p v-else-if="items.length === 0" data-test="empty">暂无员工数据</p>
    <table v-else data-test="table">
      <thead>
        <tr>
          <th>ID</th>
          <th>姓名</th>
          <th>邮箱</th>
          <th>角色</th>
          <th>状态</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="s in items" :key="s.id">
          <td>{{ s.id }}</td>
          <td>{{ s.real_name || '-' }}</td>
          <td>{{ s.email || '-' }}</td>
          <td>{{ s.role || '-' }}</td>
          <td>{{ s.status === 1 ? '启用' : '禁用' }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
