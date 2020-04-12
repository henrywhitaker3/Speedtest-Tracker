import _extends from "@babel/runtime/helpers/esm/extends";
import _objectWithoutPropertiesLoose from "@babel/runtime/helpers/esm/objectWithoutPropertiesLoose";
import classNames from 'classnames';
import React from 'react';
import { useUncontrolled } from 'uncontrollable';
import useEventCallback from '@restart/hooks/useEventCallback';
import createWithBsPrefix from './createWithBsPrefix';
import divWithClassName from './divWithClassName';
import { useBootstrapPrefix } from './ThemeProvider';
import Fade from './Fade';
import CloseButton from './CloseButton';
import SafeAnchor from './SafeAnchor';
var defaultProps = {
  show: true,
  transition: Fade,
  closeLabel: 'Close alert'
};
var controllables = {
  show: 'onClose'
};
var Alert = React.forwardRef(function (uncontrolledProps, ref) {
  var _useUncontrolled = useUncontrolled(uncontrolledProps, controllables),
      bsPrefix = _useUncontrolled.bsPrefix,
      show = _useUncontrolled.show,
      closeLabel = _useUncontrolled.closeLabel,
      className = _useUncontrolled.className,
      children = _useUncontrolled.children,
      variant = _useUncontrolled.variant,
      onClose = _useUncontrolled.onClose,
      dismissible = _useUncontrolled.dismissible,
      Transition = _useUncontrolled.transition,
      props = _objectWithoutPropertiesLoose(_useUncontrolled, ["bsPrefix", "show", "closeLabel", "className", "children", "variant", "onClose", "dismissible", "transition"]);

  var prefix = useBootstrapPrefix(bsPrefix, 'alert');
  var handleClose = useEventCallback(function (e) {
    onClose(false, e);
  });
  var alert = React.createElement("div", _extends({
    role: "alert"
  }, Transition ? props : undefined, {
    ref: ref,
    className: classNames(className, prefix, variant && prefix + "-" + variant, dismissible && prefix + "-dismissible")
  }), dismissible && React.createElement(CloseButton, {
    onClick: handleClose,
    label: closeLabel
  }), children);
  if (!Transition) return show ? alert : null;
  return React.createElement(Transition, _extends({
    unmountOnExit: true
  }, props, {
    in: show
  }), alert);
});
var DivStyledAsH4 = divWithClassName('h4');
DivStyledAsH4.displayName = 'DivStyledAsH4';
Alert.displayName = 'Alert';
Alert.defaultProps = defaultProps;
Alert.Link = createWithBsPrefix('alert-link', {
  Component: SafeAnchor
});
Alert.Heading = createWithBsPrefix('alert-heading', {
  Component: DivStyledAsH4
});
export default Alert;