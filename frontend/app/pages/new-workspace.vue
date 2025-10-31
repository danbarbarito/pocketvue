<template>
  <UContainer class="flex h-dvh max-w-4xl flex-col items-center justify-center">
    <div class="w-full max-w-md">
      <div class="mb-12 text-center">
        <p class="text-lg font-semibold">Create a new workspace</p>
        <p class="text-sm text-neutral-500">
          A workspace is a common space for you to work on your projects.
        </p>
      </div>

      <UForm
        :schema="schema"
        :state="state"
        class="space-y-4"
        @submit="onSubmit"
      >
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

        <UFormField label="Logo" name="logo">
          <AccountAvatarUploader
            v-model="state.logo"
            @file-selected="handleFileSelected"
          />
        </UFormField>

        <UButton
          label="Create workspace"
          class="mt-6"
          block
          size="lg"
          type="submit"
          :loading="loading"
          :disabled="loading"
        />
      </UForm>
    </div>
  </UContainer>
</template>

<script setup lang="ts">
import * as z from 'zod'
import type { FormSubmitEvent } from '@nuxt/ui'
import type { WorkspacesResponse } from '@@/types/pocketbase'
import { isSlugBlocked } from '~/constants/blockedSlugs'

const { user } = useAuth()
const { pb } = usePocketbase()
const { workspaces } = useWorkspaces()
const loading = ref(false)
const selectedFile = ref<File | null>(null)
const toast = useToast()

definePageMeta({
  middleware: ['auth']
})

const schema = z.object({
  name: z.string().min(1, 'Workspace name is required'),
  domain: z
    .string()
    .min(1, 'Website URL is required')
    .regex(
      /^https?:\/\/[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(\/.*)?$/,
      'Please enter a valid URL (e.g., https://example.com)'
    ),
  logo: z.string().optional(),
  slug: z
    .string()
    .min(1, 'Slug is required')
    .regex(
      /^[a-z0-9-]+$/,
      'Slug can only contain lowercase letters, numbers, and hyphens'
    )
    .refine((slug: string) => !isSlugBlocked(slug), {
      message: 'This slug is already taken, try a different one'
    })
})

type Schema = z.output<typeof schema>

const state = reactive<Partial<Schema>>({
  name: '',
  domain: '',
  logo: '',
  slug: ''
})

const handleFileSelected = (file: File | null) => {
  selectedFile.value = file
  if (!file) {
    state.logo = ''
  }
}

const slugError = ref<string | null>(null)

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

const onSubmit = async (event: FormSubmitEvent<Schema>) => {
  try {
    loading.value = true
    const formData = new FormData()

    // Add form fields to FormData
    formData.append('name', event.data.name)
    formData.append('domain', event.data.domain)
    formData.append('slug', event.data.slug)
    formData.append('user', user.value?.id as string)

    // Add logo if a file was selected
    if (selectedFile.value) {
      formData.append('logo', selectedFile.value)
    }

    // Create workspace in PocketBase
    const workspace = await pb
      .collection('workspaces')
      .create<WorkspacesResponse>(formData)

    // Add the new workspace to the workspaces state to prevent "Workspace not found" error
    workspaces.value = [...workspaces.value, workspace]

    // Show success toast
    toast.add({
      title: 'Workspace created',
      description: `${workspace.name} workspace was created successfully`,
      color: 'success'
    })

    // Redirect to workspace dashboard
    await navigateTo(`/${workspace.slug}/dashboard`, { replace: true })
  } catch (error: unknown) {
    console.error(error)
    const pbError = error as PocketbaseError
    if (pbError?.data?.data?.slug?.code === 'validation_not_unique') {
      slugError.value = 'URL is already taken, try a different one'
    }

    toast.add({
      title: 'Error creating workspace',
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
