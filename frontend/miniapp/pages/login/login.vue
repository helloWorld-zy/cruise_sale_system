<!-- miniapp/pages/login/login.vue — 小程序端登录页面 -->
<!-- 当前为轻量入口，复用 PrimaryButton 组件承载后续登录交互 -->
<script setup lang="ts">
import { ref } from 'vue'
import PrimaryButton from '../../components/PrimaryButton.vue'
import { request } from '../../src/utils/request'
import { useAuthStore } from '../../src/stores/auth'

// 小程序登录：短信验证码发送、倒计时与登录提交。
const authStore = useAuthStore()
const phone = ref('')
const code = ref('')
const loading = ref(false)
const sending = ref(false)
const countdown = ref(0)
const error = ref('')
const success = ref('')

function isValidPhone(value: string) {
  return /^1\d{10}$/.test(value)
}

function startCountdown() {
  countdown.value = 60
  const timer = setInterval(() => {
    countdown.value -= 1
    if (countdown.value <= 0) {
      clearInterval(timer)
    }
  }, 1000)
}

async function sendCode() {
  if (sending.value || countdown.value > 0) return
  error.value = ''
  success.value = ''
  if (!isValidPhone(phone.value)) {
    error.value = '请输入 11 位手机号'
    return
  }
  sending.value = true
  try {
    await request('/users/sms-code', {
      method: 'POST',
      data: { phone: phone.value },
    })
    success.value = '验证码已发送'
    startCountdown()
  } catch (e: any) {
    error.value = e?.message ?? '发送验证码失败'
  } finally {
    sending.value = false
  }
}

async function handleLogin() {
  if (loading.value) return
  error.value = ''
  success.value = ''
  if (!isValidPhone(phone.value)) {
    error.value = '请输入 11 位手机号'
    return
  }
  if (!code.value) {
    error.value = '请输入验证码'
    return
  }
  loading.value = true
  try {
    const res = await request('/users/login', {
      method: 'POST',
      data: { phone: phone.value, code: code.value },
    })
    const token = res?.token ?? res?.data?.token
    if (!token) {
      throw new Error('登录响应缺少 token')
    }
    authStore.setToken(token)
    success.value = '登录成功'
  } catch (e: any) {
    error.value = e?.message ?? '登录失败'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <view class="page">
    <text class="title">用户登录</text>
    <input v-model="phone" type="number" placeholder="手机号" :disabled="loading || sending" />
    <input v-model="code" type="number" placeholder="验证码" :disabled="loading" />
    <button class="sub-btn" :disabled="sending || countdown > 0" @click="sendCode">
      {{ countdown > 0 ? `${countdown}s` : '发送验证码' }}
    </button>
    <PrimaryButton @click="handleLogin">{{ loading ? '提交中...' : '登录' }}</PrimaryButton>
    <text v-if="error" class="error">{{ error }}</text>
    <text v-if="success" class="success">{{ success }}</text>
  </view>
</template>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
  padding: 40rpx;
  min-height: 100vh;
  background: linear-gradient(180deg, #eff7ff 0%, #f9fcff 100%);
}

.title {
  font-size: 44rpx;
  font-weight: 700;
  color: #16395b;
}

input {
  background: #fff;
  border: 2rpx solid #c9dff5;
  border-radius: 16rpx;
  padding: 20rpx;
}

.sub-btn {
  border-radius: 16rpx;
  background: #e8f0ff;
  color: #2a4d74;
}

.error {
  color: #d13e5b;
}

.success {
  color: #0f8b62;
}
</style>
