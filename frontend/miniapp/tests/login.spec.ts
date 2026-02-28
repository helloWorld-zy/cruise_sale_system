import { describe, it, expect, beforeEach, vi } from 'vitest'
import { render } from '@testing-library/vue'
import { createPinia, setActivePinia } from 'pinia'
import Login from '../pages/login/login.vue'

const mockRequest = vi.fn().mockResolvedValue({})
vi.mock('../src/utils/request', () => ({
    request: (...args: any[]) => mockRequest(...args),
}))

beforeEach(() => {
    setActivePinia(createPinia())
    mockRequest.mockClear()
})

describe('Miniapp Login', () => {
    it('shows login title', () => {
        const { getByText } = render(Login)
        expect(getByText('Login')).toBeTruthy()
    })
})
