<!-- miniapp/components/NavBar.vue — 统一小程序顶部导航栏 -->
<!-- 参照微信小程序标准导航栏规范，所有页面共用 -->
<script setup lang="ts">
import { ChevronLeft } from 'lucide-vue-next'

defineProps<{
  /** 页面标题 */
  title: string
  /** 是否显示返回按钮（详情页传 true，Tab 页传 false） */
  showBack?: boolean
  /** 透明模式（用于有顶部大图的页面） */
  transparent?: boolean
}>()

const emit = defineEmits<{
  (e: 'back'): void
}>()
</script>

<template>
  <header
    class="navbar"
    :class="transparent ? 'navbar--transparent' : 'navbar--solid'"
  >
    <!-- 状态栏占位 -->
    <div class="navbar__status-bar"></div>

    <!-- 内容区：44px 标准高度 -->
    <div class="navbar__content">
      <!-- 左侧：返回按钮 -->
      <div class="navbar__left">
        <button v-if="showBack" class="navbar__back" @click="emit('back')">
          <ChevronLeft class="w-5 h-5" />
        </button>
      </div>

      <!-- 居中标题 -->
      <div class="navbar__title">{{ title }}</div>

      <!-- 右侧：胶囊按钮占位（模拟微信原生胶囊） -->
      <div class="navbar__right">
        <div class="navbar__capsule">
          <span class="navbar__capsule-dots">
            <span></span><span></span><span></span>
          </span>
          <span class="navbar__capsule-divider"></span>
          <span class="navbar__capsule-avatar"></span>
        </div>
      </div>
    </div>
  </header>
  <!-- 占位块，避免页面内容被 fixed 的 navbar 遮挡 -->
  <div v-if="!transparent" class="navbar__placeholder"></div>
</template>

<style scoped>
/* ====== 结构 ====== */
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 50;
}

.navbar__status-bar {
  height: 12px; /* 模拟状态栏 */
}

.navbar__content {
  display: flex;
  align-items: center;
  height: 44px;
  padding: 0 12px;
}

.navbar__left,
.navbar__right {
  width: 80px;
  flex-shrink: 0;
}

.navbar__right {
  display: flex;
  justify-content: flex-end;
}

.navbar__title {
  flex: 1;
  text-align: center;
  font-size: 17px;
  font-weight: 600;
  line-height: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.navbar__placeholder {
  height: 56px; /* status-bar 12 + content 44 */
}

/* ====== 返回按钮 ====== */
.navbar__back {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: 0;
  border-radius: 50%;
  background: transparent;
  cursor: pointer;
  padding: 0;
  transition: background 0.15s;
}

.navbar__back:active {
  background: rgba(0, 0, 0, 0.05);
}

/* ====== 胶囊按钮 ====== */
.navbar__capsule {
  display: flex;
  align-items: center;
  height: 30px;
  border-radius: 15px;
  padding: 0 10px;
  gap: 6px;
  border: 1px solid rgba(0, 0, 0, 0.12);
  background: rgba(255, 255, 255, 0.6);
}

.navbar__capsule-dots {
  display: flex;
  gap: 3px;
  align-items: center;
}

.navbar__capsule-dots span {
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: #353535;
}

.navbar__capsule-divider {
  width: 1px;
  height: 16px;
  background: rgba(0, 0, 0, 0.12);
}

.navbar__capsule-avatar {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  border: 1.5px solid #353535;
  position: relative;
}

.navbar__capsule-avatar::after {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #353535;
}

/* ====== 实心模式（白底） ====== */
.navbar--solid {
  background: #fff;
  box-shadow: 0 1px 0 rgba(0, 0, 0, 0.06);
}

.navbar--solid .navbar__title {
  color: #1a1a1a;
}

.navbar--solid .navbar__back {
  color: #1a1a1a;
}

/* ====== 透明模式（浮在大图上方） ====== */
.navbar--transparent {
  background: transparent;
}

.navbar--transparent .navbar__title {
  color: #fff;
  text-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
}

.navbar--transparent .navbar__back {
  color: #fff;
}

.navbar--transparent .navbar__capsule {
  border-color: rgba(255, 255, 255, 0.5);
  background: rgba(0, 0, 0, 0.15);
}

.navbar--transparent .navbar__capsule-dots span {
  background: #fff;
}

.navbar--transparent .navbar__capsule-divider {
  background: rgba(255, 255, 255, 0.4);
}

.navbar--transparent .navbar__capsule-avatar {
  border-color: #fff;
}

.navbar--transparent .navbar__capsule-avatar::after {
  background: #fff;
}
</style>
