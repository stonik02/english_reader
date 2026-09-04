import {
  useInfiniteQuery,
  useMutation,
  useQueryClient,
} from '@tanstack/react-query'

import {
  addToMyLibrary,
  deleteBook,
  listCatalog,
  listMyBooks,
  removeFromMyLibrary,
} from '../../api/library'
import { uploadBook } from '../../api/uploadBook'

const pageSize = 20

export function useCatalogBooks(enabled: boolean) {
  return useInfiniteQuery({
    queryKey: ['books', 'catalog'],
    queryFn: ({ pageParam }) =>
      listCatalog({ cursor: pageParam, limit: pageSize }),
    initialPageParam: '',
    getNextPageParam: (lastPage) => lastPage.getNextCursor() || undefined,
    enabled,
  })
}

export function useAddToMyLibrary() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: addToMyLibrary,
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: ['books', 'my-library'] }),
  })
}

export function useMyLibraryBooks(enabled: boolean, limit = pageSize) {
  return useInfiniteQuery({
    queryKey: ['books', 'my-library'],
    queryFn: ({ pageParam }) => listMyBooks({ cursor: pageParam, limit }),
    initialPageParam: '',
    getNextPageParam: (lastPage) => lastPage.getNextCursor() || undefined,
    enabled,
  })
}

export function useRemoveFromMyLibrary() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: removeFromMyLibrary,
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: ['books', 'my-library'] }),
  })
}

export function useUploadBook() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({
      file,
      onProgress,
    }: {
      file: File
      onProgress(percent: number): void
    }) => uploadBook(file, onProgress),
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: ['books', 'catalog'] }),
  })
}

export function useDeleteBook() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: deleteBook,
    onSuccess: () =>
      Promise.all([
        queryClient.invalidateQueries({ queryKey: ['books', 'catalog'] }),
        queryClient.invalidateQueries({ queryKey: ['books', 'my-library'] }),
      ]),
  })
}
