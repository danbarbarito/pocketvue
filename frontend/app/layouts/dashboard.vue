<template>
  <UDashboardGroup unit="rem">
    <UDashboardSidebar
      id="default"
      v-model:open="open"
      collapsible
      resizable
      mode="drawer"
      class="@container"
      :ui="{
        footer: 'lg:border-t lg:border-none p-2',
        body: 'gap-2 px-2 @max-[4rem]:p-3 @max-[4rem]:pt-2 md:pt-2',
        header: 'px-2 pb-0 h-10 pt-2',
        root: 'min-w-14'
      }"
    >
      <template #header="{ collapsed }">
        <WorkspacesMenu :collapsed="collapsed" />
      </template>

      <template #default="{ collapsed }">
        <UDashboardSearchButton
          :collapsed="collapsed"
          class="bg-transparent py-1.5 pr-1.5 ring-neutral-300
            dark:ring-white/5"
          size="md"
          :ui="{
            leadingIcon: 'size-4'
          }"
          label="Go to"
        />

        <UNavigationMenu
          :collapsed="collapsed"
          :items="links[0]"
          orientation="vertical"
          tooltip
          popover
          color="neutral"
          :ui="{
            linkTrailingIcon: 'text-dimmed'
          }"
        />
      </template>

      <template #footer="{ collapsed }">
        <UserMenu :collapsed="collapsed" />
      </template>
    </UDashboardSidebar>
    <UDashboardSearch
      :groups="groups"
      :ui="{ close: 'text-muted', modal: 'sm:max-w-xl' }"
    />
    <slot />
  </UDashboardGroup>
</template>

<script setup lang="ts">
import { getWorkspaceAvatar, getWorkspacePath } from '~/utils/workspace'

const { workspaces } = useWorkspaces()
const { navigationLinks } = useNavigationLinks()
const { getFileUrl } = usePocketbaseFiles()
const { getAvatarUrl } = useAvatar()
const router = useRouter()
const open = ref(false)

const links = computed(() => {
  return navigationLinks.value.map(group =>
    group.map(link => ({
      ...link,
      onSelect: () => {
        open.value = false
        if (link.to && !link.target) {
          router.push(link.to)
        }
      }
    }))
  )
})

const groups = computed(() => {
  const navigationItems = navigationLinks.value.flat().map(link => ({
    label: link.label,
    icon: link.icon,
    to: link.to,
    target: link.target,
    exact: link.exact,
    onSelect: () => {
      if (link.to && !link.target) {
        router.push(String(link.to))
      } else if (link.to && link.target) {
        window.open(String(link.to), link.target)
      }
    }
  }))

  const baseGroups = [
    {
      id: 'links',
      label: 'Navigate to',
      items: navigationItems.flat()
    }
  ]

  if (workspaces.value.length > 0) {
    const workspaceItems = workspaces.value.map(workspace => {
      // Extract current sub-route to preserve it when switching workspaces
      const currentPath = router.currentRoute.value.path
      const workspacePattern = /^\/[^\/]+\/dashboard(.*)$/
      const match = currentPath.match(workspacePattern)
      const subRoute = match ? match[1] : ''

      return {
        exact: true,
        label: workspace.name,
        suffix: workspace.slug,
        to: getWorkspacePath(workspace.slug),
        icon: '',
        target: undefined,
        avatar: getWorkspaceAvatar(workspace, getFileUrl, getAvatarUrl),
        onSelect: () => {
          router.push(getWorkspacePath(workspace.slug, subRoute))
        }
      }
    })

    baseGroups.push({
      id: 'workspaces',
      label: 'Switch workspace',
      items: [
        ...workspaceItems,
        {
          label: 'Create new workspace',
          icon: 'i-lucide-circle-plus',
          to: '/new-workspace',
          target: undefined,
          exact: true,
          onSelect: () => {
            router.push('/new-workspace')
          }
        }
      ]
    })
  }

  return baseGroups
})
</script>
