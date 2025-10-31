export interface PolarProduct {
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

export interface PolarProductFeature {
  icon: string
  label: string
}

export const useProducts = () => {
  const products = useState<PolarProduct[]>('polar-products', () => [])

  const normalizeProduct = (product: PolarProduct): PolarProduct => ({
    ...product,
    price_currency: product.price_currency?.toLowerCase() || 'usd',
    recurring_interval: product.recurring_interval || 'month',
    recurring_interval_count: product.recurring_interval_count || 1,
    features: product.features || []
  })

  const setProducts = (newProducts: PolarProduct[]) => {
    products.value = newProducts.map(normalizeProduct)
  }

  const formatInterval = (product: PolarProduct) => {
    const count = product.recurring_interval_count || 1
    const interval = product.recurring_interval || 'month'

    if (count <= 1) {
      return interval
    }

    const plural = interval.endsWith('s') ? interval : `${interval}s`
    return `${count} ${plural}`
  }

  const getProductById = (id: string) => {
    return products.value.find(p => p.id === id)
  }

  const formatPrice = (product: PolarProduct) => {
    const amount = product.price_amount / 100
    const formatted = new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: product.price_currency?.toUpperCase() || 'USD'
    }).format(amount)

    return `${formatted} / ${formatInterval(product)}`
  }

  return {
    products: readonly(products),
    setProducts,
    getProductById,
    formatPrice,
    formatInterval
  }
}
