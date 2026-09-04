import { DictionaryServiceClient } from '@reader/proto/reader/v1/DictionaryServiceClientPb'

import { grpcWebUrl } from './config'

export const dictionaryClient = new DictionaryServiceClient(grpcWebUrl, null, {
  withCredentials: true,
})
