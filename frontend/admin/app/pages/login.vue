<!-- admin/app/pages/login.vue — 管理员登录页面 -->
<!-- 提供用户名和密码输入表单，验证通过后获取 JWT 令牌 -->
<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50">
    <div class="w-full max-w-md bg-white border rounded-xl p-6">
      <h1 class="text-xl font-semibold mb-4">管理员登录</h1>
      <!-- 登录表单：阻止默认提交行为，改用 handleLogin 处理 -->
      <form class="space-y-3" @submit.prevent="handleLogin">
        <UInput v-model="username" placeholder="用户名" />
        <!-- HI-05 FIX: 使用 v-model="password" 而非 v-model="password.value" -->
        <UInput v-model="password" type="password" placeholder="密码" />
        <!-- 错误提示信息 -->
        <p v-if="error" class="text-red-500 text-sm">{{ error }}</p>
        <UButton class="w-full" color="primary" type="submit" :loading="loading">登录</UButton>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

// 表单状态
const username = ref('')   // 用户名
const password = ref('')   // 密码
const error = ref('')      // 错误提示信息
const loading = ref(false) // 登录中状态

/**
 * handleLogin 处理登录表单提交。
 * 验证必填项 → 调用认证 API → 跳转到仪表盘。
 */
async function handleLogin() {
  // 验证必填项
  if (!username.value || !password.value) {
    error.value = '请填写用户名和密码'
    return
  }
  loading.value = true
  error.value = ''
  // TODO: 在 Sprint 2 集成阶段对接真实认证 API
  await new Promise(resolve => setTimeout(resolve, 500))
  // navigateTo('/dashboard')
  loading.value = false
}
</script>
