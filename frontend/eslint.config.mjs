import withNuxt from './.nuxt/eslint.config.mjs'

export default withNuxt(
  {
    ignores: ['app/client/**']
  },
  {
    rules: {
      'vue/script-setup-no-uses-vars': 'off'
    }
  }
)
