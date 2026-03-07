<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

const { request } = useApi()

const loading = ref(false)
const saving = ref(false)
const error = ref('')
const success = ref('')
const form = ref({
  name: '',
  logo: '',
  contact_phone: '',
  contact_email: '',
  company_desc: '',
  service_desc: '',
  icp_number: '',
})

const isEmpty = computed(() => {
  const f = form.value
  return !f.name && !f.contact_phone && !f.contact_email && !f.icp_number
})

async function loadShopInfo() {
  loading.value = true
  error.value = ''
  try {
    const res = await request('/shop-info')
    const payload = res?.data ?? res ?? {}
    form.value = {
      name: payload?.name || '',
      logo: payload?.logo || '',
      contact_phone: payload?.contact_phone || '',
      contact_email: payload?.contact_email || '',
      company_desc: payload?.company_desc || '',
      service_desc: payload?.service_desc || '',
      icp_number: payload?.icp_number || '',
    }
  } catch (e: any) {
    error.value = e?.message ?? '加载店铺信息失败'
  } finally {
    loading.value = false
  }
}

async function saveShopInfo() {
  saving.value = true
  success.value = ''
  error.value = ''
  try {
    await request('/shop-info', { method: 'PUT', body: form.value })
    success.value = '保存成功'
  } catch (e: any) {
    error.value = e?.message ?? '保存失败'
  } finally {
    saving.value = false
  }
}

onMounted(loadShopInfo)
</script>

<template>
  <div class="admin-page">
    <AdminPageHeader title="店铺设置" subtitle="维护店铺基础信息与联系方式" />
    <AdminFormCard>
      <p v-if="loading" data-test="loading" class="text-sm text-slate-600">加载中...</p>
      <p v-else-if="error" data-test="error" class="text-sm text-rose-500">{{ error }}</p>
      <div v-else>
        <p v-if="isEmpty" data-test="empty" class="mb-3 text-sm text-slate-600">暂无店铺配置，请先填写基础信息</p>
        <form class="grid max-w-3xl gap-3" @submit.prevent="saveShopInfo">
          <label class="text-sm text-slate-600">店铺名称
            <input v-model="form.name" data-test="name" type="text" class="mt-1 h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
          </label>
          <label class="text-sm text-slate-600">联系电话
            <input v-model="form.contact_phone" data-test="phone" type="text" class="mt-1 h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
          </label>
          <label class="text-sm text-slate-600">联系邮箱
            <input v-model="form.contact_email" data-test="email" type="text" class="mt-1 h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
          </label>
          <label class="text-sm text-slate-600">备案号
            <input v-model="form.icp_number" data-test="icp" type="text" class="mt-1 h-10 w-full rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
          </label>
          <AdminActionBar>
            <button data-test="submit" type="submit" class="rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-500" :disabled="saving">{{ saving ? '保存中...' : '保存' }}</button>
          </AdminActionBar>
        </form>
        <p v-if="success" data-test="success" class="text-sm text-emerald-600">{{ success }}</p>
      </div>
    </AdminFormCard>
  </div>
</template>
