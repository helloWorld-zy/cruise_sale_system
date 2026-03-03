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
  <div class="page">
    <h1>店铺设置</h1>
    <p v-if="loading" data-test="loading">加载中...</p>
    <p v-else-if="error" data-test="error" class="error">{{ error }}</p>
    <div v-else>
      <p v-if="isEmpty" data-test="empty">暂无店铺配置，请先填写基础信息</p>
      <form @submit.prevent="saveShopInfo">
        <label>
          店铺名称
          <input v-model="form.name" data-test="name" type="text" />
        </label>
        <label>
          联系电话
          <input v-model="form.contact_phone" data-test="phone" type="text" />
        </label>
        <label>
          联系邮箱
          <input v-model="form.contact_email" data-test="email" type="text" />
        </label>
        <label>
          备案号
          <input v-model="form.icp_number" data-test="icp" type="text" />
        </label>
        <button data-test="submit" type="submit" :disabled="saving">{{ saving ? '保存中...' : '保存' }}</button>
      </form>
      <p v-if="success" data-test="success">{{ success }}</p>
    </div>
  </div>
</template>
