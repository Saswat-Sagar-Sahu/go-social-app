const TOKEN_KEY = 'go_social_auth_token'

export function saveToken(token) {
  localStorage.setItem(TOKEN_KEY, token)
}

export function getToken() {
  return localStorage.getItem(TOKEN_KEY)
}

export function clearToken() {
  localStorage.removeItem(TOKEN_KEY)
}

export function isAuthenticated() {
  return Boolean(getToken())
}

function decodeBase64Url(value) {
  const normalized = value.replace(/-/g, '+').replace(/_/g, '/')
  const padding = '='.repeat((4 - (normalized.length % 4)) % 4)
  return atob(normalized + padding)
}

export function getUserIdFromToken() {
  const token = getToken()
  if (!token) {
    return null
  }

  try {
    const [, payload] = token.split('.')
    if (!payload) {
      return null
    }
    const decoded = decodeBase64Url(payload)
    const claims = JSON.parse(decoded)
    const userId = Number(claims.sub)
    return Number.isFinite(userId) ? userId : null
  } catch {
    return null
  }
}
