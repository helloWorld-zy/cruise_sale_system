<!-- web/components/LoginForm.vue — 短信验证码登录表单组件 -->
<!-- H-01 修复：完整 SMS 登录流程（手机号验证、倒计时、防重按、错误展示） -->
<script setup lang="ts">
import { ref, computed } from 'vue'

// 短信验证码登录表单：手机号输入 → 获取验证码 → 登录

const phone = ref('')
const code = ref('')
const loading = ref(false)
const sending = ref(false)
const countdown = ref(0)
const errorMsg = ref('')
const successMsg = ref('')

// 手机号格式验证（1[3-9] + 9位）
const isValidPhone = computed(() => /^1[3-9]\d{9}$/.test(phone.value))
const canSendCode = computed(() => isValidPhone.value && countdown.value === 0 && !sending.value)
const canSubmit = computed(() => isValidPhone.value && code.value.length === 6 && !loading.value)

let timer: ReturnType<typeof setInterval> | null = null

async function sendCode() {
  if (!canSendCode.value) return
  errorMsg.value = ''
  sending.value = true
  try {
    await $fetch('/api/v1/users/sms-code', {
      method: 'POST',
      body: { phone: phone.value },
    })
    countdown.value = 60
    timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(timer!)
        timer = null
      }
    }, 1000)
  } catch (e: any) {
    errorMsg.value = e?.data?.message ?? e?.message ?? '发送失败，请重试'
  } finally {
    sending.value = false
  }
}

async function handleSubmit() {
  if (!canSubmit.value) return
  errorMsg.value = ''
  loading.value = true
  try {
    const res = await $fetch<{ data: { token: string } }>('/api/v1/users/login', {
      method: 'POST',
      body: { phone: phone.value, code: code.value },
    })
    successMsg.value = '登录成功'
    // 安全地将 Token 存储到 sessionStorage（避免 localStorage XSS 风险）
    sessionStorage.setItem('auth_token', res?.data?.token ?? '')
    emit('logged-in', res?.data?.token ?? '')
  } catch (e: any) {
    errorMsg.value = e?.data?.message ?? e?.message ?? '登录失败，请重试'
  } finally {
    loading.value = false
  }
}

const emit = defineEmits<{ 'logged-in': [token: string] }>()
</script>

<template>
  <form class="login-form" data-testid="login-form" @submit.prevent="handleSubmit">
    <div class="field">
      <label for="phone">手机号</label>
      <input
        id="phone"
        v-model="phone"
        type="tel"
        inputmode="numeric"
        placeholder="请输入手机号"
        autocomplete="tel"
        :disabled="loading"
      />
      <span v-if="phone && !isValidPhone" class="hint">请输入正确的手机号</span>
    </div>

    <div class="field code-row">
      <label for="code">验证码</label>
      <input
        id="code"
        v-model="code"
        type="text"
        inputmode="numeric"
        maxlength="6"
        placeholder="6位验证码"
        autocomplete="one-time-code"
        :disabled="loading"
      />
      <button
        type="button"
        class="send-btn"
        :disabled="!canSendCode"
        @click="sendCode"
      >
        {{ sending ? '发送中…' : countdown > 0 ? `${countdown}s 后重新发送` : '获取验证码' }}
      </button>
    </div>

    <p v-if="errorMsg" class="error" role="alert">{{ errorMsg }}</p>
    <p v-if="successMsg" class="success">{{ successMsg }}</p>

    <button type="submit" class="submit-btn" :disabled="!canSubmit">
      {{ loading ? '登录中…' : '登录' }}
    </button>
  </form>
</template>
