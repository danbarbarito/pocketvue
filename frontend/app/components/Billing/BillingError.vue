<template>
  <UAlert
    v-if="error"
    color="error"
    variant="soft"
    :title="error"
    :close-button="{
      icon: 'i-heroicons-x-mark-20-solid',
      color: 'error',
      variant: 'link'
    }"
    @close="$emit('close')"
  />
</template>

<script lang="ts" setup>
interface Props {
  error: string | null
}

defineProps<Props>()

defineEmits<{
  close: []
}>()
</script>
