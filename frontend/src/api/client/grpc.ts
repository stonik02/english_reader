import type { Metadata } from 'grpc-web'

import { sessionToken } from './sessionToken'

export function authorizationMetadata(): Metadata {
  const token = sessionToken.get()

  return token === null ? {} : { authorization: `Bearer ${token}` }
}
