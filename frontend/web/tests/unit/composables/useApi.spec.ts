import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useApi } from '../../../app/composables/useApi'

vi.stubGlobal('useRuntimeConfig', () => ({
    public: { apiBase: '/api/v1' }
}))

const mockFetch = vi.fn().mockResolvedValue({ data: 'ok' })
vi.stubGlobal('$fetch', mockFetch)
vi.stubGlobal('sessionStorage', {
    getItem: vi.fn(() => null)
})

describe('useApi', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        mockFetch.mockClear()
    })

    it('exposes baseUrl', () => {
        const api = useApi()
        expect(api.baseUrl).toBe('/api/v1')
    })

    it('makes request', async () => {
        const api = useApi()
        await api.request('/test', { method: 'POST' })
        expect(mockFetch).toHaveBeenCalledWith('/api/v1/test', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        })
    })

    it('injects token when sessionStorage has auth_token', async () => {
        vi.stubGlobal('sessionStorage', {
            getItem: vi.fn(() => 'token-1')
        })
        const api = useApi()
        await api.request('/secure')
        expect(mockFetch).toHaveBeenCalledWith('/api/v1/secure', {
            headers: {
                'Content-Type': 'application/json',
                Authorization: 'Bearer token-1'
            }
        })
    })
})
