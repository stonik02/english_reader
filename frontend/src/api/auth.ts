import { httpApiUrl } from './client'

export type AuthUser = {
  id: string
  email: string
  created_at: string
}

type AuthResponse = {
  user: AuthUser
  access_token: string
}

type Credentials = {
  email: string
  password: string
}

export class AuthRequestError extends Error {
  constructor(
    message: string,
    readonly status: number,
  ) {
    super(message)
    this.name = 'AuthRequestError'
  }
}

async function request<T>(path: string, options: RequestInit = {}) {
  const response = await fetch(`${httpApiUrl}${path}`, {
    ...options,
    credentials: 'include',
  })

  if (!response.ok) {
    const body = (await response.json().catch(() => null)) as {
      error?: string
    } | null
    throw new AuthRequestError(
      body?.error ?? 'Не удалось выполнить запрос к серверу.',
      response.status,
    )
  }

  if (response.status === 204) {
    return undefined as T
  }

  return (await response.json()) as T
}

function credentialsBody(credentials: Credentials) {
  return JSON.stringify({
    ...credentials,
    deviceLabel: 'English Reader web',
  })
}

export function register(credentials: Credentials) {
  return request<AuthResponse>('/api/v1/auth/register', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: credentialsBody(credentials),
  })
}

export function login(credentials: Credentials) {
  return request<AuthResponse>('/api/v1/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: credentialsBody(credentials),
  })
}

export function refresh() {
  return request<AuthResponse>('/api/v1/auth/refresh', { method: 'POST' })
}

export function logout() {
  return request<void>('/api/v1/auth/logout', { method: 'POST' })
}

export function getMe(accessToken: string) {
  return request<AuthUser>('/api/v1/auth/me', {
    headers: { Authorization: `Bearer ${accessToken}` },
  })
}
