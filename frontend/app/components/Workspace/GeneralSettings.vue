<template>
  <div>
    <UCard>
      <template #header>
        <p class="font-medium">Workspace Details</p>
      </template>
      <UForm
        :schema="schema"
        :state="state"
        class="space-y-4"
        @submit="onSubmit"
      >
        <UFormField label="Logo" name="logo">
          <AccountAvatarUploader
            v-model="state.logo"
            @file-selected="handleFileSelected"
          />
        </UFormField>
        <UFormField label="Workspace name" name="name">
          <UInput
            v-model="state.name"
            placeholder="My awesome workspace"
            size="lg"
            variant="soft"
            class="w-full"
          />
        </UFormField>

        <UFormField
          label="Unique workspace URL"
          name="slug"
          hint="Lowercase alphanumeric only"
          :error="slugError || undefined"
          :class="{ shake: slugError }"
        >
          <UInput
            v-model="state.slug"
            placeholder="my-workspace"
            size="lg"
            variant="soft"
            class="w-full"
          />
        </UFormField>

        <UFormField label="Website URL" name="domain">
          <UInput
            v-model="state.domain"
            placeholder="mycompany.com"
            size="lg"
            variant="soft"
            class="w-full"
          />
        </UFormField>

        <div class="flex justify-end gap-2 pt-4">
          <UButton label="Cancel" variant="soft" size="lg" @click="resetForm" />
          <UButton
            label="Save Changes"
            size="lg"
            type="submit"
            :loading="loading"
            :disabled="loading"
          />
        </div>
      </UForm>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import * as z from 'zod'
import type { FormSubmitEvent } from '@nuxt/ui'
import type { WorkspacesResponse } from '@@/types/pocketbase'
import type { RecordModel } from 'pocketbase'

const { activeWorkspace } = useWorkspaces()
const { pb } = usePocketbase()
const { getFileUrl } = usePocketbaseFiles()
const loading = ref(false)
const selectedFile = ref<File | null>(null)
const toast = useToast()
const slugError = ref<string | null>(null)

// Define form schema
const schema = z.object({
  name: z.string().min(1, 'Workspace name is required'),
  domain: z.url('Invalid URL'),
  logo: z.string().optional(),
  slug: z
    .string()
    .min(1, 'Slug is required')
    .regex(
      /^[a-z0-9-]+$/,
      'Slug can only contain lowercase letters, numbers, and hyphens'
    )
})

type Schema = z.output<typeof schema>

// Get the logo URL from PocketBase if available
const logoUrl = computed<string>(() => {
  if (activeWorkspace.value?.id && activeWorkspace.value?.logo) {
    // Type assertion to RecordModel since WorkspacesResponse has all the required properties
    const url = getFileUrl(
      activeWorkspace.value as unknown as RecordModel,
      'logo'
    )
    return url || ''
  }
  return ''
})

// Initialize form state with active workspace data
const state = reactive<Partial<Schema>>({
  name: activeWorkspace.value?.name || '',
  domain: activeWorkspace.value?.domain || '',
  logo: logoUrl.value,
  slug: activeWorkspace.value?.slug || ''
})

// Handle file selection for logo
const handleFileSelected = (file: File | null) => {
  selectedFile.value = file
  if (!file) {
    // If no file is selected, clear the logo
    state.logo = ''
  }
  // When a file is selected, the component will set the model value to the blob URL
}

// Reset form to original values
const resetForm = () => {
  state.name = activeWorkspace.value?.name || ''
  state.domain = activeWorkspace.value?.domain || ''
  state.logo = logoUrl.value
  state.slug = activeWorkspace.value?.slug || ''
  selectedFile.value = null
  slugError.value = null
}

interface PocketbaseError {
  data?: {
    data?: {
      slug?: {
        code?: string
      }
    }
  }
  message?: string
}

// Handle form submission
const onSubmit = async (event: FormSubmitEvent<Schema>) => {
  try {
    loading.value = true

    // Only proceed if we have a workspace ID
    if (!activeWorkspace.value?.id) {
      throw new Error('Workspace not found')
    }

    const formData = new FormData()

    // Add form fields to FormData
    formData.append('name', event.data.name)
    formData.append('domain', event.data.domain)
    formData.append('slug', event.data.slug)

    // Handle logo field
    if (selectedFile.value) {
      // If a new file was selected, add it to FormData
      formData.append('logo', selectedFile.value)
    } else if (state.logo === '') {
      // If the logo was explicitly removed, set it to empty
      formData.append('logo', '')
    }
    // If no change to logo, don't include it in FormData

    // Update workspace in PocketBase
    await pb
      .collection('workspaces')
      .update<WorkspacesResponse>(activeWorkspace.value.id, formData)

    // Show success toast
    toast.add({
      title: 'Workspace updated',
      description: 'Workspace details were updated successfully',
      color: 'success'
    })

    // Redirect to new slug URL if slug was changed
    if (event.data.slug !== activeWorkspace.value.slug) {
      await navigateTo(`/${event.data.slug}/dashboard/settings`, {
        replace: true
      })
    } else {
      // Refresh the workspace data
      const { fetchWorkspaces } = useWorkspaces()
      await fetchWorkspaces()
    }
  } catch (error: unknown) {
    console.error('Update workspace error:', error)
    const pbError = error as PocketbaseError

    if (pbError?.data?.data?.slug?.code === 'validation_not_unique') {
      slugError.value = 'URL is already taken, try a different one'
    }

    toast.add({
      title: 'Error updating workspace',
      description:
        error instanceof Error ? error.message : 'An unexpected error occurred',
      color: 'error'
    })
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.shake {
  animation: shake 0.82s cubic-bezier(0.36, 0.07, 0.19, 0.97) both;
  transform: translate3d(0, 0, 0);
}

@keyframes shake {
  10%,
  90% {
    transform: translate3d(-1px, 0, 0);
  }
  20%,
  80% {
    transform: translate3d(2px, 0, 0);
  }
  30%,
  50%,
  70% {
    transform: translate3d(-4px, 0, 0);
  }
  40%,
  60% {
    transform: translate3d(4px, 0, 0);
  }
}
</style>
