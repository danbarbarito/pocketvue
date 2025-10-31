export const useWorkspacePreferences = () => {
  const lastWorkspaceSlug = useCookie('lastWorkspaceSlug', {
    maxAge: 60 * 60 * 24 * 365,
    path: '/'
  })

  const setLastUsedWorkspace = (slug: string) => {
    lastWorkspaceSlug.value = slug
  }

  const getLastUsedWorkspace = () => {
    return lastWorkspaceSlug.value
  }

  return {
    setLastUsedWorkspace,
    getLastUsedWorkspace
  }
}
