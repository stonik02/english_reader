import { describe, expect, it } from 'vitest'

import { contextWindow } from './contextWindow'

describe('contextWindow', () => {
  it('keeps three words on each side of the selected word', () => {
    expect(
      contextWindow('hummed', 'One two three hummed four five six seven.'),
    ).toBe('One two three hummed four five six')
  })

  it('shortens the context at a text boundary', () => {
    expect(contextWindow('Putting', 'Putting it all together now.')).toBe(
      'Putting it all together',
    )
  })
})
