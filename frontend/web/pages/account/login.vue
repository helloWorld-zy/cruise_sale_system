<!-- web/pages/account/login.vue — 前台用户登录页面 -->
<!-- 组合 LoginForm 组件承载短信验证码登录流程 -->
<script setup lang="ts">
import LoginForm from '../../components/LoginForm.vue'
</script>

<template>
  <div class="page">
    <h1>Login</h1>
    <LoginForm />
  </div>
</template>
