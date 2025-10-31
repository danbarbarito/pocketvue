import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

const currentDir = dirname(fileURLToPath(import.meta.url))

export default defineNuxtConfig({
  modules: ['@nuxt/ui', '@vueuse/nuxt'],
  compatibilityDate: '2025-07-15',
  devtools: { enabled: false },
  ssr: false,
  css: ['~/assets/css/main.css'],
  fonts: {
    families: [
      {
        name: 'Inter',
        provider: 'google',
        weights: ['400', '500', '600', '700', '800']
      }
    ]
  },
  icon: {
    customCollections: [
      {
        prefix: 'custom',
        dir: './app/assets/icons'
      }
    ]
  },
  nitro: {
    output: {
      publicDir: resolve(currentDir, '../backend/ui/dist')
    }
  }
})
