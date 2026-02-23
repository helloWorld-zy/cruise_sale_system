// admin/app/stores/auth.ts — 认证状态管理
// 使用 Pinia Composition API 风格管理用户认证状态

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

// ME-04：使用 Pinia Composition API 风格，与 <script setup> 保持一致
export const useAuthStore = defineStore('auth', () => {
    // 当前用户的 JWT 认证令牌
    const token = ref('')
    // 当前登录用户的个人信息（null 表示未登录）
    const profile = ref<null | { id: number; username: string; roles: string[] }>(null)

    // setToken 设置认证令牌（登录成功后调用）
    function setToken(t: string) {
        token.value = t
    }

    // setProfile 设置用户个人信息
    function setProfile(p: { id: number; username: string; roles: string[] }) {
        profile.value = p
    }

    // logout 清空认证状态（退出登录时调用）
    function logout() {
        token.value = ''
        profile.value = null
    }

    // isAuthenticated 计算属性：判断用户是否已登录
    const isAuthenticated = computed(() => !!token.value)
    // roles 计算属性：获取当前用户的角色列表
    const roles = computed(() => profile.value?.roles ?? [])

    return { token, profile, setToken, setProfile, logout, isAuthenticated, roles }
})
