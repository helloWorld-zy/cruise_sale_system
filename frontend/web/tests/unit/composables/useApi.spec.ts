import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useApi } from '../../../app/composables/useApi'

vi.stubGlobal('useRuntimeConfig', () => ({
    public: { apiBase: '/api/v1' }
}))

const mockFetch = vi.fn().mockResolvedValue({ data: 'ok' })
vi.stubGlobal('$fetch', mockFetch)

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
        expect(mockFetch).toHaveBeenCalledWith('/api/v1/test', { method: 'POST' })
    })
})
