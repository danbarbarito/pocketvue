import type { NotesRecord, Create } from '@@/types/pocketbase'

export const useNotes = () => {
  const { pb } = usePocketbase()
  const toast = useToast()
  const notes = useState<NotesRecord[]>('notes', () => [])
  let unsubscribeFn: (() => void) | null = null

  const subscribeNotes = async () => {
    unsubscribeFn = await pb.collection('notes').subscribe('*', e => {
      console.log('Realtime event:', e.action, e.record.id)

      if (!notes.value) return

      switch (e.action) {
        case 'create':
          console.log('Creating note:', e.record.title)
          notes.value = [e.record, ...notes.value]
          break
        case 'update':
          console.log('Updating note:', e.record.title)
          const updateIndex = notes.value.findIndex(
            note => note.id === e.record.id
          )
          if (updateIndex !== -1) {
            notes.value = [
              ...notes.value.slice(0, updateIndex),
              e.record,
              ...notes.value.slice(updateIndex + 1)
            ]
          }
          break
        case 'delete':
          console.log('Deleting note:', e.record.id)
          notes.value = notes.value.filter(note => note.id !== e.record.id)
          break
      }
    })
  }

  const unsubscribeNotes = () => {
    if (unsubscribeFn) {
      unsubscribeFn()
      unsubscribeFn = null
    }
  }

  const fetchNotes = async () => {
    return await pb.collection('notes').getFullList<NotesRecord>({
      sort: '-created'
    })
  }

  const createNote = async (note: Create<'notes'>) => {
    try {
      await pb.collection('notes').create(note)
      toast.add({
        title: 'Note created',
        description: 'Your note has been created successfully',
        color: 'success'
      })
    } catch (error) {
      console.error('Error creating note:', error)
      toast.add({
        title: 'Error creating note',
        description: 'Please try again later',
        color: 'error'
      })
    }
  }

  const updateNote = async (note: NotesRecord) => {
    try {
      await pb.collection('notes').update(note.id, note)
      toast.add({
        title: 'Note updated',
        description: 'Your note has been updated successfully',
        icon: 'i-lucide-check',
        color: 'success'
      })
    } catch (error) {
      console.error('Error updating note:', error)
      toast.add({
        title: 'Error updating note',
        description: 'Please try again later',
        color: 'error'
      })
    }
  }

  const deleteNote = async (noteId: string) => {
    try {
      await pb.collection('notes').delete(noteId)
      toast.add({
        title: 'Note deleted',
        description: 'Your note has been deleted successfully',
        color: 'success'
      })
    } catch (error) {
      console.error('Error deleting note:', error)
      toast.add({
        title: 'Error deleting note',
        description: 'Please try again later',
        color: 'error'
      })
    }
  }

  return {
    notes,
    subscribeNotes,
    unsubscribeNotes,
    createNote,
    updateNote,
    deleteNote,
    fetchNotes
  }
}
