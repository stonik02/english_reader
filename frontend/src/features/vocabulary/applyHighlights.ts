type Highlight = {
  texts: string[]
}

export function applyHighlights(root: HTMLElement, highlights: Highlight[]) {
  root.querySelectorAll('mark.vocabulary-highlight').forEach((mark) => {
    mark.replaceWith(...mark.childNodes)
  })
  const forms = [...new Set(highlights.flatMap((highlight) => highlight.texts))]
    .filter(Boolean)
    .sort((left, right) => right.length - left.length)
  if (forms.length === 0) return

  const expression = new RegExp(
    `(?<![\\p{L}'’])(?:${forms.map(escapeRegExp).join('|')})(?![\\p{L}'’])`,
    'giu',
  )
  const walker = document.createTreeWalker(root, NodeFilter.SHOW_TEXT)
  const nodes: Text[] = []
  while (walker.nextNode()) {
    nodes.push(walker.currentNode as Text)
  }
  for (const node of nodes) {
    const text = node.textContent ?? ''
    const matches = [...text.matchAll(expression)]
    for (const match of matches.reverse()) {
      if (match.index === undefined) continue
      const range = document.createRange()
      range.setStart(node, match.index)
      range.setEnd(node, match.index + match[0].length)
      const mark = document.createElement('mark')
      mark.className = 'vocabulary-highlight'
      range.surroundContents(mark)
    }
  }
}

function escapeRegExp(value: string) {
  return value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}
