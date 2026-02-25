import { describe, it, expect, vi } from 'vitest'
import { request, buildUrl } from '../../src/utils/request'

describe('request utils', () => {
    it('builds full url', () => {
        const url = buildUrl('/test')
        expect(url).toBe('http://localhost:8080/api/v1/test')
    })

    it('makes a successful request', async () => {
        const mockUni = {
            request: vi.fn((opts) => {
                opts.success({ data: { success: true } })
            })
        }
        vi.stubGlobal('uni', mockUni)

        const res = await request('/test', { method: 'POST' })
        expect(res).toEqual({ success: true })
        expect(mockUni.request).toHaveBeenCalledWith(expect.objectContaining({
            url: 'http://localhost:8080/api/v1/test',
            method: 'POST'
        }))
    })

    it('handles request failure', async () => {
        const mockUni = {
            request: vi.fn((opts) => {
                opts.fail(new Error('Network error'))
            })
        }
        vi.stubGlobal('uni', mockUni)

        await expect(request('/test')).rejects.toThrow('Network error')
    })
})
