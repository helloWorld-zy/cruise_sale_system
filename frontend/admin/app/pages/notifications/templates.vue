<script setup lang="ts">
import { onMounted, ref } from 'vue'

const { request } = useApi()

type TemplateItem = {
  id: number
  event_type: string
  channel: string
  template: string
  enabled?: boolean
}

const loading = ref(false)
const saving = ref(false)
const error = ref('')
const items = ref<TemplateItem[]>([])
const form = ref({
  event_type: '',
  channel: 'sms',
  template: '',
  enabled: true,
})

async function loadItems() {
  loading.value = true
  error.value = ''
  try {
    const res = await request('/notification-templates')
    const payload = res?.data ?? res
    items.value = Array.isArray(payload) ? payload : (payload?.list ?? [])
  } catch (e: any) {
    error.value = e?.message ?? '加载模板失败'
  } finally {
    loading.value = false
  }
}

async function createTemplate() {
  saving.value = true
  error.value = ''
  try {
    await request('/notification-templates', {
      method: 'POST',
      body: form.value,
    })
    form.value = { event_type: '', channel: 'sms', template: '', enabled: true }
    await loadItems()
  } catch (e: any) {
    error.value = e?.message ?? '创建模板失败'
  } finally {
    saving.value = false
  }
}

onMounted(loadItems)
</script>

<template>
  <div class="page">
    <h1>通知模板</h1>
    <p v-if="loading" data-test="loading">加载中...</p>
    <p v-else-if="error" data-test="error" class="error">{{ error }}</p>
    <p v-else-if="items.length === 0" data-test="empty">暂无模板数据</p>
    <table v-else data-test="table">
      <thead>
        <tr>
          <th>事件</th>
          <th>渠道</th>
          <th>模板内容</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="it in items" :key="it.id">
          <td>{{ it.event_type }}</td>
          <td>{{ it.channel }}</td>
          <td>{{ it.template }}</td>
        </tr>
      </tbody>
    </table>

    <form style="margin-top: 12px" @submit.prevent="createTemplate">
      <input v-model="form.event_type" data-test="event" type="text" placeholder="事件类型" />
      <select v-model="form.channel" data-test="channel">
        <option value="sms">sms</option>
        <option value="wechat_subscribe">wechat_subscribe</option>
        <option value="wechat_template">wechat_template</option>
        <option value="in_app">in_app</option>
      </select>
      <input v-model="form.template" data-test="template" type="text" placeholder="模板内容" />
      <button data-test="create" type="submit" :disabled="saving">{{ saving ? '提交中...' : '新增模板' }}</button>
    </form>
  </div>
</template>
