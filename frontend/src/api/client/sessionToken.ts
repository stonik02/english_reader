let accessToken: string | null = null
let renewToken: (() => Promise<string | null>) | null = null
let pendingRenewal: Promise<string | null> | null = null

export const sessionToken = {
  clear() {
    accessToken = null
  },
  get() {
    return accessToken
  },
  set(token: string) {
    accessToken = token
  },
  setRenewal(handler: (() => Promise<string | null>) | null) {
    renewToken = handler
  },
  async renew() {
    if (renewToken === null) {
      return null
    }
    if (pendingRenewal === null) {
      pendingRenewal = renewToken().finally(() => {
        pendingRenewal = null
      })
    }
    return pendingRenewal
  },
}
