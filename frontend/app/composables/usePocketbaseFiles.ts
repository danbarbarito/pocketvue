import type { RecordModel, FileOptions } from 'pocketbase'

export const usePocketbaseFiles = () => {
  const { pb } = usePocketbase()

  const getFileUrl = (
    record: RecordModel,
    filenameField: string = 'avatar',
    queryParams?: FileOptions
  ): string | null => {
    const filename = record[filenameField]
    if (filename && typeof filename === 'string') {
      return pb.files.getURL(record, filename, queryParams)
    }
    return null
  }

  return { getFileUrl }
}
