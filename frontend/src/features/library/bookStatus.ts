export type BookStatus = 'processing' | 'ready' | 'failed'

type BookStatusPresentation = {
  className: string
  isReady: boolean
  label: string
  unavailableMessage?: string
}

const presentations: Record<BookStatus, BookStatusPresentation> = {
  processing: {
    className: 'book-status-processing',
    isReady: false,
    label: 'Подготавливается',
    unavailableMessage: 'Чтение станет доступно после обработки',
  },
  ready: {
    className: 'book-status-ready',
    isReady: true,
    label: 'Готова к чтению',
  },
  failed: {
    className: 'book-status-failed',
    isReady: false,
    label: 'Ошибка обработки',
    unavailableMessage: 'Файл не прошёл обработку.',
  },
}

export function bookStatusPresentation(status: string): BookStatusPresentation {
  return (
    presentations[status.toLowerCase() as BookStatus] ?? {
      className: 'book-status-unknown',
      isReady: false,
      label: status || 'Неизвестный статус',
      unavailableMessage: 'Чтение пока недоступно.',
    }
  )
}
