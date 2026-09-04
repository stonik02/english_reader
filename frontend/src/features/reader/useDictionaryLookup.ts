import { useMutation } from '@tanstack/react-query'

import { lookupWord, translateText } from '../../api/dictionary'

export function useDictionaryLookup() {
  return useMutation({ mutationFn: lookupWord })
}

export function useTextTranslation() {
  return useMutation({ mutationFn: translateText })
}
