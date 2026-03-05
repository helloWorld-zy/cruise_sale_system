<template>
  <div class="min-h-screen bg-slate-50 p-4 md:p-6">
    <div class="mx-auto max-w-3xl rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
      <h1 class="mb-6 text-xl font-semibold text-slate-900">编辑邮轮公司</h1>
      <p v-if="loading" class="mb-3 text-sm text-slate-600">加载中...</p>
      <p v-else-if="empty" data-test="empty" class="mb-3 text-sm text-slate-600">暂无公司数据</p>
      <form class="space-y-4" @submit.prevent="handleSubmit">
        <label class="block space-y-1 text-sm text-slate-600">
          <span>中文名</span>
          <input v-model="form.name" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
        </label>
        <label class="block space-y-1 text-sm text-slate-600">
          <span>英文名</span>
          <input v-model="form.english_name" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
        </label>
        <label class="block space-y-1 text-sm text-slate-600">
          <span>Logo 上传</span>
          <input type="file" accept="image/*" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" @change="onLogoFileChange" />
        </label>
        <div>
          <div v-if="form.logo_url" class="company-logo-preview-frame">
            <img :src="form.logo_url" alt="logo-preview" class="company-logo-preview-image" />
          </div>
          <button v-if="form.logo_url" type="button" class="ml-3 text-xs text-slate-500 hover:text-slate-700" @click="form.logo_url = ''">清除 Logo</button>
        </div>
        <label class="block space-y-1 text-sm text-slate-600">
          <span>文字介绍</span>
          <textarea v-model="form.description" class="min-h-[120px] w-full rounded-md border border-slate-200 px-3 py-2 outline-none ring-indigo-500 focus:ring-2" />
        </label>

        <div class="flex items-center justify-end gap-3 border-t border-slate-200 pt-4">
          <AdminActionLink to="/companies">取消</AdminActionLink>
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
const empty = ref(false)

const id = computed(() => {
  const value = Number(route.params.id)
  return Number.isFinite(value) && value > 0 ? value : 0
})

const form = ref({
  name: '',
  english_name: '',
  logo_url: '',
  description: '',
})

async function loadDetail() {
  loading.value = true
  error.value = null
  empty.value = false
  try {
    const res = await request('/companies')
    const payload = res?.data ?? res ?? {}
    const list = Array.isArray(payload) ? payload : payload?.list ?? []
    const detail = list.find((item: Record<string, any>) => Number(item.id ?? item.ID) === id.value)
    if (!detail) {
      empty.value = true
      return
    }
    form.value = {
      name: detail.name || detail.Name || '',
      english_name: detail.english_name || detail.EnglishName || '',
      logo_url: detail.logo_url || detail.LogoURL || '',
      description: detail.description || detail.Description || '',
    }
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load company'
  } finally {
    loading.value = false
  }
}

async function handleSubmit() {
  if (!id.value || loading.value) return
  loading.value = true
  error.value = null
  try {
    await request(`/companies/${id.value}`, {
      method: 'PUT',
      body: {
        name: form.value.name,
        english_name: form.value.english_name,
        logo_url: form.value.logo_url,
        description: form.value.description,
      },
    })
    await navigateTo('/companies')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to update company'
  } finally {
    loading.value = false
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

onMounted(loadDetail)
</script>

<style scoped>
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
