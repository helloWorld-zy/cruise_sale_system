<script setup lang="ts">
import { onMounted, ref } from 'vue'
import AdminConfirmDialog from '../../components/AdminConfirmDialog.vue'
import { useAdminDeleteDialog } from '../../composables/useAdminDeleteDialog'

const route = useRoute()
const { request } = useApi()
const id = Number(route.params.id)

const loading = ref(false)
const saving = ref(false)
const error = ref<string | null>(null)
const empty = ref(false)
const amenityOptions = ['浴缸', '迷你吧', '智能电视', '胶囊咖啡机', '沙发床', '独立衣帽间']
const {
  visible: deleteDialogVisible,
  submitting: deleting,
  open: openDeleteDialog,
  close: closeDeleteDialog,
  run: runDelete,
} = useAdminDeleteDialog()

const form = ref({
  voyage_id: 0,
  cabin_type_id: 0,
  code: '',
  deck: '',
  area: 0,
  max_guests: 2,
  position: 'mid',
  orientation: 'port',
  has_window: false,
  has_balcony: false,
  bed_type: '',
  amenities: [] as string[],
  status: 1,
})

async function loadDetail() {
  loading.value = true
  error.value = null
  empty.value = false
  try {
    const res = await request(`/cabins/${id}`)
    const data = res?.data ?? res ?? {}
    if (Object.keys(data).length === 0) {
      empty.value = true
      return
    }
    form.value = {
      voyage_id: Number(data.voyage_id ?? 0),
      cabin_type_id: Number(data.cabin_type_id ?? 0),
      code: data.code ?? '',
      deck: data.deck ?? '',
      area: Number(data.area ?? 0),
      max_guests: Number(data.max_guests ?? 2),
      position: data.position || 'mid',
      orientation: data.orientation || 'port',
      has_window: Boolean(data.has_window),
      has_balcony: Boolean(data.has_balcony),
      bed_type: data.bed_type ?? '',
      amenities: splitCsv(data.amenities),
      status: Number(data.status ?? 1),
    }
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load cabin detail'
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  if (saving.value) return
  saving.value = true
  error.value = null
  try {
    await request(`/cabins/${id}`, {
      method: 'PUT',
      body: {
        voyage_id: Number(form.value.voyage_id),
        cabin_type_id: Number(form.value.cabin_type_id),
        code: form.value.code,
        deck: form.value.deck,
        area: Number(form.value.area),
        max_guests: Number(form.value.max_guests),
        position: form.value.position,
        orientation: form.value.orientation,
        has_window: Boolean(form.value.has_window),
        has_balcony: Boolean(form.value.has_balcony),
        bed_type: form.value.bed_type,
        amenities: form.value.amenities.join(','),
        status: Number(form.value.status),
      },
    })
    await navigateTo('/cabins')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to update cabin'
  } finally {
    saving.value = false
  }
}

function splitCsv(raw: unknown) {
  if (typeof raw !== 'string') return []
  return raw.split(',').map((part) => part.trim()).filter(Boolean)
}

function toggleAmenity(item: string, checked: boolean) {
  const next = new Set(form.value.amenities)
  if (checked) next.add(item)
  else next.delete(item)
  form.value.amenities = Array.from(next)
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
      await request(`/cabins/${id}`, { method: 'DELETE' })
      await navigateTo('/cabins')
    })
  } catch (e: any) {
    error.value = e?.message ?? '删除舱房失败，请稍后重试。'
  }
}

onMounted(loadDetail)
</script>

