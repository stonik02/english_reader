const wordExpression = /[\p{L}]+(?:['’][\p{L}]+)*/gu

export function decorateInteractiveWords(root: HTMLElement) {
  const walker = document.createTreeWalker(root, NodeFilter.SHOW_TEXT)
  const nodes: Text[] = []
  while (walker.nextNode()) {
    const node = walker.currentNode as Text
    if (node.parentElement?.closest('script, style, pre, code, [data-reader-word]')) {
      continue
    }
    nodes.push(node)
  }

  for (const node of nodes) {
    const matches = [...(node.textContent ?? '').matchAll(wordExpression)]
    for (const match of matches.reverse()) {
      if (match.index === undefined) continue
      const range = document.createRange()
      range.setStart(node, match.index)
      range.setEnd(node, match.index + match[0].length)
      const word = document.createElement('span')
      word.className = 'reader-word'
      word.dataset.readerWord = match[0]
      range.surroundContents(word)
    }
  }
}
