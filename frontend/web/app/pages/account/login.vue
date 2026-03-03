<!-- web/pages/account/login.vue — 前台用户登录页面 -->
<!-- 组合 LoginForm 组件承载短信验证码登录流程 -->
<script setup lang="ts">
import LoginForm from '../../../components/LoginForm.vue'
</script>

<template>
  <div class="grid gap-6 md:grid-cols-[1.1fr_1fr]">
    <section class="rounded-3xl border border-[#d7e2ea] bg-gradient-to-br from-[#0f3d5c] via-[#1d567a] to-[#2a6d91] p-8 text-white shadow-[0_24px_60px_rgba(16,47,72,0.28)]">
      <p class="mb-2 text-xs uppercase tracking-[0.3em] text-[#d4e6f5]">Member Deck</p>
      <h1 class="font-['Bodoni_MT','Times_New_Roman',serif] text-4xl leading-tight">欢迎回到你的航程中枢</h1>
      <p class="mt-4 max-w-md text-sm text-[#d9e8f5]">管理订单、收藏航线与会员权益。一次登录，继续你的海上假期计划。</p>
      <div class="mt-6 grid grid-cols-2 gap-3 text-xs">
        <div class="rounded-2xl bg-white/10 p-3">实时订单状态</div>
        <div class="rounded-2xl bg-white/10 p-3">价格与舱位提醒</div>
      </div>
    </section>

    <section class="page">
      <p class="lux-badge mb-3">Secure Sign-In</p>
      <h2 class="mb-3 font-['Bodoni_MT','Times_New_Roman',serif] text-3xl text-[#12263a]">手机号快捷登录</h2>
      <p class="mb-5 text-sm text-[#5a7189]">验证码仅用于本次登录验证，不会泄露给第三方。</p>
      <LoginForm />
    </section>
  </div>
</template>
