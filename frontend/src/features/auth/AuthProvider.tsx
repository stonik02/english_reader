import {
  type ReactNode,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from 'react'

import {
  type AuthUser,
  getMe,
  login as loginRequest,
  logout as logoutRequest,
  refresh,
  register as registerRequest,
} from '../../api/auth'
import { sessionToken } from '../../api/client'
import {
  AuthContext,
  type AuthContextValue,
  type AuthStatus,
  type Credentials,
} from './authContext'

export function AuthProvider({ children }: { children: ReactNode }) {
  const [status, setStatus] = useState<AuthStatus>('loading')
  const [user, setUser] = useState<AuthUser | null>(null)

  const clearSession = useCallback(() => {
    sessionToken.clear()
    setUser(null)
    setStatus('anonymous')
  }, [])

  const renewAccessToken = useCallback(async () => {
    try {
      const session = await refresh()
      sessionToken.set(session.access_token)
      const currentUser = await getMe(session.access_token)
      setUser(currentUser)
      setStatus('authenticated')
      return session.access_token
    } catch {
      clearSession()
      return null
    }
  }, [clearSession])

  useEffect(() => {
    sessionToken.setRenewal(renewAccessToken)
    void sessionToken.renew()

    return () => sessionToken.setRenewal(null)
  }, [renewAccessToken])

  const authenticate = useCallback(
    async (
      request: (
        credentials: Credentials,
      ) => Promise<{ access_token: string; user: AuthUser }>,
      credentials: Credentials,
    ) => {
      const session = await request(credentials)
      sessionToken.set(session.access_token)
      setUser(session.user)
      setStatus('authenticated')
    },
    [],
  )

  const value = useMemo<AuthContextValue>(
    () => ({
      status,
      user,
      login: (credentials) => authenticate(loginRequest, credentials),
      register: (credentials) => authenticate(registerRequest, credentials),
      logout: async () => {
        try {
          await logoutRequest()
        } finally {
          clearSession()
        }
      },
    }),
    [authenticate, clearSession, status, user],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}
