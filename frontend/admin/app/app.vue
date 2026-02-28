<!-- admin/app/app.vue — 应用根组件 -->
<!-- 提供 Nuxt 布局和页面渲染的入口 -->
<template>
  <div>
    <!-- NuxtLayout 提供全局布局框架 -->
    <NuxtLayout>
      <!-- NuxtPage 渲染当前路由对应的页面组件 -->
      <NuxtPage />
    </NuxtLayout>
  </div>
</template>

<style src="./assets/css/main.css"></style>
