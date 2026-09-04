import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'

import { ApiError } from '../api/client'
import { httpApiUrl, sessionToken } from '../api/client'
import { type Book, type UserBook } from '../api/gen/reader/v1/library_pb'
import {
  useAddToMyLibrary,
  useCatalogBooks,
  useDeleteBook,
  useMyLibraryBooks,
  useRemoveFromMyLibrary,
} from '../features/library/useCatalogBooks'
import { bookStatusPresentation } from '../features/library/bookStatus'
import { UploadBookDialog } from '../features/library/UploadBookDialog'

type LibraryPageProps = { kind: 'catalog' | 'my-library' }

export function LibraryPage({ kind }: LibraryPageProps) {
  return kind === 'catalog' ? <CatalogPage /> : <MyLibraryPage />
}

function CatalogPage() {
  const catalog = useCatalogBooks(true)
  const myLibrary = useMyLibraryBooks(true, 100)
  const addToMyLibrary = useAddToMyLibrary()
  const deleteBook = useDeleteBook()
  const [isUploadDialogOpen, setIsUploadDialogOpen] = useState(false)
  const [addedBookIds, setAddedBookIds] = useState<Set<string>>(() => new Set())
  const books = catalog.data?.pages.flatMap((page) => page.getBooksList()) ?? []
  const myBookIds = new Set(
    myLibrary.data?.pages
      .flatMap((page) => page.getBooksList())
      .flatMap((item) => item.getBook()?.getId() ?? '')
      .filter(Boolean) ?? [],
  )
  const progressByBookID = new Map<string, number>()
  myLibrary.data?.pages
    .flatMap((page) => page.getBooksList())
    .forEach((item) => {
      const bookID = item.getBook()?.getId()
      if (bookID) progressByBookID.set(bookID, item.getProgressPercent())
    })

  function addBook(bookId: string) {
    addToMyLibrary.mutate(bookId, {
      onSuccess: () => {
        setAddedBookIds((current) => new Set(current).add(bookId))
      },
    })
  }

  function removeBook(book: Book) {
    const title = book.getTitle() || 'Книга без названия'
    if (
      !window.confirm(
        `Удалить «${title}»? Книга исчезнет у всех членов семьи без возможности восстановления.`,
      )
    )
      return
    deleteBook.mutate(book.getId())
  }

  return (
    <section className="page-section" aria-labelledby="page-title">
      <div className="page-heading">
        <div>
          <p className="eyebrow">Библиотека</p>
          <h1 id="page-title">Общий каталог</h1>
        </div>
        <button
          className="button button-primary"
          onClick={() => setIsUploadDialogOpen(true)}
          type="button"
        >
          Загрузить EPUB
        </button>
      </div>

      {catalog.isPending && <LibraryStatus>Загружаем книги…</LibraryStatus>}
      {catalog.isError && (
        <LibraryError
          error={catalog.error}
          retry={() => void catalog.refetch()}
        />
      )}
      {catalog.isSuccess && books.length === 0 && (
        <LibraryStatus>
          В общем каталоге пока нет книг. Загрузите первый EPUB для семьи.
        </LibraryStatus>
      )}
      {books.length > 0 && (
        <>
          <div className="book-grid" aria-label="Книги общего каталога">
            {books.map((book) => (
              <BookCard
                added={
                  addedBookIds.has(book.getId()) || myBookIds.has(book.getId())
                }
                book={book}
                isAdding={
                  addToMyLibrary.isPending &&
                  addToMyLibrary.variables === book.getId()
                }
                isDeleting={
                  deleteBook.isPending && deleteBook.variables === book.getId()
                }
                deleteError={
                  deleteBook.isError && deleteBook.variables === book.getId()
                }
                progressPercent={progressByBookID.get(book.getId())}
                key={book.getId()}
                onAdd={() => addBook(book.getId())}
                onDelete={() => removeBook(book)}
              />
            ))}
          </div>
          {catalog.hasNextPage && (
            <button
              className="button button-secondary load-more"
              disabled={catalog.isFetchingNextPage}
              onClick={() => void catalog.fetchNextPage()}
              type="button"
            >
              {catalog.isFetchingNextPage ? 'Загружаем…' : 'Показать ещё'}
            </button>
          )}
        </>
      )}
      {isUploadDialogOpen && (
        <UploadBookDialog onClose={() => setIsUploadDialogOpen(false)} />
      )}
    </section>
  )
}

export function BookCard({
  added,
  book,
  isAdding,
  deleteError,
  isDeleting,
  onAdd,
  onDelete,
  progressPercent,
}: {
  added: boolean
  book: Book
  isAdding: boolean
  deleteError: boolean
  isDeleting: boolean
  onAdd(): void
  onDelete(): void
  progressPercent?: number
}) {
  const status = bookStatusPresentation(book.getStatus())
  const title = book.getTitle() || 'Книга без названия'

  return (
    <article className="book-card">
      <BookCover coverURL={book.getCoverUrl()} title={title} />
      <div className="book-card-body">
        <span className={`book-status ${status.className}`}>
          {status.label}
        </span>
        <h2>{title}</h2>
        <p>{book.getAuthor() || 'Автор не указан'}</p>
        {progressPercent !== undefined && progressPercent > 0 && (
          <div className="reading-progress">
            <span>Прочитано: {Math.round(progressPercent)}%</span>
            <progress value={progressPercent} max="100">
              {Math.round(progressPercent)}%
            </progress>
          </div>
        )}
        <div className="book-actions">
          {status.isReady ? (
            <Link
              className="button button-secondary"
              to={`/reader/${book.getId()}`}
            >
              {progressPercent !== undefined && progressPercent > 0
                ? 'Продолжить'
                : 'Читать'}
            </Link>
          ) : (
            <span className="book-unavailable">
              {status.unavailableMessage}
            </span>
          )}
          <button
            className="button button-primary"
            disabled={!status.isReady || added || isAdding}
            onClick={onAdd}
            type="button"
          >
            {added ? 'В моей библиотеке' : isAdding ? 'Добавляем…' : 'Добавить'}
          </button>
          <button
            className="button button-danger"
            disabled={isDeleting}
            onClick={onDelete}
            type="button"
          >
            {isDeleting ? 'Удаляем…' : 'Удалить книгу'}
          </button>
          {deleteError && (
            <span className="book-delete-error" role="alert">
              Не удалось удалить книгу. Повторите попытку.
            </span>
          )}
        </div>
      </div>
    </article>
  )
}

