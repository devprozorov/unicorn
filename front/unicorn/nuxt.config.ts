const isDev = process.env.NODE_ENV === 'development'

export default defineNuxtConfig({
  appDir: 'app',
  css: ['~/assets/css/main.css'],

  modules: [
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt'
  ],

  build: {
    transpile: ['vue-i18n'],
  },

  app: {
    head: {
      title: 'Unicornstar',
      meta: [
        { name: 'description', content: 'Unicornstar' }
      ],
      link: [
        { rel: 'icon', type: 'image/png', href: '/favicon.png' }, // иконка
      ]
    }
  },

  runtimeConfig: {
    apiInternalBase: isDev ? 'http://localhost:8080' : 'http://backend:8080',
    public: {
      apiBase: isDev ? '' : 'https://unicornstar.online' // В dev используем прокси (пустая строка = относительные пути)
    }
  },

  // Прокси для dev режима - решает проблему с CORS и cookies
  nitro: {
    devProxy: isDev ? {
      '/api': {
        target: 'http://localhost:8080/api',
        changeOrigin: true,
        prependPath: false,
        cookieDomainRewrite: 'localhost'
      }
    } : {}
  }

})
