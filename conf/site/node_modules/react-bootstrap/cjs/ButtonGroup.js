"use strict";

var _interopRequireDefault = require("@babel/runtime/helpers/interopRequireDefault");

exports.__esModule = true;
exports.default = void 0;

var _extends2 = _interopRequireDefault(require("@babel/runtime/helpers/extends"));

var _objectWithoutPropertiesLoose2 = _interopRequireDefault(require("@babel/runtime/helpers/objectWithoutPropertiesLoose"));

var _classnames = _interopRequireDefault(require("classnames"));

var _react = _interopRequireDefault(require("react"));

var _ThemeProvider = require("./ThemeProvider");

var defaultProps = {
  vertical: false,
  toggle: false,
  role: 'group'
};

var ButtonGroup = _react.default.forwardRef(function (props, ref) {
  var bsPrefix = props.bsPrefix,
      size = props.size,
      toggle = props.toggle,
      vertical = props.vertical,
      className = props.className,
      _props$as = props.as,
      Component = _props$as === void 0 ? 'div' : _props$as,
      rest = (0, _objectWithoutPropertiesLoose2.default)(props, ["bsPrefix", "size", "toggle", "vertical", "className", "as"]);
  var prefix = (0, _ThemeProvider.useBootstrapPrefix)(bsPrefix, 'btn-group');
  var baseClass = prefix;
  if (vertical) baseClass = prefix + "-vertical";
  return _react.default.createElement(Component, (0, _extends2.default)({}, rest, {
    ref: ref,
    className: (0, _classnames.default)(className, baseClass, size && prefix + "-" + size, toggle && prefix + "-toggle")
  }));
});

ButtonGroup.displayName = 'ButtonGroup';
ButtonGroup.defaultProps = defaultProps;
var _default = ButtonGroup;
exports.default = _default;
module.exports = exports["default"];