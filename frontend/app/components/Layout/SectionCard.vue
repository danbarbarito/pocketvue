<template>
  <div class="rounded-lg bg-neutral-100 p-[5px] dark:bg-neutral-800">
    <div class="flex items-center justify-between px-4 py-2 pe-1.5">
      <p
        class="text-xs font-semibold text-neutral-500 uppercase
          dark:text-neutral-400"
      >
        {{ title }}
      </p>
      <slot name="header-right" />
    </div>
    <div
      class="rounded-md bg-white p-4 ring-1 ring-neutral-200/70 sm:p-6
        dark:bg-neutral-950 dark:ring-neutral-800"
    >
      <slot />
    </div>
  </div>
</template>

<script lang="ts" setup>
interface Props {
  title: string
}

defineProps<Props>()
</script>
