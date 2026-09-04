import {
  type CSSProperties,
  useEffect,
  useLayoutEffect,
  useRef,
  useState,
} from 'react'
import { Link, useParams } from 'react-router-dom'

import { ApiError } from '../api/client'
import { type Chapter } from '../api/gen/reader/v1/reader_pb'
import {
  useAdjacentChapter,
  useReaderState,
  useReaderSettings,
  useSaveReadingProgress,
  useUpdateReaderSettings,
} from '../features/reader/useReader'
import { contextWindow } from '../features/reader/contextWindow'
import { useDictionaryLookup } from '../features/reader/useDictionaryLookup'
import { applyHighlights } from '../features/vocabulary/applyHighlights'
import {
  useHighlights,
  useSaveVocabularyEntry,
} from '../features/vocabulary/useVocabulary'

export function ReaderPage() {
  const { bookId } = useParams()
  const state = useReaderState(bookId)
  const adjacent = useAdjacentChapter()
  const saveProgress = useSaveReadingProgress()
  const settingsQuery = useReaderSettings()
  const updateSettings = useUpdateReaderSettings()
  const dictionary = useDictionaryLookup()
  const saveVocabulary = useSaveVocabularyEntry()
  const [chapter, setChapter] = useState<Chapter | null>(null)
  const [revision, setRevision] = useState(0)
  const [progressPercent, setProgressPercent] = useState(0)
  const [settingsOpen, setSettingsOpen] = useState(false)
  const [fontScale, setFontScale] = useState(100)
  const [theme, setTheme] = useState('system')
  const [lineHeight, setLineHeight] = useState(1.5)
  const [highlightColor, setHighlightColor] = useState('yellow')
  const [selection, setSelection] = useState<SelectionContext | null>(null)
  const [savedLemmaIds, setSavedLemmaIds] = useState<number[]>([])
  const chapterContent = useRef<HTMLElement>(null)
  const revisionRef = useRef(0)
  const chapterNavigationRef = useRef(false)
  const saveProgressRef = useRef(saveProgress.mutate)
  const saveProgressPendingRef = useRef(saveProgress.isPending)
  const wholeSentencePressTimer = useRef<number | null>(null)
  const translateWholeSentenceRef = useRef(false)
  const highlights = useHighlights(bookId, chapter?.getId())

  revisionRef.current = revision
  saveProgressRef.current = saveProgress.mutate
  saveProgressPendingRef.current = saveProgress.isPending

  useEffect(() => {
    if (state.data?.getChapter()) {
      setChapter(state.data.getChapter() ?? null)
      setRevision(state.data.getProgress()?.getRevision() ?? 0)
      setProgressPercent(state.data.getProgress()?.getProgressPercent() ?? 0)
    }
  }, [state.data])

  useEffect(() => {
    const settings = settingsQuery.data
    if (!settings) return
    setFontScale(settings.getFontScale())
    setTheme(settings.getTheme())
    setLineHeight(settings.getLineHeight())
    // A running Vite session can temporarily retain an older generated protobuf
    // constructor during an API-contract update. Keep the reader usable until
    // that dependency cache is refreshed.
    const getHighlightColor = settings.getHighlightColor
    setHighlightColor(
      typeof getHighlightColor === 'function'
        ? getHighlightColor.call(settings) || 'yellow'
        : 'yellow',
    )
  }, [settingsQuery.data])

  useEffect(() => {
    if (!selection || !dictionary.data) return
    const closeOnEscape = (event: KeyboardEvent) => {
      if (event.key === 'Escape') setSelection(null)
    }
    window.addEventListener('keydown', closeOnEscape)
    return () => window.removeEventListener('keydown', closeOnEscape)
  }, [dictionary.data, selection])

  useEffect(() => {
    if (!chapterContent.current) return
    applyHighlights(
      chapterContent.current,
      highlights.data?.getTokensList().map((token) => {
        const getTexts = token.getTextsList
        return {
          texts:
            typeof getTexts === 'function'
              ? getTexts.call(token)
              : [token.getLemma()],
        }
      }) ?? [],
    )
  }, [chapter, highlights.data, highlightColor])

  useLayoutEffect(() => {
    if (!chapterNavigationRef.current) return
    window.scrollTo({ top: 0 })
  }, [chapter])

  function saveSettings() {
    updateSettings.mutate({ fontScale, theme, lineHeight, highlightColor })
  }

  useEffect(() => {
    if (!bookId || !chapter) return
    const key = `reader-scroll:${bookId}:${chapter.getId()}`
    const navigatedBetweenChapters = chapterNavigationRef.current
    chapterNavigationRef.current = false
    const savedRatio = Number.parseFloat(localStorage.getItem(key) ?? '')
    const restore = window.requestAnimationFrame(() => {
      if (navigatedBetweenChapters) {
        window.scrollTo({ top: 0 })
        return
      }
      if (!Number.isFinite(savedRatio)) return
      const maxScroll =
        document.documentElement.scrollHeight - window.innerHeight
      window.scrollTo({ top: Math.max(0, maxScroll * savedRatio) })
    })
    let timer: number | undefined
    const persist = () => {
      window.clearTimeout(timer)
      timer = window.setTimeout(() => {
        const maxScroll =
          document.documentElement.scrollHeight - window.innerHeight
        const ratio = maxScroll > 0 ? window.scrollY / maxScroll : 0
        localStorage.setItem(key, String(Math.min(1, Math.max(0, ratio))))
        if (saveProgressPendingRef.current) return
        const totalChapters = Math.max(chapter.getTotalChapters(), 1)
        const progress = ((chapter.getSequence() + ratio) / totalChapters) * 100
        saveProgressRef.current(
          {
            bookId,
            chapterId: chapter.getId(),
            epubCfi: chapter.getStartCfi(),
            progressPercent: Math.min(100, Math.max(0, progress)),
            revision: revisionRef.current + 1,
          },
          {
            onSuccess: (saved) => {
              setRevision(saved.getRevision())
              setProgressPercent(saved.getProgressPercent())
            },
          },
        )
      }, 800)
    }
    window.addEventListener('scroll', persist, { passive: true })
    return () => {
      window.cancelAnimationFrame(restore)
      window.removeEventListener('scroll', persist)
      window.clearTimeout(timer)
      persist()
    }
  }, [bookId, chapter])

  function changeChapter(direction: -1 | 1) {
    if (!bookId || !chapter) return
    adjacent.mutate(
      { bookId, chapterId: chapter.getId(), direction },
      {
        onSuccess: (nextChapter) => {
          chapterNavigationRef.current = true
          window.scrollTo({ top: 0 })
          setChapter(nextChapter)
          const currentProgress = state.data?.getProgress()
          saveProgress.mutate(
            {
              bookId,
              chapterId: nextChapter.getId(),
              epubCfi: nextChapter.getStartCfi(),
              progressPercent: currentProgress?.getProgressPercent() ?? 0,
              revision: revision + 1,
            },
            { onSuccess: (saved) => setRevision(saved.getRevision()) },
          )
        },
      },
    )
  }

  function turnPage(direction: -1 | 1) {
    const pageHeight = Math.max(window.innerHeight * 0.85, 320)
    const atStart = window.scrollY <= 1
    const atEnd =
      window.scrollY + window.innerHeight >=
      document.documentElement.scrollHeight - 2

    if ((direction === -1 && atStart) || (direction === 1 && atEnd)) {
      changeChapter(direction)
      return
    }
    window.scrollBy({ top: direction * pageHeight, behavior: 'smooth' })
  }

  function lookUp(word: string, sentence: string) {
    if (!bookId || !chapter) return
    dictionary.mutate({
      bookId,
      chapterId: chapter.getId(),
      selectedText: word,
      sentenceText: sentence,
      epubCfi: chapter.getStartCfi(),
    })
  }

  function captureSelection(translateWholeSentence = false) {
    const browserSelection = window.getSelection()
    const word = browserSelection?.toString().trim() ?? ''
    if (!word) {
      setSelection(null)
      return
    }
    dictionary.reset()
    const text = selectionContainerText(browserSelection)
    const translationText = translateWholeSentence
      ? sentenceFor(word, text)
      : contextWindow(word, text)
    setSelection({ word, sentence: translationText })
    lookUp(word, translationText)
  }

  function startTranslationPress() {
    cancelTranslationPress()
    translateWholeSentenceRef.current = false
    wholeSentencePressTimer.current = window.setTimeout(() => {
      translateWholeSentenceRef.current = true
    }, 600)
  }

  function finishTranslationPress() {
    const translateWholeSentence = translateWholeSentenceRef.current
    cancelTranslationPress()
    // Native selection is finalized after pointerup, especially on touch screens.
    window.setTimeout(() => captureSelection(translateWholeSentence), 0)
  }

  function cancelTranslationPress() {
    if (wholeSentencePressTimer.current !== null) {
      window.clearTimeout(wholeSentencePressTimer.current)
      wholeSentencePressTimer.current = null
    }
  }

  function saveSelectedWord() {
    if (!selection || !dictionary.data) return
    const lemmaId = dictionary.data.getLemmaId()
    if (!lemmaId) return
    saveVocabulary.mutate(
      {
        lemmaId,
        chosenSenseId: dictionary.data.getSensesList()[0]?.getId(),
        sourceForm: selection.word,
      },
      { onSuccess: () => setSavedLemmaIds((ids) => [...ids, lemmaId]) },
    )
  }

  if (state.isPending) return <ReaderStatus>Открываем книгу…</ReaderStatus>
  if (state.isError) {
    const message =
      state.error instanceof ApiError
        ? state.error.message
        : 'Не удалось открыть книгу.'
    return <ReaderStatus error={message} retry={() => void state.refetch()} />
  }
  if (!chapter) return <ReaderStatus>В книге нет доступных глав.</ReaderStatus>
  const progress = Math.round(progressPercent)

  return (
    <main className={`reader-page reader-theme-${theme}`}>
      <header className="reader-header">
        <Link className="back-link" to="/library">
          ← В библиотеку
        </Link>
        <span className="reader-title">Ридер</span>
        <button
          className="text-button"
          type="button"
          onClick={() => setSettingsOpen((open) => !open)}
        >
          Настройки текста
        </button>
      </header>
      {settingsOpen && (
        <section className="reader-settings" aria-label="Настройки текста">
          <label>
            Размер шрифта: {fontScale}%
            <input
              type="range"
              min="80"
              max="200"
              value={fontScale}
              onChange={(event) => setFontScale(Number(event.target.value))}
            />
          </label>
          <label>
            Межстрочный интервал: {lineHeight.toFixed(1)}
            <input
              type="range"
              min="1"
              max="3"
              step="0.1"
              value={lineHeight}
              onChange={(event) => setLineHeight(Number(event.target.value))}
            />
          </label>
          <label>
            Тема
            <select
              value={theme}
              onChange={(event) => setTheme(event.target.value)}
            >
              <option value="system">Системная</option>
              <option value="light">Светлая</option>
              <option value="dark">Тёмная</option>
            </select>
          </label>
          <fieldset className="highlight-color-picker">
            <legend>Цвет подсветки словаря</legend>
            {[
              ['yellow', 'Жёлтый'],
              ['blue', 'Синий'],
              ['green', 'Зелёный'],
              ['pink', 'Розовый'],
              ['orange', 'Оранжевый'],
              ['purple', 'Фиолетовый'],
            ].map(([value, label]) => (
              <label key={value}>
                <input
                  type="radio"
                  name="highlight-color"
                  value={value}
                  checked={highlightColor === value}
                  onChange={() => setHighlightColor(value)}
                />
                <span className={`highlight-color-swatch highlight-${value}`} />
                {label}
              </label>
            ))}
          </fieldset>
          <button
            className="button button-secondary"
            type="button"
            disabled={updateSettings.isPending}
            onClick={saveSettings}
          >
            {updateSettings.isPending ? 'Сохраняем…' : 'Сохранить настройки'}
          </button>
          {updateSettings.isError && (
            <p className="reader-error" role="alert">
              Не удалось сохранить настройки.
            </p>
          )}
        </section>
      )}
      <section className="reader-content" aria-labelledby="reader-title">
        <p className="eyebrow">Глава {chapter.getSequence() + 1}</p>
        <h1 id="reader-title">Чтение книги</h1>
        <p className="reader-progress">Текущий прогресс: {progress}%</p>
        <article
          ref={chapterContent}
          className="chapter-content"
          style={
            {
              fontSize: `${fontScale}%`,
              lineHeight,
              '--vocabulary-highlight-color': `var(--highlight-${highlightColor})`,
            } as CSSProperties
          }
          dangerouslySetInnerHTML={{ __html: chapter.getSanitizedHtml() }}
          onPointerCancel={cancelTranslationPress}
          onPointerDown={startTranslationPress}
          onPointerUp={finishTranslationPress}
        />
        {selection && (
          <div
            className="modal-backdrop"
            onMouseDown={() => setSelection(null)}
          >
            <section
              className="translation-panel"
              role="dialog"
              aria-modal="true"
              aria-labelledby="translation-title"
              onMouseDown={(event) => event.stopPropagation()}
            >
              <h2 id="translation-title">{selection.word}</h2>
              {dictionary.isPending && <p>Ищем перевод…</p>}
              {dictionary.isError && (
                <div className="reader-error" role="alert">
                  Не удалось найти перевод.
                  <button
                    className="text-button"
                    type="button"
                    onClick={() => lookUp(selection.word, selection.sentence)}
                  >
                    Повторить
                  </button>
                </div>
              )}
              {dictionary.data && (
                <>
                  <h3>{dictionary.data.getNormalizedLemma()}</h3>
                  {groupSenses(dictionary.data.getSensesList()).map((sense) => (
                    <div key={sense.partOfSpeech}>
                      <p>
                        <strong>{sense.partOfSpeech}</strong>:{' '}
                        {sense.translations.join(', ')}
                      </p>
                    </div>
                  ))}
                  {dictionary.data
                    .getSentenceTranslation()
                    ?.getTranslatedText() && (
                    <div className="translation-context">
                      <p>
                        <ContextPreview
                          context={selection.sentence}
                          selectedWord={selection.word}
                        />
                      </p>
                      <p>
                        <em>
                          {dictionary.data
                            .getSentenceTranslation()
                            ?.getTranslatedText()}
                        </em>
                      </p>
                    </div>
                  )}
                  {dictionary.data
                    .getSentenceTranslation()
                    ?.getProviderError() && (
                    <p className="reader-error">Перевод временно недоступен.</p>
                  )}
                  {groupSenses(dictionary.data.getSensesList()).map((sense) => (
                    <div key={`examples-${sense.partOfSpeech}`}>
                      {sense.examples.map((example) => (
                        <div key={`${example.en}-${example.ru}`}>
                          {example.en && <p>{example.en}</p>}
                          {example.ru && <p>{example.ru}</p>}
                        </div>
                      ))}
                    </div>
                  ))}
                  {dictionary.data.getAlreadySaved() ||
                  savedLemmaIds.includes(dictionary.data.getLemmaId()) ? (
                    <p>Слово уже в вашем словаре.</p>
                  ) : (
                    <button
                      className="button button-secondary"
                      type="button"
                      disabled={
                        saveVocabulary.isPending ||
                        !dictionary.data.getLemmaId()
                      }
                      onClick={saveSelectedWord}
                    >
                      {saveVocabulary.isPending
                        ? 'Добавляем…'
                        : 'В мой словарь'}
                    </button>
                  )}
                  {saveVocabulary.isError && (
                    <p className="reader-error" role="alert">
                      Не удалось добавить слово в словарь.
                    </p>
                  )}
                </>
              )}
              <button
                className="text-button"
                type="button"
                onClick={() => setSelection(null)}
              >
                Закрыть
              </button>
            </section>
          </div>
        )}
        <div className="reader-controls" aria-label="Навигация по книге">
          <button
            className="button button-secondary"
            type="button"
            disabled={adjacent.isPending || saveProgress.isPending}
            onClick={() => turnPage(-1)}
          >
            ← Предыдущая страница
          </button>
          <span>
            {adjacent.isPending
              ? 'Загружаем главу…'
              : `Глава ${chapter.getSequence() + 1}`}
          </span>
          <button
            className="button button-secondary"
            type="button"
            disabled={adjacent.isPending || saveProgress.isPending}
            onClick={() => turnPage(1)}
          >
            Следующая страница →
          </button>
        </div>
        {adjacent.isError && (
          <p className="reader-error" role="alert">
            Следующей главы нет или она недоступна.
          </p>
        )}
        {saveProgress.isError && (
          <p className="reader-error" role="alert">
            Не удалось сохранить новую главу. Повторите переход позже.
          </p>
        )}
      </section>
      <nav className="chapter-switcher" aria-label="Переход между главами">
        <button
          className="chapter-switcher-button chapter-switcher-previous"
          type="button"
          disabled={adjacent.isPending || saveProgress.isPending}
          onClick={() => changeChapter(-1)}
          aria-label="Предыдущая глава"
        >
          ‹
        </button>
        <button
          className="chapter-switcher-button chapter-switcher-next"
          type="button"
          disabled={adjacent.isPending || saveProgress.isPending}
          onClick={() => changeChapter(1)}
          aria-label="Следующая глава"
        >
          ›
        </button>
      </nav>
    </main>
  )
}

