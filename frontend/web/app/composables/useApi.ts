// web/app/composables/useApi.ts — 前台 API 请求组合式函数
// 封装 $fetch 请求，自动注入 API 基础路径（无需认证 Token）

// 声明 Nuxt 框架自动导入的全局函数
declare const useRuntimeConfig: any
declare const $fetch: any

/**
 * useApi 提供统一的后端 API 请求方法。
 * 前台接口默认无需身份认证，但支持从 sessionStorage 注入可选 token。
 */
export const useApi = () => {
    const config = useRuntimeConfig()
    // 后端 API 基础路径
    const baseUrl = config.public.apiBase

    /**
     * request 发起 HTTP 请求到后端 API。
     * @param path - 接口路径（如 "/cabins"）
     * @param options - 请求配置项
     * @returns API 响应数据
     */
    const request = async <T>(path: string, options: any = {}) => {
        // 统一处理请求头并按需附带登录令牌。
        const headers: Record<string, string> = {
            'Content-Type': 'application/json',
            ...(options.headers || {}),
        }
        try {
            const token = globalThis?.sessionStorage?.getItem('auth_token')
            if (token && !headers.Authorization) {
                headers.Authorization = `Bearer ${token}`
            }
        } catch {
            // SSR 或不可访问 sessionStorage 场景下忽略 token 注入。
        }

        return await $fetch<T>(`${baseUrl}${path}`, {
            ...options,
            headers,
        })
    }
    return { baseUrl, request }
}
