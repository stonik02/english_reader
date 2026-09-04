import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'

import { App } from './app/App'
import './app/global.css'
import { AppProviders } from './app/providers'

const rootElement = document.getElementById('root')

if (rootElement === null) {
  throw new Error('Не найден элемент #root для запуска приложения.')
}

createRoot(rootElement).render(
  <StrictMode>
    <AppProviders>
      <App />
    </AppProviders>
  </StrictMode>,
)
