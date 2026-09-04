import { useEffect, useId, useRef, useState } from 'react'

import { UploadBookError } from '../../api/uploadBook'
import { useUploadBook } from './useCatalogBooks'

const maxFileSize = 50 * 1024 * 1024

function formatFileSize(bytes: number) {
  return `${(bytes / 1024 / 1024).toFixed(bytes < 10 * 1024 * 1024 ? 1 : 0)} МБ`
}

function validationMessage(file: File) {
  if (!file.name.toLowerCase().endsWith('.epub')) {
    return 'Выберите файл с расширением .epub.'
  }
  if (file.size === 0) {
    return 'Этот файл пустой. Выберите EPUB ещё раз.'
  }
  if (file.size > maxFileSize) {
    return 'Размер EPUB не должен превышать 50 МБ.'
  }
  return null
}

type UploadBookDialogProps = {
  onClose(): void
}

export function UploadBookDialog({ onClose }: UploadBookDialogProps) {
  const inputId = useId()
  const titleId = useId()
  const inputRef = useRef<HTMLInputElement>(null)
  const upload = useUploadBook()
  const [file, setFile] = useState<File | null>(null)
  const [validationError, setValidationError] = useState<string | null>(null)
  const [progress, setProgress] = useState(0)
  const [isDragging, setIsDragging] = useState(false)

  const isUploading = upload.isPending
  const isComplete = upload.isSuccess
  const error =
    validationError ??
    (upload.error instanceof UploadBookError
      ? upload.error.message
      : upload.isError
        ? 'Не удалось загрузить книгу. Попробуйте ещё раз.'
        : null)

  useEffect(() => {
    function closeOnEscape(event: KeyboardEvent) {
      if (event.key === 'Escape' && !isUploading) {
        onClose()
      }
    }

    window.addEventListener('keydown', closeOnEscape)
    return () => window.removeEventListener('keydown', closeOnEscape)
  }, [isUploading, onClose])

  function selectFile(candidate: File | undefined) {
    if (candidate === undefined) {
      return
    }

    const message = validationMessage(candidate)
    setValidationError(message)
    setFile(message === null ? candidate : null)
    setProgress(0)
    upload.reset()
  }

  function startUpload() {
    if (file === null) {
      setValidationError('Сначала выберите EPUB-файл.')
      return
    }

    setValidationError(null)
    setProgress(0)
    upload.mutate({ file, onProgress: setProgress })
  }

  return (
    <div
      aria-labelledby={titleId}
      aria-modal="true"
      className="modal-backdrop"
      onMouseDown={(event) => {
        if (event.target === event.currentTarget && !isUploading) {
          onClose()
        }
      }}
      role="dialog"
    >
      <section className="upload-dialog">
        <div className="upload-dialog-header">
          <div>
            <p className="eyebrow">Общий каталог</p>
            <h2 id={titleId}>Загрузить EPUB</h2>
          </div>
          <button
            aria-label="Закрыть окно загрузки"
            className="modal-close"
            disabled={isUploading}
            onClick={onClose}
            type="button"
          >
            ×
          </button>
        </div>

        {isComplete ? (
          <div className="upload-success" role="status">
            <span aria-hidden="true">✓</span>
            <div>
              <h3>Книга загружена</h3>
              <p>Она уже появилась в каталоге и подготавливается к чтению.</p>
            </div>
            <button
              className="button button-primary"
              onClick={onClose}
              type="button"
            >
              Готово
            </button>
          </div>
        ) : (
          <>
            <input
              accept=".epub,application/epub+zip"
              className="visually-hidden"
              id={inputId}
              onChange={(event) => selectFile(event.target.files?.[0])}
              ref={inputRef}
              type="file"
            />
            <button
              className={`upload-dropzone ${isDragging ? 'upload-dropzone-active' : ''}`}
              disabled={isUploading}
              onClick={() => inputRef.current?.click()}
              onDragEnter={(event) => {
                event.preventDefault()
                setIsDragging(true)
              }}
              onDragLeave={(event) => {
                event.preventDefault()
                setIsDragging(false)
              }}
              onDragOver={(event) => event.preventDefault()}
              onDrop={(event) => {
                event.preventDefault()
                setIsDragging(false)
                selectFile(event.dataTransfer.files[0])
              }}
              type="button"
            >
              <span aria-hidden="true" className="upload-icon">
                ↑
              </span>
              <strong>Перетащите EPUB сюда</strong>
              <span>или выберите файл на компьютере</span>
              <small>Только .epub, до 50 МБ</small>
            </button>

            {file !== null && (
              <div className="selected-file">
                <span aria-hidden="true">▤</span>
                <div>
                  <strong>{file.name}</strong>
                  <span>{formatFileSize(file.size)}</span>
                </div>
                {!isUploading && (
                  <button
                    onClick={() => {
                      setFile(null)
                      setValidationError(null)
                      setProgress(0)
                      upload.reset()
                    }}
                    type="button"
                  >
                    Убрать
                  </button>
                )}
              </div>
            )}

            {isUploading && (
              <div aria-live="polite" className="upload-progress">
                <div>
                  <span>Отправляем файл…</span>
                  <strong>{progress}%</strong>
                </div>
                <progress max="100" value={progress} />
              </div>
            )}
            {error !== null && (
              <p className="form-error" role="alert">
                {error}
              </p>
            )}
            <div className="upload-dialog-actions">
              <button
                className="button button-secondary"
                disabled={isUploading}
                onClick={onClose}
                type="button"
              >
                Отмена
              </button>
              <button
                className="button button-primary"
                disabled={file === null || isUploading}
                onClick={startUpload}
                type="button"
              >
                {isUploading ? 'Загружаем…' : 'Загрузить'}
              </button>
            </div>
          </>
        )}
      </section>
    </div>
  )
}
