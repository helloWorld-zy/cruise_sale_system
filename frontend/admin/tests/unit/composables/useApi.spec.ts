import { describe, test, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useApi } from '../../../app/composables/useApi'
import { useAuthStore } from '../../../app/stores/auth'

// Mock globals
vi.stubGlobal('useRuntimeConfig', () => ({
    public: { apiBase: '/api/v1' }
}))

const mockFetch = vi.fn().mockResolvedValue({ success: true })
vi.stubGlobal('$fetch', mockFetch)

describe('useApi', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        mockFetch.mockClear()
    })

    test('exposes baseUrl', () => {
        const api = useApi()
        expect(api.baseUrl).toContain('/api/v1')
    })

    test('makes request with token', async () => {
        const auth = useAuthStore()
        auth.setToken('dummy-token')

        const api = useApi()
        const res = await api.request('/test', { method: 'POST' })

        expect(mockFetch).toHaveBeenCalledWith('/api/v1/test', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer dummy-token'
            }
        })
        expect(res).toEqual({ success: true })
    })

    test('makes request without token', async () => {
        const api = useApi()
        await api.request('/test')

        expect(mockFetch).toHaveBeenCalledWith('/api/v1/test', {
            headers: {
                'Content-Type': 'application/json'
            }
        })
    })
})
