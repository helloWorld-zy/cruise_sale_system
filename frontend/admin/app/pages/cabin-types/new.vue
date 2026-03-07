<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'

const { request } = useApi()

const loading = ref(false)
const error = ref<string | null>(null)
const companies = ref<Record<string, any>[]>([])
const cruises = ref<Record<string, any>[]>([])
const categories = ref<Record<string, any>[]>([])
const selectedCruiseIds = ref<number[]>([])

const form = ref({
  company_id: 0,
  category_id: 0,
  name: '',
  english_name: '',
  code: '',
  area_min: 0,
  area_max: 0,
  occupancy: 2,
  capacity: 2,
  max_capacity: 2,
  intro: '',
  tags: '',
  amenities: '',
  sort_order: 0,
  status: 1,
})

const canSubmit = computed(() => {
  return selectedCruiseIds.value.length > 0 && form.value.category_id > 0 && form.value.name.trim().length > 0
})

async function loadCompanies() {
  try {
    const res = await request('/companies', { query: { page: 1, page_size: 100 } })
    const payload = res?.data ?? res ?? {}
    companies.value = Array.isArray(payload) ? payload : payload?.list ?? []
    if (companies.value.length > 0) {
      const first = Number(companies.value[0]?.id)
      if (Number.isFinite(first) && first > 0) form.value.company_id = first
    }
  } catch {
    companies.value = []
  }
}

async function loadCruises() {
  try {
    const query: Record<string, any> = { page: 1, page_size: 200 }
    if (form.value.company_id > 0) query.company_id = form.value.company_id
    const res = await request('/cruises', { query })
    const payload = res?.data ?? res ?? {}
    cruises.value = Array.isArray(payload) ? payload : payload?.list ?? []
    const allowed = new Set(cruises.value.map((item) => Number(item.id)))
    selectedCruiseIds.value = selectedCruiseIds.value.filter((id) => allowed.has(id))
  } catch {
    cruises.value = []
    selectedCruiseIds.value = []
  }
}

async function loadCategories() {
  try {
    const res = await request('/cabin-type-categories')
    const payload = res?.data ?? res ?? {}
    const list = Array.isArray(payload) ? payload : payload?.list ?? []
    categories.value = list.filter((item: Record<string, any>) => Number(item.status ?? 1) !== 0)
    if (categories.value.length > 0 && form.value.category_id <= 0) {
      form.value.category_id = Number(categories.value[0]?.id) || 0
    }
  } catch {
    categories.value = []
  }
}

function toggleCruise(cruiseID: number, checked: boolean) {
  const next = new Set(selectedCruiseIds.value)
  if (checked) next.add(cruiseID)
  else next.delete(cruiseID)
  selectedCruiseIds.value = Array.from(next)
}

async function handleSubmit() {
  if (!canSubmit.value || loading.value) return
  loading.value = true
  error.value = null
  try {
    await request('/cabin-types/batch-create', {
      method: 'POST',
      body: {
        // Back-end batch request embeds CabinTypeRequest where cruise_id is still required.
        cruise_id: Number(selectedCruiseIds.value[0] ?? 0),
        cruise_ids: selectedCruiseIds.value,
        category_id: Number(form.value.category_id),
        name: form.value.name.trim(),
        english_name: form.value.english_name.trim(),
        code: form.value.code.trim(),
        area_min: Number(form.value.area_min),
        area_max: Number(form.value.area_max),
        occupancy: Number(form.value.occupancy),
        capacity: Number(form.value.capacity),
        max_capacity: Number(form.value.max_capacity),
        intro: form.value.intro.trim(),
        description: form.value.intro.trim(),
        tags: form.value.tags.trim(),
        amenities: form.value.amenities.trim(),
        sort_order: Number(form.value.sort_order),
        floor_plan_url: '',
        deck: '',
        bed_type: '',
        status: Number(form.value.status),
      },
    })
    await navigateTo('/cabin-types')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to create cabin type'
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await loadCompanies()
  await Promise.all([loadCruises(), loadCategories()])
})

watch(
  () => form.value.company_id,
  async () => {
    await loadCruises()
  },
)
</script>

