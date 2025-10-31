import type { WorkspacesResponse } from '@@/types/pocketbase'
import type { RecordModel } from 'pocketbase'

/**
 * Find workspace by slug from a list of workspaces
 */
export function findWorkspaceBySlug(
  workspaces: WorkspacesResponse[],
  slug: string
): WorkspacesResponse | undefined {
  return workspaces.find(workspace => workspace.slug === slug)
}

/**
 * Get default workspace based on last used preference
 */
export function getDefaultWorkspace(
  workspaces: WorkspacesResponse[],
  lastUsedSlug?: string | null
): WorkspacesResponse | undefined {
  if (!workspaces.length) return undefined

  if (lastUsedSlug) {
    const lastWorkspace = findWorkspaceBySlug(workspaces, lastUsedSlug)
    if (lastWorkspace) return lastWorkspace
  }

  return workspaces[0]
}

/**
 * Build workspace dashboard path
 */
export function getWorkspacePath(slug: string, subPath = ''): string {
  const cleanSubPath = subPath.startsWith('/') ? subPath : `/${subPath}`
  return `/${slug}/dashboard${cleanSubPath}`
}

/**
 * Get workspace logo URL
 */
export function getWorkspaceLogoUrl(
  workspace: WorkspacesResponse,
  getFileUrl: (record: RecordModel, field: string) => string | null
): string | null {
  return workspace.logo
    ? getFileUrl(workspace as unknown as RecordModel, 'logo')
    : null
}

/**
 * Get workspace avatar props for components
 */
export function getWorkspaceAvatar(
  workspace: WorkspacesResponse,
  getFileUrl: (record: RecordModel, field: string) => string | null,
  getAvatarUrl: (
    url: string | null | undefined,
    fallback: string | null | undefined
  ) => string
) {
  const logoUrl = getWorkspaceLogoUrl(workspace, getFileUrl)
  return {
    src: getAvatarUrl(logoUrl, workspace.name),
    alt: workspace.name
  }
}
