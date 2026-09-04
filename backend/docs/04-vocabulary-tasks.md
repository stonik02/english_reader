# Задачи реализации: личный словарь и подсветка

Этот план декомпозирует [VocabularyService](04-vocabulary-plan.md). Новые
сценарии следуют обязательной цепочке:

```text
gRPC handler → request-specific usecase → optional service → repository
```

Каждая задача отмечается готовой только после тестов и генерации protobuf.

## V01. Контракт VocabularyService — готово

**Результат:** frontend может вызывать личный словарь через gRPC-Web.

- Дополнить `dictionary.proto` методами `SaveEntry`, `ListEntries`,
  `DeleteEntry`, `GetHighlights` либо вынести их в отдельный
  `vocabulary.proto` с `VocabularyService`.
- Описать request/response: lemma ID, optional sense ID, source form,
  cursor/limit/query, entry ID/lemma ID, book/chapter ID.
- В `SaveEntryResponse` добавить `already_saved`.
- `HighlightToken` содержит lemma ID, lemma, optional positions и match kind
  (`lemma` или `exact_fallback`), но не HTML главы.
- Запустить `make proto` и зарегистрировать contract-first gRPC service.

**Проверка:** старые номера полей не меняются; `make proto`, `buf lint` и
`go test ./...` проходят.

## V02. Repository и миграции — готово

**Результат:** личные записи хранятся изолированно и идемпотентно.

- Проверить существующую миграцию `vocabulary_entries`; при необходимости
  добавить новую, не редактируя применённую.
- Создать `internal/repository/postgres/vocabulary`.
- Реализовать lookup lemma и проверку, что selected sense принадлежит lemma.
- Реализовать `INSERT ... ON CONFLICT (user_id, lemma_id) DO NOTHING` с
  чтением существующей записи и флагом `already_saved`.
- Реализовать cursor pagination `(created_at, id)` и search по lemma.
- Реализовать идемпотентный delete только по subject пользователя.

**Проверка:** integration tests на двойное сохранение, чужую запись, неверный
sense и cursor pagination в `make test-repository`.

## V03. SaveEntry usecase и gRPC handler — готово

**Результат:** слово из LookupWord сохраняется в личный словарь.

- Создать `usecase/vocabulary/saveentry` и `handler/grpc/vocabulary/saveentry`.
- Handler получает subject только из JWT access token.
- Usecase валидирует source form, lemma ID и optional sense ID; не позволяет
  произвольный lemma text вместо dictionary lemma ID.
- Вернуть entry и `already_saved`; повторный запрос не меняет chosen sense.
- Добавить GoMock tests: success, duplicate, unknown lemma, invalid sense pair.

**Проверка:** другой пользователь не может создать/увидеть связь владельца.

## V04. ListEntries и DeleteEntry — готово

**Результат:** frontend показывает и удаляет только собственный словарь.

- Создать отдельные usecase/handler packages `listentries`, `deleteentry`.
- Ограничить limit, валидировать cursor и нормализовать search query через
  общий word normalizer.
- List response включает lemma, source form, chosen sense для UI и timestamps.
- Delete принимает entry ID либо lemma ID, удаляет только `(user_id, …)` и
  возвращает успешный результат при повторе.
- Добавить GoMock и repository integration tests для изоляции двух аккаунтов.

**Проверка:** никакой endpoint не принимает `user_id` от клиента.

## V05. Общие normalizer и tokenizer — готово

**Результат:** карточка перевода и подсветка одинаково обрабатывают формы слов.

- Оставить `service/wordnormalizer` единым местом normalisation.
- Создать `service/tokenizer` с `dependencies.go`: токены из server-side
  `plain_text`, apostrophe внутри слова, исходные offsets при необходимости.
- Использовать `morphology` только как best-effort; сохранять `exact_fallback`,
  если lemma не определена.
- Добавить table-driven tests на punctuation, Unicode apostrophe, irregular
  forms и лимиты входа.

**Проверка:** tokenizer не принимает XHTML и не использует текст клиента.

## V06. GetHighlights — готово

**Результат:** подсветка выдаёт пересечение словаря пользователя и текущей главы.

- Создать `usecase/vocabulary/gethighlights` и отдельный gRPC handler.
- Repository/reader port проверяет `READY`, book/chapter ownership и выбирает
  только `plain_text` указанной chapter.
- Выбрать lemmas текущего пользователя, токенизировать plain text, сопоставить
  exact/lemma form и вернуть только совпавшие tokens.
- Ограничить размер главы, число vocabulary lemmas и количество highlights.
- Не возвращать HTML, полный список словаря или слова иных пользователей.
- Добавить GoMock tests: processing book, чужой chapter, exact fallback,
  несколько совпадений и пустой результат.

**Проверка:** XHTML в `book_chapters` остаётся неизменённым.

## V07. Интеграция и эксплуатация — готово

**Результат:** VocabularyService готов для семейного Compose окружения.

- Добавить Dictionary/Vocabulary fixture в `test/integration`.
- Проверить два пользователя с разными entries, idempotent SaveEntry, delete,
  highlights для разных форм одной lemma и отсутствие межпользовательской утечки.
- Обновить README, session journal и OpenAPI только при наличии HTTP fallback.
- Измерить простой plain-text highlights flow; `book_chapter_tokens` добавлять
  только при подтверждённой проблеме производительности.

**Definition of done:** все V01–V07 выполнены; личное слово доступно только
сохранившему его пользователю, повторное сохранение не создаёт дубликат,
подсветка ограничена текущей готовой главой.

## Итог проверки

Выполнены `make proto`, `make mocks`, `go test ./...`, `go vet ./...` и
`make test-repository`. Repository-интеграция использует временную PostgreSQL
на порту `5434` и проверяет idempotent save, неверную пару lemma/sense,
изоляцию двух пользователей, удаление, поиск и cursor pagination. Подсветка
дополнительно получает главу только через связь `user_books`, проверяет статус
`READY` и принадлежность главы книге.
