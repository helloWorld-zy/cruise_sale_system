<template>
  <div class="min-h-screen bg-slate-50 p-4 md:p-6">
    <div class="mx-auto max-w-7xl">
      <div class="mb-4 flex items-center justify-between">
        <h1 class="text-xl font-semibold text-slate-900">邮轮公司管理</h1>
      </div>

      <div class="mb-4 rounded-lg border border-slate-200 bg-white p-4 shadow-sm">
        <form class="grid gap-3 md:grid-cols-2" @submit.prevent="createCompany">
          <label class="text-sm text-slate-600">中文名
            <input v-model="form.name" class="mt-1 h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
          </label>
          <label class="text-sm text-slate-600">英文名
            <input v-model="form.english_name" class="mt-1 h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
          </label>
          <label class="text-sm text-slate-600 md:col-span-2">Logo 上传
            <input type="file" accept="image/*" class="mt-1 h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" @change="onLogoFileChange" />
          </label>
          <div class="md:col-span-2">
            <div v-if="form.logo_url" class="company-logo-preview-frame">
              <img :src="form.logo_url" alt="logo-preview" class="company-logo-preview-image" />
            </div>
            <button v-if="form.logo_url" type="button" class="ml-3 text-xs text-slate-500 hover:text-slate-700" @click="form.logo_url = ''">清除 Logo</button>
          </div>
          <label class="text-sm text-slate-600 md:col-span-2">文字介绍
            <textarea v-model="form.description" class="mt-1 min-h-[100px] w-full rounded-md border border-slate-200 px-3 py-2 outline-none ring-indigo-500 focus:ring-2" />
          </label>
          <div class="md:col-span-2">
            <button type="submit" class="rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-500" :disabled="submitting">
              {{ submitting ? '提交中...' : '新增公司' }}
            </button>
          </div>
        </form>
      </div>

      <div class="rounded-lg border border-slate-200 bg-white shadow-sm">
        <table class="w-full text-sm">
          <thead class="bg-slate-50 text-left text-slate-600">
            <tr>
              <th class="p-3">Logo</th>
              <th class="p-3">中文名</th>
              <th class="p-3">英文名</th>
              <th class="p-3">文字介绍</th>
              <th class="p-3">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading"><td class="p-3" colspan="5">加载中...</td></tr>
            <tr v-else-if="error"><td class="p-3 text-rose-500" colspan="5">{{ error }}</td></tr>
            <tr v-else-if="items.length === 0"><td class="p-3" colspan="5">暂无数据</td></tr>
            <tr v-for="item in normalizedItems" v-else :key="item.id">
              <td class="p-3">
                <div
                  v-if="item.logo_url"
                  class="company-logo-cell"
                >
                  <img :src="item.logo_url" alt="logo" class="company-logo-image" />
                </div>
                <span v-else class="text-slate-400">-</span>
              </td>
              <td class="p-3 text-slate-900">{{ item.name || '-' }}</td>
              <td class="p-3 text-slate-600">{{ item.english_name || '-' }}</td>
              <td class="p-3 text-slate-600">{{ item.description || '-' }}</td>
              <td class="p-3">
                <AdminActionLink :to="`/companies/${item.id}`">编辑</AdminActionLink>
                <button type="button" class="ml-3 text-rose-500 hover:text-rose-400" @click="askRemoveCompany(item)">删除</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <AdminConfirmDialog
        :visible="deleteDialogVisible"
        title="确认删除公司"
        :message="`确认删除公司「${deleteTarget?.name || `#${deleteTarget?.id ?? ''}`}」吗？删除后不可恢复。`"
        hint="若该公司下存在邮轮，将无法删除。"
        :loading="deleteSubmitting"
        loading-text="删除中..."
        @close="closeDeleteDialog"
        @confirm="confirmRemoveCompany"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import AdminConfirmDialog from '../../components/AdminConfirmDialog.vue'

const { request } = useApi()
const loading = ref(false)
const submitting = ref(false)
const deleteSubmitting = ref(false)
const error = ref<string | null>(null)
const items = ref<any[]>([])
const normalizedItems = ref<Array<{ id: number; name: string; english_name: string; logo_url: string; description: string }>>([])
const deleteDialogVisible = ref(false)
const deleteTarget = ref<{ id: number; name: string } | null>(null)

