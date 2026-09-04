# gRPC-контракты и локальный запуск

## Транспорт

Клиент — браузер, поэтому использует **gRPC-Web**, а Envoy конвертирует его в
обычный gRPC для Go API. Официальный grpc-web проект описывает необходимость
специального proxy и приводит Envoy как стандартный вариант. Все методы v1
unary: browser-side bidirectional streaming не нужен.
[grpc-web README](https://github.com/grpc/grpc-web)

Пока Go-сервис один, server-to-server вызовов между нашими сервисами нет.
Техническое требование gRPC выполняется для связи client-server. LibreTranslate
— сторонний self-hosted HTTP API; его протокол не контролируется проектом и
поэтому не является нарушением внутреннего gRPC-контракта.

## Публичный API v1

| Proto service | Unary RPC | Назначение |
| --- | --- | --- |
| `AuthService` | `Register`, `Login`, `Refresh`, `Logout`, `GetMe` | Семейные аккаунты и сессии. |
| `LibraryService` | `UploadEpub`, `ListCatalog`, `GetBook`, `AddToMyLibrary`, `ListMyBooks`, `RemoveFromMyLibrary` | Общий EPUB-каталог и личные списки книг. |
| `ReaderService` | `GetReadingState`, `SaveReadingProgress`, `GetChapter` | Последняя позиция, настройки и подготовленная глава. |
| `DictionaryService` | `LookupWord`, `SaveEntry`, `ListEntries`, `DeleteEntry`, `GetHighlights` | Перевод, личный словарь и подсветка. |

`UploadEpub` принимает файл потоковой загрузкой, но сервер ограничивает его
размер (стартовое значение: 50 MiB). При необходимости gRPC-Web upload окажется
неудобным, `LibraryService` будет возвращать одноразовый HTTP upload URL от
самого Go-приложения; это не меняет модель хранения и не требует MinIO.

После загрузки EPUB появляется в общем каталоге со статусом `PROCESSING`,
`READY` или `FAILED`. `AddToMyLibrary` идемпотентен; `GetReadingState` или
первый переход к `GetChapter` автоматически выполняют то же добавление.
`RemoveFromMyLibrary` удаляет только связь текущего пользователя с книгой,
а не сам общий EPUB. `SaveReadingProgress` содержит `revision`: повторные и устаревшие запросы не
перетирают более новую позицию. Все методы, кроме регистрации и входа, требуют
`authorization: Bearer <access-token>`; доступ к каталогу доступен членам семьи,
а на каждой личной операции проверяется владелец записи/словаря.

`LookupWord` возвращает отдельно `dictionary_senses[]` (часть речи, переводы,
пример, атрибуция) и `sentence_translation`. Это не смешивает словарные данные
с машинным переводом и позволяет заменить provider впоследствии.

## Контракты и генерация

```text
protos/reader/v1/*.proto
  |-- buf lint + breaking
  |-- protoc-gen-go / protoc-gen-go-grpc -> backend/gen/go
  `-- protoc-gen-grpc-web + TypeScript -> frontend/src/api/gen
```

Proto-поля только добавляются; номера не переиспользуются. Ошибки — стандартные
gRPC codes: `UNAUTHENTICATED`, `PERMISSION_DENIED`, `NOT_FOUND`,
`INVALID_ARGUMENT`, `FAILED_PRECONDITION` (EPUB ещё готовится), `RESOURCE_EXHAUSTED`
и `UNAVAILABLE` (переводчик временно недоступен).

## Локальный запуск

```text
frontend (5173) -> envoy (8080) -> reader-api (9090 gRPC, 8081 health)
                                      |-> postgres (5432)
                                      `-> libretranslate (5000, доступен только внутри сети Compose)
```

`docker compose up --build` поднимает четыре контейнера и два persistent
volume: `postgres-data` и `book-files`. `reader-api` запускает миграции и
возобновляет записи `PROCESSING`. LibreTranslate получает только модели `en`
и `ru`; web UI и внешний порт для него выключены.

Перед запуском разработчик копирует `.env.example` в `.env`, задаёт JWT secret
и пароль PostgreSQL, затем выполняет `docker compose up --build`, `buf lint` и
тесты Go/frontend. В production доступ извне открыт только к frontend/Envoy;
PostgreSQL, volume и LibreTranslate остаются во внутренней сети.
