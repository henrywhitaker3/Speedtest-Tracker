module.exports = {
  presets: [
    [
      '@babel/env',
      {
        targets: {
          node: 6
        }
      }
    ]
  ],
  env: {
    test: {
      plugins: ['istanbul']
    }
  }
};
