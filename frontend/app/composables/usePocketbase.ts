import type { TypedPocketBase } from '@@/types/pocketbase'
import type { RecordModel } from 'pocketbase'

// Global reactive user state using useState
const useGlobalUser = () => useState<RecordModel | null>('pb-user', () => null)
const useIsInitialized = () => useState<boolean>('pb-initialized', () => false)

export const usePocketbase = () => {
  const { $pb, $buildApiUrl } = useNuxtApp()
  const pb = $pb as TypedPocketBase
  const buildApiUrl = $buildApiUrl as (path: string) => string

  const globalUser = useGlobalUser()
  const isInitialized = useIsInitialized()

  // Initialize global user state once
  if (!isInitialized.value) {
    globalUser.value = pb.authStore.record

    // Listen for auth store changes and update global state
    pb.authStore.onChange(() => {
      globalUser.value = pb.authStore.record
    })

    isInitialized.value = true
  }

  const user = computed(() => globalUser.value)

  const refreshUser = async () => {
    try {
      await pb.collection('users').authRefresh()
    } catch (error) {
      console.error('Failed to refresh user:', error)
      throw error
    }
  }

  return {
    pb,
    user,
    buildApiUrl,
    refreshUser
  }
}
