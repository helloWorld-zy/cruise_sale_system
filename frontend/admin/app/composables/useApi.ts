// admin/app/composables/useApi.ts — API 请求组合式函数
// 封装 $fetch 请求，自动注入认证 Token 和 API 基础路径

import { useAuthStore } from '../stores/auth'

// 声明 Nuxt 框架自动导入的全局函数
declare const useRuntimeConfig: any
declare const $fetch: any

/**
 * useApi 提供统一的后端 API 请求方法。
 * - 自动从运行时配置中读取 API 基础路径
 * - 自动附加 JSON Content-Type 请求头
 * - 若用户已登录，自动附加 Bearer 令牌
 */
export const useApi = () => {
    const config = useRuntimeConfig()
    const auth = useAuthStore()

    // 后端 API 基础路径（如 http://localhost:8080/api/v1）
    const baseUrl = config.public.apiBase

    /**
     * request 发起 HTTP 请求到后端 API。
     * @param path - 接口路径（如 "/cruises"）
     * @param options - 请求配置项（method、body 等）
     * @returns API 响应数据
     */
    const request = async <T>(path: string, options: any = {}) => {
        // 构建请求头，默认 JSON 格式
        const headers: Record<string, string> = {
            'Content-Type': 'application/json',
        }
        // 若已登录，附加认证令牌
        if (auth.token) headers.Authorization = `Bearer ${auth.token}`

        return await $fetch<T>(`${baseUrl}${path}`, {
            ...options,
            headers: { ...headers, ...(options.headers || {}) },
        })
    }

    return { baseUrl, request }
}
