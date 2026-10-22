<template>
  <UApp>
    <UHeader>
      <template #left>
        <NuxtLink
          to="/"
          class="outline-primary/25 -ms-1 flex items-center gap-2 rounded-md p-1 focus-visible:outline-3"
        >
          <AppLogo class="h-10 w-auto shrink-0" />
          <span class="text-xl">
            Speedtest Tracker
          </span>
        </NuxtLink>
      </template>

      <template #right>
        <RunSpeedtest />
        <UColorModeButton />

        <UButton
          to="https://github.com/henrywhitaker3/Speedtest-Tracker"
          target="_blank"
          icon="i-simple-icons-github"
          aria-label="GitHub"
          color="neutral"
          variant="ghost"
        />
      </template>
    </UHeader>

    <UMain>
      <NuxtPage />
    </UMain>

    <UFooter>
      <p class="text-muted text-sm">
        Speedtest Tracker {{ version }} • © {{ year }}
      </p>
    </UFooter>
  </UApp>
</template>

<script setup lang="ts">
import appleTouchIcon from '~/assets/icons/apple-icon-180x180.png'
import favicon from '~/assets/icons/favicon.ico'
import favicon16 from '~/assets/icons/favicon-16x16.png'
import favicon32 from '~/assets/icons/favicon-32x32.png'

useHead({
  meta: [
    { name: 'viewport', content: 'width=device-width, initial-scale=1' }
  ],
  link: [
    { rel: 'icon', href: favicon },
    { rel: 'icon', type: 'image/png', sizes: '16x16', href: favicon16 },
    { rel: 'icon', type: 'image/png', sizes: '32x32', href: favicon32 },
    { rel: 'apple-touch-icon', sizes: '180x180', href: appleTouchIcon }
  ],
  htmlAttrs: {
    lang: 'en'
  }
})

const title = 'Speedtest-Tracker'
const description = 'Continuously track your internet speed'

useSeoMeta({
  title,
  description,
  ogTitle: title,
  ogDescription: description
})

const year = new Date().getFullYear()

const version = useRuntimeConfig().public.version
</script>
