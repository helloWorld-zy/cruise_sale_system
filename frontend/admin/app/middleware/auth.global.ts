import { useAuthStore } from '../stores/auth'

export default defineNuxtRouteMiddleware((to) => {
  // Admin token is persisted in localStorage, so server-side middleware cannot
  // reliably determine auth state during hard reload.
  if (import.meta.server) {
    return
  }

  const auth = useAuthStore()

  // Ensure token is restored before route guard checks after page reload.
  if (!auth.token && typeof window !== 'undefined') {
    const cachedToken = typeof window.localStorage?.getItem === 'function'
      ? window.localStorage.getItem('admin_token') || ''
      : ''
    if (cachedToken) {
      auth.setToken(cachedToken)
    }
  }

  if (to.path === '/login') {
    if (auth.isAuthenticated) {
      return navigateTo('/dashboard', { replace: true })
    }
    return
  }

  if (!auth.isAuthenticated) {
    return navigateTo('/login', { replace: true })
  }
})
