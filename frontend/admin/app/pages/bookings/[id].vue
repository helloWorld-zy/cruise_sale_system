<!-- admin/app/pages/bookings/[id].vue — 订单详情页面 -->
<!-- 根据路由参数加载并展示单个订单的详细信息 -->
<script setup lang="ts">
import { onMounted, ref } from 'vue'
import AdminConfirmDialog from '../../components/AdminConfirmDialog.vue'
import { useAdminDeleteDialog } from '../../composables/useAdminDeleteDialog'

declare const useApi: any
declare const navigateTo: any

// 根据路由参数中的订单 ID 加载并展示订单详情
const route = useRoute()
const { request } = useApi()

const loading = ref(false)
const saving = ref(false)
const error = ref<string | null>(null)
const booking = ref<{ id: number; status: string; total_cents: number } | null>(null)
const status = ref('')
const {
  visible: deleteDialogVisible,
  submitting: deleting,
  open: openDeleteDialog,
  close: closeDeleteDialog,
  run: runDelete,
} = useAdminDeleteDialog()

onMounted(async () => {
  loading.value = true
  error.value = null
  try {
    const id = String(route.params.id ?? '')
    const res = await request(`/bookings/${id}`)
    booking.value = res?.data ?? res ?? null
    status.value = booking.value?.status ?? ''
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load booking detail'
  } finally {
    loading.value = false
  }
})

async function handleSave() {
  if (!booking.value || saving.value) return
  saving.value = true
  error.value = null
  try {
    await request(`/bookings/${booking.value.id}`, {
      method: 'PUT',
      body: {
        status: status.value,
      },
    })
    booking.value.status = status.value
  } catch (e: any) {
    error.value = e?.message ?? 'failed to update booking'
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  if (!booking.value || deleting.value) return
  openDeleteDialog()
}

async function confirmDelete() {
  if (!booking.value || deleting.value) return
  error.value = null
  try {
    await runDelete(async () => {
      await request(`/bookings/${booking.value.id}`, { method: 'DELETE' })
      await navigateTo('/bookings')
    })
  } catch (e: any) {
    error.value = e?.message ?? '删除订单失败，请稍后重试。'
  }
}
</script>

<template>
  <div class="admin-page">
    <AdminPageHeader :title="`订单详情 #${route.params.id}`" />
    <AdminFormCard title="订单详情">
      <p v-if="loading" class="text-sm text-slate-600">Loading...</p>
      <p v-else-if="error" class="text-sm text-rose-500">{{ error }}</p>
      <div v-else-if="booking" class="space-y-4">
        <div class="grid grid-cols-1 gap-3 md:grid-cols-3">
          <div class="rounded-md border border-slate-200 bg-slate-50 p-3">
            <p class="text-xs text-slate-500">订单 ID</p>
            <p class="text-base font-semibold text-slate-900">{{ booking.id }}</p>
          </div>
          <div class="rounded-md border border-slate-200 bg-slate-50 p-3">
            <p class="text-xs text-slate-500">当前状态</p>
            <p class="text-base font-semibold text-slate-900">{{ booking.status }}</p>
          </div>
          <div class="rounded-md border border-slate-200 bg-slate-50 p-3">
            <p class="text-xs text-slate-500">订单总额（分）</p>
            <p class="text-base font-semibold text-slate-900">{{ booking.total_cents }}</p>
          </div>
        </div>

        <label class="block max-w-sm space-y-1 text-sm text-slate-600">
          <span>新状态</span>
          <select v-model="status" :disabled="saving || deleting" class="h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2">
            <option value="pending_payment">待支付</option>
            <option value="paid">已支付</option>
            <option value="cancelled">已取消</option>
            <option value="refunding">退改中</option>
            <option value="refunded">已退款</option>
          </select>
        </label>

        <AdminActionBar>
          <AdminActionLink to="/bookings" class="admin-btn admin-btn--secondary">返回列表</AdminActionLink>
          <button type="button" class="admin-btn" :disabled="saving || deleting" @click="handleSave">{{ saving ? '保存中...' : '保存状态' }}</button>
          <button type="button" class="admin-btn admin-btn--danger" :disabled="saving || deleting" @click="handleDelete">{{ deleting ? '删除中...' : '删除订单' }}</button>
        </AdminActionBar>
      </div>
      <p v-else data-test="empty" class="text-sm text-slate-600">暂无订单数据</p>
    </AdminFormCard>

    <AdminConfirmDialog
      :visible="deleteDialogVisible"
      title="确认删除订单"
      :message="`确认删除订单 #${booking?.id ?? ''} 吗？删除后不可恢复。`"
      :loading="deleting"
      loading-text="删除中..."
      @close="closeDeleteDialog"
      @confirm="confirmDelete"
    />
  </div>
</template>
