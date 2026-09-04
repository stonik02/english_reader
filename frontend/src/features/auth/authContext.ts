import { createContext } from 'react'

import type { AuthUser } from '../../api/auth'

export type Credentials = { email: string; password: string }
export type AuthStatus = 'anonymous' | 'authenticated' | 'loading'

export type AuthContextValue = {
  login(credentials: Credentials): Promise<void>
  logout(): Promise<void>
  register(credentials: Credentials): Promise<void>
  status: AuthStatus
  user: AuthUser | null
}

export const AuthContext = createContext<AuthContextValue | null>(null)
