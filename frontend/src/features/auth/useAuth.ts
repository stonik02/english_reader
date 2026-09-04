import { useContext } from 'react'

import { AuthContext } from './authContext'

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === null) {
    throw new Error('useAuth должен использоваться внутри AuthProvider.')
  }
  return context
}
