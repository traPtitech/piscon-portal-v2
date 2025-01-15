import createClient from 'openapi-fetch'
import type { paths } from '@/api/openapi'

export const apiBaseUrl = '/api'

export const api = createClient<paths>({
  baseUrl: apiBaseUrl,
})
