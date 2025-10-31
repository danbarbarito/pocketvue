import type { NavigationMenuItem } from '@nuxt/ui'
import { getWorkspacePath } from '~/utils/workspace'

export const useNavigationLinks = () => {
  const { activeWorkspaceSlug } = useWorkspaces()

  const navigationLinks = useState<NavigationMenuItem[][]>(
    'navigationLinks',
    () => []
  )

  const initializeLinks = () => {
    if (!activeWorkspaceSlug.value) return

    navigationLinks.value = [
      [
        {
          label: 'Home',
          icon: 'i-custom-home',
          to: getWorkspacePath(activeWorkspaceSlug.value),
          exact: true
        },
        {
          label: 'Notes',
          icon: 'i-lucide-notebook-pen',
          to: getWorkspacePath(activeWorkspaceSlug.value, '/notes')
        },
        {
          label: 'Settings',
          icon: 'i-lucide-folder-cog',
          to: getWorkspacePath(activeWorkspaceSlug.value, '/settings')
        }
      ],
      [
        {
          label: 'Got Feedback?',
          icon: 'i-lucide-send'
        },
        {
          label: 'Need Help?',
          icon: 'i-lucide-life-buoy'
        }
      ]
    ]
  }

  watch(
    activeWorkspaceSlug,
    () => {
      initializeLinks()
    },
    { immediate: true }
  )

  return {
    navigationLinks,
    initializeLinks
  }
}