<template>
  <div class="admin-page">
    <AdminPageHeader :title="`编辑舱位 #${id}`" />
    <AdminFormCard title="舱位配置">
      <p v-if="loading" class="text-sm text-slate-600">加载中...</p>
      <p v-else-if="empty" data-test="empty" class="text-sm text-slate-600">暂无舱位数据</p>
      <form v-else class="space-y-4" @submit.prevent="handleSave">
        <section class="rounded-lg border border-slate-200 p-4">
          <h2 class="mb-3 text-sm font-semibold text-slate-700">基本信息</h2>
          <div class="grid grid-cols-1 gap-4 md:grid-cols-3">
            <label class="space-y-1 text-sm text-slate-600"><span>航线 ID</span><input v-model.number="form.voyage_id" type="number" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" /></label>
            <label class="space-y-1 text-sm text-slate-600"><span>舱型 ID</span><input v-model.number="form.cabin_type_id" type="number" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" /></label>
            <label class="space-y-1 text-sm text-slate-600"><span>编号</span><input v-model="form.code" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" /></label>
            <label class="space-y-1 text-sm text-slate-600"><span>甲板</span><input v-model="form.deck" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" /></label>
            <label class="space-y-1 text-sm text-slate-600"><span>面积</span><input v-model.number="form.area" type="number" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" /></label>
            <label class="space-y-1 text-sm text-slate-600"><span>人数上限</span><input v-model.number="form.max_guests" type="number" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" /></label>
          </div>
        </section>

        <section class="rounded-lg border border-slate-200 p-4">
          <h2 class="mb-3 text-sm font-semibold text-slate-700">位置属性</h2>
          <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
            <div>
              <p class="mb-2 text-sm text-slate-600">位置</p>
              <div class="inline-flex rounded-md border border-slate-200 p-1">
                <button type="button" class="rounded px-3 py-1 text-sm" :class="form.position === 'fore' ? 'bg-indigo-600 text-white' : 'text-slate-600'" @click="form.position = 'fore'">前段</button>
                <button type="button" class="rounded px-3 py-1 text-sm" :class="form.position === 'mid' ? 'bg-indigo-600 text-white' : 'text-slate-600'" @click="form.position = 'mid'">中段</button>
                <button type="button" class="rounded px-3 py-1 text-sm" :class="form.position === 'aft' ? 'bg-indigo-600 text-white' : 'text-slate-600'" @click="form.position = 'aft'">后段</button>
              </div>
            </div>
            <div>
              <p class="mb-2 text-sm text-slate-600">朝向</p>
              <div class="inline-flex rounded-md border border-slate-200 p-1">
                <button type="button" class="rounded px-3 py-1 text-sm" :class="form.orientation === 'port' ? 'bg-indigo-600 text-white' : 'text-slate-600'" @click="form.orientation = 'port'">左舷</button>
                <button type="button" class="rounded px-3 py-1 text-sm" :class="form.orientation === 'starboard' ? 'bg-indigo-600 text-white' : 'text-slate-600'" @click="form.orientation = 'starboard'">右舷</button>
              </div>
            </div>
            <label class="flex items-center gap-2 text-sm text-slate-700"><input v-model="form.has_window" type="checkbox" /><span>有窗户</span></label>
            <label class="flex items-center gap-2 text-sm text-slate-700"><input v-model="form.has_balcony" type="checkbox" /><span>有阳台</span></label>
          </div>
        </section>

        <section class="rounded-lg border border-slate-200 p-4">
          <h2 class="mb-3 text-sm font-semibold text-slate-700">设施配置</h2>
          <label class="mb-3 block space-y-1 text-sm text-slate-600"><span>床型</span><input v-model="form.bed_type" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" /></label>
          <div class="grid grid-cols-2 gap-2 md:grid-cols-3">
            <label v-for="item in amenityOptions" :key="item" class="flex items-center gap-2 rounded border border-slate-200 px-3 py-2 text-sm text-slate-700">
              <input type="checkbox" :checked="form.amenities.includes(item)" @change="toggleAmenity(item, ($event.target as HTMLInputElement).checked)" />
              <span>{{ item }}</span>
            </label>
          </div>
          <label class="mt-3 block space-y-1 text-sm text-slate-600"><span>状态</span><select v-model.number="form.status" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2"><option :value="1">上架</option><option :value="2">维护中</option><option :value="0">下架</option></select></label>
        </section>

        <AdminActionBar>
          <AdminActionLink to="/cabins" class="admin-btn admin-btn--secondary">取消</AdminActionLink>
          <button type="submit" class="admin-btn" :disabled="saving || deleting">{{ saving ? '保存中...' : '保存' }}</button>
          <button type="button" class="admin-btn admin-btn--danger" :disabled="saving || deleting" @click="handleDelete">{{ deleting ? '删除中...' : '删除' }}</button>
        </AdminActionBar>
        <p v-if="error" class="text-sm text-rose-500">{{ error }}</p>
      </form>

      <AdminConfirmDialog
        :visible="deleteDialogVisible"
        title="确认删除舱房"
        :message="`确认删除舱房 #${id} 吗？删除后不可恢复。`"
        :loading="deleting"
        loading-text="删除中..."
        @close="closeDeleteDialog"
        @confirm="confirmDelete"
      />
    </AdminFormCard>
  </div>
</template>
