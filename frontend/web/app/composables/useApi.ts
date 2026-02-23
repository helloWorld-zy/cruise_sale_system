// web/app/composables/useApi.ts — 前台 API 请求组合式函数
// 封装 $fetch 请求，自动注入 API 基础路径（无需认证 Token）

// 声明 Nuxt 框架自动导入的全局函数
declare const useRuntimeConfig: any
declare const $fetch: any

/**
 * useApi 提供统一的后端 API 请求方法。
 * 与 admin 版本不同，前台接口无需身份认证。
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
        return await $fetch<T>(`${baseUrl}${path}`, options)
    }
    return { baseUrl, request }
}
