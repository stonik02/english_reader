import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

import {
  deleteVocabularyEntry,
  getHighlights,
  listVocabularyEntries,
  saveVocabularyEntry,
} from '../../api/vocabulary'

export function useVocabularyEntries(query: string, cursor = '') {
  return useQuery({
    queryKey: ['vocabulary', 'entries', query, cursor],
    queryFn: () => listVocabularyEntries({ query, cursor }),
  })
}

export function useSaveVocabularyEntry() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: saveVocabularyEntry,
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: ['vocabulary'] }),
  })
}

export function useDeleteVocabularyEntry() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: deleteVocabularyEntry,
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: ['vocabulary'] }),
  })
}

export function useHighlights(
  bookId: string | undefined,
  chapterId: string | undefined,
) {
  return useQuery({
    queryKey: ['vocabulary', 'highlights', bookId, chapterId],
    queryFn: () => getHighlights(bookId ?? '', chapterId ?? ''),
    enabled: Boolean(bookId && chapterId),
  })
}
