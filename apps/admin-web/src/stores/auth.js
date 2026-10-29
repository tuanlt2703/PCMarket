import { defineStore } from 'pinia'
import { authApi } from '@/services/authApi'

const SESSION_KEY = 'pcmarket-admin-session'

function readSession() {
  try {
    return JSON.parse(sessionStorage.getItem(SESSION_KEY)) ?? null
  } catch {
    return null
  }
}

export const useAuthStore = defineStore('auth', {
  state: () => ({ session: readSession() }),
  getters: {
    isAuthenticated: (s) => !!s.session?.token,
    user: (s) => s.session?.user ?? null,
  },
  actions: {
    async login(username, password) {
      this.session = await authApi.login(username, password)
      sessionStorage.setItem(SESSION_KEY, JSON.stringify(this.session))
    },
    logout() {
      this.session = null
      sessionStorage.removeItem(SESSION_KEY)
    },
  },
})
