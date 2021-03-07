"use strict";

var _interopRequireDefault = require("@babel/runtime/helpers/interopRequireDefault");

exports.__esModule = true;
exports.default = transitionEndListener;

var _transitionEnd = _interopRequireDefault(require("dom-helpers/transitionEnd"));

function transitionEndListener(element, handler) {
  var remove = (0, _transitionEnd.default)(element, function (e) {
    if (e.target === element) {
      remove();
      handler(e);
    }
  });
}

module.exports = exports["default"];