<template>
  <div ref="rootRef" class="port-city-select">
    <div class="port-city-select-field">
      <input
        :value="searchText"
        :placeholder="placeholder"
        :disabled="disabled"
        class="port-city-select-input"
        :data-test="`${testIdBase}-input`"
        @focus="openPanel"
        @input="onInput"
        @compositionstart="onCompositionStart"
        @compositionend="onCompositionEnd"
      />
      <button
        v-if="searchText"
        type="button"
        class="port-city-select-clear"
        :disabled="disabled"
        :data-test="`${testIdBase}-clear`"
        @click="clearValue"
      >
        清除
      </button>
    </div>

    <div v-if="open" class="port-city-select-panel">
      <div v-if="loading" class="port-city-select-state">搜索中...</div>
      <button
        v-for="(item, index) in options"
        :key="`${item.label}-${index}`"
        type="button"
        class="port-city-select-option"
        :data-test="`${testIdBase}-option-${index}`"
        @click="selectOption(item.label)"
      >
        <span class="port-city-select-option-label">{{ item.label }}</span>
        <span v-if="item.is_special" class="port-city-select-option-meta">特殊站点</span>
      </button>
      <div v-if="!loading && options.length === 0" class="port-city-select-state">未找到可选城市，请继续输入关键词</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, ref, watch } from 'vue'
import { useApi } from '../composables/useApi'

type PortCityOption = {
  label: string
  city_name?: string
  country_name?: string
  is_special?: boolean
}

const props = withDefaults(defineProps<{
  modelValue: string
  disabled?: boolean
  placeholder?: string
  testIdBase: string
}>(), {
  disabled: false,
  placeholder: '搜索城市或输入海上巡游',
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const { request } = useApi()
const rootRef = ref<HTMLElement | null>(null)
const open = ref(false)
const loading = ref(false)
const searchText = ref(props.modelValue || '')
const options = ref<PortCityOption[]>([])
const isComposing = ref(false)
let lastIssuedKeyword = ''
let searchTimer: ReturnType<typeof setTimeout> | null = null

watch(() => props.modelValue, (value) => {
  searchText.value = value || ''
})

async function fetchOptions(keyword: string) {
  const trimmed = keyword.trim()
  if (!trimmed) {
    options.value = []
    loading.value = false
    return
  }
  loading.value = true
  lastIssuedKeyword = trimmed
  try {
    const res = await request('/port-cities', { query: { keyword: trimmed } })
    if (lastIssuedKeyword !== trimmed) return
    const payload = res?.data ?? res ?? []
    options.value = Array.isArray(payload) ? payload : []
  } catch {
    if (lastIssuedKeyword === trimmed) {
      options.value = []
    }
  } finally {
    if (lastIssuedKeyword === trimmed) {
      loading.value = false
    }
  }
}

function scheduleFetch(keyword: string) {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    void fetchOptions(keyword)
  }, 320)
}

function openPanel() {
  if (props.disabled) return
  open.value = true
  scheduleFetch(searchText.value)
}

function onInput(event: Event) {
  const target = event.target as HTMLInputElement
  searchText.value = target.value
  open.value = true
  if (isComposing.value) return
  scheduleFetch(searchText.value)
}

function onCompositionStart() {
  isComposing.value = true
}

function onCompositionEnd(event: CompositionEvent) {
  isComposing.value = false
  const target = event.target as HTMLInputElement
  searchText.value = target.value
  open.value = true
  scheduleFetch(searchText.value)
}

function selectOption(label: string) {
  emit('update:modelValue', label)
  searchText.value = label
  open.value = false
}

function clearValue() {
  emit('update:modelValue', '')
  searchText.value = ''
  options.value = []
  open.value = false
  loading.value = false
  if (searchTimer) clearTimeout(searchTimer)
}

function onDocumentClick(event: MouseEvent) {
  if (!open.value) return
  const target = event.target as Node | null
  if (rootRef.value && target && !rootRef.value.contains(target)) {
    searchText.value = props.modelValue || ''
    open.value = false
  }
}

document.addEventListener('click', onDocumentClick)

onBeforeUnmount(() => {
  if (searchTimer) clearTimeout(searchTimer)
  document.removeEventListener('click', onDocumentClick)
})

function isSeaCruiseKeyword(keyword: string) {
  return ['海上', '巡游', '巡航'].some((item) => keyword.includes(item))
}
</script>

<style scoped>
.port-city-select {
  position: relative;
}

.port-city-select-field {
  display: flex;
  gap: 8px;
  align-items: center;
}

.port-city-select-input {
  flex: 1;
  min-width: 0;
}

.port-city-select-clear {
  border: 1px solid #cbd5e1;
  background: #fff;
  border-radius: 8px;
  padding: 8px 12px;
  cursor: pointer;
}

.port-city-select-panel {
  position: absolute;
  z-index: 40;
  left: 0;
  right: 0;
  top: calc(100% + 6px);
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.14);
  max-height: 240px;
  overflow: auto;
}

.port-city-select-option {
  width: 100%;
  border: 0;
  background: transparent;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  cursor: pointer;
  text-align: left;
}

.port-city-select-option:hover {
  background: #f8fafc;
}

.port-city-select-option-label {
  color: #0f172a;
}

.port-city-select-option-meta,
.port-city-select-state {
  color: #64748b;
  font-size: 12px;
}

.port-city-select-state {
  padding: 10px 12px;
}
</style>