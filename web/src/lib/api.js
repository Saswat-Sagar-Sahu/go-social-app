import { getToken } from './auth'

const API_BASE = import.meta.env.VITE_API_BASE || ''

async function parseResponse(response) {
  const text = await response.text()

  if (!text) {
    return null
  }

  try {
    return JSON.parse(text)
  } catch {
    return null
  }
}

export async function apiRequest(path, options = {}) {
  const headers = {
    'Content-Type': 'application/json',
    ...(options.headers || {}),
  }

  if (options.auth) {
    const token = getToken()
    if (token) {
      headers.Authorization = `Bearer ${token}`
    }
  }

  const response = await fetch(`${API_BASE}${path}`, {
    method: options.method || 'GET',
    headers,
    body: options.body ? JSON.stringify(options.body) : undefined,
  })

  const data = await parseResponse(response)

  if (!response.ok) {
    const message = data?.error || `Request failed with status ${response.status}`
    const error = new Error(message)
    error.status = response.status
    throw error
  }

  return data
}
