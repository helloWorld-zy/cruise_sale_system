<template>
  <div class="admin-page">
    <AdminPageHeader title="编辑邮轮公司" />
    <AdminFormCard title="公司资料维护">
      <h1 class="mb-3 text-xl font-semibold text-slate-900">编辑邮轮公司</h1>
      <p v-if="loading" class="mb-3 text-sm text-slate-600">加载中...</p>
      <p v-else-if="empty" data-test="empty" class="mb-3 text-sm text-slate-600">暂无公司数据</p>
      <form class="admin-cruise-form" @submit.prevent="handleSubmit">
        <section class="admin-cruise-form__intro">
          <h2 class="admin-cruise-form__intro-title">更新公司信息</h2>
          <p class="admin-cruise-form__intro-desc">公司信息会关联该公司下的全部邮轮展示和筛选结果，请保持命名一致性。</p>
        </section>

        <section class="admin-cruise-form__section">
          <h3 class="admin-cruise-form__section-title">基础信息</h3>
          <div class="admin-cruise-form__grid">
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label"><span>中文名</span><span class="admin-cruise-form__field-hint">必填</span></span>
              <input v-model="form.name" :class="['admin-cruise-form__control', fieldErrors.name && 'admin-cruise-form__control--error']" />
              <p v-if="fieldErrors.name" class="admin-form-error-text">{{ fieldErrors.name }}</p>
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label"><span>英文名</span><span class="admin-cruise-form__field-hint">选填</span></span>
              <input v-model="form.english_name" class="admin-cruise-form__control" />
            </label>
          </div>

          <label class="admin-cruise-form__field">
            <span class="admin-cruise-form__field-label"><span>Logo 上传</span><span class="admin-cruise-form__field-hint">支持图片文件</span></span>
            <input type="file" accept="image/*" class="admin-cruise-form__control" @change="onLogoFileChange" />
          </label>
          <div>
            <div v-if="form.logo_url" class="company-logo-preview-frame">
              <img :src="form.logo_url" alt="logo-preview" class="company-logo-preview-image" />
            </div>
            <button v-if="form.logo_url" type="button" class="ml-3 text-xs text-slate-500 hover:text-slate-700" @click="form.logo_url = ''">清除 Logo</button>
          </div>

          <label class="admin-cruise-form__field">
            <span class="admin-cruise-form__field-label"><span>文字介绍</span><span class="admin-cruise-form__field-hint">选填</span></span>
            <textarea v-model="form.description" class="admin-cruise-form__control admin-cruise-form__control--textarea" />
          </label>
        </section>

        <div class="admin-cruise-form__actions">
          <AdminActionLink to="/companies" class="admin-btn admin-btn--secondary">取消</AdminActionLink>
          <button type="submit" class="admin-btn" :disabled="loading">{{ loading ? '提交中...' : '保存' }}</button>
        </div>
        <p v-if="error" class="text-sm text-rose-500">{{ error }}</p>
      </form>
    </AdminFormCard>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

const route = useRoute()
const { request } = useApi()
const loading = ref(false)
const error = ref<string | null>(null)
const fieldErrors = ref<Record<string, string>>({})
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

function validateForm() {
  const nextErrors: Record<string, string> = {}
  if (!String(form.value.name || '').trim()) {
    nextErrors.name = '请填写公司中文名'
  }
  fieldErrors.value = nextErrors
  return Object.keys(nextErrors).length === 0
}

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
  if (!validateForm()) {
    error.value = '请先修正表单校验错误'
    return
  }
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
