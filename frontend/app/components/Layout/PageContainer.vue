<template>
  <UDashboardPanel :ui="{ body: bodyPadding }">
    <template #header>
      <UDashboardNavbar :title="title" :ui="{ root: 'sm:px-4' }">
        <template #leading>
          <UDashboardSidebarCollapse :ui="{ leadingIcon: 'size-4' }" />
          <USeparator orientation="vertical" class="mr-2 h-4" />
          <slot name="header-leading" />
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <slot />
    </template>
  </UDashboardPanel>
</template>

<script lang="ts" setup>
interface Props {
  title: string
  noPadding?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  noPadding: false
})

const bodyPadding = computed(() => (props.noPadding ? 'sm:p-0 p-0' : ''))
</script>
