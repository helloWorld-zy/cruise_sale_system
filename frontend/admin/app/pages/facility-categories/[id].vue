<template>
  <div class="admin-page">
    <AdminPageHeader title="编辑设施分类" subtitle="调整分类名称、图标与显示顺序，变更后将实时影响设施筛选与展示。" />
    <AdminFormCard title="编辑设施分类">
      <p v-if="loading" class="mb-3 text-sm text-slate-600">加载中...</p>
      <p v-else-if="empty" data-test="empty" class="mb-3 text-sm text-slate-600">暂无分类数据</p>
      <form class="admin-cruise-form max-w-3xl" @submit.prevent="handleSubmit">
        <section class="facility-category-meta-panel" aria-live="polite">
          <div class="facility-category-meta-panel__item">
            <span class="facility-category-meta-panel__label">分类 ID</span>
            <strong class="facility-category-meta-panel__value">#{{ id || '-' }}</strong>
          </div>
          <div class="facility-category-meta-panel__item">
            <span class="facility-category-meta-panel__label">当前状态</span>
            <strong class="facility-category-meta-panel__value">{{ Number(form.status) === 1 ? '启用' : '停用' }}</strong>
          </div>
          <p class="facility-category-meta-panel__hint">建议：直接按图标样式选择分类视觉，排序值越小越靠前。</p>
        </section>

        <section class="admin-cruise-form__section">
          <h2 class="admin-cruise-form__section-title">基础信息</h2>
          <p class="admin-cruise-form__section-subtitle">维护分类在前台的名称、图标和显示策略。</p>

          <div class="admin-cruise-form__grid">
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">名称</span>
              <input v-model="form.name" class="admin-cruise-form__control" placeholder="请输入分类名称" />
            </label>

            <div class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">图标</span>
              <div class="facility-category-icon-current-card" data-test="facility-category-edit-icon-current" :data-icon="form.icon || ''">
                <span class="facility-category-icon-current-card__preview">
                  <FacilityCategoryIcon :name="form.icon" />
                </span>
                <div class="facility-category-icon-current-card__text">
                  <strong>{{ getFacilityCategoryIconLabel(form.icon) }}</strong>
                  <span>{{ form.icon || '未选择图标' }}</span>
                </div>
              </div>
            </div>

            <div class="facility-category-icon-presets" role="group" aria-label="图标快捷选择">
              <button
                v-for="icon in iconOptions"
                :key="icon.value"
                type="button"
                class="facility-category-icon-presets__option"
                :class="form.icon === icon.value ? 'facility-category-icon-presets__option--active' : ''"
                :data-test="`facility-category-edit-icon-option-${icon.value}`"
                @click="selectIcon(icon.value)"
              >
                <FacilityCategoryIcon :name="icon.value" />
                <span>{{ icon.label }}</span>
              </button>
            </div>

            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">排序</span>
              <input v-model.number="form.sort_order" type="number" class="admin-cruise-form__control" />
            </label>

            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">状态</span>
              <select v-model.number="form.status" class="admin-cruise-form__control">
                <option :value="1">启用</option>
                <option :value="0">停用</option>
              </select>
            </label>
          </div>
        </section>

        <AdminActionBar>
          <button type="button" class="admin-btn admin-btn--danger" :disabled="loading || deleting" @click="handleDelete">{{ deleting ? '删除中...' : '删除' }}</button>
          <AdminActionLink to="/facility-categories" class="admin-btn admin-btn--secondary">取消</AdminActionLink>
          <button type="submit" class="admin-btn" :disabled="loading">{{ loading ? '提交中...' : '保存' }}</button>
        </AdminActionBar>
        <p v-if="error" class="text-sm text-rose-500">{{ error }}</p>
      </form>

      <AdminConfirmDialog
        :visible="deleteDialogVisible"
        title="确认删除设施分类"
        :message="`确认删除分类 #${id} 吗？删除后不可恢复。`"
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
import FacilityCategoryIcon from '../../components/facility-categories/FacilityCategoryIcon.vue'
import { useAdminDeleteDialog } from '../../composables/useAdminDeleteDialog'
import { ensureFacilityCategoryIconOptions, getFacilityCategoryIconLabel } from '../../constants/facilityCategoryIcons'

