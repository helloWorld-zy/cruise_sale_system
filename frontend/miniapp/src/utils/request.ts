// miniapp/src/utils/request.ts — 小程序端 HTTP 请求工具
// 封装 uni.request 为 Promise 风格，提供统一的 API 调用方式
// 自动注入 Authorization header（从 uni.getStorageSync 读取 token）

declare const uni: any

const TOKEN_KEY = 'auth_token'

/** 从本地存储读取 token，优先 uni，降级 localStorage */
function getToken(): string {
    try {
        if (typeof uni !== 'undefined' && typeof uni.getStorageSync === 'function') {
            return uni.getStorageSync(TOKEN_KEY) ?? ''
        }
    } catch { /* ignored */ }
    try {
        return globalThis.localStorage?.getItem(TOKEN_KEY) ?? ''
    } catch { /* ignored */ }
    return ''
}

/**
 * buildUrl 拼接完整的 API 请求地址。
 * @param path - 接口路径（如 "/cruises"）
 * @returns 完整的请求 URL
 */
export const buildUrl = (path: string) => {
    // 后端 API 基础路径（开发环境配置）
    const baseUrl = 'http://localhost:8080/api/v1' // DEV
    return `${baseUrl}${path}`
}

export const buildAssetUrl = (path?: string) => {
    if (!path) return ''
    if (/^https?:\/\//i.test(path)) return path
    const apiBaseUrl = 'http://localhost:8080/api/v1'
    const origin = apiBaseUrl.replace(/\/api\/v1\/?$/, '')
    const normalizedPath = path.startsWith('/') ? path : `/${path}`
    return `${origin}${normalizedPath}`
}

/**
 * request 发起 HTTP 请求（基于 uni.request 封装）。
 * 将回调风格的 uni.request 包装为 Promise，方便 async/await 使用。
 * 自动在 header 中注入 Authorization（Bearer token），无需调用方手动传入。
 * @param path - 接口路径
 * @param options - 请求配置（method、data、header 等）
 * @returns Promise<T> - 响应数据
 */
export const request = <T>(path: string, options: any = {}) => {
    const method = String(options?.method || 'GET').toUpperCase()
    const userHeaders = options?.header || options?.headers || {}
    const data = options?.data

    // 自动注入 Authorization header（已有则不覆盖）
    const token = getToken()
    const headers: Record<string, string> = { ...userHeaders }
    if (token && !headers['Authorization']) {
        headers['Authorization'] = `Bearer ${token}`
    }

    // 浏览器预览环境没有 uni 对象，降级为 fetch，避免整页崩溃。
    if (typeof (globalThis as any).uni?.request !== 'function') {
        const init: RequestInit = {
            method,
            headers: {
                'Content-Type': 'application/json',
                ...headers,
            },
        }
        if (data != null && method !== 'GET') {
            init.body = typeof data === 'string' ? data : JSON.stringify(data)
        }
        return fetch(buildUrl(path), init).then(async (res) => {
            if (!res.ok) {
                throw new Error(`HTTP ${res.status}`)
            }
            return (await res.json()) as T
        })
    }

    return new Promise<T>((resolve, reject) => {
        ;(globalThis as any).uni.request({
            url: buildUrl(path),
            method,
            data,
            header: {
                'Content-Type': 'application/json',
                ...headers,
            },
            success: (res: any) => resolve(res.data as T), // 成功时返回响应数据
            fail: reject, // 失败时拒绝 Promise
        })
    })
}
