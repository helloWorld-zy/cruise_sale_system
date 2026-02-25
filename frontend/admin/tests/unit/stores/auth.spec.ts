import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '../../../app/stores/auth'

describe('Auth Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
    })

    it('manages token and profile', () => {
        const auth = useAuthStore()
        expect(auth.isAuthenticated).toBe(false)
        expect(auth.roles).toEqual([])

        auth.setToken('test-token')
        expect(auth.token).toBe('test-token')
        expect(auth.isAuthenticated).toBe(true)

        auth.setProfile({ id: 1, username: 'admin', roles: ['admin'] })
        expect(auth.profile?.username).toBe('admin')
        expect(auth.roles).toEqual(['admin'])

        auth.logout()
        expect(auth.token).toBe('')
        expect(auth.profile).toBeNull()
        expect(auth.isAuthenticated).toBe(false)
    })
})
