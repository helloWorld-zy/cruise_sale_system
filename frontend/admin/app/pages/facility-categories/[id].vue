<template>
  <div class="min-h-screen bg-slate-50 p-4 md:p-6">
    <div class="mx-auto max-w-3xl rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
      <h1 class="mb-6 text-xl font-semibold text-slate-900">编辑设施分类</h1>
      <form class="space-y-4" @submit.prevent="handleSubmit">
        <label class="block space-y-1 text-sm text-slate-600">
          <span>名称</span>
          <input v-model="form.name" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
        </label>
        <label class="block space-y-1 text-sm text-slate-600">
          <span>图标</span>
          <input v-model="form.icon" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
        </label>
        <label class="block space-y-1 text-sm text-slate-600">
          <span>排序</span>
          <input v-model.number="form.sort_order" type="number" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
        </label>
        <label class="block space-y-1 text-sm text-slate-600">
          <span>状态</span>
          <select v-model.number="form.status" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2">
            <option :value="1">启用</option>
            <option :value="0">停用</option>
          </select>
        </label>

        <div class="flex items-center justify-end gap-3 border-t border-slate-200 pt-4">
          <button type="button" class="rounded-md border border-rose-200 px-4 py-2 text-sm text-rose-600 hover:bg-rose-50" :disabled="loading" @click="handleDelete">删除</button>
          <NuxtLink to="/facility-categories" class="rounded-md border border-slate-200 px-4 py-2 text-sm text-slate-700 hover:bg-slate-50">取消</NuxtLink>
          <button type="submit" class="rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-500" :disabled="loading">{{ loading ? '提交中...' : '保存' }}</button>
        </div>
        <p v-if="error" class="text-sm text-rose-500">{{ error }}</p>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

const route = useRoute()
const { request } = useApi()
const loading = ref(false)
const error = ref<string | null>(null)

const id = computed(() => {
  const value = Number(route.params.id)
  return Number.isFinite(value) && value > 0 ? value : 0
})

const form = ref({
  name: '',
  icon: '',
  sort_order: 0,
  status: 1,
})

async function loadDetail() {
  loading.value = true
  error.value = null
  try {
    const res = await request('/facility-categories')
    const payload = res?.data ?? res ?? []
    const list = Array.isArray(payload) ? payload : payload?.list ?? []
    const detail = list.find((item: Record<string, any>) => Number(item.id) === id.value)
    if (!detail) {
      error.value = '未找到该分类'
      return
    }
    form.value = {
      name: detail.name || '',
      icon: detail.icon || '',
      sort_order: Number(detail.sort_order || 0),
      status: Number(detail.status ?? 1),
    }
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load facility category'
  } finally {
    loading.value = false
  }
}

async function handleSubmit() {
  if (!id.value || loading.value) return
  loading.value = true
  error.value = null
  try {
    await request(`/facility-categories/${id.value}`, {
      method: 'PUT',
      body: {
        name: form.value.name,
        icon: form.value.icon,
        sort_order: Number(form.value.sort_order || 0),
        status: Number(form.value.status ?? 1),
      },
    })
    await navigateTo('/facility-categories')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to update facility category'
  } finally {
    loading.value = false
  }
}

async function handleDelete() {
  if (!id.value || loading.value) return
  if (!confirm(`确认删除分类 #${id.value} 吗？`)) return
  loading.value = true
  error.value = null
  try {
    await request(`/facility-categories/${id.value}`, { method: 'DELETE' })
    await navigateTo('/facility-categories')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to delete facility category'
  } finally {
    loading.value = false
  }
}

onMounted(loadDetail)
</script>
