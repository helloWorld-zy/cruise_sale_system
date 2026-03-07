import { ref } from 'vue'

export function useAdminDeleteDialog<T = never>() {
  const visible = ref(false)
  const submitting = ref(false)
  const target = ref<T | null>(null)

  function open(nextTarget?: T) {
    if (submitting.value) return
    if (arguments.length > 0) {
      target.value = (nextTarget ?? null) as T | null
    }
    visible.value = true
  }

  function close(force = false) {
    if (submitting.value && !force) return false
    visible.value = false
    target.value = null
    return true
  }

  async function run(action: () => Promise<void>) {
    if (submitting.value) return false
    submitting.value = true
    try {
      await action()
      close(true)
      return true
    } finally {
      submitting.value = false
    }
  }

  return {
    visible,
    submitting,
    target,
    open,
    close,
    run,
  }
}