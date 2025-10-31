<template>
  <NuxtLoadingIndicator color="#000" />
  <NuxtRouteAnnouncer />
  <UApp
    :toaster="{
      position: 'top-center',
      progress: false
    }"
  >
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </UApp>
</template>

<script setup lang="ts">
import { withoutTrailingSlash } from 'ufo'
const route = useRoute()

const { seo } = useAppConfig()

useSeoMeta({
  title: seo.title,
  description: seo.description,
  ogTitle: seo.title,
  ogDescription: seo.description,
  ogImage: seo.image,
  ogUrl: seo.url,
  twitterTitle: seo.title,
  twitterDescription: seo.description,
  twitterImage: seo.image,
  twitterCard: 'summary_large_image'
})

useHead({
  meta: [
    {
      name: 'viewport',
      content:
        'width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no'
    }
  ],
  link: [
    { rel: 'icon', href: '/favicon.ico', sizes: '48x48' },
    {
      rel: 'canonical',
      href: `https://piqo.app${withoutTrailingSlash(route.path)}`
    }
  ],
  htmlAttrs: {
    lang: 'en'
  }
})
</script>