type SelectionContext = {
  word: string
  sentence: string
}

type DictionarySenseLike = {
  getPartOfSpeech(): string
  getTranslationsList(): string[]
  getExampleEn(): string
  getExampleRu(): string
}

type DisplaySense = {
  partOfSpeech: string
  translations: string[]
  examples: Array<{ en: string; ru: string }>
}

function groupSenses(senses: DictionarySenseLike[]): DisplaySense[] {
  const grouped = new Map<string, DisplaySense>()
  for (const sense of senses) {
    const partOfSpeech = sense.getPartOfSpeech() || 'other'
    const group = grouped.get(partOfSpeech) ?? {
      partOfSpeech,
      translations: [],
      examples: [],
    }
    for (const translation of sense.getTranslationsList()) {
      if (!group.translations.includes(translation)) {
        group.translations.push(translation)
      }
    }
    const example = { en: sense.getExampleEn(), ru: sense.getExampleRu() }
    if (
      (example.en || example.ru) &&
      !group.examples.some(
        (value) => value.en === example.en && value.ru === example.ru,
      )
    ) {
      group.examples.push(example)
    }
    grouped.set(partOfSpeech, group)
  }
  return [...grouped.values()]
}

function ContextPreview({
  context,
  selectedWord,
}: {
  context: string
  selectedWord: string
}) {
  const expression = new RegExp(`(${escapeRegExp(selectedWord)})`, 'gi')
  return context
    .split(expression)
    .map((part, index) =>
      part.toLocaleLowerCase() === selectedWord.toLocaleLowerCase() ? (
        <strong key={index}>{part}</strong>
      ) : (
        part
      ),
    )
}

function escapeRegExp(value: string) {
  return value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}

function selectionContainerText(selection: Selection | null) {
  const node = selection?.rangeCount
    ? selection.getRangeAt(0).commonAncestorContainer
    : null
  const element =
    node?.nodeType === Node.ELEMENT_NODE
      ? (node as Element)
      : node?.parentElement
  return element?.closest('p, li, h1, h2, h3, h4, h5, h6, blockquote')
    ?.textContent
}

function sentenceFor(word: string, text: string | null | undefined) {
  const sentences = (text ?? word).match(/[^.!?]+[.!?]?/g) ?? []
  const normalizedWord = word.toLocaleLowerCase()
  return (
    sentences
      .find((sentence) => sentence.toLocaleLowerCase().includes(normalizedWord))
      ?.trim() || word
  )
}

function ReaderStatus({
  children,
  error,
  retry,
}: {
  children?: string
  error?: string
  retry?: () => void
}) {
  return (
    <main className="reader-page">
      <section className="reader-placeholder">
        <p>{error ?? children}</p>
        {retry && (
          <button
            className="button button-secondary"
            onClick={retry}
            type="button"
          >
            Повторить
          </button>
        )}
      </section>
    </main>
  )
}
