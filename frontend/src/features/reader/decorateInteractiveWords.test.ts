import { afterEach, describe, expect, it, vi } from 'vitest'

import {
  decorateInteractiveWords,
  decorateInteractiveWordsInBatches,
} from './decorateInteractiveWords'

afterEach(() => vi.useRealTimers())

describe('decorateInteractiveWords', () => {
  it('wraps readable words and leaves code untouched', () => {
    const root = document.createElement('article')
    root.innerHTML = '<p>Hello, world!</p><code>const value = 1</code>'

    decorateInteractiveWords(root)

    expect(
      [...root.querySelectorAll<HTMLElement>('[data-reader-word]')].map(
        (word) => word.dataset.readerWord,
      ),
    ).toEqual(['Hello', 'world'])
    expect(root.querySelector('code')?.textContent).toBe('const value = 1')
  })

  it('does not wrap a word twice', () => {
    const root = document.createElement('article')
    root.textContent = 'Putting'

    decorateInteractiveWords(root)
    decorateInteractiveWords(root)

    expect(root.querySelectorAll('[data-reader-word]')).toHaveLength(1)
  })

  it('adds wrappers in deferred batches', () => {
    vi.useFakeTimers()
    const root = document.createElement('article')
    root.textContent = 'One two three'

    decorateInteractiveWordsInBatches(root, 1)
    expect(root.querySelectorAll('[data-reader-word]')).toHaveLength(0)

    vi.runAllTimers()
    expect(root.querySelectorAll('[data-reader-word]')).toHaveLength(3)
  })
})