<template>
  <div class="admin-page">
    <AdminPageHeader title="新建舱型" subtitle="一次可选择多个邮轮批量创建同类舱型，减少重复录入。">
      <template #actions>
        <AdminActionLink to="/cabin-types">返回列表</AdminActionLink>
      </template>
    </AdminPageHeader>

    <AdminFormCard title="舱型基础信息">
      <form class="admin-cruise-form" @submit.prevent="handleSubmit">
        <section class="admin-cruise-form__intro">
          <h2 class="admin-cruise-form__intro-title">批量创建舱型</h2>
          <p class="admin-cruise-form__intro-desc">先配置通用参数，再勾选适用邮轮，一次提交即可为多个邮轮生成舱型。</p>
        </section>

        <section class="admin-cruise-form__section">
          <h3 class="admin-cruise-form__section-title">识别信息</h3>
          <p class="admin-cruise-form__section-subtitle">用于后台检索、运营筛选和前台展示。</p>
          <div class="admin-cruise-form__grid">
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">所属公司</span>
              <select v-model.number="form.company_id" class="admin-cruise-form__control">
                <option :value="0">请选择公司</option>
                <option v-for="company in companies" :key="company.id" :value="Number(company.id)">{{ company.name || `公司 #${company.id}` }}</option>
              </select>
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">舱型大类</span>
              <select v-model.number="form.category_id" class="admin-cruise-form__control">
                <option :value="0">请选择大类</option>
                <option v-for="category in categories" :key="category.id" :value="Number(category.id)">{{ category.name || `分类 #${category.id}` }}</option>
              </select>
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">舱型名称</span>
              <input v-model="form.name" class="admin-cruise-form__control" />
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">代码</span>
              <input v-model="form.code" class="admin-cruise-form__control" />
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">英文名</span>
              <input v-model="form.english_name" class="admin-cruise-form__control" />
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">排序权重</span>
              <input v-model.number="form.sort_order" type="number" class="admin-cruise-form__control" />
            </label>
          </div>
        </section>

        <section class="admin-cruise-form__section">
          <h3 class="admin-cruise-form__section-title">规格参数</h3>
          <p class="admin-cruise-form__section-subtitle">建议按真实舱型参数填写，避免后续库存和价格管理偏差。</p>
          <div class="admin-cruise-form__grid">
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">面积最小值</span>
              <input v-model.number="form.area_min" type="number" min="0" step="0.1" class="admin-cruise-form__control" />
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">面积最大值</span>
              <input v-model.number="form.area_max" type="number" min="0" step="0.1" class="admin-cruise-form__control" />
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">默认入住人数</span>
              <input v-model.number="form.occupancy" type="number" min="1" class="admin-cruise-form__control" />
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">标准容量</span>
              <input v-model.number="form.capacity" type="number" min="1" class="admin-cruise-form__control" />
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">最大容量</span>
              <input v-model.number="form.max_capacity" type="number" min="1" class="admin-cruise-form__control" />
            </label>
          </div>
        </section>

        <section class="admin-cruise-form__section">
          <h3 class="admin-cruise-form__section-title">介绍与标签</h3>
          <p class="admin-cruise-form__section-subtitle">标签和设施支持逗号分隔，便于前台卡片展示。</p>
          <label class="admin-cruise-form__field">
            <span class="admin-cruise-form__field-label">简介</span>
            <textarea v-model="form.intro" rows="4" class="admin-cruise-form__control admin-cruise-form__control--textarea" />
          </label>
          <div class="admin-cruise-form__grid">
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">标签（逗号分隔）</span>
              <input v-model="form.tags" class="admin-cruise-form__control" />
            </label>
            <label class="admin-cruise-form__field">
              <span class="admin-cruise-form__field-label">设施（逗号分隔）</span>
              <input v-model="form.amenities" class="admin-cruise-form__control" />
            </label>
          </div>
        </section>

        <section class="admin-cruise-form__section">
          <h3 class="admin-cruise-form__section-title">适用邮轮</h3>
          <p class="admin-cruise-form__section-subtitle">可多选，系统会按每个邮轮独立创建舱型记录。</p>
          <div class="cabin-type-cruise-list">
            <label v-for="cruise in cruises" :key="cruise.id" class="cabin-type-cruise-list__item">
              <input
                type="checkbox"
                :checked="selectedCruiseIds.includes(Number(cruise.id))"
                @change="toggleCruise(Number(cruise.id), ($event.target as HTMLInputElement).checked)"
              />
              <span>{{ cruise.name || `邮轮 #${cruise.id}` }}</span>
            </label>
            <p v-if="cruises.length === 0" class="text-sm text-slate-500">当前公司下暂无邮轮</p>
          </div>
        </section>

        <AdminActionBar>
          <AdminActionLink to="/cabin-types" class="admin-btn admin-btn--secondary">取消</AdminActionLink>
          <button
            type="submit"
            :disabled="loading || !canSubmit"
            class="admin-btn disabled:cursor-not-allowed disabled:opacity-60"
          >
            {{ loading ? '提交中...' : '保存' }}
          </button>
        </AdminActionBar>
        <p v-if="error" class="text-sm text-rose-500">{{ error }}</p>
      </form>
    </AdminFormCard>
  </div>
</template>

<style scoped>
.cabin-type-cruise-list {
  display: grid;
  gap: 10px;
  grid-template-columns: repeat(1, minmax(0, 1fr));
  border: 1px solid #dbe4f0;
  border-radius: 8px;
  background: #f8fbff;
  padding: 12px;
}

.cabin-type-cruise-list__item {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: #334155;
}

@media (min-width: 768px) {
  .cabin-type-cruise-list {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
</style>
