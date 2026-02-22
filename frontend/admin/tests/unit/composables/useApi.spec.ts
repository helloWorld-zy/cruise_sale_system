import { describe, test, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useApi } from '../../../app/composables/useApi'

// Mock globals
vi.stubGlobal('useRuntimeConfig', () => ({
    public: { apiBase: '/api/v1' }
}))

describe('useApi', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
    })

    test('exposes baseUrl', () => {
        const api = useApi()
        expect(api.baseUrl).toContain('/api/v1')
    })
})
