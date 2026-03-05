<template>
  <div ref="rootRef" class="company-select">
    <button
      type="button"
      class="company-select-trigger"
      :disabled="disabled"
      data-test="company-select-trigger"
      @click="toggleOpen"
    >
      <span class="company-select-value">
        <img
          v-if="selectedOption?.logo_url"
          :src="selectedOption.logo_url"
          alt="company-logo"
          class="company-select-logo"
        />
        <span class="company-select-text">{{ selectedLabel }}</span>
      </span>
      <span class="company-select-arrow">▾</span>
    </button>

    <div v-if="open" class="company-select-panel">
      <button
        type="button"
        class="company-select-item"
        data-test="company-option-0"
        @click="selectOption(0)"
      >
        <span class="company-select-text">{{ placeholder }}</span>
      </button>
      <button
        v-for="item in options"
        :key="item.id"
        type="button"
        class="company-select-item"
        :data-test="`company-option-${item.id}`"
        @click="selectOption(item.id)"
      >
        <img
          v-if="item.logo_url"
          :src="item.logo_url"
          alt="company-logo"
          class="company-select-logo"
        />
        <span class="company-select-text">{{ item.name || item.english_name || `#${item.id}` }}</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'

type CompanyOption = {
  id: number
  name: string
  english_name?: string
  logo_url?: string
}

const props = withDefaults(defineProps<{
  modelValue: number
  options: CompanyOption[]
  disabled?: boolean
  placeholder?: string
}>(), {
  disabled: false,
  placeholder: '请选择公司',
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: number): void
}>()

const open = ref(false)
const rootRef = ref<HTMLElement | null>(null)

const selectedOption = computed(() => {
  const currentId = Number(props.modelValue)
  return props.options.find((item) => Number(item.id) === currentId)
})

const selectedLabel = computed(() => {
  if (!selectedOption.value) return props.placeholder
  return selectedOption.value.name || selectedOption.value.english_name || `#${selectedOption.value.id}`
})

function toggleOpen() {
  if (props.disabled) return
  open.value = !open.value
}

function selectOption(id: number) {
  emit('update:modelValue', Number(id) || 0)
  open.value = false
}

function onDocumentClick(event: MouseEvent) {
  if (!open.value) return
  const root = rootRef.value
  const target = event.target as Node | null
  if (root && target && !root.contains(target)) {
    open.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', onDocumentClick)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', onDocumentClick)
})
</script>

<style scoped>
.company-select {
  position: relative;
}

.company-select-trigger {
  width: 100%;
  min-height: 40px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  background: #ffffff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 8px 12px;
  cursor: pointer;
}

.company-select-trigger:disabled {
  cursor: not-allowed;
  opacity: 0.7;
}

.company-select-value {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.company-select-logo {
  width: 22px;
  height: 22px;
  border-radius: 6px;
  object-fit: cover;
  border: 1px solid #e2e8f0;
  flex-shrink: 0;
}

.company-select-text {
  color: #0f172a;
  font-size: 14px;
  line-height: 1.4;
  text-align: left;
  word-break: break-all;
}

.company-select-arrow {
  color: #64748b;
  font-size: 12px;
  flex-shrink: 0;
}

.company-select-panel {
  position: absolute;
  z-index: 40;
  left: 0;
  right: 0;
  top: calc(100% + 6px);
  max-height: 240px;
  overflow: auto;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  background: #ffffff;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.16);
}

.company-select-item {
  width: 100%;
  border: 0;
  background: transparent;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  cursor: pointer;
}

.company-select-item:hover {
  background: #f8fafc;
}
</style>
