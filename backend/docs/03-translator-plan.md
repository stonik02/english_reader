# Backend-план: встроенный переводчик

## Цель и границы

По выделенному английскому слову backend возвращает модальной карточке:

1. все найденные варианты перевода на русский и части речи;
2. пример употребления;
3. машинный перевод именно того предложения, в котором выбрано слово.

Словарный ответ и перевод предложения имеют разное происхождение. Первый
строится из локального открытого словарного набора, производного от Wiktionary;
второй получает локально запущенный LibreTranslate с моделью English → Russian.
Это оставляет сервис бесплатным и не отправляет текст книги стороннему API.
Кнопка добавления слова не реализуется здесь: она относится к словарю и
использует `SaveEntry` из следующего документа.

## Контент и качество данных

Для словаря не следует скрейпить веб-страницы по запросу. На этапе импорта
используется versioned dump/структурированный экспорт Wiktionary, затем в
PostgreSQL создаются `dictionary_lemmas` и `dictionary_senses`. Для каждой
смысловой записи сохраняются lemma, part of speech, переводы, пример,
source URL, атрибуция и версия источника. Необходимые условия лицензии CC BY-SA
и атрибуции проверяются до публикации импортированных данных.

Лемматизация — best effort: сначала exact normalized lookup, затем английская
lemma через локальный morphological analyzer. Если lemma не определена,
карточка показывает результаты exact form и не выдаёт догадку за факт.
Система не выбирает единственный «правильный» смысл из контекста в MVP: она
возвращает все подходящие senses в сохранённом порядке, как требует БТ.

LibreTranslate запускается в Compose внутри сети, с минимальным набором моделей
`en`, `ru`. Вызов через Go adapter имеет timeout, лимит размера одного
предложения, малое число повторов и circuit breaker. Недоступность переводчика
не должна делать недоступными словарные senses и возвращается как отдельный
статус поля `sentence_translation`.

## Сценарий LookupWord

```text
Frontend выделяет token + предложение + CFI
  → DictionaryService.LookupWord
  → validation и normalisation
  → dictionary cache / PostgreSQL senses
  → translation cache / LibreTranslate (только предложение)
  → senses[] + sentence translation + source metadata
```

Входной запрос содержит `book_id`, `chapter_id`, `selected_text`,
`sentence_text`, `epub_cfi` и язык `en`. Backend проверяет JWT, существование
готовой книги и принадлежность chapter этой книге. Он не обязан требовать, чтобы
книга уже была добавлена в личную библиотеку: сам факт чтения создаёт такую
связь через Reader use case.

Ограничения запроса: одно слово/короткая фраза до согласованной длины,
предложение до безопасного лимита, Unicode normalisation, без HTML. Если CFI
или клиентское предложение не соответствует сохранённому тексту главы,
сервер предпочитает извлечь контекст из `book_chapters.plain_text`; если
точно извлечь нельзя, отдаёт пометку `context_unverified`, но не использует
текст другой книги.

## Данные и кэш

| Таблица | Назначение |
| --- | --- |
| `dictionary_lemmas` | `id`, language, normalized lemma, source и source version; уникальность по этим полям. |
| `dictionary_senses` | lemma FK, POS, список русских переводов, EN/RU пример, source URL, attribution, position. |
| `translation_cache` | hash нормализованного предложения + язык + версия модели, перевод, expiry. |
| `dictionary_import_runs` | версия источника, start/finish, количество записей, ошибки и checksum; позволяет повторяемый импорт. |

Кэш находится в PostgreSQL и является лишь оптимизацией: при промахе ответ
получается вновь. TTL вводится в конфигурацию. В `translation_cache` не
попадают user id, CFI или книга; в логах не пишутся выбранное слово и
предложение, только hash/длина/статус provider.

## gRPC-контракт

```proto
rpc LookupWord(LookupWordRequest) returns (LookupWordResponse);
```

Ответ включает `normalized_lemma`, `repeated DictionarySense senses`,
`SentenceTranslation sentence_translation`, `context_verified` и `source`
metadata. `DictionarySense` содержит `part_of_speech`, повторяемые
`translations`, `example_en`, `example_ru`, `attribution`. У
`SentenceTranslation` есть oneof `translated_text`/`provider_error`; так UI
может показать словарь даже при ошибке LibreTranslate.

Ошибки запроса: `INVALID_ARGUMENT` (пустое/слишком длинное выделение),
`NOT_FOUND` (книга/chapter), `FAILED_PRECONDITION` (не `READY`). Отсутствие
словарной статьи — валидный успешный ответ с пустыми `senses`, а не 404.

## План работ

1. Описать proto для lookup и models ответа, включая отдельное состояние
   ошибки перевода предложения; сгенерировать stubs.
2. Подготовить versioned importer словарного набора, миграции lemmas/senses/
   import runs и документ атрибуции; выполнить маленький тестовый импорт.
3. Реализовать normalizer и morphological port с table-driven тестами для
   регистра, пунктуации, irregular forms и отсутствующей lemma.
4. Реализовать PostgreSQL repository: exact/lemma lookup, порядок senses и
   translation cache с TTL.
5. Поднять LibreTranslate `en`/`ru` в Compose и реализовать HTTP adapter с
   timeout/retry/circuit breaker; добавить fake adapter для unit-тестов.
6. Реализовать `LookupWord` use case с проверкой book/chapter/context и
   частично успешным ответом при сбое машинного перевода.
7. Добавить integration tests: несколько POS/переводов, кэш, пустой результат,
   provider timeout и запрет обратиться к неготовой книге.

## Definition of Done

- Один lookup может отдать несколько частей речи и переводов с примером и
  атрибуцией.
- Переводится только предложение, а не вся глава; provider не видит identity
  пользователя или полный EPUB.
- Ошибка/медленный LibreTranslate не ломает словарную часть карточки.
- Импорт можно повторить с новой версией без ручного изменения production-таблиц.
