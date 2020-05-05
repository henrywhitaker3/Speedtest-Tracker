import _extends from "@babel/runtime/helpers/esm/extends";
import _objectWithoutPropertiesLoose from "@babel/runtime/helpers/esm/objectWithoutPropertiesLoose";
import classNames from 'classnames';
import React, { useContext } from 'react';
import { useDropdownMenu } from 'react-overlays/DropdownMenu';
import useMergedRefs from '@restart/hooks/useMergedRefs';
import NavbarContext from './NavbarContext';
import { useBootstrapPrefix } from './ThemeProvider';
import useWrappedRefWithWarning from './useWrappedRefWithWarning';
import usePopperMarginModifiers from './usePopperMarginModifiers';
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
      showProps = _ref.show,
      renderOnMount = _ref.renderOnMount,
      _ref$as = _ref.as,
      Component = _ref$as === void 0 ? 'div' : _ref$as,
      _ref$popperConfig = _ref.popperConfig,
      popperConfig = _ref$popperConfig === void 0 ? {} : _ref$popperConfig,
      props = _objectWithoutPropertiesLoose(_ref, ["bsPrefix", "className", "alignRight", "rootCloseEvent", "flip", "show", "renderOnMount", "as", "popperConfig"]);

  var isNavbar = useContext(NavbarContext);
  var prefix = useBootstrapPrefix(bsPrefix, 'dropdown-menu');

  var _usePopperMarginModif = usePopperMarginModifiers(),
      popperRef = _usePopperMarginModif[0],
      marginModifiers = _usePopperMarginModif[1];

  var _useDropdownMenu = useDropdownMenu({
    flip: flip,
    rootCloseEvent: rootCloseEvent,
    show: showProps,
    alignEnd: alignRight,
    usePopper: !isNavbar,
    popperConfig: _extends({}, popperConfig, {
      modifiers: marginModifiers.concat(popperConfig.modifiers || [])
    })
  }),
      hasShown = _useDropdownMenu.hasShown,
      placement = _useDropdownMenu.placement,
      show = _useDropdownMenu.show,
      alignEnd = _useDropdownMenu.alignEnd,
      close = _useDropdownMenu.close,
      menuProps = _useDropdownMenu.props;

  menuProps.ref = useMergedRefs(popperRef, useMergedRefs(useWrappedRefWithWarning(ref, 'DropdownMenu'), menuProps.ref));
  if (!hasShown && !renderOnMount) return null; // For custom components provide additional, non-DOM, props;

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

  return /*#__PURE__*/React.createElement(Component, _extends({}, props, menuProps, {
    style: style,
    className: classNames(className, prefix, show && 'show', alignEnd && prefix + "-right")
  }));
});
DropdownMenu.displayName = 'DropdownMenu';
DropdownMenu.defaultProps = defaultProps;
export default DropdownMenu;