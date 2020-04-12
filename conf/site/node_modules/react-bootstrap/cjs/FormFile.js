"use strict";

var _interopRequireWildcard = require("@babel/runtime/helpers/interopRequireWildcard");

var _interopRequireDefault = require("@babel/runtime/helpers/interopRequireDefault");

exports.__esModule = true;
exports.default = void 0;

var _extends2 = _interopRequireDefault(require("@babel/runtime/helpers/extends"));

var _objectWithoutPropertiesLoose2 = _interopRequireDefault(require("@babel/runtime/helpers/objectWithoutPropertiesLoose"));

var _classnames = _interopRequireDefault(require("classnames"));

var _react = _interopRequireWildcard(require("react"));

var _all = _interopRequireDefault(require("prop-types-extra/lib/all"));

var _Feedback = _interopRequireDefault(require("./Feedback"));

var _FormFileInput = _interopRequireDefault(require("./FormFileInput"));

var _FormFileLabel = _interopRequireDefault(require("./FormFileLabel"));

var _FormContext = _interopRequireDefault(require("./FormContext"));

var _ThemeProvider = require("./ThemeProvider");

var defaultProps = {
  disabled: false,
  isValid: false,
  isInvalid: false
};

var FormFile = _react.default.forwardRef(function (_ref, ref) {
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
      props = (0, _objectWithoutPropertiesLoose2.default)(_ref, ["id", "bsPrefix", "bsCustomPrefix", "disabled", "isValid", "isInvalid", "feedback", "className", "style", "label", "children", "custom", "lang", "data-browse", "as", "inputAs"]);
  bsPrefix = custom ? (0, _ThemeProvider.useBootstrapPrefix)(bsCustomPrefix, 'custom') : (0, _ThemeProvider.useBootstrapPrefix)(bsPrefix, 'form-file');
  var type = 'file';

  var _useContext = (0, _react.useContext)(_FormContext.default),
      controlId = _useContext.controlId;

  var innerFormContext = (0, _react.useMemo)(function () {
    return {
      controlId: id || controlId,
      custom: custom
    };
  }, [controlId, custom, id]);
  var hasLabel = label != null && label !== false && !children;

  var input = _react.default.createElement(_FormFileInput.default, (0, _extends2.default)({}, props, {
    ref: ref,
    isValid: isValid,
    isInvalid: isInvalid,
    disabled: disabled,
    as: inputAs,
    lang: lang
  }));

  return _react.default.createElement(_FormContext.default.Provider, {
    value: innerFormContext
  }, _react.default.createElement(Component, {
    style: style,
    className: (0, _classnames.default)(className, bsPrefix, custom && "custom-" + type)
  }, children || _react.default.createElement(_react.default.Fragment, null, custom ? _react.default.createElement(_react.default.Fragment, null, input, hasLabel && _react.default.createElement(_FormFileLabel.default, {
    "data-browse": dataBrowse
  }, label)) : _react.default.createElement(_react.default.Fragment, null, hasLabel && _react.default.createElement(_FormFileLabel.default, null, label), input), (isValid || isInvalid) && _react.default.createElement(_Feedback.default, {
    type: isValid ? 'valid' : 'invalid'
  }, feedback))));
});

FormFile.displayName = 'FormFile';
FormFile.defaultProps = defaultProps;
FormFile.Input = _FormFileInput.default;
FormFile.Label = _FormFileLabel.default;
var _default = FormFile;
exports.default = _default;
module.exports = exports["default"];