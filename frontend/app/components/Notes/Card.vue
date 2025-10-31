<template>
  <UContextMenu
    :items="items"
    :ui="{
      content: 'w-48'
    }"
  >
    <div
      class="relative flex cursor-pointer flex-col gap-3 overflow-hidden
        rounded-lg bg-white p-4 text-left shadow-md transition-all duration-200
        focus:outline-none dark:bg-neutral-800 w-full"
    >
      <p class="text-sm font-bold">
        {{ note.title }}
      </p>
      <div
        v-html="note.content"
        class="prose prose-neutral prose-sm text-muted flex-1"
      />
      <USeparator
        :ui="{
          border: 'border-black/5 dark:border-white/5'
        }"
      />
      <div class="flex justify-end">
        <span class="text-dimmed text-[10px]">
          {{ useDateFormat(note.created, 'MMM D, YYYY').value }}
        </span>
      </div>
    </div>
  </UContextMenu>
</template>

<script lang="ts" setup>
import type { NotesRecord } from '@@/types/pocketbase'

const props = defineProps<{
  note: NotesRecord
}>()

const { deleteNote } = useNotes()

const handleDeleteNote = async (noteId: string) => {
  try {
    await deleteNote(noteId)
  } catch (error) {
    console.error('Error deleting note:', error)
  }
}

import type { ContextMenuItem } from '@nuxt/ui'

const items = ref<ContextMenuItem[]>([
  {
    label: 'Delete',
    icon: 'i-lucide-trash-2',
    color: 'error',
    onSelect: () => handleDeleteNote(props.note.id)
  }
])
</script>
