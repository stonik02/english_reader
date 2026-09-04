import {
  DeleteEntryRequest,
  GetHighlightsRequest,
  ListEntriesRequest,
  SaveEntryRequest,
} from '@reader/proto/reader/v1/vocabulary_pb'

import {
  authorizationMetadata,
  sessionToken,
  unaryCall,
  vocabularyClient,
} from './client'

export async function saveVocabularyEntry({
  lemmaId,
  chosenSenseId,
  sourceForm,
}: {
  lemmaId: number
  chosenSenseId?: number
  sourceForm: string
}) {
  return unaryCall(() => {
    const request = new SaveEntryRequest()
    request.setLemmaId(lemmaId)
    if (chosenSenseId) request.setChosenSenseId(chosenSenseId)
    request.setSourceForm(sourceForm)
    request.setAccessToken(sessionToken.get() ?? '')
    return vocabularyClient.saveEntry(request, authorizationMetadata())
  })
}

export async function listVocabularyEntries({
  cursor = '',
  limit = 20,
  query = '',
}: {
  cursor?: string
  limit?: number
  query?: string
} = {}) {
  return unaryCall(() => {
    const request = new ListEntriesRequest()
    request.setCursor(cursor)
    request.setLimit(limit)
    request.setQuery(query)
    request.setAccessToken(sessionToken.get() ?? '')
    return vocabularyClient.listEntries(request, authorizationMetadata())
  })
}

export async function deleteVocabularyEntry(entryId: string) {
  return unaryCall(() => {
    const request = new DeleteEntryRequest()
    request.setEntryId(entryId)
    request.setAccessToken(sessionToken.get() ?? '')
    return vocabularyClient.deleteEntry(request, authorizationMetadata())
  })
}

export async function getHighlights(bookId: string, chapterId: string) {
  return unaryCall(() => {
    const request = new GetHighlightsRequest()
    request.setBookId(bookId)
    request.setChapterId(chapterId)
    request.setAccessToken(sessionToken.get() ?? '')
    return vocabularyClient.getHighlights(request, authorizationMetadata())
  })
}
