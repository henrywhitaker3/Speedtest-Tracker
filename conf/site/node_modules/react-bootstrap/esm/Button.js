import _extends from "@babel/runtime/helpers/esm/extends";
import _objectWithoutPropertiesLoose from "@babel/runtime/helpers/esm/objectWithoutPropertiesLoose";
import classNames from 'classnames';
import React from 'react';
import { useBootstrapPrefix } from './ThemeProvider';
import SafeAnchor from './SafeAnchor';
var defaultProps = {
  variant: 'primary',
  active: false,
  disabled: false,
  type: 'button'
};
var Button = React.forwardRef(function (_ref, ref) {
  var bsPrefix = _ref.bsPrefix,
      variant = _ref.variant,
      size = _ref.size,
      active = _ref.active,
      className = _ref.className,
      block = _ref.block,
      type = _ref.type,
      as = _ref.as,
      props = _objectWithoutPropertiesLoose(_ref, ["bsPrefix", "variant", "size", "active", "className", "block", "type", "as"]);

  var prefix = useBootstrapPrefix(bsPrefix, 'btn');
  var classes = classNames(className, prefix, active && 'active', prefix + "-" + variant, block && prefix + "-block", size && prefix + "-" + size);

  if (props.href) {
    return React.createElement(SafeAnchor, _extends({}, props, {
      as: as,
      ref: ref,
      className: classNames(classes, props.disabled && 'disabled')
    }));
  }

  if (ref) {
    props.ref = ref;
  }

  if (!as) {
    props.type = type;
  }

  var Component = as || 'button';
  return React.createElement(Component, _extends({}, props, {
    className: classes
  }));
});
Button.displayName = 'Button';
Button.defaultProps = defaultProps;
export default Button;