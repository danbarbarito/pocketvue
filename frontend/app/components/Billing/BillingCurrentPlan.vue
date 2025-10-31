<template>
  <section class="w-full space-y-4">
    <!-- Current Plan Card -->
    <LayoutSectionCard title="Current Plan">
      <template #header-right>
        <UBadge
          :color="getStatusColor(subscriptionStatus)"
          variant="subtle"
          class="rounded-[6px] capitalize"
        >
          {{ subscriptionStatus }}
        </UBadge>
      </template>

      <div class="space-y-4">
        <!-- Plan Name and Price -->
        <div class="space-y-2">
          <h2 class="text-highlighted font-bold">
            {{ currentPlan?.name || 'Custom Plan' }} -
            {{ formatPrice(currentPlan) }}
          </h2>
          <p class="text-muted text-sm">
            {{
              currentPlan?.description ||
              'This is the plan tied to your workspace.'
            }}
          </p>
        </div>

        <!-- Billing Period Info -->
        <div
          v-if="user?.subscription_current_period_end"
          class="space-y-1 text-sm"
        >
          <p class="text-muted text-sm">
            {{ subscriptionStatus === 'active' ? 'Renews on' : 'Expires on' }}
            <span class="font-semibold">
              {{ formatDate(user?.subscription_current_period_end) }}
            </span>
          </p>
        </div>

        <!-- Cancellation Warning -->
        <UAlert
          v-if="user?.subscription_cancel_at_period_end"
          color="warning"
          variant="soft"
          icon="i-lucide-alert-triangle"
          title="Subscription Canceling"
          :description="`Your subscription will cancel at the end of the current billing period. You'll keep access until ${formatDate(user?.subscription_current_period_end)}.`"
        />

        <!-- Features Grid -->
        <div
          v-if="currentPlan?.features?.length"
          class="mt-6 space-y-3 text-sm"
        >
          <div
            v-for="(feature, i) in currentPlan.features"
            :key="i"
            class="flex items-start gap-3"
          >
            <UIcon
              :name="feature.icon || 'i-lucide-check'"
              class="mt-0.5 h-4 w-4 shrink-0"
            />
            <span class="text-neutral-700 dark:text-neutral-300">{{
              feature.label
            }}</span>
          </div>
        </div>
      </div>
    </LayoutSectionCard>

    <!-- Manage Subscription Card -->
    <LayoutSectionCard title="Manage Subscription">
      <div class="space-y-4">
        <div class="space-y-2">
          <h3 class="text-sm font-medium">Update Your Subscription</h3>
          <p class="text-muted text-sm">
            Change your plan, update payment methods, or cancel when you need
            to.
          </p>
        </div>
        <UButton
          :loading="isLoadingPortal"
          @click="$emit('manage-subscription')"
        >
          Manage via Polar
        </UButton>
      </div>
    </LayoutSectionCard>
  </section>
</template>

<script lang="ts" setup>
interface Props {
  user: any
  currentPlan: any
  subscriptionStatus: string
  subscriptionPeriodEnd: string
  isLoadingPortal: boolean
  formatPrice: (product: any) => string
}

defineProps<Props>()

defineEmits<{
  'manage-subscription': []
}>()

const getStatusColor = (status: string) => {
  switch (status) {
    case 'active':
      return 'success'
    case 'canceled':
      return 'warning'
    case 'past_due':
      return 'error'
    case 'trialing':
      return 'info'
    default:
      return 'neutral'
  }
}

const formatDate = (dateString: string | undefined) => {
  if (!dateString) return ''
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}
</script>
