// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from '@tailwindcss/vite'

export default defineNuxtConfig({
  compatibilityDate: '2025-02-24',
  devtools: { enabled: true },
  css: ['~/tailwind.css'],
  vite: {
    plugins: [tailwindcss() as any],
  },
  imports: {
    dirs: ['types'],
  },
  runtimeConfig: {
    public: {
      apiBaseUrl: 'http://localhost:7080',
    },
  },
})
