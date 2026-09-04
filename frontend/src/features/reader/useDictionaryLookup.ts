import { useMutation } from '@tanstack/react-query'

import { lookupWord } from '../../api/dictionary'

export function useDictionaryLookup() {
  return useMutation({ mutationFn: lookupWord })
}
