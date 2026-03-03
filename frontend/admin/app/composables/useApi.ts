// admin/app/composables/useApi.ts — API 请求组合式函数
// 封装 $fetch 请求，自动注入认证 Token 和 API 基础路径

import { useAuthStore } from '../stores/auth'

/**
 * useApi 提供统一的后端 API 请求方法。
 * - 自动从运行时配置中读取 API 基础路径
 * - 自动附加 JSON Content-Type 请求头
 * - 若用户已登录，自动附加 Bearer 令牌
 */
export const useApi = () => {
    const runtimeConfigFactory = (globalThis as any).useRuntimeConfig
    const config = typeof runtimeConfigFactory === 'function'
        ? runtimeConfigFactory()
        : { public: { apiBase: '' } }
    let tokenFromStore = ''
    try {
        const auth = useAuthStore()
        tokenFromStore = auth.token || ''
    } catch {
        tokenFromStore = ''
    }
    const fetcher = (globalThis as any).$fetch

    // 后端 API 基础路径（如 http://localhost:8080/api/v1）
    const baseUrl = config.public.apiBase

    /**
     * request 发起 HTTP 请求到后端 API。
     * @param path - 接口路径（如 "/cruises"）
     * @param options - 请求配置项（method、body 等）
     * @returns API 响应数据
     */
    const request = async <T>(path: string, options: any = {}) => {
        // 后台页面默认走 /admin 命名空间，避免漏写导致 404。
        const normalizedPath =
            path.startsWith('/admin/') ||
            path.startsWith('/users') ||
            path.startsWith('/pay') ||
            path.startsWith('/refund')
                ? path
                : `/admin${path}`

        // 构建请求头，默认 JSON 格式
        const headers: Record<string, string> = {
            'Content-Type': 'application/json',
        }
        // 若已登录，附加认证令牌
        const localStorageToken = typeof window !== 'undefined' && typeof (window as any).localStorage?.getItem === 'function'
            ? (window as any).localStorage.getItem('admin_token') || ''
            : ''
        const token = tokenFromStore || localStorageToken
        if (token) headers.Authorization = `Bearer ${token}`

        try {
            if (typeof fetcher !== 'function') {
                throw new Error('$fetch is not available')
            }
            return await fetcher(`${baseUrl}${normalizedPath}`, {
                ...options,
                headers: { ...headers, ...(options.headers || {}) },
            })
        } catch (err: any) {
            const message =
                err?.data?.message ||
                err?.response?._data?.message ||
                err?.message ||
                'request failed'
            throw new Error(message)
        }
    }

    return { baseUrl, request }
}
