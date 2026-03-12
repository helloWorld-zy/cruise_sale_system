// miniapp/src/stores/auth.ts — 小程序端认证状态管理
// 使用 Pinia Composition API 风格管理用户认证状态
// 通过 uni.setStorageSync / uni.getStorageSync 实现 token 持久化

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

declare const uni: any

const TOKEN_KEY = 'auth_token'
const PROFILE_KEY = 'auth_profile'

/** 读取本地存储，优先 uni，降级 localStorage */
function storageGet(key: string): string {
    try {
        if (typeof uni !== 'undefined' && typeof uni.getStorageSync === 'function') {
            return uni.getStorageSync(key) ?? ''
        }
    } catch { /* ignored */ }
    try {
        return globalThis.localStorage?.getItem(key) ?? ''
    } catch { /* ignored */ }
    return ''
}

/** 写入本地存储，优先 uni，降级 localStorage */
function storageSet(key: string, value: string) {
    try {
        if (typeof uni !== 'undefined' && typeof uni.setStorageSync === 'function') {
            uni.setStorageSync(key, value)
            return
        }
    } catch { /* ignored */ }
    try {
        globalThis.localStorage?.setItem(key, value)
    } catch { /* ignored */ }
}

/** 删除本地存储，优先 uni，降级 localStorage */
function storageRemove(key: string) {
    try {
        if (typeof uni !== 'undefined' && typeof uni.removeStorageSync === 'function') {
            uni.removeStorageSync(key)
        }
    } catch { /* ignored */ }
    try {
        globalThis.localStorage?.removeItem(key)
    } catch { /* ignored */ }
}

/**
 * useAuthStore 管理小程序端的用户认证信息。
 * - token: 认证令牌（由微信小程序登录接口或短信登录获取）
 * - profile: 用户基本信息
 * token 和 profile 会持久化到本地存储，页面刷新/冷启动后自动恢复。
 */
export const useAuthStore = defineStore('auth', () => {
    // 从本地存储恢复 token
    const token = ref(storageGet(TOKEN_KEY))
    // 从本地存储恢复 profile
    const savedProfile = storageGet(PROFILE_KEY)
    const profile = ref<null | { id: number; username: string; roles: string[] }>(
        savedProfile ? JSON.parse(savedProfile) : null,
    )

    // setToken 设置认证令牌并持久化
    function setToken(t: string) {
        token.value = t
        if (t) {
            storageSet(TOKEN_KEY, t)
        } else {
            storageRemove(TOKEN_KEY)
        }
    }

    // setProfile 设置用户信息并持久化
    function setProfile(p: { id: number; username: string; roles: string[] }) {
        profile.value = p
        storageSet(PROFILE_KEY, JSON.stringify(p))
    }

    // logout 清空认证状态并清除本地存储
    function logout() {
        token.value = ''
        profile.value = null
        storageRemove(TOKEN_KEY)
        storageRemove(PROFILE_KEY)
    }

    // isAuthenticated 判断是否已登录
    const isAuthenticated = computed(() => !!token.value)
    // roles 获取当前用户角色列表
    const roles = computed(() => profile.value?.roles ?? [])

    return { token, profile, setToken, setProfile, logout, isAuthenticated, roles }
})
