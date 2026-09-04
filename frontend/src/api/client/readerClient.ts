import { ReaderServiceClient } from '@reader/proto/reader/v1/ReaderServiceClientPb'

import { grpcWebUrl } from './config'

export const readerClient = new ReaderServiceClient(grpcWebUrl, null, {
  withCredentials: true,
})
