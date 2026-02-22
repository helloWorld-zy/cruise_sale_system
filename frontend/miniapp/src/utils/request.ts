export const buildUrl = (path: string) => {
    const baseUrl = 'http://localhost:8080/api/v1' // DEV
    return `${baseUrl}${path}`
}

declare const uni: any

export const request = <T>(path: string, options: any = {}) => {
    return new Promise<T>((resolve, reject) => {
        uni.request({
            url: buildUrl(path),
            ...options,
            success: (res: any) => resolve(res.data as T),
            fail: reject,
        })
    })
}
