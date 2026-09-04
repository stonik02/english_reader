# Импорт словарного набора

Приложение не скачивает словарные данные во время пользовательского запроса и
не использует платный внешний Dictionary API. Оператор самостоятельно выбирает
юридически допустимый открытый dump, проверяет условия лицензии и импортирует
его как новую версию.

## Формат JSONL

Каждая строка описывает один sense и содержит следующие поля:

```json
{
  "lemma": "go",
  "language": "en",
  "part_of_speech": "verb",
  "translations": ["идти", "ехать"],
  "example_en": "We go home.",
  "example_ru": "Мы идём домой.",
  "source_url": "https://source.example/entry/go",
  "attribution": "Источник и текст лицензии",
  "position": 0
}
```

`source_url` и `attribution` обязательны для каждого sense: frontend показывает
их в карточке слова. Для Wiktionary необходимо соблюдать условия CC BY-SA;
подробности — в [DICTIONARY_ATTRIBUTION.md](DICTIONARY_ATTRIBUTION.md).

## Импорт новой версии

Сначала примените миграции и подготовьте файл вне git-репозитория. Затем
задайте явные источник и версию:

```bash
cd backend
make dictionary-import \
  DICTIONARY_FILE=/absolute/path/english-russian.jsonl \
  DICTIONARY_SOURCE=wiktionary \
  DICTIONARY_VERSION=2026-09-03
```

Команда создаёт запись в `dictionary_import_runs` с source, version и checksum.
Повторный импорт той же версии идемпотентен для `(language, lemma, source,
source_version)`; новая версия сохраняется рядом со старой. Перед production
импортом проверьте размер dump, лицензию и наличие attribution у каждой строки.

Для локальной разработки есть только небольшой fixture
`test/fixtures/dictionary-mini.jsonl`; он не является словарным набором для
семейного использования.

## Kaikki/Wiktextract

Полная английская выгрузка Kaikki называется `raw-wiktextract-data.jsonl.gz`.
Она содержит записи многих языков; используйте
`tools/convert_kaikki_en_ru.py`, чтобы потоково оставить только English lemma
с Russian translations. Конвертер не распаковывает архив целиком и сохраняет
source URL, attribution, часть речи и первый английский пример.
