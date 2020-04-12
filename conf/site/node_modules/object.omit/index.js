/*!
 * object.omit <https://github.com/jonschlinkert/object.omit>
 *
 * Copyright (c) 2014-2017, Jon Schlinkert.
 * Released under the MIT License.
 */

'use strict';

var isObject = require('is-extendable');

module.exports = function omit(obj, props, fn) {
  if (!isObject(obj)) return {};

  if (typeof props === 'function') {
    fn = props;
    props = [];
  }

  if (typeof props === 'string') {
    props = [props];
  }

  var isFunction = typeof fn === 'function';
  var keys = Object.keys(obj);
  var res = {};

  for (var i = 0; i < keys.length; i++) {
    var key = keys[i];
    var val = obj[key];

    if (!props || (props.indexOf(key) === -1 && (!isFunction || fn(val, key, obj)))) {
      res[key] = val;
    }
  }
  return res;
};
