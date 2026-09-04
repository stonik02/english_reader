import { toApiError } from './errors'
import { grpcCode } from './errors'
import { sessionToken } from './sessionToken'

export async function unaryCall<Response>(call: () => Promise<Response>) {
  try {
    return await call()
  } catch (error) {
    const apiError = toApiError(error)
    if (apiError.code !== grpcCode.unauthenticated) {
      throw apiError
    }

    const token = await sessionToken.renew()
    if (token === null) {
      throw apiError
    }

    try {
      return await call()
    } catch (retryError) {
      throw toApiError(retryError)
    }
  }
}
