<!-- admin/app/pages/cabins/new.vue -->
<script setup lang="ts">
import { ref } from 'vue'
import { useApi } from '../../composables/useApi'
const { request } = useApi()
const loading = ref(false)
const error = ref<string | null>(null)
const form = ref({
  voyage_id: 0,
  cabin_type_id: 0,
  code: '',
  max_guests: 2,
  deck: '',
})

async function handleSubmit() {
  if (loading.value) return
  loading.value = true
  error.value = null
  try {
    await request('/cabins', {
      method: 'POST',
      body: {
        voyage_id: Number(form.value.voyage_id),
        cabin_type_id: Number(form.value.cabin_type_id),
        code: form.value.code,
        max_guests: Number(form.value.max_guests),
        deck: form.value.deck,
      },
    })
    await navigateTo('/cabins')
  } catch (e: any) {
    error.value = e?.message ?? 'failed to create cabin'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="page">
    <h1>新建舱房</h1>
    <form @submit.prevent="handleSubmit">
      <input v-model.number="form.voyage_id" type="number" min="1" placeholder="Voyage ID" :disabled="loading" />
      <input v-model.number="form.cabin_type_id" type="number" min="1" placeholder="Cabin Type ID" :disabled="loading" />
      <input v-model="form.code" placeholder="Code" :disabled="loading" />
      <input v-model.number="form.max_guests" type="number" min="1" placeholder="Max Guests" :disabled="loading" />
      <input v-model="form.deck" placeholder="Deck" :disabled="loading" />
      <p v-if="error" class="error">{{ error }}</p>
      <button type="button" :disabled="loading" @click="handleSubmit">{{ loading ? '提交中...' : '创建' }}</button>
    </form>
  </div>
</template>

