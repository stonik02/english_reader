# Frontend: системный анализ и рекомендуемый стек

Дата: 2026-09-03. Это закрытое домашнее web-приложение для владельца и семьи,
с менее чем 50 пользователями. Оно поддерживает только EPUB; PDF, подписки,
биллинг, SEO и публичные страницы отсутствуют.

## Рекомендация

**React + TypeScript + Vite + React Router + TanStack Query + gRPC-Web через
Envoy**. Это лёгкий SPA-стек: React отвечает исходному предпочтению, TypeScript
сохраняет типы API, Vite не добавляет ненужный серверный runtime, а TanStack
Query кэширует список книг, прогресс и словарь. SSR и Next.js здесь не нужны:
у приложения нет поискового трафика и публичного контента.

React рекомендует начинать новый production-проект с framework; в данном случае
Vite выбран осознанно как минимальный framework/tooling для закрытого SPA.
[React: Creating a React App](https://react.dev/learn/creating-a-react-app)

## gRPC в браузере

```text
React components → generated TypeScript gRPC-Web client → Envoy → Go reader-api
```

Браузер посылает binary gRPC-Web запросы к Envoy; Envoy преобразует их в обычный
gRPC. gRPC-Web требует такой proxy, а Envoy является стандартным вариантом.
Все browser RPC остаются unary. Загрузка EPUB — исключение: browser gRPC-Web
не поддерживает client-streaming, поэтому UI передаёт файл через существующий
`POST /api/v1/library/books` multipart endpoint, а затем получает статус книги
через gRPC-Web.
[grpc-web README](https://github.com/grpc/grpc-web)

`buf generate` создаёт TypeScript clients из `../protos`; hand-written JSON API
и дублированные типы не допускаются. Interceptor добавляет access token,
реагирует на `UNAUTHENTICATED`, один раз обновляет сессию и повторяет запрос.
Refresh token хранится в `HttpOnly`, `Secure`, `SameSite` cookie, не в
`localStorage`.

## Архитектура проекта

```text
frontend/src/
  app/             # router, providers, global styles
  api/gen/         # generated protobuf; не редактировать вручную
  api/client/      # grpc client, auth interceptor, error mapping
  features/
    auth/          # регистрация, вход, session hook
    library/       # список, upload EPUB, статус обработки
    reader/        # EPUB.js adapter, навигация, выделение
    dictionary/    # modal, словарь, подсветка
  pages/           # Login, Library, Reader, Vocabulary, NotFound
  shared/          # UI primitives, a11y, utilities
```

UI-состояние (модалка, выделение, локальное применение шрифта) остаётся в
React state/context. Server-state ведёт TanStack Query: query keys разделяют
кэш `books`, `readingState`, `chapter` и `vocabulary`; мутации инвалидируют
только затронутые ключи.

## Реализация требований в UI

| Функция | Экран и взаимодействие |
| --- | --- |
| Аккаунты | Вход, регистрация, route guard и `GetMe`; каждый член семьи работает со своим профилем. |
| Библиотека | Общий семейный каталог EPUB с HTTP multipart-загрузкой и статусами `PROCESSING`/`READY`/`FAILED`; отдельная вкладка «Моя библиотека» и кнопка добавления через gRPC-Web. |
| Ридер | EPUB.js рендерит главы и CFI; навигация назад/вперёд, размер шрифта, тема, автоматическое сохранение позиции. Первое открытие автоматически добавляет книгу в «Мою библиотеку». |
| Перевод | Выделение мышью/пальцем или клавиатурой открывает доступную modal: senses, части речи, пример и бесплатный перевод предложения. |
| Словарь | Кнопка `+`, постраничный личный список, поиск/удаление, безопасная подсветка слов текущей главы. |

Модальное окно поддерживает focus trap, Escape, touch и aria-label для `+`.
Загрузка и перевод имеют loading/error/retry состояния и не блокируют чтение.
Подготовленный backend XHTML отображается изолированно; выделение слов работает
по text nodes, без небезопасной вставки текста через `innerHTML`.

## Проверка и поставка

- TypeScript strict, ESLint и Prettier; Vitest/Testing Library для компонентов.
- Playwright сценарий: вход → EPUB upload → книга готова → открыть → изменить
  шрифт → закрыть/открыть на прежнем месте → добавить слово в словарь; отдельно
  проверяется явное добавление и автоматическое добавление книги при чтении.
- `docker compose up --build` поднимает frontend, Envoy, Go API, PostgreSQL и
  LibreTranslate; последнему не нужен публичный порт.
- Поддерживаются актуальные desktop/mobile браузеры. Отправка метрик и ошибок
  не содержит EPUB-текст, выделения или токены.

## Что подтвердить до кода

Для frontend-кода достаточно подтвердить предложенный стек и приемлемость
LibreTranslate для переводов English → Russian. Если качество окажется слабым,
мы заменим только backend translation adapter, не меняя UI или gRPC контракт.
