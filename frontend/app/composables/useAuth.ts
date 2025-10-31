import type { RecordModel } from 'pocketbase'
import { parsePocketBaseError } from '~/utils/pocketbaseErrorHandler'
import { getDefaultWorkspace, getWorkspacePath } from '~/utils/workspace'

export const useAuth = () => {
  const { pb, user } = usePocketbase()
  const toast = useToast()

  const isLoggedIn = computed(() => pb.authStore.isValid)
  const { getFileUrl } = usePocketbaseFiles()

  const avatarUrl = computed(() =>
    user.value ? getFileUrl(user.value as RecordModel, 'avatar') : null
  )

  const token = computed(() => pb.authStore.token)

  const signOut = async () => {
    await pb.authStore.clear()
    return navigateTo('/')
  }

  const login = async (email: string, password: string) => {
    try {
      await pb.collection('users').authWithPassword(email, password)
      await redirectAfterLogin()
      return true
    } catch (error) {
      console.error(error)
      const parsedError = parsePocketBaseError(error, 'login')
      toast.add({
        title: parsedError.title,
        description: parsedError.message,
        color: 'error'
      })
      return false
    }
  }

  const register = async (email: string, password: string, name: string) => {
    try {
      await pb.collection('users').create({
        email,
        password,
        passwordConfirm: password,
        name
      })

      await pb.collection('users').requestVerification(email)
      toast.add({
        title: 'Registration successful',
        description:
          'An email was sent with a verification link for account verification',
        color: 'success'
      })
      await navigateTo('/')
      return true
    } catch (error) {
      console.error(error)
      const parsedError = parsePocketBaseError(error, 'register')
      toast.add({
        title: parsedError.title,
        description: parsedError.message,
        color: 'error'
      })
      return false
    }
  }

  const oauthLogin = async (provider: 'github' | 'google') => {
    try {
      await pb.collection('users').authWithOAuth2({ provider })
      await new Promise(resolve => setTimeout(resolve, 300))
      await redirectAfterLogin()
      return true
    } catch (error) {
      console.error('OAuth login error:', error)
      const parsedError = parsePocketBaseError(error, 'login')
      toast.add({
        title: parsedError.title,
        description: parsedError.message,
        color: 'error'
      })
      return false
    }
  }

  const redirectAfterLogin = async () => {
    const { fetchWorkspaces, workspaces } = useWorkspaces()
    const { getLastUsedWorkspace } = useWorkspacePreferences()

    await fetchWorkspaces()

    if (workspaces.value.length === 0) {
      return navigateTo('/new-workspace')
    }

    const defaultWorkspace = getDefaultWorkspace(
      workspaces.value,
      getLastUsedWorkspace()
    )

    if (defaultWorkspace) {
      return navigateTo(getWorkspacePath(defaultWorkspace.slug))
    }

    return navigateTo('/new-workspace')
  }

  return {
    user,
    isLoggedIn,
    avatarUrl,
    token,
    signOut,
    login,
    register,
    oauthLogin,
    redirectAfterLogin
  }
}
