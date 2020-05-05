"use strict";

var _interopRequireWildcard = require("@babel/runtime/helpers/interopRequireWildcard");

var _interopRequireDefault = require("@babel/runtime/helpers/interopRequireDefault");

exports.__esModule = true;
exports.default = void 0;

var _extends2 = _interopRequireDefault(require("@babel/runtime/helpers/extends"));

var _objectWithoutPropertiesLoose2 = _interopRequireDefault(require("@babel/runtime/helpers/objectWithoutPropertiesLoose"));

var _react = _interopRequireWildcard(require("react"));

var _classnames = _interopRequireDefault(require("classnames"));

var _useTimeout = _interopRequireDefault(require("@restart/hooks/useTimeout"));

var _Fade = _interopRequireDefault(require("./Fade"));

var _ToastHeader = _interopRequireDefault(require("./ToastHeader"));

var _ToastBody = _interopRequireDefault(require("./ToastBody"));

var _ThemeProvider = require("./ThemeProvider");

var _ToastContext = _interopRequireDefault(require("./ToastContext"));

var defaultProps = {
  animation: true,
  autohide: false,
  delay: 3000,
  show: true,
  transition: _Fade.default
};

var Toast = _react.default.forwardRef(function (_ref, ref) {
  var bsPrefix = _ref.bsPrefix,
      className = _ref.className,
      children = _ref.children,
      Transition = _ref.transition,
      show = _ref.show,
      animation = _ref.animation,
      delay = _ref.delay,
      autohide = _ref.autohide,
      onClose = _ref.onClose,
      props = (0, _objectWithoutPropertiesLoose2.default)(_ref, ["bsPrefix", "className", "children", "transition", "show", "animation", "delay", "autohide", "onClose"]);
  bsPrefix = (0, _ThemeProvider.useBootstrapPrefix)('toast');
  var delayRef = (0, _react.useRef)(delay);
  var onCloseRef = (0, _react.useRef)(onClose);
  (0, _react.useEffect)(function () {
    // We use refs for these, because we don't want to restart the autohide
    // timer in case these values change.
    delayRef.current = delay;
    onCloseRef.current = onClose;
  }, [delay, onClose]);
  var autohideTimeout = (0, _useTimeout.default)();
  var autohideFunc = (0, _react.useCallback)(function () {
    if (!(autohide && show)) {
      return;
    }

    onCloseRef.current();
  }, [autohide, show]);
  autohideTimeout.set(autohideFunc, delayRef.current);
  var hasAnimation = (0, _react.useMemo)(function () {
    return Transition && animation;
  }, [Transition, animation]);

  var toast = /*#__PURE__*/_react.default.createElement("div", (0, _extends2.default)({}, props, {
    ref: ref,
    className: (0, _classnames.default)(bsPrefix, className, !hasAnimation && (show ? 'show' : 'hide')),
    role: "alert",
    "aria-live": "assertive",
    "aria-atomic": "true"
  }), children);

  var toastContext = {
    onClose: onClose
  };
  return /*#__PURE__*/_react.default.createElement(_ToastContext.default.Provider, {
    value: toastContext
  }, hasAnimation ? /*#__PURE__*/_react.default.createElement(Transition, {
    in: show,
    unmountOnExit: true
  }, toast) : toast);
});

Toast.defaultProps = defaultProps;
Toast.displayName = 'Toast';
Toast.Body = _ToastBody.default;
Toast.Header = _ToastHeader.default;
var _default = Toast;
exports.default = _default;
module.exports = exports["default"];