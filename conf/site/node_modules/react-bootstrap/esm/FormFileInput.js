import _extends from "@babel/runtime/helpers/esm/extends";
import _objectWithoutPropertiesLoose from "@babel/runtime/helpers/esm/objectWithoutPropertiesLoose";
import classNames from 'classnames';
import React, { useContext } from 'react';
import FormContext from './FormContext';
import { useBootstrapPrefix } from './ThemeProvider';
var FormFileInput = React.forwardRef(function (_ref, ref) {
  var id = _ref.id,
      bsPrefix = _ref.bsPrefix,
      bsCustomPrefix = _ref.bsCustomPrefix,
      className = _ref.className,
      isValid = _ref.isValid,
      isInvalid = _ref.isInvalid,
      lang = _ref.lang,
      _ref$as = _ref.as,
      Component = _ref$as === void 0 ? 'input' : _ref$as,
      props = _objectWithoutPropertiesLoose(_ref, ["id", "bsPrefix", "bsCustomPrefix", "className", "isValid", "isInvalid", "lang", "as"]);

  var _useContext = useContext(FormContext),
      controlId = _useContext.controlId,
      custom = _useContext.custom;

  var type = 'file';
  bsPrefix = custom ? useBootstrapPrefix(bsCustomPrefix, 'custom-file-input') : useBootstrapPrefix(bsPrefix, 'form-control-file');
  return React.createElement(Component, _extends({}, props, {
    ref: ref,
    id: id || controlId,
    type: type,
    lang: lang,
    className: classNames(className, bsPrefix, isValid && 'is-valid', isInvalid && 'is-invalid')
  }));
});
FormFileInput.displayName = 'FormFileInput';
export default FormFileInput;