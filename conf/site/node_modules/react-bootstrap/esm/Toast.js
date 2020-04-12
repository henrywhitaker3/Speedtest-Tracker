import _extends from "@babel/runtime/helpers/esm/extends";
import _objectWithoutPropertiesLoose from "@babel/runtime/helpers/esm/objectWithoutPropertiesLoose";
import React, { useEffect, useRef, useMemo, useCallback } from 'react';
import classNames from 'classnames';
import useTimeout from '@restart/hooks/useTimeout';
import Fade from './Fade';
import Header from './ToastHeader';
import Body from './ToastBody';
import { useBootstrapPrefix } from './ThemeProvider';
import ToastContext from './ToastContext';
var defaultProps = {
  animation: true,
  autohide: false,
  delay: 3000,
  show: true,
  transition: Fade
};
var Toast = React.forwardRef(function (_ref, ref) {
  var bsPrefix = _ref.bsPrefix,
      className = _ref.className,
      children = _ref.children,
      Transition = _ref.transition,
      show = _ref.show,
      animation = _ref.animation,
      delay = _ref.delay,
      autohide = _ref.autohide,
      onClose = _ref.onClose,
      props = _objectWithoutPropertiesLoose(_ref, ["bsPrefix", "className", "children", "transition", "show", "animation", "delay", "autohide", "onClose"]);

  bsPrefix = useBootstrapPrefix('toast');
  var delayRef = useRef(delay);
  var onCloseRef = useRef(onClose);
  useEffect(function () {
    // We use refs for these, because we don't want to restart the autohide
    // timer in case these values change.
    delayRef.current = delay;
    onCloseRef.current = onClose;
  }, [delay, onClose]);
  var autohideTimeout = useTimeout();
  var autohideFunc = useCallback(function () {
    if (!(autohide && show)) {
      return;
    }

    onCloseRef.current();
  }, [autohide, show]);
  autohideTimeout.set(autohideFunc, delayRef.current);
  var useAnimation = useMemo(function () {
    return Transition && animation;
  }, [Transition, animation]);
  var toast = React.createElement("div", _extends({}, props, {
    ref: ref,
    className: classNames(bsPrefix, className, !useAnimation && show && 'show'),
    role: "alert",
    "aria-live": "assertive",
    "aria-atomic": "true"
  }), children);
  var toastContext = {
    onClose: onClose
  };
  return React.createElement(ToastContext.Provider, {
    value: toastContext
  }, useAnimation ? React.createElement(Transition, {
    in: show
  }, toast) : toast);
});
Toast.defaultProps = defaultProps;
Toast.displayName = 'Toast';
Toast.Body = Body;
Toast.Header = Header;
export default Toast;