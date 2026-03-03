// miniapp/src/utils/request.ts — 小程序端 HTTP 请求工具
// 封装 uni.request 为 Promise 风格，提供统一的 API 调用方式

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

// 声明 uni-app 框架的全局 API 对象
declare const uni: any

/**
 * request 发起 HTTP 请求（基于 uni.request 封装）。
 * 将回调风格的 uni.request 包装为 Promise，方便 async/await 使用。
 * @param path - 接口路径
 * @param options - 请求配置（method、data、header 等）
 * @returns Promise<T> - 响应数据
 */
export const request = <T>(path: string, options: any = {}) => {
    const method = String(options?.method || 'GET').toUpperCase()
    const headers = options?.header || options?.headers || {}
    const data = options?.data

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
            header: headers,
            success: (res: any) => resolve(res.data as T), // 成功时返回响应数据
            fail: reject, // 失败时拒绝 Promise
        })
    })
}
