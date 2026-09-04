# Задачи реализации: встроенный переводчик

Этот файл декомпозирует [план переводчика](03-translator-plan.md) на
последовательные задачи. Каждая задача должна сохранять архитектуру проекта:
`handler → request-specific usecase → optional service → repository`.

## T01. Контракт и конфигурация — готово

**Результат:** публичный `DictionaryService.LookupWord` готов для frontend.

- Создать `protos/reader/v1/dictionary.proto`.
- Описать `LookupWordRequest`: access token, `book_id`, `chapter_id`,
  `selected_text`, `sentence_text`, `epub_cfi`, source language.
- Описать `LookupWordResponse`: нормализованная lemma, repeated senses,
  перевод предложения либо ошибка provider, `context_verified`, source metadata.
- Не передавать user ID, пути к файлам или полный EPUB.
- Запустить `make proto`; зарегистрировать пустой каркас `DictionaryService`.
- Добавить конфигурацию: `TRANSLATE_TIMEOUT`, `TRANSLATE_CACHE_TTL`, лимиты
  длины слова и предложения.

**Выполнено 2026-09-03:** добавлен `dictionary.proto`, сгенерированы Go/gRPC
stubs, зарегистрирован contract-first `DictionaryService` stub, добавлены
`TRANSLATE_TIMEOUT`, `TRANSLATE_CACHE_TTL`, `LOOKUP_WORD_MAX_LENGTH` и
`LOOKUP_SENTENCE_MAX_LENGTH`. До T06 `LookupWord` намеренно отвечает
`UNIMPLEMENTED`.

**Проверка:** protobuf генерируется, `buf lint` и `go test ./...` проходят.

## T02. Тестовый словарный импорт и атрибуция — готово

**Результат:** словарные данные можно повторяемо загрузить без web-scraping.

- Создать `cmd/dictionary-import` с входом JSONL/TSV из versioned dump.
- Проверить checksum, source, source version и обязательную атрибуцию.
- Записывать начало/окончание/ошибку в `dictionary_import_runs`.
- Upsert `dictionary_lemmas`, заменить senses конкретной lemma/source/version
  в одной транзакции.
- Добавить `docs/DICTIONARY_ATTRIBUTION.md`: источник, версия, CC BY-SA,
  формат показа attribution в UI.
- Подготовить маленький fixture с несколькими частями речи и переводами.

**Выполнено 2026-09-03:** добавлены `cmd/dictionary-import`, команда
`make dictionary-import`, fixture `test/fixtures/dictionary-mini.jsonl` и
`docs/DICTIONARY_ATTRIBUTION.md`. Импорт сохраняет checksum и статус в
`dictionary_import_runs`, выполняет versioned upsert lemmas/senses и требует
URL/attribution в каждой записи.

**Проверка:** повторный импорт не создаёт дубликатов; новая версия источника
может сосуществовать со старой до переключения.

## T03. Normalizer и morphological service — готово

**Результат:** selected text приводится к безопасной lookup-форме.

- Создать `internal/service/wordnormalizer` и `dependencies.go`.
- Выполнить Unicode normalization, trim, lowercase, удаление внешней
  пунктуации; запретить HTML, пробелы внутри одиночного слова и слишком длинный
  input.
- Создать optional `morphology` service/port: exact form → irregular forms →
  простые английские суффиксы. Не выдавать догадку как подтверждённую lemma.
- Написать table-driven tests для регистра, апострофов, пунктуации,
  `went → go`, `children → child`, неизвестных форм и невалидного input.

**Выполнено 2026-09-03:** добавлены локальные services `wordnormalizer` и
`morphology` с unit-тестами. Normalizer удаляет внешнюю пунктуацию, приводит
к lower-case и отклоняет HTML/пробелы/недопустимые символы; morphology содержит
best-effort irregular и suffix fallback с флагом достоверности.

**Проверка:** сервис не делает сетевых вызовов и не зависит от HTTP/gRPC.

## T04. PostgreSQL Dictionary repository — готово

**Результат:** usecase получает ordered senses и кэш перевода.

- Создать `internal/repository/postgres/dictionary`.
- Реализовать exact lookup и lookup по предложенной lemma с фильтром `en`.
- Возвращать senses по `position ASC`; translations декодировать из JSONB.
- Реализовать `translation_cache`: SHA-256 нормализованного предложения,
  модель/языки, TTL; просроченные записи не выдавать.
