# Backend-план: библиотека и аккаунты

## Цель и границы

Раздел реализует четыре требования библиотеки: семейные аккаунты, загрузку
EPUB в общий каталог, явное/автоматическое добавление в личную библиотеку и
просмотр личного списка. PDF, общий публичный доступ, роли, подписки и
самостоятельное удаление общей книги не входят в MVP.

Общий каталог и личная библиотека — разные сущности. Один EPUB существует один
раз в `books` и виден всем аутентифицированным членам семьи. Запись
`user_books` означает «эта книга добавлена данному пользователю» и не является
копией файла. Это защищает от ситуации, когда удаление из личного списка
удаляет книгу у других.

## Сценарии и правила

| Сценарий | Правило backend |
| --- | --- |
| Регистрация | Email нормализуется, должен быть уникален; пароль хранится только как Argon2id-хэш. Первый пользователь создаётся обычным способом, отдельной роли администратора нет. |
| Вход и сессия | `Login` выдаёт короткий access JWT и refresh session; refresh token в БД хранится только в виде хэша и может быть отозван. |
| Загрузка EPUB | Любой аутентифицированный член семьи загружает EPUB. После проверки появляется запись общего каталога со статусом `PROCESSING`; имя загрузившего сохраняется для аудита. |
| Тот же EPUB | Сравнение SHA-256 предотвращает дублирование. Повторная загрузка возвращает существующую общую книгу, а не создаёт второй файл/job. |
| Добавление кнопкой | `AddToMyLibrary` создаёт `(user_id, book_id)` один раз. Повторный запрос успешен и возвращает имеющуюся связь. |
| Первое чтение | Reader use case вызывает тот же внутренний `EnsureBookInMyLibrary`, указывая причину `first_read`. |
| Моя библиотека | Только `user_books` текущего JWT; сортировка по `added_at DESC`, курсорная пагинация. |
| Удаление из моей библиотеки | Удаляется только `user_books`. Если есть прогресс, он удаляется каскадно. Общий EPUB, главы и записи других пользователей остаются. |

Книга со статусом `PROCESSING` или `FAILED` видна в общем каталоге вместе с
безопасным статусом, но не добавляется чтением и не открывается. Книгу можно
добавить кнопкой только в `READY`; ошибка возвращается как
`FAILED_PRECONDITION`.

## Данные и инварианты

```text
users ──< auth_sessions
users ──< user_books >── books ──< book_chapters
                              └──< ingestion_jobs
```

Нужные миграции:

- `users`: UUID, email `citext UNIQUE`, `password_hash`, timestamps;
- `auth_sessions`: UUID, `user_id`, `refresh_token_hash UNIQUE`, expiry/revocation;
- `books`: UUID, `uploaded_by_user_id`, title/author, `format = epub`, status,
  SHA-256 `UNIQUE`, пути на private volume, error code, timestamps;
- `book_chapters`: `book_id`, sequence, href, CFI, sanitized XHTML/plain text;
- `ingestion_jobs`: одна активная job на книгу, attempt/lock/error;
- `user_books`: составной PK `(user_id, book_id)`, `added_at`,
  `added_via = button|first_read`.

File path не приходит от клиента и не выдаётся в публичном API. Он формируется
из server-generated `book_id`. Все запросы изменяют данные в транзакции; для
создания связи применяется `INSERT ... ON CONFLICT DO NOTHING`.

## gRPC-контракт

| Метод | Request / response | Проверки и ошибки |
| --- | --- | --- |
| `Register` | email, password → account/session | `INVALID_ARGUMENT`, `ALREADY_EXISTS` |
| `Login`, `Refresh`, `Logout`, `GetMe` | учётные данные/session → account | `UNAUTHENTICATED` для невалидной сессии |
| `UploadEpub` | bytes или upload ticket → `Book` | JWT, размер/MIME/ZIP; `INVALID_ARGUMENT`, `RESOURCE_EXHAUSTED` |
| `ListCatalog` | cursor, limit → `BookPage` | JWT; не выдаёт файловые пути |
| `GetBook` | `book_id` → `Book` | JWT, `NOT_FOUND` |
| `AddToMyLibrary` | `book_id`, idempotency key → `UserBook` | JWT, `FAILED_PRECONDITION`, `NOT_FOUND` |
| `ListMyBooks` | cursor, limit → `UserBookPage` | JWT; только текущий пользователь |
| `RemoveFromMyLibrary` | `book_id` → empty | JWT; идемпотентный успех при отсутствии связи |

Для browser upload предпочтителен отдельный короткоживущий HTTP upload endpoint
самого Go-приложения, если gRPC-Web не поддержит нужный размер/прогресс. Все
остальные client-server операции остаются gRPC-Web → Envoy → gRPC.

## План работ

1. Создать proto-пакет `reader/v1`, общие сообщения (`Book`, cursor, errors) и
   контракты auth/library; включить `buf lint` и breaking check.
2. Реализовать миграции `users`, `auth_sessions`, `books`, `user_books`,
   `ingestion_jobs`, `book_chapters` и repository-интерфейсы.
3. Реализовать auth use cases, JWT interceptor и ownership middleware; покрыть
   регистрацию, login, refresh, logout и изоляцию пользователей тестами.
4. Реализовать безопасное принятие EPUB во временный каталог, SHA-256
   дедупликацию, создание книги/job и статусы.
5. Реализовать idempotent add/list/remove личной библиотеки и внутренний
   `EnsureBookInMyLibrary` для reader module.
6. Добавить фоновой обработчик EPUB и выдачу статуса; подробная обработка файла
   следует [документу о хранении книг](../../docs/BOOK_STORAGE_AND_DELIVERY.md).
7. Поднять интеграционный тест Compose: два аккаунта, один общий EPUB,
   добавление/удаление у первого без изменения списка второго.

## Definition of Done

- Одну и ту же книгу можно увидеть в общем каталоге с двух аккаунтов.
- Личная библиотека пользователя содержит книгу после кнопки либо первого
  чтения, но не содержит не добавленные книги из каталога.
- Повторные upload/add/remove не создают дубликаты или ошибочные состояния.
- Невалидный/слишком большой EPUB не появляется в каталоге и не оставляет
  доступных файлов; пользователь B не может читать сессию пользователя A.
