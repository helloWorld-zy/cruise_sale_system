<!-- miniapp/pages/login/login.vue — 小程序端登录页面 -->
<!-- 支持微信一键登录（uni.login）和短信验证码登录双通道 -->
<script setup lang="ts">
import { ref, computed } from 'vue'
import PrimaryButton from '../../components/PrimaryButton.vue'
import NavBar from '../../components/NavBar.vue'
import { request } from '../../src/utils/request'
import { useAuthStore } from '../../src/stores/auth'

declare const uni: any

const emit = defineEmits<{ (e: 'back'): void }>()

const authStore = useAuthStore()
const phone = ref('')
const code = ref('')
const loading = ref(false)
const sending = ref(false)
const countdown = ref(0)
const error = ref('')
const success = ref('')

/** 是否处于微信小程序环境（有 uni.login 可用） */
const isWxMiniApp = computed(() => {
  try {
    return typeof uni !== 'undefined' && typeof uni.login === 'function'
  } catch {
    return false
  }
})

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

/** 短信验证码登录 */
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

/**
 * 微信一键登录：
 * 1. uni.login() 获取微信临时 code
 * 2. 将 code 发送到后端 /users/wx-login 换取业务 token
 * 3. 存入 auth store（自动持久化）
 */
async function handleWxLogin() {
  if (loading.value) return
  error.value = ''
  success.value = ''
  loading.value = true
  try {
    // 第一步：调用 uni.login 获取微信 code
    const loginRes = await new Promise<{ code: string }>((resolve, reject) => {
      uni.login({
        provider: 'weixin',
        success: (res: any) => {
          if (res.code) {
            resolve({ code: res.code })
          } else {
            reject(new Error(res.errMsg || '微信登录失败'))
          }
        },
        fail: (err: any) => reject(new Error(err?.errMsg || '微信登录失败')),
      })
    })

    // 第二步：将 code 发送后端换取 token
    const res = await request('/users/wx-login', {
      method: 'POST',
      data: { code: loginRes.code },
    })
    const token = res?.token ?? res?.data?.token
    if (!token) {
      throw new Error('微信登录响应缺少 token')
    }
    authStore.setToken(token)

    // 可选：设置用户 profile
    const profile = res?.profile ?? res?.data?.profile
    if (profile) {
      authStore.setProfile(profile)
    }

    success.value = '微信登录成功'
  } catch (e: any) {
    error.value = e?.message ?? '微信登录失败'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <view class="page">
    <NavBar title="用户登录" show-back @back="emit('back')" />
    <view class="hero">
      <text class="eyebrow">Azure Deck Member</text>
      <text class="subtitle">输入手机号与验证码，继续你的海上假期计划。</text>
    </view>

    <view class="panel">
      <!-- 微信一键登录（仅在小程序环境可用） -->
      <button v-if="isWxMiniApp" class="wx-btn" :disabled="loading" @click="handleWxLogin">
        {{ loading ? '登录中...' : '微信一键登录' }}
      </button>

      <view v-if="isWxMiniApp" class="divider">
        <text class="divider-line"></text>
        <text class="divider-text">或使用手机号登录</text>
        <text class="divider-line"></text>
      </view>

      <!-- 短信验证码登录 -->
      <input v-model="phone" type="number" placeholder="手机号" :disabled="loading || sending" />
      <view class="code-row">
        <input v-model="code" type="number" placeholder="验证码" :disabled="loading" />
        <button class="sub-btn" :disabled="sending || countdown > 0" @click="sendCode">
          {{ countdown > 0 ? `${countdown}s` : '发送验证码' }}
        </button>
      </view>
      <PrimaryButton @click="handleLogin">{{ loading ? '提交中...' : '登录' }}</PrimaryButton>
      <text v-if="error" class="error">{{ error }}</text>
      <text v-if="success" class="success">{{ success }}</text>
    </view>
  </view>
</template>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  gap: 18rpx;
  padding: 30rpx;
  min-height: 100vh;
  background:
    radial-gradient(circle at 8% 4%, #d9e9f5 0, transparent 30%),
    linear-gradient(180deg, #f3f8fb 0%, #eef3f7 100%);
}

.hero {
  padding: 22rpx 4rpx 6rpx;
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.eyebrow {
  font-size: 20rpx;
  letter-spacing: 0.2em;
  text-transform: uppercase;
  color: #8a6f3f;
}

.title {
  font-size: 48rpx;
  font-weight: 700;
  color: #122b42;
}

.subtitle {
  font-size: 24rpx;
  color: #5a728a;
}

.panel {
  display: flex;
  flex-direction: column;
  gap: 14rpx;
  padding: 24rpx;
  background: #fff;
  border: 1rpx solid #d4e0ea;
  border-radius: 24rpx;
  box-shadow: 0 16rpx 36rpx rgba(16, 47, 72, 0.12);
}

input {
  background: #f9fcff;
  border: 2rpx solid #d0dee8;
  border-radius: 16rpx;
  padding: 20rpx;
}

.code-row {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 10rpx;
}

.sub-btn {
  border-radius: 16rpx;
  padding: 0 20rpx;
  background: #e9f1f8;
  color: #204666;
  font-size: 24rpx;
}

.error {
  color: #c53f57;
}

.success {
  color: #0f8a60;
}

.wx-btn {
  border: 0;
  border-radius: 16rpx;
  padding: 22rpx;
  background: #07c160;
  color: #fff;
  font-size: 28rpx;
  font-weight: 700;
}

.wx-btn[disabled] {
  opacity: 0.6;
}

.divider {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin: 8rpx 0;
}

.divider-line {
  flex: 1;
  height: 1rpx;
  background: #d0dee8;
}

.divider-text {
  font-size: 22rpx;
  color: #9ab;
  white-space: nowrap;
}
</style>
