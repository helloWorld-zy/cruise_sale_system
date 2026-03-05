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
        Object.defineProperty(window, 'localStorage', {
            value: {
                getItem: () => null,
                setItem: () => undefined,
                removeItem: () => undefined,
            },
            configurable: true,
        })
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

        expect(mockFetch).toHaveBeenCalledWith('/api/v1/admin/test', {
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

        expect(mockFetch).toHaveBeenCalledWith('/api/v1/admin/test', {
            headers: {
                'Content-Type': 'application/json'
            }
        })
    })

    test('does not force JSON content-type for FormData body', async () => {
        const api = useApi()
        const form = new FormData()
        form.append('file', new Blob(['x'], { type: 'image/png' }), 'logo.png')

        await api.request('/upload/image', {
            method: 'POST',
            body: form,
        })

        expect(mockFetch).toHaveBeenCalledWith('/api/v1/admin/upload/image', {
            method: 'POST',
            body: form,
            headers: {},
        })
    })
})
