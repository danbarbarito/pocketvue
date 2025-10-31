<template>
  <UAlert
    v-if="showSuccessMessage"
    color="success"
    variant="soft"
    title="Subscription Successful!"
    description="Your subscription has been activated."
    :close-button="{
      icon: 'i-heroicons-x-mark-20-solid',
      color: 'success',
      variant: 'link'
    }"
    @close="$emit('close')"
  />
</template>

<script lang="ts" setup>
interface Props {
  showSuccessMessage: boolean
}

defineProps<Props>()

defineEmits<{
  close: []
}>()
</script>
