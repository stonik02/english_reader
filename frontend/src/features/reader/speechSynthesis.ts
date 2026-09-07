const maximumChunkLength = 1_200

export type SpeechPlaybackOptions = {
  onEnd: () => void
  onError: () => void
  onStart: () => void
  rate: number
  voiceURI: string
}

export function speechSynthesisSupported() {
  return typeof window !== 'undefined' && 'speechSynthesis' in window
}

export function englishVoices() {
  if (!speechSynthesisSupported()) return []
  return window.speechSynthesis
    .getVoices()
    .filter((voice) => voice.lang.toLowerCase().startsWith('en'))
}

export function subscribeEnglishVoices(
  onChange: (voices: SpeechSynthesisVoice[]) => void,
) {
  if (!speechSynthesisSupported()) return () => {}
  const update = () => onChange(englishVoices())
  update()
  window.speechSynthesis.addEventListener('voiceschanged', update)
  return () =>
    window.speechSynthesis.removeEventListener('voiceschanged', update)
}

export function preferredVoice(
  voices: SpeechSynthesisVoice[],
  voiceURI: string,
) {
  return (
    voices.find((voice) => voice.voiceURI === voiceURI) ??
    voices.find((voice) => voice.lang.toLowerCase() === 'en-us') ??
    voices[0]
  )
}

export function splitSpeechText(value: string, limit = maximumChunkLength) {
  const text = value.replace(/\s+/g, ' ').trim()
  if (!text) return []
  const sentences = text.match(/[^.!?]+(?:[.!?]+|$)/g) ?? [text]
  const chunks: string[] = []
  let current = ''
  for (const sentence of sentences) {
    const normalized = sentence.trim()
    if (!normalized) continue
    if (normalized.length > limit) {
      if (current) chunks.push(current)
      current = ''
      chunks.push(...splitLongText(normalized, limit))
      continue
    }
    const next = current ? `${current} ${normalized}` : normalized
    if (next.length > limit) {
      chunks.push(current)
      current = normalized
    } else {
      current = next
    }
  }
  if (current) chunks.push(current)
  return chunks
}

export function startSpeech(text: string, options: SpeechPlaybackOptions) {
  if (!speechSynthesisSupported()) return null
  const voice = preferredVoice(englishVoices(), options.voiceURI)
  if (!voice) return null
  const chunks = splitSpeechText(text)
  if (chunks.length === 0) return null

  const synth = window.speechSynthesis
  let active = true
  let index = 0
  synth.cancel()
  options.onStart()

  const speakNext = () => {
    if (!active) return
    const utterance = new SpeechSynthesisUtterance(chunks[index])
    utterance.lang = 'en-US'
    utterance.voice = voice
    utterance.volume = 1
    utterance.rate = options.rate
    utterance.pitch = 1
    utterance.onend = () => {
      if (!active) return
      index += 1
      if (index < chunks.length) {
        speakNext()
      } else {
        active = false
        options.onEnd()
      }
    }
    utterance.onerror = () => {
      if (!active) return
      active = false
      options.onError()
    }
    synth.speak(utterance)
  }
  speakNext()

  return () => {
    if (!active) return
    active = false
    synth.cancel()
    options.onEnd()
  }
}

function splitLongText(text: string, limit: number) {
  const chunks: string[] = []
  let current = ''
  for (const word of text.split(' ')) {
    const next = current ? `${current} ${word}` : word
    if (next.length > limit && current) {
      chunks.push(current)
      current = word
    } else {
      current = next
    }
  }
  if (current) chunks.push(current)
  return chunks
}
