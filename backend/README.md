# English Reader Backend

Backend закрытого семейного EPUB-ридера. Реализован на Go 1.25, PostgreSQL
через `pgx`, gRPC/gRPC-Web, HTTP health API, Swagger и локальном LibreTranslate.

## Архитектура кода

Для каждого HTTP и gRPC запроса используется одна цепочка:

`handler/<transport>/<feature>/<request> → usecase/<feature>/<request> → service (при необходимости) → repository`.

Каждый handler и usecase находится в собственном package. Повторно используемая
логика Auth оформлена как services `password`, `token`, `refreshtoken`; доступ к
PostgreSQL — concrete repository `repository/postgres/auth`. Интерфейсы зависимостей
находятся у потребителя в `dependencies.go`, а конкретные реализации собираются
только в `cmd/server`. Полные правила — в [AGENTS.md](AGENTS.md).

## Требования

- Docker Desktop / Docker Engine с Compose v2;
- Go 1.25 для запуска и тестов вне Docker;
- `make`;
- `protoc` и `buf` нужны только для изменения protobuf-контрактов.

## Быстрый запуск в Docker

```bash
cd backend
cp .env.example .env
make infra-up
make migrate-up
make app-up
```

Первый старт LibreTranslate скачивает языковые модели English и Russian, поэтому
может занять несколько минут. Их состояние сохраняется в Docker volume.

Проверить контейнеры и логи:

```bash
make ps
make logs
```

Остановить сервисы, сохранив данные:

```bash
make app-down
```

## Запуск API на хосте

Сначала поднимите зависимости, затем примените миграции. `Makefile` читает
переменные из `.env` и экспортирует их для команды Go.

```bash
cd backend
cp .env.example .env
make infra-up
make migrate-up
make run
```

Для development используйте минимум 32 случайных символа в `JWT_SECRET`.
Значения в `.env.example` подходят только для локальной разработки.

## Порты и адреса

| Сервис | Адрес с хоста | Адрес в сети Docker | Назначение |
| --- | --- | --- | --- |
| PostgreSQL | `localhost:5433` | `postgres:5432` | База данных `reader`. Порт можно изменить через `POSTGRES_PORT`. |
| API HTTP | `http://localhost:8083` | `api:8081` | Health checks, REST Auth API и Swagger. Меняется через `API_HTTP_PORT`. |
| Swagger UI | `http://localhost:8083/swagger/index.html` | `http://api:8081/swagger/index.html` | Интерактивная документация HTTP API. |
| OpenAPI | `http://localhost:8083/openapi.yaml` | `http://api:8081/openapi.yaml` | OpenAPI 3.0 YAML-спецификация. |
| Envoy / gRPC-Web | `http://localhost:8082` | `envoy:8080` | Вход браузерных gRPC-Web запросов. Меняется через `ENVOY_PORT`. |
| API gRPC | не публикуется | `api:9090` | Обычный gRPC для Envoy и будущих внутренних вызовов. |
| LibreTranslate | не публикуется | `libretranslate:5000` | Локальный бесплатный English ↔ Russian переводчик. |

Порты `5433`, `8082` и `8083` выбраны по умолчанию, чтобы не конфликтовать с
часто занятыми `5432`, `8080` и `8081`. В production HTTP API и PostgreSQL не
следует публиковать напрямую: наружу должен смотреть только защищённый gateway.

Для браузерной регистрации и обновления сессии HTTP API принимает credentialed
CORS-запросы только с `FRONTEND_ORIGIN` (по умолчанию
`http://localhost:5173`). Если Vite запускается на другом адресе, задайте этот
же origin в `backend/.env` и перезапустите `make app-up`.

## Проверка API

```bash
curl http://localhost:8083/health/live
curl http://localhost:8083/health/ready
```

Регистрация тестового пользователя:

```bash
curl -i -X POST http://localhost:8083/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  --data '{
    "email": "family@example.com",
    "password": "correct horse battery staple",
    "device_label": "local verification"
  }'
```

Ответ содержит access JWT, а refresh token устанавливается как `HttpOnly`
cookie. Не добавляйте реальные пароли или токены в shell history, логи и git.

## Тесты и качество кода

```bash
make test        # все Go-тесты
make test-race   # тесты с race detector
make test-docker # тесты в образе Go 1.25
make vet         # go vet
make fmt         # gofmt
make tidy        # обновить go.mod и go.sum
```

## Миграции

```bash
make migrate-up    # применить все новые миграции
make migrate-down  # откатить последнюю миграцию
```

Миграции находятся в `internal/database/migrations/` и встраиваются в бинарник
`reader-migrate`. Не редактируйте уже применённую миграцию: для изменения схемы
создавайте следующий нумерованный файл `*.up.sql` и соответствующий `*.down.sql`.