function BookCover({ coverURL, title }: { coverURL: string; title: string }) {
  const [imageURL, setImageURL] = useState<string | null>(null)

  useEffect(() => {
    if (!coverURL) {
      setImageURL(null)
      return
    }
    const controller = new AbortController()
    let objectURL: string | null = null

    void fetch(`${httpApiUrl}${coverURL}`, {
      headers: { Authorization: `Bearer ${sessionToken.get() ?? ''}` },
      signal: controller.signal,
    })
      .then(async (response) => {
        if (!response.ok) throw new Error('cover is unavailable')
        objectURL = URL.createObjectURL(await response.blob())
        setImageURL(objectURL)
      })
      .catch(() => setImageURL(null))

    return () => {
      controller.abort()
      if (objectURL) URL.revokeObjectURL(objectURL)
    }
  }, [coverURL])

  if (imageURL) {
    return (
      <img
        className="book-cover book-cover-image"
        src={imageURL}
        alt={`Обложка: ${title}`}
      />
    )
  }
  return (
    <div className="book-cover" aria-hidden="true">
      {title.slice(0, 1).toUpperCase()}
    </div>
  )
}

function MyLibraryPage() {
  const library = useMyLibraryBooks(true)
  const remove = useRemoveFromMyLibrary()
  const books = library.data?.pages.flatMap((page) => page.getBooksList()) ?? []

  function removeBook(item: UserBook) {
    const book = item.getBook()
    if (!book) return
    const title = book.getTitle() || 'Книга без названия'
    if (
      !window.confirm(
        `Убрать «${title}» из моей библиотеки? Книга останется в общем каталоге.`,
      )
    )
      return
    remove.mutate(book.getId())
  }

  return (
    <section className="page-section" aria-labelledby="page-title">
      <div className="page-heading">
        <div>
          <p className="eyebrow">Библиотека</p>
          <h1 id="page-title">Моя библиотека</h1>
          <p className="page-description">
            Здесь появятся книги, которые вы добавили или начали читать.
          </p>
        </div>
        <Link className="button button-primary" to="/catalog">
          Открыть каталог
        </Link>
      </div>
      {library.isPending && (
        <LibraryStatus>Загружаем вашу библиотеку…</LibraryStatus>
      )}
      {library.isError && (
        <LibraryError
          error={library.error}
          retry={() => void library.refetch()}
        />
      )}
      {library.isSuccess && books.length === 0 && (
        <LibraryStatus>
          Ваша библиотека пока пуста. Добавьте готовую книгу из общего каталога.
        </LibraryStatus>
      )}
      {books.length > 0 && (
        <div className="book-grid" aria-label="Книги моей библиотеки">
          {books.map((item) => (
            <MyBookCard
              item={item}
              isRemoving={
                remove.isPending && remove.variables === item.getBook()?.getId()
              }
              key={item.getBook()?.getId()}
              onRemove={() => removeBook(item)}
              removeError={
                remove.isError && remove.variables === item.getBook()?.getId()
              }
            />
          ))}
        </div>
      )}
    </section>
  )
}

function MyBookCard({
  item,
  isRemoving,
  onRemove,
  removeError,
}: {
  item: UserBook
  isRemoving: boolean
  onRemove(): void
  removeError: boolean
}) {
  const book = item.getBook()
  if (!book) return null
  const title = book.getTitle() || 'Книга без названия'
  const progress = Math.round(item.getProgressPercent())
  return (
    <article className="book-card">
      <BookCover coverURL={book.getCoverUrl()} title={title} />
      <div className="book-card-body">
        <span className="book-status book-status-ready">В моей библиотеке</span>
        <h2>{title}</h2>
        <p>{book.getAuthor() || 'Автор не указан'}</p>
        <div className="reading-progress">
          <span>Прочитано: {progress}%</span>
          <progress value={progress} max="100">
            {progress}%
          </progress>
        </div>
        <div className="book-actions">
          <Link
            className="button button-secondary"
            to={`/reader/${book.getId()}`}
          >
            Читать
          </Link>
          <button
            className="button button-danger"
            disabled={isRemoving}
            onClick={onRemove}
            type="button"
          >
            {isRemoving ? 'Убираем…' : 'Убрать из моей библиотеки'}
          </button>
          {removeError && (
            <span className="book-delete-error" role="alert">
              Не удалось обновить библиотеку. Повторите попытку.
            </span>
          )}
        </div>
      </div>
    </article>
  )
}

function LibraryStatus({ children }: { children: string }) {
  return (
    <div className="empty-state">
      <span className="empty-state-icon" aria-hidden="true">
        ▤
      </span>
      <p>{children}</p>
    </div>
  )
}

function LibraryError({ error, retry }: { error: unknown; retry(): void }) {
  const message =
    error instanceof ApiError ? error.message : 'Не удалось загрузить каталог.'

  return (
    <div className="empty-state" role="alert">
      <h2>Каталог недоступен</h2>
      <p>{message}</p>
      <button className="button button-secondary" onClick={retry} type="button">
        Повторить
      </button>
    </div>
  )
}
