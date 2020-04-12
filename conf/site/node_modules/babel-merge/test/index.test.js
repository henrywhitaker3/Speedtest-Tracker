import assert from 'assert';

import * as babel from '@babel/core';
import omit from 'object.omit';

import babelMerge from '../src';

function formatBabelConfig({ file, options }) {
  return options ? [file.resolved, options] : file.resolved;
}

describe('babel-merge', () => {
  it('should deeply merge preset options', () => {
    assert.deepStrictEqual(
      babelMerge(
        {
          presets: [
            ['@babel/env', {
              targets: {
                browsers: [
                  'latest 1 Chrome'
                ]
              }
            }]
          ]
        },
        {
          presets: [
            ['@babel/env', {
              targets: {
                browsers: [
                  'latest 1 Firefox'
                ]
              }
            }]
          ]
        }
      ),
      {
        presets: [
          [require.resolve('@babel/preset-env'), {
            targets: {
              browsers: [
                'latest 1 Firefox'
              ]
            }
          }]
        ]
      }
    );
  });

  it('should merge by resolved name', () => {
    assert.deepStrictEqual(
      babelMerge(
        {
          presets: [
            [require.resolve('@babel/preset-env'), {
              targets: {
                browsers: [
                  'latest 1 Chrome'
                ]
              }
            }]
          ]
        },
        {
          presets: [
            ['@babel/env', {
              targets: {
                browsers: [
                  'latest 1 Firefox'
                ]
              }
            }]
          ]
        }
      ),
      {
        presets: [
          [require.resolve('@babel/preset-env'), {
            targets: {
              browsers: [
                'latest 1 Firefox'
              ]
            }
          }]
        ]
      }
    );
  });

  it('should merge env options', () => {
    assert.deepStrictEqual(
      babelMerge(
        {
          env: {
            development: {
              presets: [
                [require.resolve('@babel/preset-env'), {
                  targets: {
                    browsers: [
                      'latest 1 Chrome'
                    ]
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
                ['@babel/env', {
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
      ),
      {
        env: {
          development: {
            presets: [
              [require.resolve('@babel/preset-env'), {
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
    );
  });

  it('should preserve plugin / preset order', () => {
    assert.deepStrictEqual(
      babelMerge(
        {
          presets: [
            './test/local-preset'
          ],
          plugins: [
            'module:fast-async',
            '@babel/plugin-syntax-dynamic-import',
            './test/local-plugin'
          ]
        },
        {
          presets: [
            '@babel/env'
          ],
          plugins: [
            ['./test/local-plugin', { foo: 'bar' }],
            '@babel/plugin-proposal-object-rest-spread',
            ['module:fast-async', { spec: true }],
            '@babel/plugin-proposal-class-properties'
          ]
        }
      ),
      {
        presets: [
          require.resolve('./local-preset'),
          require.resolve('@babel/preset-env')
        ],
        plugins: [
          [require.resolve('fast-async'), { 'spec': true }],
          require.resolve('@babel/plugin-syntax-dynamic-import'),
          [require.resolve('./local-plugin'), { foo: 'bar' }],
          require.resolve('@babel/plugin-proposal-object-rest-spread'),
          require.resolve('@babel/plugin-proposal-class-properties')
        ]
      }
    );
  });

  it('should merge an array of config objects', () => {
    assert.deepStrictEqual(
      babelMerge.all([
        {
          presets: [
            require.resolve('@babel/preset-env')
          ]
        },
        {
          presets: [
            '@babel/preset-env'
          ]
        },
        {
          presets: [
            '@babel/env'
          ]
        }
      ]),
      {
        presets: [
          require.resolve('@babel/preset-env')
        ]
      }
    );
  });

  it('should dedupe merged arrays', () => {
    assert.deepStrictEqual(
      babelMerge.all([
        {
          presets: [
            [require.resolve('@babel/preset-env'), {
              targets: {
                browsers: [
                  'latest 1 Chrome'
                ]
              }
            }]
          ]
        },
        {
          presets: [
            ['@babel/preset-env', {
              targets: {
                browsers: [
                  'latest 1 Chrome'
                ]
              }
            }]
          ]
        },
        {
          presets: [
            ['@babel/env', {
              targets: {
                browsers: [
                  'latest 1 Chrome'
                ]
              }
            }]
          ]
        }
      ]),
      {
        presets: [
          [require.resolve('@babel/preset-env'), {
            targets: {
              browsers: [
                'latest 1 Chrome'
              ]
            }
          }]
        ]
      }
    );
  });

  it('should support ES6+ data structures', () => {
    const a = {
      Map: new Map([['a', 'a']]),
      Set: new Set(['a']),
      WeakMap: new WeakMap([[{ a: true }, 'a']]),
      WeakSet: new WeakSet([{ a: true }])
    };

    const b = {
      Map: new Map([['b', 'b']]),
      Set: new Set(['b']),
      WeakMap: new WeakMap([[{ b: true }, 'b']]),
      WeakSet: new WeakSet([{ b: true }])
    };

    const c = {
      Map: new Map([['c', 'c']]),
      Set: new Set(['c']),
      WeakMap: new WeakMap([[{ c: true }, 'c']]),
      WeakSet: new WeakSet([{ c: true }])
    };

    assert.deepStrictEqual(
      babelMerge.all([
        { presets: [[require.resolve('@babel/preset-env'), a]] },
        { presets: [['@babel/preset-env', b]] },
        { presets: [['@babel/env', c]] }
      ]),
      {
        presets: [
          [require.resolve('@babel/preset-env'), c]
        ]
      }
    );
  });

  it('should support deepmerge option overrides', () => {
    assert.deepStrictEqual(
      babelMerge(
        {
          presets: [
            ['@babel/env', {
              targets: {
                browsers: new Set()
              }
            }]
          ]
        },
        undefined,
        { isMergeableObject: () => true }
      ),
      {
        presets: [
          [require.resolve('@babel/preset-env'), {
            targets: {
              browsers: {}
            }
          }]
        ]
      }
    );

    assert.deepStrictEqual(
      babelMerge.all(
        [{
          presets: [
            ['@babel/env', {
              targets: {
                browsers: new Set()
              }
            }]
          ]
        }],
        { isMergeableObject: () => true }
      ),
      {
        presets: [
          [require.resolve('@babel/preset-env'), {
            targets: {
              browsers: {}
            }
          }]
        ]
      }
    );
  });

  it("should mirror babel's merge behavior", () => {
    function getOverrides() {
      return {
        presets: [
          ['./test/local-preset', { foo: 'bar' }],
          [
            '@babel/env',
            {
              targets: {
                browsers: ['>= 0.25%', 'not dead']
              }
            }
          ]
        ],
        plugins: [
          '@babel/plugin-proposal-object-rest-spread',
          ['module:fast-async', { spec: true }],
          '@babel/plugin-proposal-class-properties'
        ]
      };
    }

    const { options: { presets, plugins } } = babel.loadPartialConfig({
      ...getOverrides(),
      configFile: require.resolve('./.babelrc.test')
    });

    delete require.cache[require.resolve('./.babelrc.test')];

    const babelrc = require('./.babelrc.test');

    assert.deepStrictEqual(
      {
        presets: presets.map(formatBabelConfig),
        plugins: plugins.map(formatBabelConfig)
      },
      omit(babelMerge.all([babelrc, babelrc.env.test, getOverrides()]), ['env'])
    );
  });
});
