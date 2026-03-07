<template>
  <section class="login-page">
    <div class="login-card">
      <h1 class="login-card__title">Cruise Booking Admin</h1>
      <p class="login-card__subtitle">系统登录</p>
      <form class="login-form" @submit.prevent="handleLogin">
        <label class="login-form__label" for="username">用户名</label>
        <input id="username" v-model="username" placeholder="请输入用户名" autocomplete="username" />
        <label class="login-form__label" for="password">密码</label>
        <input id="password" v-model="password" type="password" placeholder="请输入密码" autocomplete="current-password" />
        <p v-if="error" class="login-form__error">{{ error }}</p>
        <button class="login-form__submit" type="submit" :disabled="loading">{{ loading ? '登录中...' : '登录' }}</button>
      </form>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useApi } from '../composables/useApi'

definePageMeta({
  layout: 'auth',
})

const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)
const authStore = useAuthStore()
const { request } = useApi()

async function handleLogin() {
  if (!username.value || !password.value) {
    error.value = '请填写用户名和密码'
    return
  }
  loading.value = true
  error.value = ''
  try {
    const res = await request('/admin/auth/login', {
      method: 'POST',
      body: {
        username: username.value,
        password: password.value,
      },
    })
    const token = res?.token ?? res?.data?.token
    if (!token) {
      throw new Error('登录响应缺少 token')
    }
    authStore.setToken(token)
    await navigateTo('/cruises')
  } catch (e: any) {
    error.value = e?.message ?? '登录失败'
  } finally {
    loading.value = false
  }
}
</script>

