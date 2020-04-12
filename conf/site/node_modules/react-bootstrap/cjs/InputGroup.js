"use strict";

var _interopRequireDefault = require("@babel/runtime/helpers/interopRequireDefault");

exports.__esModule = true;
exports.default = void 0;

var _extends2 = _interopRequireDefault(require("@babel/runtime/helpers/extends"));

var _objectWithoutPropertiesLoose2 = _interopRequireDefault(require("@babel/runtime/helpers/objectWithoutPropertiesLoose"));

var _classnames = _interopRequireDefault(require("classnames"));

var _react = _interopRequireDefault(require("react"));

var _createWithBsPrefix = _interopRequireDefault(require("./createWithBsPrefix"));

var _ThemeProvider = require("./ThemeProvider");

/**
 *
 * @property {InputGroupAppend} Append
 * @property {InputGroupPrepend} Prepend
 * @property {InputGroupText} Text
 * @property {InputGroupRadio} Radio
 * @property {InputGroupCheckbox} Checkbox
 */
var InputGroup = _react.default.forwardRef(function (_ref, ref) {
  var bsPrefix = _ref.bsPrefix,
      size = _ref.size,
      className = _ref.className,
      _ref$as = _ref.as,
      Component = _ref$as === void 0 ? 'div' : _ref$as,
      props = (0, _objectWithoutPropertiesLoose2.default)(_ref, ["bsPrefix", "size", "className", "as"]);
  bsPrefix = (0, _ThemeProvider.useBootstrapPrefix)(bsPrefix, 'input-group');
  return _react.default.createElement(Component, (0, _extends2.default)({
    ref: ref
  }, props, {
    className: (0, _classnames.default)(className, bsPrefix, size && bsPrefix + "-" + size)
  }));
});

var InputGroupAppend = (0, _createWithBsPrefix.default)('input-group-append');
var InputGroupPrepend = (0, _createWithBsPrefix.default)('input-group-prepend');
var InputGroupText = (0, _createWithBsPrefix.default)('input-group-text', {
  Component: 'span'
});

var InputGroupCheckbox = function InputGroupCheckbox(props) {
  return _react.default.createElement(InputGroupText, null, _react.default.createElement("input", (0, _extends2.default)({
    type: "checkbox"
  }, props)));
};

var InputGroupRadio = function InputGroupRadio(props) {
  return _react.default.createElement(InputGroupText, null, _react.default.createElement("input", (0, _extends2.default)({
    type: "radio"
  }, props)));
};

InputGroup.displayName = 'InputGroup';
InputGroup.Text = InputGroupText;
InputGroup.Radio = InputGroupRadio;
InputGroup.Checkbox = InputGroupCheckbox;
InputGroup.Append = InputGroupAppend;
InputGroup.Prepend = InputGroupPrepend;
var _default = InputGroup;
exports.default = _default;
module.exports = exports["default"];