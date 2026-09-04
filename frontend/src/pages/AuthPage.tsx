import { type FormEvent, useState } from 'react'
import { Link, Navigate, useLocation, useNavigate } from 'react-router-dom'

import { AuthRequestError } from '../api/auth'
import { useAuth } from '../features/auth/useAuth'

type AuthPageProps = { mode: 'login' | 'register' }

export function AuthPage({ mode }: AuthPageProps) {
  const isLogin = mode === 'login'
  const { login, register, status } = useAuth()
  const location = useLocation()
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [passwordConfirmation, setPasswordConfirmation] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const destination =
    (location.state as { from?: string } | null)?.from ?? '/library'

  if (status === 'authenticated') {
    return <Navigate to={destination} replace />
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError(null)

    if (!isLogin && password !== passwordConfirmation) {
      setError('Пароли не совпадают.')
      return
    }

    setIsSubmitting(true)
    try {
      await (isLogin ? login : register)({ email, password })
      navigate(destination, { replace: true })
    } catch (requestError) {
      setError(
        requestError instanceof AuthRequestError
          ? requestError.message
          : 'Не удалось связаться с сервером. Попробуйте ещё раз.',
      )
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <main className="auth-page">
      <section className="auth-card" aria-labelledby="auth-title">
        <Link className="brand" to="/library">
          <span className="brand-mark" aria-hidden="true">
            ER
          </span>
          <span>English Reader</span>
        </Link>
        <p className="eyebrow">Семейная библиотека</p>
        <h1 id="auth-title">
          {isLogin ? 'Рады видеть снова' : 'Создайте аккаунт'}
        </h1>
        <p className="page-description">
          {isLogin
            ? 'Войдите, чтобы продолжить читать с сохранённого места.'
            : 'Отдельный профиль сохранит вашу библиотеку, прогресс и словарь.'}
        </p>
        <form className="placeholder-form" onSubmit={handleSubmit}>
          <label>
            Email
            <input
              autoComplete="email"
              name="email"
              onChange={(event) => setEmail(event.target.value)}
              placeholder="name@example.com"
              required
              type="email"
              value={email}
            />
          </label>
          <label>
            Пароль
            <input
              autoComplete={isLogin ? 'current-password' : 'new-password'}
              minLength={12}
              name="password"
              onChange={(event) => setPassword(event.target.value)}
              required
              type="password"
              value={password}
            />
          </label>
          {!isLogin && (
            <label>
              Повторите пароль
              <input
                autoComplete="new-password"
                minLength={12}
                name="passwordConfirmation"
                onChange={(event) =>
                  setPasswordConfirmation(event.target.value)
                }
                required
                type="password"
                value={passwordConfirmation}
              />
            </label>
          )}
          {error !== null && (
            <p className="form-error" role="alert">
              {error}
            </p>
          )}
          <button
            className="button button-primary"
            disabled={isSubmitting}
            type="submit"
          >
            {isSubmitting
              ? 'Подождите…'
              : isLogin
                ? 'Войти'
                : 'Создать аккаунт'}
          </button>
        </form>
        <p className="form-note">
          {isLogin ? 'Нет аккаунта?' : 'Уже есть аккаунт?'}{' '}
          <Link to={isLogin ? '/register' : '/login'}>
            {isLogin ? 'Зарегистрироваться' : 'Войти'}
          </Link>
        </p>
      </section>
    </main>
  )
}
