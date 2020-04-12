import _extends from "@babel/runtime/helpers/esm/extends";
import _objectWithoutPropertiesLoose from "@babel/runtime/helpers/esm/objectWithoutPropertiesLoose";
import classNames from 'classnames';
import React, { useContext, useMemo } from 'react';
import all from 'prop-types-extra/lib/all';
import Feedback from './Feedback';
import FormFileInput from './FormFileInput';
import FormFileLabel from './FormFileLabel';
import FormContext from './FormContext';
import { useBootstrapPrefix } from './ThemeProvider';
var defaultProps = {
  disabled: false,
  isValid: false,
  isInvalid: false
};
var FormFile = React.forwardRef(function (_ref, ref) {
  var id = _ref.id,
      bsPrefix = _ref.bsPrefix,
      bsCustomPrefix = _ref.bsCustomPrefix,
      disabled = _ref.disabled,
      isValid = _ref.isValid,
      isInvalid = _ref.isInvalid,
      feedback = _ref.feedback,
      className = _ref.className,
      style = _ref.style,
      label = _ref.label,
      children = _ref.children,
      custom = _ref.custom,
      lang = _ref.lang,
      dataBrowse = _ref['data-browse'],
      _ref$as = _ref.as,
      Component = _ref$as === void 0 ? 'div' : _ref$as,
      _ref$inputAs = _ref.inputAs,
      inputAs = _ref$inputAs === void 0 ? 'input' : _ref$inputAs,
      props = _objectWithoutPropertiesLoose(_ref, ["id", "bsPrefix", "bsCustomPrefix", "disabled", "isValid", "isInvalid", "feedback", "className", "style", "label", "children", "custom", "lang", "data-browse", "as", "inputAs"]);

  bsPrefix = custom ? useBootstrapPrefix(bsCustomPrefix, 'custom') : useBootstrapPrefix(bsPrefix, 'form-file');
  var type = 'file';

  var _useContext = useContext(FormContext),
      controlId = _useContext.controlId;

  var innerFormContext = useMemo(function () {
    return {
      controlId: id || controlId,
      custom: custom
    };
  }, [controlId, custom, id]);
  var hasLabel = label != null && label !== false && !children;
  var input = React.createElement(FormFileInput, _extends({}, props, {
    ref: ref,
    isValid: isValid,
    isInvalid: isInvalid,
    disabled: disabled,
    as: inputAs,
    lang: lang
  }));
  return React.createElement(FormContext.Provider, {
    value: innerFormContext
  }, React.createElement(Component, {
    style: style,
    className: classNames(className, bsPrefix, custom && "custom-" + type)
  }, children || React.createElement(React.Fragment, null, custom ? React.createElement(React.Fragment, null, input, hasLabel && React.createElement(FormFileLabel, {
    "data-browse": dataBrowse
  }, label)) : React.createElement(React.Fragment, null, hasLabel && React.createElement(FormFileLabel, null, label), input), (isValid || isInvalid) && React.createElement(Feedback, {
    type: isValid ? 'valid' : 'invalid'
  }, feedback))));
});
FormFile.displayName = 'FormFile';
FormFile.defaultProps = defaultProps;
FormFile.Input = FormFileInput;
FormFile.Label = FormFileLabel;
export default FormFile;