import { describe, expect, it } from 'vitest'

import { preferredVoice, splitSpeechText } from './speechSynthesis'

describe('speech synthesis helpers', () => {
  it('normalizes and splits long text at sentence boundaries', () => {
    expect(splitSpeechText(' One.   Two! Three? ', 10)).toEqual([
      'One. Two!',
      'Three?',
    ])
  })

  it('splits a sentence that is longer than the limit at word boundaries', () => {
    expect(splitSpeechText('one two three four', 7)).toEqual([
      'one two',
      'three',
      'four',
    ])
  })

  it('prefers the saved voice then en-US then any English voice', () => {
    const british = {
      lang: 'en-GB',
      voiceURI: 'british',
    } as SpeechSynthesisVoice
    const american = {
      lang: 'en-US',
      voiceURI: 'american',
    } as SpeechSynthesisVoice

    expect(preferredVoice([british, american], 'british')).toBe(british)
    expect(preferredVoice([british, american], '')).toBe(american)
  })
})