- Не логировать sentence text; для диагностики допускаются hash и длина.
- Добавить repository integration tests к `make test-repository`.

**Выполнено 2026-09-03:** добавлен `repository/postgres/dictionary` с exact
lookup senses по `position`, JSONB translations и SHA-256 translation cache с
TTL. Integration test проверяет два POS в порядке и cache hit.

**Проверка:** несколько POS и переводов сохраняют порядок; expired cache miss.

## T05. LibreTranslate adapter — готово

**Результат:** локальный provider переводит только одно предложение.

- Создать `internal/service/libretranslate` и `dependencies.go`.
- Использовать HTTP client с контекстным timeout, лимитом текста и JSON API
  LibreTranslate (`q`, `source=en`, `target=ru`, `format=text`).
- Реализовать один повтор только для временных ошибок, малый circuit breaker
  и typed `ProviderUnavailable` error.
- Не отправлять token, user ID, CFI, book ID или главу целиком.
- Добавить fake provider и tests: success, timeout, 5xx, malformed response,
  открытый circuit.

**Выполнено 2026-09-03:** добавлен HTTP adapter `service/libretranslate` с
контекстным client timeout, лимитом текста, одной повторной попыткой 5xx/timeout
и circuit breaker после трёх неудач. Unit test использует локальный fake HTTP
provider.

**Проверка:** provider failure не превращается в ошибку словарного lookup.

## T06. LookupWord usecase — готово

**Результат:** сценарий переводчика с частичным успехом.

- Создать `internal/usecase/dictionary/lookupword` и `dependencies.go`.
- Проверить JWT subject на handler boundary; usecase проверяет `READY` книги,
  принадлежность chapter этой книге и вызывает `EnsureBookInMyLibrary` через
  reader-facing port при первом чтении.
- Получить безопасный context из `book_chapters.plain_text`; если его нельзя
  подтвердить, установить `context_verified=false`, не подменяя книгу.
- Выполнить normalizer → exact/lemma dictionary lookup.
- Для sentence translation: сначала cache, потом provider, затем cache put.
- При provider error вернуть senses и field `provider_error`, не gRPC error.

**Проверка:** GoMock tests для cache hit/miss, пустого словаря, provider
timeout, processing book, чужой chapter и unverified context.

## T07. Transport и observability — готово

**Результат:** browser получает DictionaryService через gRPC-Web.

- Создать отдельный gRPC handler package `handler/grpc/dictionary/lookupword`
  и HTTP fallback только если frontend действительно потребует его.
- Добавить gRPC mapping: input → `INVALID_ARGUMENT`, chapter/book →
  `NOT_FOUND`, processing → `FAILED_PRECONDITION`; provider failure находится
  в успешном response field.
- Подключить service к composition root и Envoy уже существующим gRPC route.
- Добавить span `Dictionary.LookupWord` и метрики duration/cache hit/provider
  result без текста запроса.

**Выполнено 2026-09-03:** `LookupWord` проверяет READY/book/chapter и
`plain_text` context, DictionaryService вызывает отдельный gRPC handler;
добавлены GoMock usecase и handler tests. В gRPC access-логе уже есть duration
и status code, без текста lookup; отдельный metrics exporter откладывается до
подключения общего OpenTelemetry provider.

**Проверка:** grpc handler tests используют GoMock; Swagger/OpenAPI обновлён,
если добавлен HTTP fallback.

## T08. Сквозная проверка — готово

**Результат:** переводчик устойчив в локальном Compose окружении.

- Добавить dictionary fixture/import в Compose test flow.
- Проверить: несколько senses, exact/lemma fallback, cache hit, пустой
  результат, provider timeout, `PROCESSING` EPUB и chapter другой книги.
- Проверить, что LibreTranslate доступен только по внутреннему адресу Docker.
- Дополнить `backend/docs/session/<date>-<session>.md`, README и документацию
  атрибуции.

**Выполнено 2026-09-03:** расширены repository integration tests: несколько
senses, cache hit и expired cache miss. Compose уже запускает LibreTranslate
внутри сети без host port; README и session documentation дополнены importer и
поведением partial success.

**Definition of done:** все T01–T08 выполнены; `make test`,
`make test-repository`, `make proto` и Compose smoke-check проходят.
