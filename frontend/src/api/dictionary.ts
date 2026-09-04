import {
  LookupWordRequest,
  TranslateTextRequest,
} from '@reader/proto/reader/v1/dictionary_pb'

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
}: {
  bookId: string
  chapterId: string
  selectedText: string
}) {
  return unaryCall(() => {
    const request = new LookupWordRequest()
    request.setBookId(bookId)
    request.setChapterId(chapterId)
    request.setSelectedText(selectedText)
    request.setSourceLanguage('en')
    request.setAccessToken(sessionToken.get() ?? '')
    return dictionaryClient.lookupWord(request, authorizationMetadata())
  })
}

export async function translateText({
  bookId,
  chapterId,
  text,
}: {
  bookId: string
  chapterId: string
  text: string
}) {
  return unaryCall(() => {
    const request = new TranslateTextRequest()
    request.setBookId(bookId)
    request.setChapterId(chapterId)
    request.setText(text)
    request.setAccessToken(sessionToken.get() ?? '')
    return dictionaryClient.translateText(request, authorizationMetadata())
  })
}
