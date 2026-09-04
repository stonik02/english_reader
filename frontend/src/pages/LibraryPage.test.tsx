import { cleanup, render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { afterEach, describe, expect, it } from 'vitest'

import { Book } from '../api/gen/reader/v1/library_pb'
import { BookCard } from './LibraryPage'

afterEach(cleanup)

function renderCard(status: string) {
  const book = new Book()
  book.setId('book-1')
  book.setTitle('Harry Potter')
  book.setStatus(status)

  render(
    <MemoryRouter>
      <BookCard
        added={false}
        book={book}
        isAdding={false}
        isDeleting={false}
        deleteError={false}
        onAdd={() => {}}
        onDelete={() => {}}
      />
    </MemoryRouter>,
  )
}

describe('BookCard', () => {
  it('shows a read link for a lowercase ready status', () => {
    renderCard('ready')

    expect(screen.getByRole('link', { name: 'Читать' })).toHaveAttribute(
      'href',
      '/reader/book-1',
    )
    expect(screen.getByText('Готова к чтению')).toBeInTheDocument()
  })

  it('explains that a failed file did not pass processing', () => {
    renderCard('failed')

    expect(screen.getByText('Ошибка обработки')).toBeInTheDocument()
    expect(screen.getByText('Файл не прошёл обработку.')).toBeInTheDocument()
    expect(
      screen.queryByRole('link', { name: 'Читать' }),
    ).not.toBeInTheDocument()
  })
})