## gRPC-контракты

Proto-файлы находятся в `../protos/reader/v1`. AuthService описан в
`auth.proto`, LibraryService — в `library.proto`, ReaderService — в
`reader.proto`, DictionaryService — в `dictionary.proto`, а личный словарь и
подсветка — в `vocabulary.proto`. Сгенерированный Go-код хранится в
`gen/reader/v1`. После изменения контрактов выполните `make proto` и добавьте
сгенерированные файлы в изменение. `make proto-lint` запускает `buf lint`, если
`buf` установлен.

## Словарь и переводчик

`DictionaryService.LookupWord` возвращает только senses из PostgreSQL, поэтому
словарная карточка не ждёт внешний переводчик. `DictionaryService.TranslateText`
переводит явно переданный фрагмент и возвращает `provider_error`, если локальный
LibreTranslate недоступен. В Compose заданы один запрос без повтора, тайм-аут
`TRANSLATE_TIMEOUT=2s` и лимит `LOOKUP_SENTENCE_MAX_LENGTH=360`; при
необходимости их можно переопределить в `.env`. Текст книги, JWT и user ID во
внешний provider не уходят.

Тестовый словарь можно импортировать после запуска PostgreSQL:

```bash
make dictionary-import \
  DICTIONARY_FILE=test/fixtures/dictionary-mini.jsonl \
  DICTIONARY_VERSION=fixture-v1
```

Атрибуция импортированных данных описана в [DICTIONARY_ATTRIBUTION.md](docs/DICTIONARY_ATTRIBUTION.md).
Инструкция по подготовке настоящего открытого набора и versioned import — в
[DICTIONARY_IMPORT.md](docs/DICTIONARY_IMPORT.md). Fixture предназначен только
для разработки и не заменяет словарь для использования семьёй.

Для выгрузки Kaikki/Wiktextract сначала создайте компактный `en → ru` JSONL:

```bash
make dictionary-convert \
  KAIKKI_FILE='en-extract/Wiktextract Data.jsonl.gz' \
  DICTIONARY_FILE=/private/tmp/en-ru-wiktionary.jsonl \
  DICTIONARY_VERSION=2026-08-05

make dictionary-import \
  DICTIONARY_FILE=/private/tmp/en-ru-wiktionary.jsonl \
  DICTIONARY_SOURCE=wiktionary \
  DICTIONARY_VERSION=2026-08-05
```

## gRPC-Web LibraryService

Через Envoy gRPC-Web (`http://localhost:8082`) доступны `ListCatalog`,
`GetBook`, `AddToMyLibrary`, `ListMyBooks` и `RemoveFromMyLibrary`.

`UploadBook` — native gRPC client-streaming вызов. Первое сообщение содержит
`UploadBookMetadata` с `access_token` и именем EPUB, последующие — `chunk`.
Сервер передаёт chunks напрямую в защищённое временное хранилище, не собирая
весь файл в памяти; применяются те же 50 MiB, проверка расширения и ZIP magic
bytes, что и у HTTP fallback. Browser gRPC-Web не поддерживает
client-streaming, поэтому frontend загружает EPUB через `POST
/api/v1/library/books` multipart endpoint; остальные library-сценарии
используют gRPC-Web.

## Личный словарь и подсветка

`VocabularyService` доступен через Envoy gRPC-Web (`http://localhost:8082`) и
требует access JWT в поле `access_token` gRPC-запроса. В нём доступны
`SaveEntry`, `ListEntries`, `DeleteEntry` и `GetHighlights`. Запись словаря
всегда привязана к subject из JWT: клиент не передаёт `user_id`. Повторное
сохранение той же lemma идемпотентно и не изменяет выбранный sense.

`GetHighlights` принимает только идентификаторы книги и главы. Сервер сам
проверяет доступ к готовой книге, получает `plain_text` главы, токенизирует его
и возвращает только совпавшие позиции из словаря текущего пользователя. XHTML,
текст главы и записи других пользователей в ответ не попадают.

Проверка PostgreSQL-репозиториев, включая изоляцию словаря, выполняется в
отдельной временной базе:

```bash
make test-repository
make test-db-down
```

## Полезные документы

- [Инициализация и текущий статус](docs/00-init.md)
- [План библиотеки](docs/01-library-plan.md)
- [План личного словаря](docs/04-vocabulary-plan.md)
- [Задачи личного словаря](docs/04-vocabulary-tasks.md)
- [Сессия изменений](docs/session/2026-09-03-01.md)
- [Архитектура и хранение книг](../docs/BACKEND_ARCHITECTURE.md)
