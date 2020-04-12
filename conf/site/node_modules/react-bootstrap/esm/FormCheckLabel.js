import _extends from "@babel/runtime/helpers/esm/extends";
import _objectWithoutPropertiesLoose from "@babel/runtime/helpers/esm/objectWithoutPropertiesLoose";
import classNames from 'classnames';
import React, { useContext } from 'react';
import FormContext from './FormContext';
import { useBootstrapPrefix } from './ThemeProvider';
var FormCheckLabel = React.forwardRef(function (_ref, ref) {
  var bsPrefix = _ref.bsPrefix,
      bsCustomPrefix = _ref.bsCustomPrefix,
      className = _ref.className,
      htmlFor = _ref.htmlFor,
      props = _objectWithoutPropertiesLoose(_ref, ["bsPrefix", "bsCustomPrefix", "className", "htmlFor"]);

  var _useContext = useContext(FormContext),
      controlId = _useContext.controlId,
      custom = _useContext.custom;

  bsPrefix = custom ? useBootstrapPrefix(bsCustomPrefix, 'custom-control-label') : useBootstrapPrefix(bsPrefix, 'form-check-label');
  return React.createElement("label", _extends({}, props, {
    ref: ref,
    htmlFor: htmlFor || controlId,
    className: classNames(className, bsPrefix)
  }));
});
FormCheckLabel.displayName = 'FormCheckLabel';
export default FormCheckLabel;