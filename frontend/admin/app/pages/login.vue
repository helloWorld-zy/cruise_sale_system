<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50">
    <div class="w-full max-w-md bg-white border rounded-xl p-6">
      <h1 class="text-xl font-semibold mb-4">管理员登录</h1>
      <form class="space-y-3" @submit.prevent="handleLogin">
        <UInput v-model="username" placeholder="用户名" />
        <!-- HI-05 FIX: use v-model="password" not v-model="password.value" -->
        <UInput v-model="password" type="password" placeholder="密码" />
        <p v-if="error" class="text-red-500 text-sm">{{ error }}</p>
        <UButton class="w-full" color="primary" type="submit" :loading="loading">登录</UButton>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function handleLogin() {
  if (!username.value || !password.value) {
    error.value = '请填写用户名和密码'
    return
  }
  loading.value = true
  error.value = ''
  try {
    // TODO: call auth API in Sprint 2 integration
    await new Promise(resolve => setTimeout(resolve, 500))
    // navigateTo('/dashboard')
  } catch (e: any) {
    error.value = e.message || '登录失败，请重试'
  } finally {
    loading.value = false
  }
}
</script>
