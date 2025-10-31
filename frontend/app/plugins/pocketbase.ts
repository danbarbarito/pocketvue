import type { TypedPocketBase } from '@@/types/pocketbase'
import PocketBase from 'pocketbase'

export default defineNuxtPlugin(nuxtApp => {
  const pocketbaseUrl = import.meta.dev ? 'http://localhost:8090' : '/'
  const pb = new PocketBase(pocketbaseUrl as string) as TypedPocketBase

  // Helper to build API URLs correctly in both dev and prod
  const buildApiUrl = (path: string) => {
    const cleanPath = path.startsWith('/') ? path.slice(1) : path
    const base = pb.baseURL.endsWith('/') ? pb.baseURL.slice(0, -1) : pb.baseURL
    return `${base}/${cleanPath}`
  }

  nuxtApp.provide('pb', pb)
  nuxtApp.provide('buildApiUrl', buildApiUrl)
})
