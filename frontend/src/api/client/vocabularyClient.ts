import { VocabularyServiceClient } from '@reader/proto/reader/v1/VocabularyServiceClientPb'

import { grpcWebUrl } from './config'

export const vocabularyClient = new VocabularyServiceClient(grpcWebUrl, null, {
  withCredentials: true,
})
