<template>
  <Teleport to="body">
    <div v-if="visible" class="admin-confirm-overlay" @click="emit('close')">
      <div class="admin-confirm-card" @click.stop>
        <h2 class="admin-confirm-title">{{ title }}</h2>
        <p class="admin-confirm-message">{{ message }}</p>
        <p v-if="hint" class="admin-confirm-hint">{{ hint }}</p>
        <div class="admin-confirm-actions">
          <button
            type="button"
            class="admin-confirm-cancel"
            :disabled="loading"
            @click="emit('close')"
          >
            {{ cancelText }}
          </button>
          <button
            type="button"
            class="admin-confirm-submit"
            :disabled="loading"
            @click="emit('confirm')"
          >
            {{ loading ? loadingText : confirmText }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  visible: boolean
  title: string
  message: string
  hint?: string
  loading?: boolean
  confirmText?: string
  cancelText?: string
  loadingText?: string
}>(), {
  hint: '',
  loading: false,
  confirmText: '确认删除',
  cancelText: '取消',
  loadingText: '处理中...',
})

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'confirm'): void
}>()
</script>

<style scoped>
.admin-confirm-overlay {
  position: fixed;
  inset: 0;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(15, 23, 42, 0.45);
  padding: 16px;
}

.admin-confirm-card {
  width: min(520px, 100%);
  border-radius: 12px;
  background: #ffffff;
  padding: 20px;
  box-shadow: 0 20px 48px rgba(15, 23, 42, 0.24);
}

.admin-confirm-title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #0f172a;
}

.admin-confirm-message {
  margin-top: 10px;
  font-size: 14px;
  line-height: 1.6;
  color: #334155;
}

.admin-confirm-hint {
  margin-top: 6px;
  font-size: 13px;
  line-height: 1.6;
  color: #b45309;
}

.admin-confirm-actions {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.admin-confirm-cancel {
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  background: #ffffff;
  color: #334155;
  padding: 8px 16px;
  font-size: 14px;
  cursor: pointer;
}

.admin-confirm-cancel:hover {
  background: #f8fafc;
}

.admin-confirm-submit {
  border: 1px solid #dc2626;
  border-radius: 8px;
  background: #dc2626;
  color: #ffffff;
  padding: 8px 16px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
}

.admin-confirm-submit:hover {
  background: #b91c1c;
  border-color: #b91c1c;
}

.admin-confirm-cancel:disabled,
.admin-confirm-submit:disabled {
  cursor: not-allowed;
  opacity: 0.7;
}
</style>
