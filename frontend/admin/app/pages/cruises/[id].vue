<!-- admin/app/pages/cruises/[id].vue -->
<script setup lang="ts">
import { onMounted, ref } from 'vue'
import AdminConfirmDialog from '../../components/AdminConfirmDialog.vue'
import AdminCompanySelect from '../../components/AdminCompanySelect.vue'
import { useAdminDeleteDialog } from '../../composables/useAdminDeleteDialog'

declare const useApi: any
declare const navigateTo: any

const route = useRoute()
const { request } = useApi()

const loading = ref(false)
const saving = ref(false)
const error = ref<string | null>(null)
const fieldErrors = ref<Record<string, string>>({})
const empty = ref(false)
const companies = ref<Array<{ id: number; name: string; english_name?: string; logo_url?: string }>>([])
const {
  visible: deleteDialogVisible,
  submitting: deleting,
  open: openDeleteDialog,
  close: closeDeleteDialog,
  run: runDelete,
} = useAdminDeleteDialog()
const form = ref({
  id: 0,
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

const id = Number(route.params.id)

function validateForm() {
  const nextErrors: Record<string, string> = {}
  if (!String(form.value.name || '').trim()) {
    nextErrors.name = '请填写邮轮名称'
  }
  if (!Number.isFinite(Number(form.value.company_id)) || Number(form.value.company_id) <= 0) {
    nextErrors.company_id = '请选择所属公司'
  }

  const buildYear = Number(form.value.build_year)
  const refurbishYear = Number(form.value.refurbish_year)
  if (buildYear > 0 && (buildYear < 1900 || buildYear > 2100)) {
    nextErrors.build_year = '建造年份需在 1900-2100 之间'
  }
  if (refurbishYear > 0 && (refurbishYear < 1900 || refurbishYear > 2100)) {
    nextErrors.refurbish_year = '翻新年份需在 1900-2100 之间'
  }
  if (buildYear > 0 && refurbishYear > 0 && refurbishYear < buildYear) {
    nextErrors.refurbish_year = '翻新年份不能小于建造年份'
  }

  fieldErrors.value = nextErrors
  return Object.keys(nextErrors).length === 0
}

async function loadDetail() {
  loading.value = true
  error.value = null
  empty.value = false
  try {
    const res = await request(`/cruises/${id}`)
    const data = res?.data ?? res ?? {}
    if (Object.keys(data).length === 0) {
      empty.value = true
      return
    }
    form.value = {
      id: Number(data.id ?? id),
      name: data.name ?? '',
      english_name: data.english_name ?? '',
      code: data.code ?? '',
      company_id: Number(data.company_id ?? 0),
      tonnage: Number(data.tonnage ?? 0),
      passenger_capacity: Number(data.passenger_capacity ?? 0),
      crew_count: Number(data.crew_count ?? 0),
      build_year: Number(data.build_year ?? 0),
      refurbish_year: Number(data.refurbish_year ?? 0),
      length: Number(data.length ?? 0),
      width: Number(data.width ?? 0),
      deck_count: Number(data.deck_count ?? 0),
      description: data.description ?? '',
      sort_order: Number(data.sort_order ?? 0),
      status: Number(data.status ?? 1),
    }
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load cruise detail'
  } finally {
    loading.value = false
  }
}

async function loadCompanies() {
  try {
    const res = await request('/companies')
    const payload = res?.data ?? res ?? {}
    companies.value = (Array.isArray(payload) ? payload : payload?.list ?? []).map((item: any) => ({
      id: Number(item.id),
      name: item.name || '',
      english_name: item.english_name || '',
      logo_url: item.logo_url || '',
    }))
  } catch {
    companies.value = []
  }
}

async function handleSave() {
  if (saving.value) return
  if (!validateForm()) {
    error.value = '请先修正表单校验错误'
    return
  }
  saving.value = true
  error.value = null
  try {
    await request(`/cruises/${id}`, {
      method: 'PUT',
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
    error.value = e?.message ?? 'failed to update cruise'
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  if (deleting.value) return
  openDeleteDialog()
}

async function confirmDelete() {
  if (deleting.value) return
  error.value = null
  try {
    await runDelete(async () => {
      await request(`/cruises/${id}`, { method: 'DELETE' })
      await navigateTo('/cruises')
    })
  } catch (e: any) {
    const code = Number(e?.code ?? 0)
    const status = Number(e?.status ?? 0)
    const message = String(e?.message ?? '')
    if (code === 42204 || (status === 409 && message.includes('cruise has voyages'))) {
      error.value = '删除失败：该邮轮下存在航次，请先处理关联航次后再删除。'
    } else if (code === 42201 || (status === 409 && message.includes('cruise has cabins'))) {
      error.value = '删除失败：该邮轮下存在舱房类型，请先处理关联舱房后再删除。'
    } else {
      error.value = e?.message ?? '删除邮轮失败，请稍后重试。'
    }
  }
}

onMounted(async () => {
  await Promise.all([loadDetail(), loadCompanies()])
})
</script>

<template>
  <div class="admin-page">
    <AdminPageHeader :title="`编辑邮轮 #${id}`" />
    <AdminFormCard title="邮轮资料维护">
      <h1 class="mb-3 text-xl font-semibold text-slate-900">编辑邮轮 #{{ id }}</h1>
      <p v-if="loading" class="text-sm text-slate-600">加载中...</p>
      <p v-else-if="empty" data-test="empty" class="text-sm text-slate-600">暂无邮轮数据</p>
      <form v-else class="admin-cruise-form" @submit.prevent="handleSave">
        <section class="admin-cruise-form__intro">
          <h2 class="admin-cruise-form__intro-title">维护邮轮档案</h2>
          <p class="admin-cruise-form__intro-desc">更新后将影响航次配置、舱型管理与前台展示，请确认关键数据准确后再保存。</p>
        </section>

        <section class="admin-cruise-form__section">
          <h3 class="admin-cruise-form__section-title">识别信息</h3>
          <p class="admin-cruise-form__section-subtitle">用于后台检索、筛选与对外展示。</p>
          <div class="admin-cruise-form__grid">
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label"><span>名称</span><span class="admin-cruise-form__field-hint">必填</span></span>
              <input v-model="form.name" :class="['admin-cruise-form__control', fieldErrors.name && 'admin-cruise-form__control--error']" :disabled="saving || deleting" />
              <p v-if="fieldErrors.name" class="admin-form-error-text">{{ fieldErrors.name }}</p>
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label"><span>英文名</span><span class="admin-cruise-form__field-hint">选填</span></span>
              <input v-model="form.english_name" class="admin-cruise-form__control" :disabled="saving || deleting" />
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label"><span>代码</span><span class="admin-cruise-form__field-hint">建议唯一</span></span>
              <input v-model="form.code" class="admin-cruise-form__control" :disabled="saving || deleting" />
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label"><span>所属公司</span><span class="admin-cruise-form__field-hint">必选</span></span>
              <AdminCompanySelect v-model="form.company_id" :options="companies" :disabled="saving || deleting" placeholder="请选择公司" />
              <p v-if="fieldErrors.company_id" class="admin-form-error-text">{{ fieldErrors.company_id }}</p>
            </label>
          </div>
        </section>

        <section class="admin-cruise-form__section">
          <h3 class="admin-cruise-form__section-title">规格参数</h3>
          <p class="admin-cruise-form__section-subtitle">按实际船舶参数填写，便于后续运营统计与展示。</p>
          <div class="admin-cruise-form__grid">
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label"><span>吨位</span><span class="admin-cruise-form__field-hint">单位: 吨</span></span>
              <input v-model.number="form.tonnage" type="number" class="admin-cruise-form__control" :disabled="saving || deleting" />
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label"><span>载客量</span><span class="admin-cruise-form__field-hint">人数</span></span>
              <input v-model.number="form.passenger_capacity" type="number" class="admin-cruise-form__control" :disabled="saving || deleting" />
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label"><span>船员数</span><span class="admin-cruise-form__field-hint">人数</span></span>
              <input v-model.number="form.crew_count" type="number" class="admin-cruise-form__control" :disabled="saving || deleting" />
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label"><span>建造年份</span><span class="admin-cruise-form__field-hint">YYYY</span></span>
              <input v-model.number="form.build_year" type="number" :class="['admin-cruise-form__control', fieldErrors.build_year && 'admin-cruise-form__control--error']" :disabled="saving || deleting" />
              <p v-if="fieldErrors.build_year" class="admin-form-error-text">{{ fieldErrors.build_year }}</p>
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label"><span>翻新年份</span><span class="admin-cruise-form__field-hint">YYYY</span></span>
              <input v-model.number="form.refurbish_year" type="number" :class="['admin-cruise-form__control', fieldErrors.refurbish_year && 'admin-cruise-form__control--error']" :disabled="saving || deleting" />
              <p v-if="fieldErrors.refurbish_year" class="admin-form-error-text">{{ fieldErrors.refurbish_year }}</p>
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label"><span>长度(m)</span><span class="admin-cruise-form__field-hint">船长</span></span>
              <input v-model.number="form.length" type="number" class="admin-cruise-form__control" :disabled="saving || deleting" />
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label"><span>宽度(m)</span><span class="admin-cruise-form__field-hint">船宽</span></span>
              <input v-model.number="form.width" type="number" class="admin-cruise-form__control" :disabled="saving || deleting" />
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label"><span>甲板数</span><span class="admin-cruise-form__field-hint">整数</span></span>
              <input v-model.number="form.deck_count" type="number" class="admin-cruise-form__control" :disabled="saving || deleting" />
            </label>
          </div>
        </section>

        <section class="admin-cruise-form__section">
          <h3 class="admin-cruise-form__section-title">运营配置</h3>
          <p class="admin-cruise-form__section-subtitle">控制排序和上架状态，描述内容用于前台展示。</p>
          <div class="admin-cruise-form__grid">
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label"><span>排序</span><span class="admin-cruise-form__field-hint">数字越小越靠前</span></span>
              <input v-model.number="form.sort_order" type="number" class="admin-cruise-form__control" :disabled="saving || deleting" />
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label"><span>状态</span><span class="admin-cruise-form__field-hint">影响前台可见性</span></span>
              <select v-model.number="form.status" class="admin-cruise-form__control" :disabled="saving || deleting">
                <option :value="1">上架</option>
                <option :value="2">维护中</option>
                <option :value="0">下架</option>
              </select>
            </label>
          </div>
          <label class="admin-cruise-form__field">
            <span class="admin-cruise-form__field-label"><span>描述</span><span class="admin-cruise-form__field-hint">建议填写亮点与定位</span></span>
            <textarea v-model="form.description" class="admin-cruise-form__control admin-cruise-form__control--textarea" :disabled="saving || deleting" />
          </label>
        </section>

        <div class="admin-cruise-form__upload">图片上传（占位，Task 16 后续接入拖拽与主图标识）</div>
        <p v-if="error" class="text-sm text-rose-500">{{ error }}</p>
        <div class="admin-cruise-form__actions">
          <button type="button" class="admin-btn admin-btn--secondary" :disabled="saving || deleting" @click="navigateTo('/cruises')">返回列表</button>
          <button type="submit" class="admin-btn" :disabled="saving || deleting">{{ saving ? '保存中...' : '保存' }}</button>
          <button type="button" class="admin-btn" style="background: var(--admin-color-danger)" :disabled="saving || deleting" @click="handleDelete">{{ deleting ? '删除中...' : '删除' }}</button>
        </div>
      </form>

      <AdminConfirmDialog
        :visible="deleteDialogVisible"
        title="确认删除邮轮"
        :message="`确认删除邮轮 #${id} 吗？删除后不可恢复。`"
        :loading="deleting"
        loading-text="删除中..."
        @close="closeDeleteDialog"
        @confirm="confirmDelete"
      />
    </AdminFormCard>
  </div>
</template>
