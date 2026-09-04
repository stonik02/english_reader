import { useEffect, useState } from 'react'

import { ApiError } from '../api/client'
import {
  useDeleteVocabularyEntry,
  useVocabularyEntries,
} from '../features/vocabulary/useVocabulary'

export function VocabularyPage() {
  const [query, setQuery] = useState('')
  const [cursor, setCursor] = useState('')
  const entries = useVocabularyEntries(query, cursor)
  const remove = useDeleteVocabularyEntry()

  useEffect(() => setCursor(''), [query])

  const message =
    entries.error instanceof ApiError
      ? entries.error.message
      : 'Не удалось загрузить словарь.'

  return (
    <section className="page-section" aria-labelledby="vocabulary-title">
      <div className="page-heading">
        <div>
          <p className="eyebrow">Личная практика</p>
          <h1 id="vocabulary-title">Мой словарь</h1>
          <p className="page-description">
            Слова, сохранённые во время чтения.
          </p>
        </div>
        <label className="search-field">
          <span className="visually-hidden">Поиск слова</span>
          <input
            placeholder="Поиск слова"
            type="search"
            value={query}
            onChange={(event) => setQuery(event.target.value)}
          />
        </label>
      </div>
      {entries.isPending && <p>Загружаем словарь…</p>}
      {entries.isError && (
        <div className="reader-error" role="alert">
          {message}{' '}
          <button
            className="text-button"
            type="button"
            onClick={() => void entries.refetch()}
          >
            Повторить
          </button>
        </div>
      )}
      {entries.data && entries.data.getEntriesList().length === 0 && (
        <div className="empty-state">
          <span className="empty-state-icon" aria-hidden="true">
            Аа
          </span>
          <h2>{query ? 'Ничего не найдено' : 'Словарь ждёт первое слово'}</h2>
          <p>
            {query
              ? 'Попробуйте изменить запрос.'
              : 'Добавьте слово из окна перевода во время чтения.'}
          </p>
        </div>
      )}
      {entries.data && entries.data.getEntriesList().length > 0 && (
        <ul className="vocabulary-list">
          {entries.data.getEntriesList().map((entry) => (
            <li key={entry.getId()} className="vocabulary-entry">
              <div>
                <h2>{entry.getLemma()}</h2>
                {entry.getChosenSense() ? (
                  <>
                    <p>
                      {entry.getChosenSense()?.getPartOfSpeech()}:{' '}
                      {entry.getChosenSense()?.getTranslationsList().join(', ')}
                    </p>
                    {entry.getChosenSense()?.getExampleEn() && (
                      <p>{entry.getChosenSense()?.getExampleEn()}</p>
                    )}
                    {entry.getChosenSense()?.getExampleRu() && (
                      <p>{entry.getChosenSense()?.getExampleRu()}</p>
                    )}
                  </>
                ) : (
                  <p>Сохранено из текста: {entry.getSourceForm()}</p>
                )}
              </div>
              <button
                className="text-button"
                type="button"
                disabled={remove.isPending}
                onClick={() => remove.mutate(entry.getId())}
              >
                Удалить
              </button>
            </li>
          ))}
        </ul>
      )}
      {remove.isError && (
        <p className="reader-error" role="alert">
          Не удалось удалить слово.
        </p>
      )}
      {entries.data?.getNextCursor() && (
        <button
          className="button button-secondary load-more"
          type="button"
          onClick={() => setCursor(entries.data?.getNextCursor() ?? '')}
        >
          Следующая страница
        </button>
      )}
      {cursor && (
        <button
          className="text-button"
          type="button"
          onClick={() => setCursor('')}
        >
          К началу списка
        </button>
      )}
    </section>
  )
}
