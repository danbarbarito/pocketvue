<template>
  <div>
    <UCard>
      <template #header>
        <p class="font-medium">Delete Workspace</p>
      </template>
      <p class="text-muted text-sm">
        Deleting your workspace will remove the workspace permenantly and all
        data associated with it. This action is irreversible.
      </p>
      <ul class="text-muted mt-2 list-inside list-disc space-y-1 text-sm">
        <li>Labels</li>
        <li>Posts</li>
      </ul>

      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <UButton
            icon="i-lucide-trash-2"
            label="Delete Workspace"
            color="error"
            size="lg"
            @click="deleteModalConfirmation = true"
          />
        </div>
      </template>
    </UCard>
    <UModal
      v-model:open="deleteModalConfirmation"
      title="Warning, destructive action"
    >
      <template #body>
        <div
          class="mb-4 rounded-md bg-red-100 p-4 text-center font-medium
            text-balance text-red-800 dark:bg-red-950 dark:text-red-100"
        >
          This action will delete everything associated with this workspace
          immidiately and this action is irreversible.
        </div>
        <p class="text-toned">This includes</p>
        <ul class="text-muted mt-2 list-inside list-disc space-y-1">
          <li>Labels</li>
          <li>Posts</li>
        </ul>
        <UFormField
          :label="`Type 'Delete ${activeWorkspace.name}' to proceed`"
          class="mt-5"
        >
          <UInput v-model="deleteConfirmation" class="w-full" size="lg" />
        </UFormField>
      </template>
      <template #footer>
        <div class="flex w-full items-center justify-end gap-2">
          <UButton
            label="Cancel"
            variant="ghost"
            @click="deleteModalConfirmation = false"
          />
          <UButton
            label="Delete this workspace permenantly"
            color="error"
            :loading="loading"
            :disabled="!deleteConfirmationCheck"
            @click="deleteWorkspace"
          />
        </div>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
const { activeWorkspace, workspaces, fetchWorkspaces } = useWorkspaces()
const deleteModalConfirmation = ref(false)
const loading = ref(false)
const { pb } = usePocketbase()
const router = useRouter()
const toast = useToast()
const { setLastUsedWorkspace } = useWorkspacePreferences()
const deleteConfirmation = ref('')

const deleteConfirmationCheck = computed(() => {
  return deleteConfirmation.value === `Delete ${activeWorkspace.value.name}`
})

const deleteWorkspace = async () => {
  loading.value = true
  try {
    await pb.collection('workspaces').delete(activeWorkspace.value.id)

    // Refetch workspaces to get the updated list
    await fetchWorkspaces()

    // Check if there are any workspaces left
    if (workspaces.value.length > 0) {
      // Redirect to the first workspace
      const firstWorkspace = workspaces.value[0]
      if (firstWorkspace && firstWorkspace.slug) {
        setLastUsedWorkspace(firstWorkspace.slug)
        await router.push(`/${firstWorkspace.slug}/dashboard`)
        toast.add({
          title: 'Workspace deleted',
          description: `Redirected to ${firstWorkspace.name || 'your workspace'}`,
          icon: 'i-lucide-check',
          color: 'success'
        })
      } else {
        // Fallback if workspace data is incomplete
        await router.push('/')
        toast.add({
          title: 'Workspace deleted',
          description: 'Redirected to homepage',
          icon: 'i-lucide-check',
          color: 'success'
        })
      }
    } else {
      // No workspaces left, redirect to new-workspace
      await router.push('/new-workspace')
      toast.add({
        title: 'Workspace deleted',
        description: 'Create a new workspace to continue',
        icon: 'i-lucide-check',
        color: 'success'
      })
    }
  } catch (error) {
    console.error('Error deleting workspace:', error)
    toast.add({
      title: 'Error deleting workspace',
      description: 'Please try again later',
      icon: 'i-lucide-alert-circle',
      color: 'error'
    })
  } finally {
    loading.value = false
    deleteModalConfirmation.value = false
  }
}
</script>
