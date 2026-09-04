import { useMutation, useQuery } from '@tanstack/react-query'

import {
  getAdjacentChapter,
  getReaderSettings,
  getReadingState,
  saveReadingProgress,
  updateReaderSettings,
} from '../../api/reader'

export function useReaderState(bookId: string | undefined) {
  return useQuery({
    queryKey: ['reader', bookId],
    queryFn: () => getReadingState(bookId ?? ''),
    enabled: Boolean(bookId),
  })
}

export function useAdjacentChapter() {
  return useMutation({
    mutationFn: ({
      bookId,
      chapterId,
      direction,
    }: {
      bookId: string
      chapterId: string
      direction: -1 | 1
    }) => getAdjacentChapter(bookId, chapterId, direction),
  })
}

export function useSaveReadingProgress() {
  return useMutation({ mutationFn: saveReadingProgress })
}

export function useReaderSettings() {
  return useQuery({
    queryKey: ['reader', 'settings'],
    queryFn: getReaderSettings,
  })
}

export function useUpdateReaderSettings() {
  return useMutation({ mutationFn: updateReaderSettings })
}
