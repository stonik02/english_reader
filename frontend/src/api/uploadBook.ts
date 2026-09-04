import { httpApiUrl, sessionToken } from './client'

type UploadedBook = {
  id: string
  title: string
  author: string
  status: string
  uploaded_by_user_id: string
  created_at: string
}

type ErrorResponse = { error?: string }

export class UploadBookError extends Error {
  constructor(
    message: string,
    readonly status: number,
  ) {
    super(message)
    this.name = 'UploadBookError'
  }
}

function messageForStatus(status: number, response: ErrorResponse | null) {
  if (status === 401) {
    return 'Сессия закончилась. Войдите снова.'
  }
  if (status === 413) {
    return 'Файл больше допустимого размера — 50 МБ.'
  }
  if (status === 400) {
    return 'Не удалось прочитать EPUB. Выберите другой файл.'
  }
  return response?.error ?? 'Не удалось загрузить книгу. Попробуйте ещё раз.'
}

function sendFile(
  file: File,
  accessToken: string,
  onProgress: (percent: number) => void,
) {
  return new Promise<UploadedBook>((resolve, reject) => {
    const request = new XMLHttpRequest()
    const body = new FormData()
    body.append('file', file, file.name)

    request.open('POST', `${httpApiUrl}/api/v1/library/books`)
    request.withCredentials = true
    request.setRequestHeader('Authorization', `Bearer ${accessToken}`)
    request.upload.addEventListener('progress', (event) => {
      if (event.lengthComputable) {
        onProgress(Math.round((event.loaded / event.total) * 100))
      }
    })
    request.addEventListener('error', () => {
      reject(
        new UploadBookError(
          'Не удалось связаться с сервером. Проверьте, что backend запущен.',
          0,
        ),
      )
    })
    request.addEventListener('load', () => {
      const response = request.responseText
        ? (JSON.parse(request.responseText) as UploadedBook | ErrorResponse)
        : null

      if (request.status >= 200 && request.status < 300) {
        resolve(response as UploadedBook)
        return
      }

      reject(
        new UploadBookError(
          messageForStatus(request.status, response as ErrorResponse | null),
          request.status,
        ),
      )
    })
    request.send(body)
  })
}

export async function uploadBook(
  file: File,
  onProgress: (percent: number) => void,
) {
  const accessToken = sessionToken.get() ?? (await sessionToken.renew())
  if (accessToken === null) {
    throw new UploadBookError('Войдите, чтобы загрузить книгу.', 401)
  }

  try {
    return await sendFile(file, accessToken, onProgress)
  } catch (error) {
    if (!(error instanceof UploadBookError) || error.status !== 401) {
      throw error
    }

    const renewedAccessToken = await sessionToken.renew()
    if (renewedAccessToken === null) {
      throw error
    }
    return sendFile(file, renewedAccessToken, onProgress)
  }
}
