<template>
  <div class="admin-page">
    <AdminPageHeader title="新建设施" subtitle="补充设施的基础信息、收费策略和适用人群，用于前台展示与筛选。" />
    <AdminFormCard title="新建设施">
      <p v-if="empty" data-test="empty" class="mb-3 text-sm text-slate-600">暂无邮轮或分类数据，无法建设施</p>
      <form class="admin-cruise-form" @submit.prevent="handleSubmit">
        <section class="admin-cruise-form__intro">
          <h2 class="admin-cruise-form__intro-title">创建设施档案</h2>
          <p class="admin-cruise-form__intro-desc">先选择邮轮与分类，再完善营业信息和收费策略，便于后续在前台准确检索。</p>
        </section>

        <section class="admin-cruise-form__section">
          <h3 class="admin-cruise-form__section-title">基本信息</h3>
          <p class="admin-cruise-form__section-subtitle">确定设施归属、名称和位置等识别信息。</p>
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
              <input v-model="form.open_hours" placeholder="如 08:00-22:00" class="admin-cruise-form__control" />
            </label>
          </div>
        </section>

        <section class="admin-cruise-form__section">
          <h3 class="admin-cruise-form__section-title">收费与人群</h3>
          <p class="admin-cruise-form__section-subtitle">用于控制付费策略与适配客群标签。</p>
          <div class="admin-cruise-form__grid">
            <label class="admin-cruise-form__field facility-form-checkbox-row">
              <span class="admin-cruise-form__field-label">是否额外收费</span>
              <span class="facility-form-checkbox-wrap">
                <input v-model="form.extra_charge" type="checkbox" />
                <span>开启后可填写收费说明</span>
              </span>
            </label>
            <label v-if="form.extra_charge" class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">收费说明</span>
              <input v-model="form.charge_price_tip" placeholder="如 ¥200/人" class="admin-cruise-form__control" />
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
          <p class="admin-cruise-form__section-subtitle">建议补充体验亮点、使用规则等说明。</p>
          <label class="admin-cruise-form__field">
            <span class="admin-cruise-form__field-label">设施描述</span>
            <textarea v-model="form.description" rows="5" class="admin-cruise-form__control admin-cruise-form__control--textarea" />
          </label>
        </section>

        <AdminActionBar>
          <AdminActionLink to="/facilities" class="admin-btn admin-btn--secondary">取消</AdminActionLink>
          <button type="submit" :disabled="loading" class="admin-btn">{{ loading ? '提交中...' : '保存' }}</button>
        </AdminActionBar>
        <p v-if="error" class="text-sm text-rose-500">{{ error }}</p>
      </form>
    </AdminFormCard>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'

const { request } = useApi()
const loading = ref(false)
const error = ref<string | null>(null)
const cruises = ref<Record<string, any>[]>([])
const categories = ref<Record<string, any>[]>([])
const empty = ref(false)
const audienceOptions = ['儿童', '家庭', '情侣', '老年', '商务']

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
    empty.value = cruises.value.length === 0 || categories.value.length === 0
    if (cruises.value.length > 0) form.value.cruise_id = Number(cruises.value[0].id) || 0
    if (categories.value.length > 0) form.value.category_id = Number(categories.value[0].id) || 0
  } catch {
    cruises.value = []
    categories.value = []
    empty.value = true
  }
}

function toggleAudience(value: string, checked: boolean) {
  const next = new Set(form.value.target_audience)
  if (checked) next.add(value)
  else next.delete(value)
  form.value.target_audience = Array.from(next)
}

async function handleSubmit() {
  if (loading.value) return
  loading.value = true
  error.value = null
  try {
    await request('/facilities', {
      method: 'POST',
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
    error.value = e?.message ?? 'failed to create facility'
  } finally {
    loading.value = false
  }
}

onMounted(loadOptions)
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
