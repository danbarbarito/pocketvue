import type { WorkspacesResponse } from '@@/types/pocketbase'

export const useWorkspaces = () => {
  const workspaces = useState<WorkspacesResponse[]>('workspaces', () => [])
  const loadingWorkspaces = ref(false)
  const { pb } = usePocketbase()
  const toast = useToast()
  const router = useRouter()
  const activeWorkspaceSlug = computed(
    () => router.currentRoute.value.params.workspaceSlug as string
  )

  const activeWorkspace = computed(() => {
    if (!activeWorkspaceSlug.value || !workspaces.value.length) {
      return workspaces.value[0] || ({} as WorkspacesResponse)
    }

    const workspace = workspaces.value.find(
      workspace => workspace.slug === activeWorkspaceSlug.value
    )

    return workspace || ({} as WorkspacesResponse)
  })

  const fetchWorkspaces = async () => {
    if (loadingWorkspaces.value) return

    loadingWorkspaces.value = true
    try {
      workspaces.value = await pb
        .collection('workspaces')
        .getFullList<WorkspacesResponse>({
          sort: '-created'
        })
    } catch (error) {
      console.error('Error fetching workspaces:', error)
      toast.add({
        title: 'Error fetching workspaces',
        description: 'Please try again later',
        icon: 'i-lucide-alert-circle',
        color: 'error'
      })
    } finally {
      loadingWorkspaces.value = false
    }
  }

  return {
    workspaces,
    loadingWorkspaces,
    activeWorkspace,
    activeWorkspaceSlug,
    fetchWorkspaces
  }
}
