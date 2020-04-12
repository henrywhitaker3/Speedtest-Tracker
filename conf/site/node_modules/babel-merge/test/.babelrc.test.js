module.exports = {
  presets: [
    [
      '@babel/env',
      {
        targets: {
          node: 'current'
        }
      }
    ],
    require.resolve('./local-preset')
  ],
  plugins: [
    require.resolve('./local-plugin'),
    '@babel/plugin-proposal-object-rest-spread',
    ['module:fast-async', {
      runtimePattern: null,
      useRuntimeModule: false
    }]
  ],
  env: {
    test: {
      plugins: [
        [require.resolve('./local-plugin'), { foo: 'bar' }]
      ]
    }
  }
};
