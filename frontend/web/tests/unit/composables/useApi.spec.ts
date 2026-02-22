import { describe, test, expect, vi } from 'vitest'
import { useApi } from '../../../app/composables/useApi'

vi.stubGlobal('useRuntimeConfig', () => ({
    public: { apiBase: '/api/v1' }
}))

describe('useApi', () => {
    test('exposes baseUrl', () => {
        const api = useApi()
        expect(api.baseUrl).toContain('/api/v1')
    })
})
