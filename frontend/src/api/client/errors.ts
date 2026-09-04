import type { RpcError } from 'grpc-web'

const grpcStatus = {
  cancelled: 1,
  unknown: 2,
  invalidArgument: 3,
  deadlineExceeded: 4,
  notFound: 5,
  alreadyExists: 6,
  permissionDenied: 7,
  resourceExhausted: 8,
  failedPrecondition: 9,
  aborted: 10,
  outOfRange: 11,
  unimplemented: 12,
  internal: 13,
  unavailable: 14,
  dataLoss: 15,
  unauthenticated: 16,
} as const

export type ApiErrorCode = (typeof grpcStatus)[keyof typeof grpcStatus]

export class ApiError extends Error {
  readonly code: ApiErrorCode | number
  readonly retryable: boolean

  constructor(code: ApiErrorCode | number, message: string) {
    super(message)
    this.name = 'ApiError'
    this.code = code
    this.retryable =
      code === grpcStatus.unavailable || code === grpcStatus.deadlineExceeded
  }
}

export function toApiError(error: unknown) {
  if (error instanceof ApiError) {
    return error
  }

  const rpcError = error as Partial<RpcError>
  const code =
    typeof rpcError.code === 'number' ? rpcError.code : grpcStatus.unknown

  switch (code) {
    case grpcStatus.unauthenticated:
      return new ApiError(code, 'Сессия закончилась. Войдите снова.')
    case grpcStatus.permissionDenied:
      return new ApiError(code, 'У вас нет доступа к этому действию.')
    case grpcStatus.notFound:
      return new ApiError(code, 'Запрошенные данные не найдены.')
    case grpcStatus.invalidArgument:
      return new ApiError(
        code,
        'Проверьте введённые данные и повторите попытку.',
      )
    case grpcStatus.failedPrecondition:
      return new ApiError(code, 'Это действие пока недоступно.')
    case grpcStatus.resourceExhausted:
      return new ApiError(
        code,
        'Слишком много запросов. Попробуйте немного позже.',
      )
    case grpcStatus.unavailable:
      return new ApiError(
        code,
        'Сервер временно недоступен. Попробуйте ещё раз.',
      )
    default:
      return new ApiError(
        code,
        rpcError.message || 'Не удалось выполнить запрос к серверу.',
      )
  }
}

export const grpcCode = grpcStatus
