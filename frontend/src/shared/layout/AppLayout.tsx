import { NavLink, Outlet } from 'react-router-dom'

import { useAuth } from '../../features/auth/useAuth'

const navigation = [
  { to: '/library', label: 'Моя библиотека' },
  { to: '/catalog', label: 'Общий каталог' },
  { to: '/vocabulary', label: 'Словарь' },
]

export function AppLayout() {
  const { logout, user } = useAuth()

  return (
    <div className="site-layout">
      <header className="site-header">
        <NavLink
          className="brand"
          to="/library"
          aria-label="English Reader — моя библиотека"
        >
          <span className="brand-mark" aria-hidden="true">
            ER
          </span>
          <span>English Reader</span>
        </NavLink>

        <nav className="main-navigation" aria-label="Основная навигация">
          {navigation.map((item) => (
            <NavLink
              className={({ isActive }) =>
                `nav-link${isActive ? ' nav-link-active' : ''}`
              }
              key={item.to}
              to={item.to}
            >
              {item.label}
            </NavLink>
          ))}
        </nav>

        <div className="user-menu" aria-label="Меню пользователя">
          <span className="user-avatar" aria-hidden="true">
            {user?.email.slice(0, 1).toUpperCase()}
          </span>
          <span className="user-name">{user?.email}</span>
          <button
            className="logout-button"
            type="button"
            onClick={() => void logout()}
          >
            Выйти
          </button>
        </div>
      </header>

      <main className="page-content">
        <Outlet />
      </main>
    </div>
  )
}
