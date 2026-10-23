// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: [
    '@nuxt/eslint',
    '@nuxt/ui',
    '@pinia/nuxt'
  ],

  ssr: false,

  devtools: {
    enabled: true
  },

  css: ['~/assets/css/main.css'],

  colorMode: {
    preference: 'dark'
  },

  runtimeConfig: {
    public: {
      version: 'dev',
      apiUrl: 'http://localhost:8080'
    }
  },

  routeRules: {
    '/': { prerender: true }
  },

  compatibilityDate: '2026-06-30',

  nitro: {
    preset: 'static'
  },

  eslint: {
    config: {
      stylistic: {
        indent: 2,
        semi: false,
        commaDangle: 'never',
        braceStyle: '1tbs'
      }
    }
  }
})
