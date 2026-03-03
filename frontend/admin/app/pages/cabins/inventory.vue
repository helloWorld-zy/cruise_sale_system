<script setup lang="ts">
import { onMounted, ref } from 'vue'

const { request } = useApi()
const route = useRoute()
const skuId = Number(route?.query?.skuId ?? 0)
const loading = ref(false)
const submitting = ref(false)
const error = ref<string | null>(null)
const inventory = ref<{ total: number; locked: number; sold: number; available?: number; alert_threshold?: number } | null>(null)
const logs = ref<Array<{ time: string; delta: number; reason: string }>>([])
const delta = ref(0)
const reason = ref('手动调整')
const threshold = ref(0)

const arrows = {
  up: '↑',
  down: '↓',
}

async function loadInventory() {
  if (!skuId) {
    error.value = '缺少 skuId 参数'
    return
  }
  loading.value = true
  error.value = null
  try {
    const res = await request(`/cabins/${skuId}/inventory`)
    const payload = res?.data ?? res ?? null
    inventory.value = payload
    threshold.value = Number(payload?.alert_threshold || 0)

    try {
      const logRes = await request(`/cabins/${skuId}/inventory/logs`)
      const logPayload = logRes?.data ?? logRes ?? []
      logs.value = Array.isArray(logPayload)
        ? logPayload.map((item: Record<string, any>) => ({
            time: item.created_at || item.time || '-',
            delta: Number(item.change || item.delta || 0),
            reason: item.reason || '-',
          }))
        : []
    } catch {
      logs.value = []
    }
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load inventory'
  } finally {
    loading.value = false
  }
}

function available() {
  if (!inventory.value) return 0
  if (inventory.value.available !== undefined) return Number(inventory.value.available)
  return Number(inventory.value.total || 0) - Number(inventory.value.locked || 0) - Number(inventory.value.sold || 0)
}

async function adjustInventory() {
  if (!skuId || submitting.value) return
  submitting.value = true
  error.value = null
  try {
    await request(`/cabins/${skuId}/inventory/adjust`, {
      method: 'POST',
      body: { delta: Number(delta.value), reason: reason.value || '手动调整' },
    })
    await loadInventory()
  } catch (e: any) {
    error.value = e?.message ?? 'failed to adjust inventory'
  } finally {
    submitting.value = false
  }
}

async function saveThreshold() {
  if (!skuId || submitting.value) return
  submitting.value = true
  error.value = null
  try {
    await request(`/cabins/${skuId}/alert-threshold`, {
      method: 'PUT',
      body: { threshold: Number(threshold.value || 0) },
    })
    await loadInventory()
  } catch (e: any) {
    error.value = e?.message ?? 'failed to save threshold'
  } finally {
    submitting.value = false
  }
}

onMounted(loadInventory)
</script>

<template>
  <div class="min-h-screen bg-slate-50 p-4 md:p-6">
    <div class="mx-auto max-w-6xl space-y-4">
      <h1 class="text-xl font-semibold text-slate-900">库存管理</h1>
      <p v-if="loading" class="text-sm text-slate-600">加载中...</p>
      <p v-else-if="error" class="text-sm text-rose-500">{{ error }}</p>
      <div v-else-if="inventory" class="space-y-4">
        <div class="grid grid-cols-2 gap-3 md:grid-cols-4">
          <div class="rounded-lg border border-slate-200 bg-white p-4 shadow-sm">
            <p class="text-sm text-slate-600">总量</p>
            <p class="text-2xl font-bold text-slate-900">{{ inventory.total }}</p>
            <p class="text-xs text-emerald-600">{{ arrows.up }} 稳定</p>
          </div>
          <div class="rounded-lg border border-slate-200 bg-white p-4 shadow-sm">
            <p class="text-sm text-slate-600">锁定</p>
            <p class="text-2xl font-bold text-slate-900">{{ inventory.locked }}</p>
            <p class="text-xs text-amber-600">{{ arrows.up }} +1</p>
          </div>
          <div class="rounded-lg border border-slate-200 bg-white p-4 shadow-sm">
            <p class="text-sm text-slate-600">已售</p>
            <p class="text-2xl font-bold text-slate-900">{{ inventory.sold }}</p>
            <p class="text-xs text-emerald-600">{{ arrows.up }} +2</p>
          </div>
          <div class="rounded-lg border border-slate-200 bg-white p-4 shadow-sm">
            <p class="text-sm text-slate-600">可用</p>
            <p class="text-2xl font-bold text-slate-900">{{ available() }}</p>
            <p class="text-xs" :class="available() < threshold ? 'text-rose-600' : 'text-emerald-600'">{{ available() < threshold ? `${arrows.down} 低于阈值` : `${arrows.up} 正常` }}</p>
          </div>
        </div>

        <div class="rounded-lg border border-slate-200 bg-white p-4 shadow-sm">
          <div class="flex flex-wrap items-center gap-3">
            <span class="text-sm text-slate-700">预警阈值</span>
            <input v-model.number="threshold" type="number" min="0" class="h-10 w-32 rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
            <button type="button" class="rounded-md bg-indigo-600 px-3 py-2 text-sm text-white hover:bg-indigo-500" :disabled="submitting" @click="saveThreshold">保存</button>
          </div>
        </div>

        <div class="rounded-lg border border-slate-200 bg-white p-4 shadow-sm">
          <div class="mb-3 flex items-center gap-2">
            <input v-model.number="delta" type="number" placeholder="调整值" class="h-10 w-32 rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
            <input v-model="reason" type="text" placeholder="原因" class="h-10 w-64 rounded-md border border-slate-200 px-3 outline-none ring-indigo-500 focus:ring-2" />
            <button type="button" class="rounded-md border border-slate-200 px-3 py-2 text-sm text-slate-700 hover:bg-slate-50" :disabled="submitting" @click="adjustInventory">提交调整</button>
          </div>
          <table class="w-full text-sm">
            <thead class="bg-slate-50 text-left text-slate-600">
              <tr>
                <th class="p-3">时间</th>
                <th class="p-3">变动量</th>
                <th class="p-3">原因</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="logs.length === 0"><td class="p-3 text-slate-500" colspan="3">暂无变动记录</td></tr>
              <tr v-for="(log, idx) in logs" v-else :key="`${log.time}-${idx}`">
                <td class="p-3 text-slate-600">{{ log.time }}</td>
                <td class="p-3 font-medium" :class="log.delta >= 0 ? 'text-emerald-600' : 'text-rose-600'">{{ log.delta >= 0 ? `+${log.delta}` : log.delta }}</td>
                <td class="p-3 text-slate-600">{{ log.reason }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

