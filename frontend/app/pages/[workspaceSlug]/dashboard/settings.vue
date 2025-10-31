<template>
  <UDashboardPanel id="workspace-settings">
    <template #header>
      <UDashboardNavbar title="Workspace Settings" :ui="{ right: 'gap-3' }">
        <template #leading>
          <UDashboardSidebarCollapse :ui="{ leadingIcon: 'size-4' }" />
          <USeparator orientation="vertical" class="mr-2 h-4" />
        </template>
      </UDashboardNavbar>

      <UDashboardToolbar>
        <template #left>
          <UNavigationMenu
            :items="links"
            highlight
            class="-mx-1 flex-1 md:-mx-3"
          />
        </template>
      </UDashboardToolbar>
    </template>

    <template #body>
      <NuxtPage />
    </template>
  </UDashboardPanel>
</template>

<script setup lang="ts">
import type { NavigationMenuItem } from '@nuxt/ui'
import { getWorkspacePath } from '~/utils/workspace'

const { activeWorkspaceSlug } = useWorkspaces()

const links = [
  [
    {
      label: 'General',
      icon: 'i-lucide-folder-cog',
      to: getWorkspacePath(activeWorkspaceSlug.value, '/settings'),
      exact: true
    },
    {
      label: 'Billing',
      icon: 'i-lucide-credit-card',
      to: getWorkspacePath(activeWorkspaceSlug.value, '/settings/billing')
    }
  ]
] satisfies NavigationMenuItem[][]
</script>
