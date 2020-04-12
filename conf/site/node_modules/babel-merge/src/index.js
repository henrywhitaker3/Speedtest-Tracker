import omit from 'object.omit';
import merge from 'deepmerge';
import { resolvePlugin, resolvePreset } from '@babel/core';

function arrayMerge(source = [], overrides = []) {
  return [...new Set([...source, ...overrides])];
}

function mergeArray(source = [], overrides = [], resolve, deepmergeOpts) {
  return [...source, ...overrides].reduce((reduction, override) => {
    const overrideName = resolve(Array.isArray(override) ? override[0] : override);
    const overrideOptions = Array.isArray(override) ? override[1] : {};
    const base = reduction.find((base) => {
      const baseName = resolve(Array.isArray(base) ? base[0] : base);
      return baseName === overrideName || baseName.includes(overrideName);
    });

    const index = reduction.includes(base) ? reduction.indexOf(base) : reduction.length;
    const baseName = base ? resolve(Array.isArray(base) ? base[0] : base) : overrideName;
    const baseOptions = Array.isArray(base) ? base[1] : {};
    const options = merge(baseOptions, overrideOptions, {
      arrayMerge,
      isMergeableObject: value => Array.isArray(value),
      ...deepmergeOpts
    });

    reduction[index] = Object.keys(options).length ? [baseName, options] : baseName;

    return reduction;
  }, []);
}

function babelMerge(source = {}, overrides = {}, deepmergeOpts) {
  const plugins = mergeArray(source.plugins, overrides.plugins, resolvePlugin, deepmergeOpts);
  const presets = mergeArray(source.presets, overrides.presets, resolvePreset, deepmergeOpts);
  const sourceEnv = source.env || {};
  const overridesEnv = overrides.env || {};
  return Object.assign(
    presets.length ? { presets } : {},
    plugins.length ? { plugins } : {},
    merge.all([
      omit(source, ['plugins', 'presets', 'env']),
      omit(overrides, ['plugins', 'presets', 'env']),
      ...[...new Set([
        ...Object.keys(sourceEnv),
        ...Object.keys(overridesEnv)
      ])].map(name => ({
        env: {
          [name]: babelMerge(sourceEnv[name], overridesEnv[name], deepmergeOpts)
        }
      }))
    ], { arrayMerge, ...deepmergeOpts })
  );
}

Object.defineProperty(babelMerge, 'all', {
  value: (values = [], deepmergeOpts) =>
    values.reduce((acc, value) => {
      if (value) {
        Object.assign(acc, babelMerge(acc, value, deepmergeOpts));
      }
      return acc;
    }, {})
});

module.exports = babelMerge;
