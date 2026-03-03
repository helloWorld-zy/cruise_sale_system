<!-- admin/app/pages/cruises/create.vue -->
<script setup lang="ts">
import { ref } from 'vue'
import { useApi } from '../../composables/useApi'

const { request } = useApi()
const loading = ref(false)
const error = ref<string | null>(null)

const form = ref({
  name: '',
  english_name: '',
  code: '',
  company_id: 0,
  tonnage: 0,
  passenger_capacity: 0,
  crew_count: 0,
  build_year: 0,
  refurbish_year: 0,
  length: 0,
  width: 0,
  deck_count: 0,
  description: '',
  sort_order: 0,
  status: 1,
})

async function handleSubmit() {
  if (loading.value) return
  loading.value = true
  error.value = null
  try {
    await request('/cruises', {
      method: 'POST',
      body: {
        ...form.value,
        company_id: Number(form.value.company_id),
        tonnage: Number(form.value.tonnage),
        passenger_capacity: Number(form.value.passenger_capacity),
        crew_count: Number(form.value.crew_count),
        build_year: Number(form.value.build_year),
        refurbish_year: Number(form.value.refurbish_year),
        length: Number(form.value.length),
        width: Number(form.value.width),
        deck_count: Number(form.value.deck_count),
        sort_order: Number(form.value.sort_order),
        status: Number(form.value.status),
      },
    })
    await navigateTo('/cruises')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to create cruise'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-slate-50 p-4 md:p-6">
    <div class="mx-auto max-w-4xl rounded-lg border border-slate-200 bg-white p-4 shadow-sm">
      <h1 class="mb-4 text-xl font-semibold text-slate-900">新建邮轮</h1>
      <form class="space-y-4" @submit.prevent="handleSubmit">
        <div class="grid gap-4 md:grid-cols-2">
          <label class="text-sm text-slate-600">名称<input v-model="form.name" class="mt-1 h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" :disabled="loading" /></label>
          <label class="text-sm text-slate-600">英文名<input v-model="form.english_name" class="mt-1 h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" :disabled="loading" /></label>
          <label class="text-sm text-slate-600">代码<input v-model="form.code" class="mt-1 h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" :disabled="loading" /></label>
          <label class="text-sm text-slate-600">公司 ID<input v-model.number="form.company_id" type="number" min="1" class="mt-1 h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" :disabled="loading" /></label>
          <label class="text-sm text-slate-600">吨位<input v-model.number="form.tonnage" type="number" class="mt-1 h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" :disabled="loading" /></label>
          <label class="text-sm text-slate-600">载客量<input v-model.number="form.passenger_capacity" type="number" class="mt-1 h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" :disabled="loading" /></label>
          <label class="text-sm text-slate-600">船员数<input v-model.number="form.crew_count" type="number" class="mt-1 h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" :disabled="loading" /></label>
          <label class="text-sm text-slate-600">建造年份<input v-model.number="form.build_year" type="number" class="mt-1 h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" :disabled="loading" /></label>
          <label class="text-sm text-slate-600">翻新年份<input v-model.number="form.refurbish_year" type="number" class="mt-1 h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" :disabled="loading" /></label>
          <label class="text-sm text-slate-600">长度(m)<input v-model.number="form.length" type="number" class="mt-1 h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" :disabled="loading" /></label>
          <label class="text-sm text-slate-600">宽度(m)<input v-model.number="form.width" type="number" class="mt-1 h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" :disabled="loading" /></label>
          <label class="text-sm text-slate-600">甲板数<input v-model.number="form.deck_count" type="number" class="mt-1 h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" :disabled="loading" /></label>
          <label class="text-sm text-slate-600">排序<input v-model.number="form.sort_order" type="number" class="mt-1 h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" :disabled="loading" /></label>
        </div>
        <label class="block text-sm text-slate-600">描述<textarea v-model="form.description" class="mt-1 min-h-[180px] w-full rounded-md border border-slate-200 bg-white px-3 py-2 outline-none ring-indigo-500 focus:ring-2" :disabled="loading" /></label>
        <label class="block text-sm text-slate-600">状态
          <select v-model.number="form.status" class="mt-1 h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" :disabled="loading">
            <option :value="1">上架</option>
            <option :value="2">维护中</option>
            <option :value="0">下架</option>
          </select>
        </label>
        <div class="rounded-lg border-2 border-dashed border-slate-300 p-4 text-sm text-slate-500">图片上传（占位，Task 16 后续接入拖拽与主图标识）</div>
        <p v-if="error" class="text-sm text-rose-500">{{ error }}</p>
        <button type="submit" class="rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-500" :disabled="loading">{{ loading ? '提交中...' : '创建' }}</button>
      </form>
    </div>
  </div>
</template>

