<template>
  <UDropdownMenu
    :items="items"
    :ui="{
      content: collapsed ? 'w-48' : 'w-(--reka-dropdown-menu-trigger-width) '
    }"
  >
    <UButton
      v-bind="{
        ...user,
        label: collapsed ? undefined : user?.name,
        trailingIcon: collapsed ? undefined : 'i-lucide-chevron-up'
      }"
      color="neutral"
      variant="ghost"
      block
      :square="collapsed"
      class="data-[state=open]:bg-elevated"
      :ui="{
        trailingIcon: 'text-dimmed size-4'
      }"
    />
  </UDropdownMenu>

  <UModal v-model:open="userSettingsModal" title="Account Settings">
    <template #body>
      <AccountSettings />
    </template>
  </UModal>
</template>

<script setup lang="ts">
import type { DropdownMenuItem } from '@nuxt/ui'

defineProps<{
  collapsed?: boolean
}>()

const userSettingsModal = ref(false)
const { user: UserData, signOut, avatarUrl } = useAuth()
const { getAvatarUrl } = useAvatar()

const colorMode = useColorMode()
const user = computed(() => ({
  name: UserData.value?.name,
  email: UserData.value?.email,
  avatar: {
    src: getAvatarUrl(avatarUrl.value, UserData.value?.name),
    alt: UserData.value?.name
  }
}))

const items = computed<DropdownMenuItem[][]>(() => [
  [
    {
      type: 'label',
      label: user.value.name,
      description: user.value.email,
      avatar: user.value.avatar
    }
  ],
  [
    {
      label: 'Account Settings',
      icon: 'i-lucide-user-cog',
      onSelect: () => {
        userSettingsModal.value = true
      }
    },
    {
      label: 'Appearance',
      icon: 'i-lucide-sun-moon',
      children: [
        {
          label: 'Light',
          icon: 'i-lucide-sun',
          type: 'checkbox',
          checked: colorMode.preference === 'light',
          onSelect(e: Event) {
            e.preventDefault()
            colorMode.preference = 'light'
          }
        },
        {
          label: 'Dark',
          icon: 'i-lucide-moon',
          type: 'checkbox',
          checked: colorMode.preference === 'dark',
          onUpdateChecked(checked: boolean) {
            if (checked) {
              colorMode.preference = 'dark'
            }
          },
          onSelect(e: Event) {
            e.preventDefault()
          }
        },
        {
          label: 'System',
          icon: 'i-lucide-monitor',
          type: 'checkbox',
          checked: colorMode.preference === 'system',
          onUpdateChecked(checked: boolean) {
            if (checked) {
              colorMode.preference = 'system'
            }
          },
          onSelect(e: Event) {
            e.preventDefault()
          }
        }
      ]
    }
  ],
  [
    {
      label: 'Log out',
      icon: 'i-lucide-log-out',
      onSelect: () => {
        signOut()
      },
      color: 'error'
    }
  ]
])
</script>
