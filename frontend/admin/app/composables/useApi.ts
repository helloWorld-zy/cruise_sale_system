import { useAuthStore } from '../stores/auth'

// Mock auto-imports for vitest environment
declare const useRuntimeConfig: any
declare const $fetch: any

export const useApi = () => {
    const config = useRuntimeConfig()
    const auth = useAuthStore()

    const baseUrl = config.public.apiBase

    const request = async <T>(path: string, options: any = {}) => {
        const headers: Record<string, string> = {
            'Content-Type': 'application/json',
        }
        if (auth.token) headers.Authorization = `Bearer ${auth.token}`

        return await $fetch<T>(`${baseUrl}${path}`, {
            ...options,
            headers: { ...headers, ...(options.headers || {}) },
        })
    }

    return { baseUrl, request }
}
