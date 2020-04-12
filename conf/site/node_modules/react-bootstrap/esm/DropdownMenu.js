import _extends from "@babel/runtime/helpers/esm/extends";
import _objectWithoutPropertiesLoose from "@babel/runtime/helpers/esm/objectWithoutPropertiesLoose";
import classNames from 'classnames';
import React, { useContext } from 'react';
import { useDropdownMenu } from 'react-overlays/DropdownMenu';
import useMergedRefs from '@restart/hooks/useMergedRefs';
import NavbarContext from './NavbarContext';
import { useBootstrapPrefix } from './ThemeProvider';
import useWrappedRefWithWarning from './useWrappedRefWithWarning';
var defaultProps = {
  alignRight: false,
  flip: true
};
var DropdownMenu = React.forwardRef(function (_ref, ref) {
  var bsPrefix = _ref.bsPrefix,
      className = _ref.className,
      alignRight = _ref.alignRight,
      rootCloseEvent = _ref.rootCloseEvent,
      flip = _ref.flip,
      popperConfig = _ref.popperConfig,
      showProps = _ref.show,
      _ref$as = _ref.as,
      Component = _ref$as === void 0 ? 'div' : _ref$as,
      props = _objectWithoutPropertiesLoose(_ref, ["bsPrefix", "className", "alignRight", "rootCloseEvent", "flip", "popperConfig", "show", "as"]);

  var isNavbar = useContext(NavbarContext);
  var prefix = useBootstrapPrefix(bsPrefix, 'dropdown-menu');

  var _useDropdownMenu = useDropdownMenu({
    flip: flip,
    popperConfig: popperConfig,
    rootCloseEvent: rootCloseEvent,
    show: showProps,
    alignEnd: alignRight,
    usePopper: !isNavbar
  }),
      hasShown = _useDropdownMenu.hasShown,
      placement = _useDropdownMenu.placement,
      show = _useDropdownMenu.show,
      alignEnd = _useDropdownMenu.alignEnd,
      close = _useDropdownMenu.close,
      menuProps = _useDropdownMenu.props;

  menuProps.ref = useMergedRefs(menuProps.ref, useWrappedRefWithWarning(ref, 'DropdownMenu'));
  if (!hasShown) return null; // For custom components provide additional, non-DOM, props;

  if (typeof Component !== 'string') {
    menuProps.show = show;
    menuProps.close = close;
    menuProps.alignRight = alignEnd;
  }

  var style = props.style;

  if (placement) {
    // we don't need the default popper style,
    // menus are display: none when not shown.
    style = _extends({}, style, {}, menuProps.style);
    props['x-placement'] = placement;
  }

  return React.createElement(Component, _extends({}, props, menuProps, {
    style: style,
    className: classNames(className, prefix, show && 'show', alignEnd && prefix + "-right")
  }));
});
DropdownMenu.displayName = 'DropdownMenu';
DropdownMenu.defaultProps = defaultProps;
export default DropdownMenu;