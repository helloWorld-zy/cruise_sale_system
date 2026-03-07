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
  <div class="admin-page">
    <AdminPageHeader title="员工管理" subtitle="查看后台员工与角色状态" />

    <AdminDataCard flush>
      <div class="overflow-x-auto">
        <p v-if="loading" data-test="loading" class="p-3 text-sm text-slate-600">加载中...</p>
        <p v-else-if="error" data-test="error" class="p-3 text-sm text-rose-500">{{ error }}</p>
        <p v-else-if="items.length === 0" data-test="empty" class="p-3 text-sm text-slate-600">暂无员工数据</p>
        <table v-else data-test="table" class="w-full min-w-[760px] text-sm">
          <thead class="bg-slate-50 text-left text-slate-600">
            <tr>
              <th class="p-3">ID</th>
              <th class="p-3">姓名</th>
              <th class="p-3">邮箱</th>
              <th class="p-3">角色</th>
              <th class="p-3">状态</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="s in items" :key="s.id">
              <td class="p-3 text-slate-900">{{ s.id }}</td>
              <td class="p-3 text-slate-900">{{ s.real_name || '-' }}</td>
              <td class="p-3 text-slate-600">{{ s.email || '-' }}</td>
              <td class="p-3 text-slate-600">{{ s.role || '-' }}</td>
              <td class="p-3">
                <AdminStatusTag :type="s.status === 1 ? 'success' : 'warning'" :text="s.status === 1 ? '启用' : '禁用'" />
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </AdminDataCard>
  </div>
</template>
