import { create } from 'zustand'
import { getCookie, setCookie, removeCookie } from '@/lib/cookies'

const AUTH_SESSION_COOKIE = 'aurora_adminui_session'

interface AuthUser {
  label: string
  tokenType: 'bootstrap' | 'apitoken'
  bootstrapExchanged: boolean
  lastLoginAt: number
}

interface AuthState {
  auth: {
    user: AuthUser | null
    setUser: (user: AuthUser | null) => void
    sessionMarker: string
    setSessionMarker: (sessionMarker: string) => void
    resetSessionMarker: () => void
    reset: () => void
  }
}

export const useAuthStore = create<AuthState>()((set) => {
  const cookieState = getCookie(AUTH_SESSION_COOKIE)
  const initToken = cookieState ? JSON.parse(cookieState) : ''
  return {
    auth: {
      user: null,
      setUser: (user) =>
        set((state) => ({ ...state, auth: { ...state.auth, user } })),
      sessionMarker: initToken,
      setSessionMarker: (sessionMarker) =>
        set((state) => {
          setCookie(AUTH_SESSION_COOKIE, JSON.stringify(sessionMarker))
          return { ...state, auth: { ...state.auth, sessionMarker } }
        }),
      resetSessionMarker: () =>
        set((state) => {
          removeCookie(AUTH_SESSION_COOKIE)
          return { ...state, auth: { ...state.auth, sessionMarker: '' } }
        }),
      reset: () =>
        set((state) => {
          removeCookie(AUTH_SESSION_COOKIE)
          return {
            ...state,
            auth: { ...state.auth, user: null, sessionMarker: '' },
          }
        }),
    },
  }
})
