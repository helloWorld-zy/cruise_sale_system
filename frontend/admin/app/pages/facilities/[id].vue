<template>
  <div class="admin-page">
    <AdminPageHeader title="编辑设施" subtitle="维护设施属性、收费规则与目标人群，变更后即时影响前台展示。" />
    <AdminFormCard title="编辑设施">
      <p v-if="loading" class="mb-3 text-sm text-slate-600">加载中...</p>
      <p v-else-if="empty" data-test="empty" class="mb-3 text-sm text-slate-600">暂无设施数据</p>
      <form class="admin-cruise-form" @submit.prevent="handleSubmit">
        <section class="admin-cruise-form__intro">
          <h2 class="admin-cruise-form__intro-title">设施信息维护</h2>
          <p class="admin-cruise-form__intro-desc">当前编辑对象 ID: #{{ id || '-' }}，请按真实业务信息更新，避免前台检索异常。</p>
        </section>

        <section class="admin-cruise-form__section">
          <h3 class="admin-cruise-form__section-title">基本信息</h3>
          <p class="admin-cruise-form__section-subtitle">更新所属邮轮、分类、名称和位置等核心属性。</p>
          <div class="admin-cruise-form__grid">
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">所属邮轮</span>
              <select v-model.number="form.cruise_id" class="admin-cruise-form__control">
                <option :value="0">请选择邮轮</option>
                <option v-for="cruise in cruises" :key="cruise.id" :value="Number(cruise.id)">{{ cruise.name || `邮轮 #${cruise.id}` }}</option>
              </select>
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">设施分类</span>
              <select v-model.number="form.category_id" class="admin-cruise-form__control">
                <option :value="0">请选择分类</option>
                <option v-for="cat in categories" :key="cat.id" :value="Number(cat.id)">{{ cat.name || `分类 #${cat.id}` }}</option>
              </select>
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">名称</span>
              <input v-model="form.name" class="admin-cruise-form__control" />
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">英文名</span>
              <input v-model="form.english_name" class="admin-cruise-form__control" />
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">位置</span>
              <input v-model="form.location" class="admin-cruise-form__control" />
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">开放时间</span>
              <input v-model="form.open_hours" class="admin-cruise-form__control" />
            </label>
          </div>
        </section>

        <section class="admin-cruise-form__section">
          <h3 class="admin-cruise-form__section-title">收费与人群</h3>
          <p class="admin-cruise-form__section-subtitle">收费控制用于结算，客群标签用于前台推荐。</p>
          <div class="admin-cruise-form__grid">
            <label class="admin-cruise-form__field facility-form-checkbox-row">
              <span class="admin-cruise-form__field-label">是否额外收费</span>
              <span class="facility-form-checkbox-wrap">
                <input v-model="form.extra_charge" type="checkbox" />
                <span>开启后请填写收费说明</span>
              </span>
            </label>
            <label v-if="form.extra_charge" class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">收费说明</span>
              <input v-model="form.charge_price_tip" class="admin-cruise-form__control" />
            </label>
            <div class="admin-cruise-form__field facility-form-audience-block">
              <span class="admin-cruise-form__field-label">适合人群</span>
              <div class="facility-form-audience-list">
                <label v-for="aud in audienceOptions" :key="aud" class="facility-form-audience-item">
                  <input type="checkbox" :checked="form.target_audience.includes(aud)" @change="toggleAudience(aud, ($event.target as HTMLInputElement).checked)" />
                  <span>{{ aud }}</span>
                </label>
              </div>
            </div>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">状态</span>
              <select v-model.number="form.status" class="admin-cruise-form__control">
                <option :value="1">开放</option>
                <option :value="0">关闭</option>
              </select>
            </label>
          </div>
        </section>

        <section class="admin-cruise-form__section">
          <h3 class="admin-cruise-form__section-title">描述</h3>
          <p class="admin-cruise-form__section-subtitle">请保持文本简洁，突出使用规则和服务亮点。</p>
          <label class="admin-cruise-form__field">
            <span class="admin-cruise-form__field-label">设施描述</span>
            <textarea v-model="form.description" rows="5" class="admin-cruise-form__control admin-cruise-form__control--textarea" />
          </label>
        </section>

        <AdminActionBar>
          <button type="button" class="admin-btn admin-btn--danger" :disabled="loading || deleting" @click="handleDelete">{{ deleting ? '删除中...' : '删除' }}</button>
          <AdminActionLink to="/facilities" class="admin-btn admin-btn--secondary">取消</AdminActionLink>
          <button type="submit" :disabled="loading" class="admin-btn">{{ loading ? '提交中...' : '保存' }}</button>
        </AdminActionBar>
        <p v-if="error" class="text-sm text-rose-500">{{ error }}</p>
      </form>

      <AdminConfirmDialog
        :visible="deleteDialogVisible"
        title="确认删除设施"
        :message="`确认删除设施 #${id} 吗？删除后不可恢复。`"
        :loading="deleting"
        loading-text="删除中..."
        @close="closeDeleteDialog"
        @confirm="confirmDelete"
      />
    </AdminFormCard>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import AdminConfirmDialog from '../../components/AdminConfirmDialog.vue'
