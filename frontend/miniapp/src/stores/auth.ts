import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', () => {
    const token = ref('')
    const profile = ref<null | { id: number; username: string; roles: string[] }>(null)

    function setToken(t: string) {
        token.value = t
    }

    function setProfile(p: { id: number; username: string; roles: string[] }) {
        profile.value = p
    }

    function logout() {
        token.value = ''
        profile.value = null
    }

    const isAuthenticated = computed(() => !!token.value)
    const roles = computed(() => profile.value?.roles ?? [])

    return { token, profile, setToken, setProfile, logout, isAuthenticated, roles }
})
