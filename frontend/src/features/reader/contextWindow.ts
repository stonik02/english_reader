const wordExpression = /[\p{L}]+(?:['’][\p{L}]+)*/gu

// Returns at most three words on each side of the selected word.
export function contextWindow(
  selectedWord: string,
  text: string | null | undefined,
) {
  const words = (text ?? selectedWord).match(wordExpression) ?? []
  const selected = normalize(selectedWord)
  const index = words.findIndex((word) => normalize(word) === selected)
  if (index < 0) return selectedWord
  return words.slice(Math.max(0, index - 3), index + 4).join(' ')
}

function normalize(value: string) {
  return value.trim().toLocaleLowerCase().replaceAll('’', "'")
}
