import { Link } from 'react-router-dom'

export function NotFoundPage() {
  return (
    <main className="centered-page">
      <section
        className="empty-state not-found"
        aria-labelledby="not-found-title"
      >
        <p className="eyebrow">Ошибка 404</p>
        <h1 id="not-found-title">Страница не найдена</h1>
        <p>Похоже, такого адреса в English Reader пока нет.</p>
        <Link className="button button-primary" to="/library">
          В мою библиотеку
        </Link>
      </section>
    </main>
  )
}
