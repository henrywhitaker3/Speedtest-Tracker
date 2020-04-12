import _extends from "@babel/runtime/helpers/esm/extends";
import _objectWithoutPropertiesLoose from "@babel/runtime/helpers/esm/objectWithoutPropertiesLoose";
import classNames from 'classnames';
import React from 'react';
import createWithBsPrefix from './createWithBsPrefix';
import { useBootstrapPrefix } from './ThemeProvider';

/**
 *
 * @property {InputGroupAppend} Append
 * @property {InputGroupPrepend} Prepend
 * @property {InputGroupText} Text
 * @property {InputGroupRadio} Radio
 * @property {InputGroupCheckbox} Checkbox
 */
var InputGroup = React.forwardRef(function (_ref, ref) {
  var bsPrefix = _ref.bsPrefix,
      size = _ref.size,
      className = _ref.className,
      _ref$as = _ref.as,
      Component = _ref$as === void 0 ? 'div' : _ref$as,
      props = _objectWithoutPropertiesLoose(_ref, ["bsPrefix", "size", "className", "as"]);

  bsPrefix = useBootstrapPrefix(bsPrefix, 'input-group');
  return React.createElement(Component, _extends({
    ref: ref
  }, props, {
    className: classNames(className, bsPrefix, size && bsPrefix + "-" + size)
  }));
});
var InputGroupAppend = createWithBsPrefix('input-group-append');
var InputGroupPrepend = createWithBsPrefix('input-group-prepend');
var InputGroupText = createWithBsPrefix('input-group-text', {
  Component: 'span'
});

var InputGroupCheckbox = function InputGroupCheckbox(props) {
  return React.createElement(InputGroupText, null, React.createElement("input", _extends({
    type: "checkbox"
  }, props)));
};

var InputGroupRadio = function InputGroupRadio(props) {
  return React.createElement(InputGroupText, null, React.createElement("input", _extends({
    type: "radio"
  }, props)));
};

InputGroup.displayName = 'InputGroup';
InputGroup.Text = InputGroupText;
InputGroup.Radio = InputGroupRadio;
InputGroup.Checkbox = InputGroupCheckbox;
InputGroup.Append = InputGroupAppend;
InputGroup.Prepend = InputGroupPrepend;
export default InputGroup;