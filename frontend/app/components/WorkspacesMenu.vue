<template>
  <UDropdownMenu
    :items="items"
    :ui="{
      content: collapsed ? 'w-48' : 'w-(--reka-dropdown-menu-trigger-width)',
      item: 'cursor-pointer'
    }"
  >
    <UButton
      v-bind="{
        ...selectedWorkspaceProps,
        label: collapsed ? undefined : selectedWorkspaceProps?.label,
        trailingIcon: collapsed ? undefined : 'i-lucide-chevron-down'
      }"
      color="neutral"
      variant="ghost"
      block
      :square="collapsed"
      class="data-[state=open]:bg-elevated"
      :class="[!collapsed && 'py-2']"
      :ui="{
        trailingIcon: 'text-dimmed size-4'
      }"
    />
  </UDropdownMenu>
</template>

<script setup lang="ts">
import type { DropdownMenuItem } from '@nuxt/ui'
import { getWorkspaceAvatar, getWorkspacePath } from '~/utils/workspace'

const { workspaces, activeWorkspace, fetchWorkspaces, activeWorkspaceSlug } =
  useWorkspaces()
const { getFileUrl } = usePocketbaseFiles()
const { getAvatarUrl } = useAvatar()
const { setLastUsedWorkspace } = useWorkspacePreferences()

defineProps<{
  collapsed?: boolean
}>()

onMounted(() => {
  if (workspaces.value.length === 0) {
    fetchWorkspaces()
  }
})

watch(
  activeWorkspaceSlug,
  newSlug => {
    if (newSlug) {
      setLastUsedWorkspace(newSlug)
    }
  },
  { immediate: true }
)

const selectedWorkspaceProps = computed(() => {
  if (!activeWorkspace.value?.name) {
    if (workspaces.value.length > 0 && !activeWorkspaceSlug.value) {
      const firstWorkspace = workspaces.value[0]!
      return {
        label: firstWorkspace.name,
        avatar: getWorkspaceAvatar(firstWorkspace, getFileUrl, getAvatarUrl)
      }
    }

    return {
      label: 'Select workspace',
      avatar: {
        src: getAvatarUrl(undefined, 'default-workspace'),
        alt: 'Select workspace'
      }
    }
  }

  return {
    label: activeWorkspace.value.name,
    avatar: getWorkspaceAvatar(activeWorkspace.value, getFileUrl, getAvatarUrl)
  }
})

const items = computed<DropdownMenuItem[][]>(() => {
  const workspaceItems = workspaces.value.map(workspace => ({
    label: workspace.name,
    avatar: getWorkspaceAvatar(workspace, getFileUrl, getAvatarUrl),
    active: workspace.slug === activeWorkspaceSlug.value,
    to: getWorkspacePath(workspace.slug),
    onSelect() {
      setLastUsedWorkspace(workspace.slug)
    }
  }))

  return [
    workspaceItems,
    [
      {
        label: 'Create workspace',
        icon: 'i-lucide-circle-plus',
        to: '/new-workspace'
      }
    ]
  ]
})
</script>
