import { beforeEach, describe, expect, it, vi } from 'vitest'

const mockUseAuthStore = vi.fn()
const mockNavigateTo = vi.fn()

vi.mock('../../../app/stores/auth', () => ({
  useAuthStore: () => mockUseAuthStore(),
}))

vi.stubGlobal('defineNuxtRouteMiddleware', (handler: unknown) => handler)
vi.stubGlobal('navigateTo', mockNavigateTo)

describe('auth.global middleware', () => {
  beforeEach(() => {
    mockNavigateTo.mockReset()
    mockUseAuthStore.mockReset()
    Object.defineProperty(window, 'localStorage', {
      value: {
        getItem: vi.fn(() => 'cached-admin-token'),
        setItem: vi.fn(),
        removeItem: vi.fn(),
      },
      configurable: true,
    })
  })

  it('allows visiting login page without restoring cached token or redirecting', async () => {
    const setToken = vi.fn()
    mockUseAuthStore.mockReturnValue({
      token: '',
      isAuthenticated: false,
      setToken,
    })

    const mod = await import('../../../app/middleware/auth.global')
    const middleware = mod.default as (to: { path: string }) => unknown
    const result = middleware({ path: '/login' })

    expect(result).toBeUndefined()
    expect(setToken).not.toHaveBeenCalled()
    expect(mockNavigateTo).not.toHaveBeenCalled()
  })

  it('restores cached token for protected routes and allows access', async () => {
    const state = {
      token: '',
      isAuthenticated: false,
      setToken: vi.fn((token: string) => {
        state.token = token
        state.isAuthenticated = true
      }),
    }
    mockUseAuthStore.mockImplementation(() => state)

    const mod = await import('../../../app/middleware/auth.global')
    const middleware = mod.default as (to: { path: string }) => unknown
    const result = middleware({ path: '/dashboard' })

    expect(result).toBeUndefined()
    expect(state.setToken).toHaveBeenCalledWith('cached-admin-token')
    expect(mockNavigateTo).not.toHaveBeenCalled()
  })

  it('redirects unauthenticated protected routes to login', async () => {
    Object.defineProperty(window, 'localStorage', {
      value: {
        getItem: vi.fn(() => ''),
        setItem: vi.fn(),
        removeItem: vi.fn(),
      },
      configurable: true,
    })
    mockUseAuthStore.mockReturnValue({
      token: '',
      isAuthenticated: false,
      setToken: vi.fn(),
    })
    mockNavigateTo.mockReturnValue('/login')

    const mod = await import('../../../app/middleware/auth.global')
    const middleware = mod.default as (to: { path: string }) => unknown
    const result = middleware({ path: '/dashboard' })

    expect(mockNavigateTo).toHaveBeenCalledWith('/login', { replace: true })
    expect(result).toBe('/login')
  })
})