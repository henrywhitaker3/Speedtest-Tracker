"use strict";

var _interopRequireWildcard = require("@babel/runtime/helpers/interopRequireWildcard");

var _interopRequireDefault = require("@babel/runtime/helpers/interopRequireDefault");

exports.__esModule = true;
exports.default = void 0;

var _extends2 = _interopRequireDefault(require("@babel/runtime/helpers/extends"));

var _objectWithoutPropertiesLoose2 = _interopRequireDefault(require("@babel/runtime/helpers/objectWithoutPropertiesLoose"));

var _classnames = _interopRequireDefault(require("classnames"));

var _react = _interopRequireWildcard(require("react"));

var _DropdownMenu = require("react-overlays/DropdownMenu");

var _useMergedRefs = _interopRequireDefault(require("@restart/hooks/useMergedRefs"));

var _NavbarContext = _interopRequireDefault(require("./NavbarContext"));

var _ThemeProvider = require("./ThemeProvider");

var _useWrappedRefWithWarning = _interopRequireDefault(require("./useWrappedRefWithWarning"));

var defaultProps = {
  alignRight: false,
  flip: true
};

var DropdownMenu = _react.default.forwardRef(function (_ref, ref) {
  var bsPrefix = _ref.bsPrefix,
      className = _ref.className,
      alignRight = _ref.alignRight,
      rootCloseEvent = _ref.rootCloseEvent,
      flip = _ref.flip,
      popperConfig = _ref.popperConfig,
      showProps = _ref.show,
      _ref$as = _ref.as,
      Component = _ref$as === void 0 ? 'div' : _ref$as,
      props = (0, _objectWithoutPropertiesLoose2.default)(_ref, ["bsPrefix", "className", "alignRight", "rootCloseEvent", "flip", "popperConfig", "show", "as"]);
  var isNavbar = (0, _react.useContext)(_NavbarContext.default);
  var prefix = (0, _ThemeProvider.useBootstrapPrefix)(bsPrefix, 'dropdown-menu');

  var _useDropdownMenu = (0, _DropdownMenu.useDropdownMenu)({
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

  menuProps.ref = (0, _useMergedRefs.default)(menuProps.ref, (0, _useWrappedRefWithWarning.default)(ref, 'DropdownMenu'));
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
    style = (0, _extends2.default)({}, style, {}, menuProps.style);
    props['x-placement'] = placement;
  }

  return _react.default.createElement(Component, (0, _extends2.default)({}, props, menuProps, {
    style: style,
    className: (0, _classnames.default)(className, prefix, show && 'show', alignEnd && prefix + "-right")
  }));
});

DropdownMenu.displayName = 'DropdownMenu';
DropdownMenu.defaultProps = defaultProps;
var _default = DropdownMenu;
exports.default = _default;
module.exports = exports["default"];