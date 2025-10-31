import {
  findWorkspaceBySlug,
  getDefaultWorkspace,
  getWorkspacePath
} from '~/utils/workspace'

export default defineNuxtRouteMiddleware(async to => {
  // Handle root path redirect for logged-in users
  if (to.path === '/' && !to.params.workspaceSlug) {
    const { isLoggedIn } = useAuth()
    if (isLoggedIn.value) {
      const { workspaces, fetchWorkspaces } = useWorkspaces()
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
    return
  }

  // Only validate workspace routes
  if (!to.params.workspaceSlug) {
    return
  }

  const { workspaces, fetchWorkspaces } = useWorkspaces()

  // Ensure workspaces are loaded
  if (workspaces.value.length === 0) {
    await fetchWorkspaces()
  }

  // No workspaces available
  if (workspaces.value.length === 0) {
    return navigateTo('/new-workspace')
  }

  // Validate workspace exists
  let workspace = findWorkspaceBySlug(
    workspaces.value,
    to.params.workspaceSlug as string
  )

  // Try refreshing if not found (stale cache)
  if (!workspace) {
    await fetchWorkspaces()
    workspace = findWorkspaceBySlug(
      workspaces.value,
      to.params.workspaceSlug as string
    )
  }

  if (!workspace) {
    throw createError({
      statusCode: 404,
      statusMessage: 'Workspace not found'
    })
  }
})
