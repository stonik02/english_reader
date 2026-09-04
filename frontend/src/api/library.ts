import {
  AddToMyLibraryRequest,
  ListCatalogRequest,
  ListMyBooksRequest,
  RemoveFromMyLibraryRequest,
} from '@reader/proto/reader/v1/library_pb'
import {
  authorizationMetadata,
  libraryClient,
  sessionToken,
  unaryCall,
} from './client'
import { httpApiUrl } from './client'

type ListCatalogParams = {
  cursor?: string
  limit?: number
}

export async function listCatalog({
  cursor = '',
  limit = 20,
}: ListCatalogParams = {}) {
  return unaryCall(() => {
    const request = new ListCatalogRequest()
    request.setCursor(cursor)
    request.setLimit(limit)
    request.setAccessToken(sessionToken.get() ?? '')

    return libraryClient.listCatalog(request, authorizationMetadata())
  })
}

export async function addToMyLibrary(bookId: string) {
  return unaryCall(() => {
    const request = new AddToMyLibraryRequest()
    request.setBookId(bookId)
    request.setAccessToken(sessionToken.get() ?? '')

    return libraryClient.addToMyLibrary(request, authorizationMetadata())
  })
}

export async function listMyBooks({
  cursor = '',
  limit = 20,
}: ListCatalogParams = {}) {
  return unaryCall(() => {
    const request = new ListMyBooksRequest()
    request.setCursor(cursor)
    request.setLimit(limit)
    request.setAccessToken(sessionToken.get() ?? '')
    return libraryClient.listMyBooks(request, authorizationMetadata())
  })
}

export async function removeFromMyLibrary(bookId: string) {
  return unaryCall(() => {
    const request = new RemoveFromMyLibraryRequest()
    request.setBookId(bookId)
    request.setAccessToken(sessionToken.get() ?? '')
    return libraryClient.removeFromMyLibrary(request, authorizationMetadata())
  })
}

export async function deleteBook(bookId: string) {
  const response = await fetch(`${httpApiUrl}/api/v1/library/books/${bookId}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${sessionToken.get() ?? ''}` },
  })
  if (!response.ok) {
    throw new Error('Не удалось удалить книгу.')
  }
}
