<template>
  <section class="space-y-8">
    <div>
      <h2 class="font-bold">Plans</h2>
      <p class="text-muted text-sm">You are currently on the free plan.</p>
    </div>

    <div v-if="productsLoading" class="flex justify-center py-12">
      <div
        class="border-primary h-10 w-10 animate-spin rounded-full border-2
          border-t-transparent"
      />
    </div>

    <div
      v-else-if="products.length"
      class="grid grid-cols-1 items-start gap-4 md:grid-cols-2 lg:grid-cols-3"
    >
      <div
        v-for="product in products"
        :key="product.id"
        class="p-[5px]"
        :class="{
          'rounded-2xl bg-neutral-100 dark:bg-neutral-800':
            product.id === recommendedPlanId
        }"
      >
        <div
          v-if="product.id === recommendedPlanId"
          class="shimmer flex h-9 items-center justify-center gap-2"
        >
          <UIcon name="solar:star-bold" class="size-4" />
          <span class="text-primary text-sm">Most Popular</span>
        </div>
        <div v-else class="h-9" />
        <div
          class="rounded-xl p-6"
          :class="{
            [`bg-white ring-1 ring-neutral-200 dark:bg-neutral-950
            dark:ring-neutral-800`]: product.id === recommendedPlanId
          }"
        >
          <div class="mb-6">
            <h3 class="text-sm font-medium">
              {{ product.name }}
            </h3>
            <div class="mt-1 text-xl font-bold">
              {{ formatPrice(product) }}
            </div>
            <p class="text-muted mt-4 text-sm">
              {{ product.description || 'Subscribe to unlock more features.' }}
            </p>
          </div>
          <UButton
            label="Subscribe"
            :loading="loadingProduct === product.id"
            :disabled="!!loadingProduct"
            block
            size="lg"
            :variant="product.id === recommendedPlanId ? 'solid' : 'outline'"
            @click="$emit('subscribe', product.id)"
          />
          <ul v-if="product.features?.length" class="text-muted mt-6 space-y-3">
            <li
              v-for="feature in product.features"
              :key="feature.label"
              class="flex items-start gap-3 text-sm"
            >
              <UIcon
                :name="feature.icon || 'i-lucide-check'"
                class="text-muted mt-0.5 size-4"
              />
              <span>{{ feature.label }}</span>
            </li>
          </ul>
        </div>
      </div>
    </div>

    <div
      v-else
      class="text-muted rounded-2xl border border-dashed border-neutral-200
        py-16 text-center"
    >
      No products available at the moment.
    </div>
  </section>
</template>

<script lang="ts" setup>
interface PolarProduct {
  id: string
  name: string
  description?: string
  price_amount: number
  price_currency: string
  recurring_interval: string
  recurring_interval_count: number
  is_recurring: boolean
  trial_interval?: string
  trial_interval_count?: number
  polar_price_id: string
  features: readonly PolarProductFeature[]
}

interface PolarProductFeature {
  icon: string
  label: string
}

interface Props {
  products: readonly PolarProduct[]
  productsLoading: boolean
  loadingProduct: string | null
  recommendedPlanId: string | null
  formatPrice: (product: PolarProduct) => string
}

defineProps<Props>()

defineEmits<{
  subscribe: [planId: string]
}>()
</script>

<style scoped></style>