const route = useRoute()
const { request } = useApi()
const loading = ref(false)
const error = ref<string | null>(null)
const empty = ref(false)
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
  name: '',
  icon: '',
  sort_order: 0,
  status: 1,
})

const iconOptions = computed(() => ensureFacilityCategoryIconOptions(form.value.icon))

function selectIcon(icon: string) {
  form.value.icon = icon
}

async function loadDetail() {
  loading.value = true
  error.value = null
  empty.value = false
  try {
    const res = await request('/facility-categories')
    const payload = res?.data ?? res ?? []
    const list = Array.isArray(payload) ? payload : payload?.list ?? []
    const detail = list.find((item: Record<string, any>) => Number(item.id) === id.value)
    if (!detail) {
      empty.value = true
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
  openDeleteDialog()
}

async function confirmDelete() {
  if (!id.value || loading.value || deleting.value) return
  error.value = null
  try {
    await runDelete(async () => {
      await request(`/facility-categories/${id.value}`, { method: 'DELETE' })
      await navigateTo('/facility-categories')
    })
  } catch (e: any) {
    error.value = e?.message ?? '删除设施分类失败，请稍后重试。'
  }
}

onMounted(loadDetail)
</script>

<style scoped>
.facility-category-meta-panel {
  display: grid;
  gap: 10px;
  grid-template-columns: repeat(1, minmax(0, 1fr));
  padding: 14px;
  border: 1px solid #d9ecff;
  background: linear-gradient(135deg, #f5faff 0%, #edf6ff 100%);
  border-radius: 8px;
}

.facility-category-meta-panel__item {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.facility-category-meta-panel__label {
  color: #64748b;
  font-size: 13px;
}

.facility-category-meta-panel__value {
  color: #1e293b;
  font-size: 14px;
}

.facility-category-meta-panel__hint {
  margin: 0;
  color: #475569;
  font-size: 13px;
  line-height: 1.6;
}

.facility-category-icon-presets {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(110px, 1fr));
  gap: 8px;
  margin-top: -2px;
}

.facility-category-icon-current-card {
  min-height: 72px;
  display: grid;
  grid-template-columns: 56px minmax(0, 1fr);
  gap: 12px;
  align-items: center;
  padding: 12px 14px;
  border: 1px solid #dbe4f0;
  border-radius: 10px;
  background: linear-gradient(135deg, #f8fbff 0%, #eef6ff 100%);
}

.facility-category-icon-current-card__preview {
  width: 56px;
  height: 56px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 14px;
  background: #dbeafe;
  color: #1d4ed8;
}

.facility-category-icon-current-card__text {
  display: grid;
  gap: 4px;
}

.facility-category-icon-current-card__text strong {
  color: #0f172a;
  font-size: 14px;
}

.facility-category-icon-current-card__text span {
  color: #64748b;
  font-size: 12px;
}

.facility-category-icon-presets__option {
  min-height: 84px;
  border: 1px solid #dbe4f0;
  border-radius: 10px;
  background: #fff;
  color: #334155;
  padding: 10px 12px;
  font-size: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  cursor: pointer;
  transition: border-color 0.2s ease, background-color 0.2s ease, color 0.2s ease, box-shadow 0.2s ease;
}

.facility-category-icon-presets__option:hover,
.facility-category-icon-presets__option:focus-visible {
  border-color: #93c5fd;
  background: #eff6ff;
  color: #1d4ed8;
  outline: none;
}

.facility-category-icon-presets__option--active {
  border-color: #60a5fa;
  background: #dbeafe;
  color: #1e40af;
  box-shadow: inset 0 0 0 1px rgba(59, 130, 246, 0.28);
}

@media (min-width: 768px) {
  .facility-category-meta-panel {
    grid-template-columns: auto auto 1fr;
    align-items: center;
  }
}

@media (prefers-reduced-motion: reduce) {
  .facility-category-icon-presets__option {
    transition: none;
  }
}
</style>
