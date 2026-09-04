# Доработки до функционального завершения backend

Этот документ закрывает выявленные после реализации планов расхождения с БТ.
В него намеренно не входят E2E-тесты, телеметрия и добавление новых тестов.
Существующие unit и repository integration tests остаются регрессионной
проверкой изменений.

## C01. Полный gRPC LibraryService — готово

**Проблема:** `library.proto` существует, но сервис не зарегистрирован; UI
может выполнять library-сценарии только по HTTP.

- Создать отдельные gRPC handlers для `ListCatalog`, `GetBook`,
  `AddToMyLibrary`, `ListMyBooks`, `RemoveFromMyLibrary`.
- Собрать их в concrete `grpcserver.LibraryService` и зарегистрировать в
  composition root.
- Сохранить HTTP endpoints как совместимый Swagger fallback.

## C02. Загрузка EPUB: native gRPC и HTTP fallback — готово с ограничением gRPC-Web

**Проблема:** в proto отсутствовал безопасный путь загрузки EPUB. Browser
gRPC-Web не поддерживает client-streaming, необходимый для безопасного файла.

- Добавить native gRPC client-streaming `UploadBook`, где первое сообщение несёт JWT и имя
  файла, последующие — chunks с ограничением размера на существующем storage.
- Создать отдельный handler, который передаёт поток в существующий upload
  usecase без буферизации полного файла в памяти.
- Зарегистрировать метод в LibraryService; HTTP multipart upload оставить
  обязательным frontend fallback и задокументировать это ограничение.

## C03. Статус сохранения lemma в карточке перевода — готово

**Проблема:** `LookupWordResponse` не сообщает UI, есть ли lemma в личном
словаре текущего пользователя.

- Добавить поле `already_saved` в ответ DictionaryService без изменения
  существующих номеров полей.
- Добавить минимальный vocabulary read-port в LookupWord usecase и concrete
  PostgreSQL-проверку `(user_id, lemma_id)`.
- Возвращать `false` для не найденной в словаре lemma и не раскрывать данные
  других пользователей.

## C04. Эксплуатационный импорт настоящего словаря — готово

**Проблема:** в репозитории лежит только небольшой fixture; production-словарь
не должен тайно подменяться тестовыми данными.

- Зафиксировать формат JSONL, обязательную атрибуцию и команду versioned import
  в README и отдельной инструкции.
- Явно указать, что оператор выбирает и скачивает юридически допустимый dump,
  затем импортирует его в PostgreSQL; приложение не требует внешнего платного
  dictionary API.

## C05. Синхронизация документации и проверка контрактов — готово

- Обновить README и session journal списком доступных gRPC services и upload
  flow.
- Выполнить `make proto`, `make mocks`, `go test ./...`, `go vet ./...` и
  `docker compose config`.

## Definition of done

Все сценарии БТ доступны frontend по gRPC-Web, кроме загрузки EPUB: для неё
используется документированный HTTP multipart fallback из-за ограничения
browser gRPC-Web. Native gRPC-клиенты загружают EPUB потоково, а карточка lookup
содержит личный статус lemma. Для полезного словаря оператор импортирует
версионированный открытый набор данных с атрибуцией.
