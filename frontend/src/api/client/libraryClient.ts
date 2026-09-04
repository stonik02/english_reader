import { LibraryServiceClient } from '@reader/proto/reader/v1/LibraryServiceClientPb'
import { grpcWebUrl } from './config'

export const libraryClient = new LibraryServiceClient(grpcWebUrl, null, {
  withCredentials: true,
})
