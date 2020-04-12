'use strict';

/*!
 * concatenate dependencies
 */

var fs = require('fs')
  , globs = require('globs');

/**
 * Simple callback generator for `concatenate()`
 *
 * @api private
 * @param {String} file
 * @return {Function}
 */
function makeCb(file) {
  return function (err, data) {
    if (err) {
      throw err;
    }

    fs.writeFileSync(file, data);
  };
}

/**
 * Concatenate files and such.  Will maintain order of files.
 *
 * Examples:
 *
 * ```js
 * // write all logs to `logs.txt`
 * concatenate([ '*.txt', '*.log' ], 'logs.txt');
 *
 * // output the contents of all local javascript files
 * concatenate('*.js', function (err, js) {
 *   if (err) throw err;
 *
 *   console.log(js);
 * });
 * ```
 *
 * @api public
 * @param {String|Array} patterns Files/wildcard patterns to concatenate
 * @param {String|Function} cb File to write or callback
 */
var concatenate = module.exports = function (patterns, cb) {
  // concatenate('*.txt', function () { ... })
  if (!Array.isArray(patterns)) {
    patterns = [ patterns ];
  }

  // concatenate([ '*.txt' ], 'all-text.js')
  if (typeof cb === 'string') {
    cb = makeCb(cb);
  }

  globs(patterns, function (err, files) {
    if (err) {
      return cb(err);
    }

    var done = false
      , res = [];

    /**
     * Add the `file` at the `index` to the stack
     *
     * Will add the next `file` at `index++`, or fire
     * `cb` if there is not another `file`.
     *
     * This is rather ugly, but maintains order
     * of the provided files.
     *
     * @api private
     * @param {Number} index
     */
    function add(index) {
      var file = files[index];
      if (!file) {
        done = true;
        return cb(new Error('Invalid file "' + file + '"'));
      }
      fs.readFile(files[index], 'utf-8', function (err, data) {
        if (done) { return; }
        if (err) {
          done = true;
          return cb(err);
        }

        res.push(data);
        index++;
        if (!files[index]) {
          done = true;
          cb(null, res.join('\n'));
        } else {
          add(index);
        }
      });
    }

    // start the loop by adding the first file
    add(0);
  });
};

/**
 * Synchronously concatenate files and such
 *
 * @api public
 * @param {String|Array} patterns
 * @param {String} [out] Output file
 */
concatenate.sync = function (patterns, out) {
  var index
    , files = globs.sync(patterns)
    , length = files.length
    , all = [];

  for (index = 0; index < length; index++) {
    all.push(fs.readFileSync(files[index], 'utf-8'));
  }
  all = all.join('\n');

  if (out) {
    // write file
    fs.writeFileSync(out, all);
  }
  // return concatenated stuff
  return all;
};
