declare const useRuntimeConfig: any
declare const $fetch: any

export const useApi = () => {
    const config = useRuntimeConfig()
    const baseUrl = config.public.apiBase
    const request = async <T>(path: string, options: any = {}) => {
        return await $fetch<T>(`${baseUrl}${path}`, options)
    }
    return { baseUrl, request }
}
