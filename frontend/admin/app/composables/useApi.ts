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
    const runtimeConfigFactory = typeof useRuntimeConfig === 'function'
        ? useRuntimeConfig
        : ((globalThis as any).useRuntimeConfig as any)
    const config = typeof runtimeConfigFactory === 'function'
        ? runtimeConfigFactory()
        : { public: { apiBase: '/api/v1' } }
    let tokenFromStore = ''
    try {
        const auth = useAuthStore()
        tokenFromStore = auth.token || ''
    } catch {
        tokenFromStore = ''
    }
    const fetcher = typeof $fetch === 'function' ? $fetch : (globalThis as any).$fetch

    // 后端 API 基础路径（默认同源 /api/v1）
    const baseUrl = config?.public?.apiBase || '/api/v1'

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

        const isFormData = typeof FormData !== 'undefined' && options?.body instanceof FormData
        // 构建请求头：JSON 请求默认附加 Content-Type；FormData 交给浏览器自动生成 boundary。
        const headers: Record<string, string> = isFormData
            ? {}
            : {
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
            const code = Number(err?.data?.code ?? err?.response?._data?.code ?? 0)
            const status = Number(err?.status ?? err?.response?.status ?? 0)
            const message =
                err?.data?.message ||
                err?.response?._data?.message ||
                err?.message ||
                'request failed'
            const wrappedError = new Error(message) as Error & { code?: number; status?: number; raw?: unknown }
            if (Number.isFinite(code) && code > 0) wrappedError.code = code
            if (Number.isFinite(status) && status > 0) wrappedError.status = status
            wrappedError.raw = err
            throw wrappedError
        }
    }

    return { baseUrl, request }
}
