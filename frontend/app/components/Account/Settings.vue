<template>
  <UForm
    :schema="schema"
    :state="state"
    class="space-y-4"
    @submit="onSubmit as any"
  >
    <UFormField label="Avatar" name="avatar">
      <AccountAvatarUploader
        v-model="state.avatarUrl"
        @file-selected="handleFileSelected"
      />
    </UFormField>
    <UFormField label="Name" name="name">
      <UInput
        v-model="state.name"
        placeholder="Name"
        class="w-full"
        size="lg"
      />
    </UFormField>
    <UFormField label="Email">
      <UInput
        :value="user?.email"
        placeholder="Email"
        class="w-full"
        disabled
        variant="subtle"
        size="lg"
      />
    </UFormField>
    <UFormField label="Account ID">
      <UInput
        :value="user?.id"
        placeholder="Account ID"
        class="w-full"
        disabled
        variant="subtle"
        size="lg"
      />
    </UFormField>
    <UButton
      color="neutral"
      :loading="loading"
      :disabled="loading"
      type="submit"
      label="Save"
    />
  </UForm>
</template>

<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types'
import { z } from 'zod'
import type { RecordModel } from 'pocketbase'
const { getFileUrl } = usePocketbaseFiles()
const { pb, user } = usePocketbase()
const loading = ref(false)
const selectedFile = ref<File | null>(null)

// Form validation schema
const schema = z.object({
  avatarUrl: z.string().optional(),
  name: z.string().min(1, 'Name is required')
})

type FormSchema = typeof schema

const handleFileSelected = (file: File | null) => {
  selectedFile.value = file
  if (!file) {
    state.avatarUrl = ''
  }
}
const toast = useToast()
const state = reactive<z.infer<FormSchema>>({
  name: user.value?.name ?? '',
  avatarUrl: getFileUrl(user.value as RecordModel, 'avatar') ?? ''
})

const onSubmit = async (event: FormSubmitEvent<z.infer<FormSchema>>) => {
  try {
    loading.value = true
    const formData = new FormData()
    formData.append('name', event.data.name)
    if (selectedFile.value) {
      formData.append('avatar', selectedFile.value)
    }

    // Update user profile in PocketBase
    const updatedUser = await pb
      .collection('users')
      .update(user.value!.id, formData)

    // Update the global auth store with the new user data
    pb.authStore.save(pb.authStore.token, updatedUser as RecordModel)

    toast.add({
      title: 'Profile updated',
      description: 'Your profile has been updated',
      color: 'success'
    })
  } catch (error) {
    console.error(error)
    toast.add({
      title: 'Update failed',
      description: 'Failed to update profile. Please try again.',
      color: 'error'
    })
  } finally {
    loading.value = false
  }
}
</script>