import { useAdminDeleteDialog } from '../../composables/useAdminDeleteDialog'

const route = useRoute()
const { request } = useApi()
const loading = ref(false)
const error = ref<string | null>(null)
const empty = ref(false)
const cruises = ref<Record<string, any>[]>([])
const categories = ref<Record<string, any>[]>([])
const audienceOptions = ['儿童', '家庭', '情侣', '老年', '商务']
const {
  visible: deleteDialogVisible,
  submitting: deleting,
  open: openDeleteDialog,
  close: closeDeleteDialog,
  run: runDelete,
} = useAdminDeleteDialog()

const id = computed(() => {
  const value = Number(route.params.id)
  return Number.isFinite(value) && value > 0 ? value : 0
})

const form = ref({
  category_id: 0,
  cruise_id: 0,
  name: '',
  english_name: '',
  location: '',
  open_hours: '',
  extra_charge: false,
  charge_price_tip: '',
  target_audience: [] as string[],
  description: '',
  status: 1,
  sort_order: 0,
})

function splitCsv(raw: unknown) {
  if (typeof raw !== 'string') return []
  return raw.split(',').map((part) => part.trim()).filter(Boolean)
}

async function loadOptions() {
  try {
    const [cruiseRes, categoryRes] = await Promise.all([
      request('/cruises', { query: { page: 1, page_size: 100 } }),
      request('/facility-categories'),
    ])
    const cruisePayload = cruiseRes?.data ?? cruiseRes ?? {}
    cruises.value = Array.isArray(cruisePayload) ? cruisePayload : cruisePayload?.list ?? []
    const categoryPayload = categoryRes?.data ?? categoryRes ?? []
    categories.value = Array.isArray(categoryPayload) ? categoryPayload : categoryPayload?.list ?? []
  } catch {
    cruises.value = []
    categories.value = []
  }
}

async function loadDetail() {
  if (!id.value) return
  loading.value = true
  error.value = null
  empty.value = false
  try {
    const res = await request(`/facilities/${id.value}`)
    const payload = res?.data ?? res ?? {}
    if (Object.keys(payload).length === 0) {
      empty.value = true
      return
    }
    form.value = {
      category_id: Number(payload.category_id || 0),
      cruise_id: Number(payload.cruise_id || 0),
      name: payload.name || '',
      english_name: payload.english_name || '',
      location: payload.location || '',
      open_hours: payload.open_hours || '',
      extra_charge: Boolean(payload.extra_charge),
      charge_price_tip: payload.charge_price_tip || '',
      target_audience: splitCsv(payload.target_audience),
      description: payload.description || '',
      status: Number(payload.status ?? 1),
      sort_order: Number(payload.sort_order || 0),
    }
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load facility detail'
  } finally {
    loading.value = false
  }
}

function toggleAudience(value: string, checked: boolean) {
  const next = new Set(form.value.target_audience)
  if (checked) next.add(value)
  else next.delete(value)
  form.value.target_audience = Array.from(next)
}

async function handleSubmit() {
  if (!id.value || loading.value) return
  loading.value = true
  error.value = null
  try {
    await request(`/facilities/${id.value}`, {
      method: 'PUT',
      body: {
        category_id: Number(form.value.category_id),
        cruise_id: Number(form.value.cruise_id),
        name: form.value.name,
        english_name: form.value.english_name,
        location: form.value.location,
        open_hours: form.value.open_hours,
        extra_charge: Boolean(form.value.extra_charge),
        charge_price_tip: form.value.extra_charge ? form.value.charge_price_tip : '',
        target_audience: form.value.target_audience.join(','),
        description: form.value.description,
        status: Number(form.value.status),
        sort_order: Number(form.value.sort_order || 0),
      },
    })
    await navigateTo('/facilities')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to update facility'
  } finally {
    loading.value = false
  }
}

async function handleDelete() {
  if (!id.value || loading.value) return
  openDeleteDialog()
}

async function confirmDelete() {
  if (!id.value || loading.value || deleting.value) return
  error.value = null
  try {
    await runDelete(async () => {
      await request(`/facilities/${id.value}`, { method: 'DELETE' })
      await navigateTo('/facilities')
    })
  } catch (e: any) {
    error.value = e?.message ?? '删除设施失败，请稍后重试。'
  }
}

onMounted(async () => {
  await loadOptions()
  await loadDetail()
})
</script>

<style scoped>
.facility-form-checkbox-wrap {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-height: 40px;
  color: #334155;
}

.facility-form-audience-block {
  grid-column: 1 / -1;
}

.facility-form-audience-list {
  display: flex;
  flex-wrap: wrap;
  gap: 14px;
  padding: 10px 0 4px;
}

.facility-form-audience-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: #334155;
}
</style>

