<template>
  <div class="admin-page">
    <AdminPageHeader title="邮轮公司管理" />

    <AdminFormCard title="新增公司">
      <form class="admin-cruise-form" @submit.prevent="createCompany">
        <section class="admin-cruise-form__intro">
          <h2 class="admin-cruise-form__intro-title">创建邮轮公司</h2>
          <p class="admin-cruise-form__intro-desc">公司信息将用于邮轮归属关系、筛选条件及前台品牌展示。</p>
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
            <button type="submit" class="admin-btn" :disabled="submitting">
              {{ submitting ? '提交中...' : '新增公司' }}
            </button>
        </div>
      </form>
    </AdminFormCard>

    <AdminDataCard flush>
      <div class="company-table-wrap overflow-x-auto">
          <table class="w-full min-w-[900px] text-sm">
          <thead class="bg-slate-50 text-left text-slate-600">
            <tr>
              <th class="p-3">Logo</th>
              <th class="p-3">中文名</th>
              <th class="p-3">英文名</th>
              <th class="p-3">文字介绍</th>
              <th class="p-3 whitespace-nowrap">操作</th>
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
              <td class="p-3 whitespace-nowrap">
                <div class="company-actions flex items-center gap-3">
                  <AdminActionLink :to="`/companies/${item.id}`" class="admin-table-action-btn">编辑</AdminActionLink>
                  <button type="button" class="admin-table-action-btn admin-table-action-btn--danger" @click="askRemoveCompany(item)">删除</button>
                </div>
              </td>
            </tr>
          </tbody>
          </table>
      </div>
    </AdminDataCard>

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
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import AdminConfirmDialog from '../../components/AdminConfirmDialog.vue'
import { useAdminDeleteDialog } from '../../composables/useAdminDeleteDialog'

const { request } = useApi()
const loading = ref(false)
const submitting = ref(false)
const error = ref<string | null>(null)
const fieldErrors = ref<Record<string, string>>({})
const items = ref<any[]>([])
const normalizedItems = ref<Array<{ id: number; name: string; english_name: string; logo_url: string; description: string }>>([])
const {
  visible: deleteDialogVisible,
  submitting: deleteSubmitting,
  target: deleteTarget,
  open: openDeleteDialog,
  close: closeDeleteDialog,
  run: runDelete,
} = useAdminDeleteDialog<{ id: number; name: string }>()

const form = ref({
  name: '',
  english_name: '',
  logo_url: '',
  description: '',
})

function validateCreateForm() {
  const nextErrors: Record<string, string> = {}
  if (!String(form.value.name || '').trim()) {
    nextErrors.name = '请填写公司中文名'
  }
  fieldErrors.value = nextErrors
  return Object.keys(nextErrors).length === 0
}

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
  if (!validateCreateForm()) {
    error.value = '请先修正表单校验错误'
    return
  }
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
    fieldErrors.value = {}
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
  openDeleteDialog({ id: targetId, name: item?.name || '' })
}

async function confirmRemoveCompany() {
  const targetId = Number(deleteTarget.value?.id ?? 0)
  if (!Number.isFinite(targetId) || targetId <= 0) {
    error.value = '无效公司 ID，无法删除'
    closeDeleteDialog()
    return
  }
  error.value = null
  try {
    await runDelete(async () => {
      await request(`/companies/${targetId}`, { method: 'DELETE' })
      normalizedItems.value = normalizedItems.value.filter((item) => Number(item.id) !== targetId)
      items.value = items.value.filter((raw) => Number(raw?.id ?? raw?.ID ?? 0) !== targetId)
      await loadItems()
    })
  } catch (e: any) {
    const code = Number(e?.code ?? 0)
    const status = Number(e?.status ?? 0)
    if (code === 42202 || status === 409) {
      error.value = '删除失败：该公司下存在邮轮，请先处理关联邮轮后再删除。'
    } else {
      error.value = '删除失败，请稍后重试。'
    }
  }
}

onMounted(loadItems)
</script>

<style scoped>
.company-table-wrap {
  scrollbar-gutter: stable;
}

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
