const defaultGrpcWebUrl = 'http://localhost:8082'
const defaultHttpApiUrl = 'http://localhost:8083'

function readGrpcWebUrl() {
  const value = import.meta.env.VITE_GRPC_WEB_URL ?? defaultGrpcWebUrl

  try {
    return new URL(value).origin
  } catch {
    throw new Error(
      'VITE_GRPC_WEB_URL должен содержать полный HTTP(S)-адрес Envoy.',
    )
  }
}

export const grpcWebUrl = readGrpcWebUrl()

function readHttpApiUrl() {
  const value = import.meta.env.VITE_HTTP_API_URL ?? defaultHttpApiUrl

  try {
    return new URL(value).origin
  } catch {
    throw new Error(
      'VITE_HTTP_API_URL должен содержать полный HTTP(S)-адрес API.',
    )
  }
}

export const httpApiUrl = readHttpApiUrl()
