<template>
  <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
    <UFormField label="Title" name="title">
      <UInput v-model="state.title" class="w-full" size="lg" variant="subtle" />
    </UFormField>

    <UFormField label="Note" name="content">
      <UTextarea
        v-model="state.content"
        type="textarea"
        class="w-full"
        :rows="10"
        size="lg"
        variant="subtle"
      />
    </UFormField>
    <div class="flex justify-end gap-2">
      <UButton type="submit" size="lg" variant="subtle"> Submit </UButton>
    </div>
  </UForm>
</template>

<script setup lang="ts">
import * as z from 'zod'
import type { FormSubmitEvent } from '@nuxt/ui'

const { user } = useAuth()
const { activeWorkspace } = useWorkspaces()
const { createNote } = useNotes()

const emit = defineEmits(['close'])
const loading = ref(false)

const schema = z.object({
  title: z.string('Title is required'),
  content: z.string('Note content is required')
})

type Schema = z.output<typeof schema>

const state = reactive<Partial<Schema>>({
  title: undefined,
  content: undefined
})

const toast = useToast()
async function onSubmit(event: FormSubmitEvent<Schema>) {
  try {
    loading.value = true
    await createNote({
      user: user.value?.id as string,
      workspace: activeWorkspace.value?.id as string,
      title: event.data.title,
      content: `<p>${event.data.content}</p>`
    })
    emit('close')
    loading.value = false
  } catch (error) {
    loading.value = false
    console.error('Error creating note:', error)
    toast.add({
      title: 'Error creating note',
      description: 'Please try again later',
      icon: 'i-lucide-alert-circle',
      color: 'error'
    })
  }
}
</script>
