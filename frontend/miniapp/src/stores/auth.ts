// miniapp/src/stores/auth.ts — 小程序端认证状态管理
// 使用 Pinia Composition API 风格管理用户认证状态

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

/**
 * useAuthStore 管理小程序端的用户认证信息。
 * - token: 认证令牌（由微信小程序登录接口获取）
 * - profile: 用户基本信息
 */
export const useAuthStore = defineStore('auth', () => {
    // 用户认证令牌
    const token = ref('')
    // 用户信息（null 表示未登录）
    const profile = ref<null | { id: number; username: string; roles: string[] }>(null)

    // setToken 设置认证令牌
    function setToken(t: string) {
        token.value = t
    }

    // setProfile 设置用户信息
    function setProfile(p: { id: number; username: string; roles: string[] }) {
        profile.value = p
    }

    // logout 清空认证状态
    function logout() {
        token.value = ''
        profile.value = null
    }

    // isAuthenticated 判断是否已登录
    const isAuthenticated = computed(() => !!token.value)
    // roles 获取当前用户角色列表
    const roles = computed(() => profile.value?.roles ?? [])

    return { token, profile, setToken, setProfile, logout, isAuthenticated, roles }
})
