const wordExpression = /[\p{L}]+(?:['’][\p{L}]+)*/gu

export function decorateInteractiveWords(root: HTMLElement) {
  for (const node of readableTextNodes(root)) {
    wrapMatches(
      node,
      [...(node.textContent ?? '').matchAll(wordExpression)].reverse(),
    )
  }
}

// Decorating every word is intentionally deferred and split into small units.
// This restores the word-hover affordance without delaying the first paint of a
// chapter, including EPUBs that contain thousands of words in one document.
export function decorateInteractiveWordsInBatches(
  root: HTMLElement,
  wordsPerBatch = 40,
) {
  const work = readableTextNodes(root).map((node) => ({
    node,
    matches: [...(node.textContent ?? '').matchAll(wordExpression)],
  }))
  let taskIndex = 0
  let cancelled = false
  let timer: number | undefined

  const process = () => {
    if (cancelled) return
    let remaining = wordsPerBatch
    while (remaining > 0 && taskIndex < work.length) {
      const task = work[taskIndex]
      const match = task.matches.pop()
      if (!match) {
        taskIndex += 1
        continue
      }
      wrapMatches(task.node, [match])
      remaining -= 1
    }
    if (taskIndex < work.length) timer = window.setTimeout(process, 16)
  }

  timer = window.setTimeout(process, 0)
  return () => {
    cancelled = true
    if (timer !== undefined) window.clearTimeout(timer)
  }
}

function readableTextNodes(root: HTMLElement) {
  const walker = document.createTreeWalker(root, NodeFilter.SHOW_TEXT)
  const nodes: Text[] = []
  while (walker.nextNode()) {
    const node = walker.currentNode as Text
    if (
      node.parentElement?.closest(
        'script, style, pre, code, [data-reader-word]',
      )
    ) {
      continue
    }
    nodes.push(node)
  }
  return nodes
}

function wrapMatches(node: Text, matches: RegExpMatchArray[]) {
  for (const match of matches) {
    if (
      match.index === undefined ||
      match.index + match[0].length > node.length
    ) {
      continue
    }
    const range = document.createRange()
    range.setStart(node, match.index)
    range.setEnd(node, match.index + match[0].length)
    const word = document.createElement('span')
    word.className = 'reader-word'
    word.dataset.readerWord = match[0]
    range.surroundContents(word)
  }
}
