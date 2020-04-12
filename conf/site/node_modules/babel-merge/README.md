# babel-merge

`babel-merge` merges multiple Babel configuration objects into a single copy.
Plugin and preset objects and arrays will be merged together.

[![NPM version][npm-image]][npm-url] [![NPM downloads][npm-downloads]][npm-url]

_Note: **options** to plugins and presets **will not be merged**, but instead
replaced by the last matching item's options. This makes the behavior consistent
with how Babel works._

## Requirements

- Node.js v6.10+
- Yarn or npm client

## Installation

`babel-merge` can be installed via the Yarn or npm clients.

#### Yarn

```bash
❯ yarn add babel-merge
```

#### npm

```bash
❯ npm install --save babel-merge
```

## Usage

* __merge(a, b, _options_)__
* __merge.all([a, b, ..., z], _options_)__

Where `a`, `b`, `z` are [Babel configuration objects](https://babeljs.io/docs/usage/api/#options) and `options` is a [deepmerge](https://github.com/KyleAMathews/deepmerge#api) options object.

```js
const merge = require('babel-merge');

const together = merge(
  {
    presets: [
      ['@babel/preset-env', {
        targets: {
          browsers: ['latest 1 Chrome']
        }
      }]
    ]
  },
  {
    presets: [
      ['@babel/preset-env', {
        targets: {
          browsers: ['latest 1 Firefox']
        }
      }]
    ]
  }
)

console.log(together);

{
  presets: [
    ['@babel/preset-env', {
      targets: {
        browsers: [
          'latest 1 Firefox'
        ]
      }
    }]
  ]
}
```

If a pathname was used in an earlier merge, you can still merge by exact name:

```js
const merge = require('babel-merge');

const together = merge(
  {
    presets: [
      [require.resolve('@babel/preset-env'), {
        targets: {
          browsers: ['latest 1 Chrome']
        }
      }]
    ]
  },
  {
    presets: [
      ['@babel/preset-env', {
        targets: {
          browsers: ['latest 1 Firefox']
        }
      }]
    ]
  }
)

console.log(together);

{
  presets: [
    ['/Users/me/code/app/node_modules/@babel/preset-env/lib/index.js', {
      targets: {
        browsers: [
          'latest 1 Firefox'
        ]
      }
    }]
  ]
}
```

Even works for plugins and presets within environments:

```js
const merge = require('babel-merge');

const together = merge(
  {
    env: {
      development: {
        presets: [
          [require.resolve('@babel/preset-env'), {
            targets: {
              browsers: ['latest 1 Chrome']
            }
          }]
        ]
      }
    }
  },
  {
    env: {
      development: {
        presets: [
          ['@babel/preset-env', {
            targets: {
              browsers: ['latest 1 Firefox']
            }
          }]
        ]
      }
    }
  }
)

console.log(together);

{
  env: {
    development: {
      presets: [
        ['/Users/me/code/app/node_modules/@babel/preset-env/lib/index.js', {
          targets: {
            browsers: [
              'latest 1 Firefox'
            ]
          }
        }]
      ]
    }
  }
}
```

Order is preserved between non-option plugins and presets and ones with options:

```js
const merge = require('babel-merge');

const together = merge(
  {
    plugins: [
      'module:fast-async',
      '@babel/plugin-syntax-dynamic-import'
    ]
  },
  {
    plugins: [
      '@babel/plugin-proposal-object-rest-spread',
      ['module:fast-async', { spec: true }],
      '@babel/plugin-proposal-class-properties'
    ]
  }
)

console.log(together);

{
  plugins: [
    ['module:fast-async', { spec: true }],
    '@babel/plugin-syntax-dynamic-import',
    '@babel/plugin-proposal-object-rest-spread',
    '@babel/plugin-proposal-class-properties'
  ]
}
```

[npm-image]: https://img.shields.io/npm/v/babel-merge.svg
[npm-downloads]: https://img.shields.io/npm/dt/babel-merge.svg
[npm-url]: https://npmjs.org/package/babel-merge
