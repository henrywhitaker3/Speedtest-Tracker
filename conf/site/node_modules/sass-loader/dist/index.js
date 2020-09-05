"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.default = void 0;

var _path = _interopRequireDefault(require("path"));

var _schemaUtils = _interopRequireDefault(require("schema-utils"));

var _loaderUtils = require("loader-utils");

var _options = _interopRequireDefault(require("./options.json"));

var _utils = require("./utils");

var _SassError = _interopRequireDefault(require("./SassError"));

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

/**
 * The sass-loader makes node-sass and dart-sass available to webpack modules.
 *
 * @this {object}
 * @param {string} content
 */
function loader(content) {
  const options = (0, _loaderUtils.getOptions)(this);
  (0, _schemaUtils.default)(_options.default, options, {
    name: 'Sass Loader',
    baseDataPath: 'options'
  });
  const implementation = (0, _utils.getSassImplementation)(options.implementation);
  const useSourceMap = typeof options.sourceMap === 'boolean' ? options.sourceMap : this.sourceMap;
  const sassOptions = (0, _utils.getSassOptions)(this, options, content, implementation, useSourceMap);
  const shouldUseWebpackImporter = typeof options.webpackImporter === 'boolean' ? options.webpackImporter : true;

  if (shouldUseWebpackImporter) {
    const {
      includePaths
    } = sassOptions;
    sassOptions.importer.push((0, _utils.getWebpackImporter)(this, implementation, includePaths));
  }

  const callback = this.async();
  const render = (0, _utils.getRenderFunctionFromSassImplementation)(implementation);
  render(sassOptions, (error, result) => {
    if (error) {
      // There are situations when the `file` property do not exist
      if (error.file) {
        // `node-sass` returns POSIX paths
        this.addDependency(_path.default.normalize(error.file));
      }

      callback(new _SassError.default(error));
      return;
    } // Modify source paths only for webpack, otherwise we do nothing


    if (result.map && useSourceMap) {
      // eslint-disable-next-line no-param-reassign
      result.map = JSON.parse(result.map); // result.map.file is an optional property that provides the output filename.
      // Since we don't know the final filename in the webpack build chain yet, it makes no sense to have it.
      // eslint-disable-next-line no-param-reassign

      delete result.map.file; // eslint-disable-next-line no-param-reassign

      result.sourceRoot = ''; // node-sass returns POSIX paths, that's why we need to transform them back to native paths.
      // This fixes an error on windows where the source-map module cannot resolve the source maps.
      // @see https://github.com/webpack-contrib/sass-loader/issues/366#issuecomment-279460722
      // eslint-disable-next-line no-param-reassign

      result.map.sources = result.map.sources.map(source => (0, _utils.absolutifySourceMapSource)(this.rootContext, source));
    }

    result.stats.includedFiles.forEach(includedFile => {
      this.addDependency(_path.default.normalize(includedFile));
    });
    callback(null, result.css.toString(), result.map);
  });
}

var _default = loader;
exports.default = _default;