const form = ref({
  name: '',
  english_name: '',
  logo_url: '',
  description: '',
})

async function loadItems() {
  loading.value = true
  error.value = null
  try {
    const res = await request('/companies')
    const payload = res?.data ?? res ?? {}
    items.value = Array.isArray(payload) ? payload : payload?.list ?? []
    normalizedItems.value = items.value.map((raw: Record<string, any>) => ({
      id: Number(raw.id ?? raw.ID ?? 0),
      name: String(raw.name ?? raw.Name ?? ''),
      english_name: String(raw.english_name ?? raw.EnglishName ?? ''),
      logo_url: String(raw.logo_url ?? raw.LogoURL ?? ''),
      description: String(raw.description ?? raw.Description ?? ''),
    }))
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load companies'
    normalizedItems.value = []
  } finally {
    loading.value = false
  }
}

async function createCompany() {
  if (submitting.value) return
  submitting.value = true
  error.value = null
  try {
    await request('/companies', {
      method: 'POST',
      body: {
        name: form.value.name,
        english_name: form.value.english_name,
        logo_url: form.value.logo_url,
        description: form.value.description,
      },
    })
    form.value = { name: '', english_name: '', logo_url: '', description: '' }
    await loadItems()
  } catch (e: any) {
    error.value = e?.message ?? 'failed to create company'
  } finally {
    submitting.value = false
  }
}

async function onLogoFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input?.files?.[0]
  if (!file) return
  if (!file.type.startsWith('image/')) {
    error.value = '请上传图片文件'
    return
  }
  const formData = new FormData()
  formData.append('file', file)
  try {
    const res = await request('/upload/image', {
      method: 'POST',
      body: formData,
    })
    const payload = res?.data ?? res ?? {}
    const url = payload?.url || payload?.data?.url || ''
    if (!url) {
      throw new Error('上传成功但未返回 URL')
    }
    form.value.logo_url = String(url)
  } catch (e: any) {
    error.value = e?.message ?? '上传失败'
  } finally {
    if (input) input.value = ''
  }
}

function askRemoveCompany(item: { id: number; name: string }) {
  const targetId = Number(item?.id)
  if (!Number.isFinite(targetId) || targetId <= 0) {
    error.value = '无效公司 ID，无法删除'
    return
  }
  deleteTarget.value = { id: targetId, name: item?.name || '' }
  deleteDialogVisible.value = true
}

function closeDeleteDialog() {
  if (deleteSubmitting.value) return
  deleteDialogVisible.value = false
  deleteTarget.value = null
}

async function confirmRemoveCompany() {
  const targetId = Number(deleteTarget.value?.id ?? 0)
  if (!Number.isFinite(targetId) || targetId <= 0) {
    error.value = '无效公司 ID，无法删除'
    closeDeleteDialog()
    return
  }
  error.value = null
  deleteSubmitting.value = true
  try {
    await request(`/companies/${targetId}`, { method: 'DELETE' })
    // Optimistically remove deleted row to avoid "no response" feeling when list reload is delayed.
    normalizedItems.value = normalizedItems.value.filter((item) => Number(item.id) !== targetId)
    items.value = items.value.filter((raw) => Number(raw?.id ?? raw?.ID ?? 0) !== targetId)
    closeDeleteDialog()
    await loadItems()
  } catch (e: any) {
    const code = Number(e?.code ?? 0)
    const status = Number(e?.status ?? 0)
    if (code === 42202 || status === 409) {
      error.value = '删除失败：该公司下存在邮轮，请先处理关联邮轮后再删除。'
    } else {
      error.value = '删除失败，请稍后重试。'
    }
  } finally {
    deleteSubmitting.value = false
  }
}

onMounted(loadItems)
</script>

<style scoped>
.company-logo-cell {
  width: 5rem;
  height: 3rem;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0.25rem;
  border: 1px solid rgb(226 232 240);
  border-radius: 0.375rem;
  background: rgb(248 250 252);
  overflow: hidden;
}

.company-logo-image {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: contain;
}

.company-logo-preview-frame {
  width: 6rem;
  height: 4rem;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.25rem;
  border: 1px solid rgb(226 232 240);
  border-radius: 0.5rem;
  background: rgb(248 250 252);
  overflow: hidden;
}

.company-logo-preview-image {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: contain;
}
</style>
