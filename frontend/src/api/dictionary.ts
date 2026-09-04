import { LookupWordRequest } from '@reader/proto/reader/v1/dictionary_pb'

import {
  authorizationMetadata,
  dictionaryClient,
  sessionToken,
  unaryCall,
} from './client'

export async function lookupWord({
  bookId,
  chapterId,
  selectedText,
  sentenceText,
  epubCfi,
}: {
  bookId: string
  chapterId: string
  selectedText: string
  sentenceText: string
  epubCfi: string
}) {
  return unaryCall(() => {
    const request = new LookupWordRequest()
    request.setBookId(bookId)
    request.setChapterId(chapterId)
    request.setSelectedText(selectedText)
    request.setSentenceText(sentenceText)
    request.setEpubCfi(epubCfi)
    request.setSourceLanguage('en')
    request.setAccessToken(sessionToken.get() ?? '')
    return dictionaryClient.lookupWord(request, authorizationMetadata())
  })
}
