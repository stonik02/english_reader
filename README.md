# English Reader

Закрытое семейное приложение для чтения английских EPUB-книг. В нём будут
общий каталог книг, личная библиотека, ридер с сохранением позиции, перевод
слов и личный словарь.

Проект состоит из двух частей:

- [`frontend/`](frontend/README.md) — React-приложение, которое открывается в
  браузере;
- [`backend/`](backend/README.md) — Go API, PostgreSQL, Envoy gRPC-Web gateway
  и локальный LibreTranslate.

## Требования

- Docker Desktop с Docker Compose v2 — для backend и базы данных;
- Node.js 22.11+ и npm 10+ — для frontend;
- свободные порты `5173`, `5433`, `8082` и `8083`.

## Запуск frontend и backend вместе

Откройте **два терминала** в корне репозитория.

### Терминал 1 — backend

```bash
cd backend
cp .env.example .env
make infra-up
make migrate-up
make app-up
```

Первый запуск LibreTranslate может занять несколько минут: контейнер загрузит
модели английского и русского языков. Дождитесь, пока API будет готов, и
проверьте это во втором окне терминала или отдельной вкладке:

```bash
curl http://localhost:8083/health/ready
```

Ответ с успешным статусом означает, что backend готов принимать запросы.

### Терминал 2 — frontend

```bash
cd frontend
npm ci
npm run dev
```

Откройте [`http://localhost:5173`](http://localhost:5173) в браузере.

## Как сервисы соединяются

```text
Браузер (frontend, :5173)
       │ gRPC-Web
       ▼
Envoy gateway (:8082)
       │ gRPC
       ▼
Go API (:9090 внутри Docker) ──► PostgreSQL и LibreTranslate
```

Envoy уже настроен принимать browser-запросы с локальных адресов Vite
`http://localhost:5173` и `http://localhost:5174`. API также публикует HTTP
health-check и Swagger на `http://localhost:8083`.

Frontend использует gRPC-Web через Envoy для каталога, ридера, перевода и
личного словаря: браузерные запросы проходят к Go API без хранения access token
в конфигурационных файлах.
Адрес Envoy можно
переопределить через `frontend/.env.local`; подробности — в
[README frontend](frontend/README.md).

Для login/register/refresh frontend обращается к HTTP API на `:8083`, потому
что backend выставляет защищённую `HttpOnly` refresh-cookie. В development API
разрешает credentialed CORS-запросы только с адреса, указанного в
`FRONTEND_ORIGIN` файла `backend/.env`. Он должен совпадать с адресом frontend
в браузере.

Загрузка EPUB также идёт в HTTP API: форма отправляет `multipart/form-data` на
`POST /api/v1/library/books`, чтобы браузер мог показать прогресс отправки.
Остальные действия с библиотекой продолжают использовать gRPC-Web через Envoy.

## Остановка

- Frontend: нажмите `Ctrl+C` в терминале, где работает `npm run dev`.
- Backend: в папке `backend` выполните `make app-down`. Данные PostgreSQL и
  загруженные книги сохраняются в Docker volumes.

## Документация

- [Архитектура frontend](docs/FRONTEND_ARCHITECTURE.md)
- [План ближайших доработок книг](docs/plans/NEXT_SESSION_PLAN_1_DONE.md)
- [Карта экранов](frontend/docs/SCREENS.md)
- [План разработки frontend](frontend/docs/IMPLEMENTATION_PLAN.md)
- [README backend](backend/README.md)
- [Развёртывание на Ubuntu рядом с Remnawave](deploy/production/README.md)

## Обновление сайта на сервере

После изменения кода: проверьте изменения локально, сделайте `git commit` и
`git push origin main`. Затем по SSH перейдите в
`/opt/english-reader/app`, выполните `git pull --ff-only origin main` и
запустите обновление контейнеров и миграций. Полная безопасная инструкция,
включая импорт обновлённого словаря и откат, находится в
[production README](deploy/production/README.md#updating-the-deployed-application).

## Поддержка инструкций

README-файлы — часть рабочего контракта проекта. Любое изменение команд,
портов, переменных окружения или механизма связи frontend/backend должно
сопровождаться обновлением соответствующего README в том же изменении.
