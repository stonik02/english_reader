import {
  GetAdjacentChapterRequest,
  GetReadingStateRequest,
  GetReaderSettingsRequest,
  SaveReadingProgressRequest,
  UpdateReaderSettingsRequest,
} from '@reader/proto/reader/v1/reader_pb'

import {
  authorizationMetadata,
  readerClient,
  sessionToken,
  unaryCall,
} from './client'

export async function getReadingState(bookId: string) {
  return unaryCall(() => {
    const request = new GetReadingStateRequest()
    request.setBookId(bookId)
    request.setAccessToken(sessionToken.get() ?? '')
    return readerClient.getReadingState(request, authorizationMetadata())
  })
}

export async function saveReadingProgress({
  bookId,
  chapterId,
  epubCfi,
  progressPercent,
  revision,
}: {
  bookId: string
  chapterId: string
  epubCfi: string
  progressPercent: number
  revision: number
}) {
  return unaryCall(() => {
    const request = new SaveReadingProgressRequest()
    request.setBookId(bookId)
    request.setChapterId(chapterId)
    request.setEpubCfi(epubCfi)
    request.setProgressPercent(progressPercent)
    request.setRevision(revision)
    request.setAccessToken(sessionToken.get() ?? '')
    return readerClient.saveReadingProgress(request, authorizationMetadata())
  })
}

export async function getAdjacentChapter(
  bookId: string,
  chapterId: string,
  direction: -1 | 1,
) {
  return unaryCall(() => {
    const request = new GetAdjacentChapterRequest()
    request.setBookId(bookId)
    request.setChapterId(chapterId)
    request.setDirection(direction)
    request.setAccessToken(sessionToken.get() ?? '')
    return readerClient.getAdjacentChapter(request, authorizationMetadata())
  })
}

export async function getReaderSettings() {
  return unaryCall(() => {
    const request = new GetReaderSettingsRequest()
    request.setAccessToken(sessionToken.get() ?? '')
    return readerClient.getReaderSettings(request, authorizationMetadata())
  })
}

export async function updateReaderSettings({
  fontScale,
  theme,
  lineHeight,
  highlightColor,
}: {
  fontScale: number
  theme: string
  lineHeight: number
  highlightColor: string
}) {
  return unaryCall(() => {
    const request = new UpdateReaderSettingsRequest()
    request.setFontScale(fontScale)
    request.setTheme(theme)
    request.setLineHeight(lineHeight)
    if (typeof request.setHighlightColor === 'function') {
      request.setHighlightColor(highlightColor)
    }
    request.setAccessToken(sessionToken.get() ?? '')
    return readerClient.updateReaderSettings(request, authorizationMetadata())
  })
}
