"use strict";

var _interopRequireWildcard = require("@babel/runtime/helpers/interopRequireWildcard");

var _interopRequireDefault = require("@babel/runtime/helpers/interopRequireDefault");

exports.__esModule = true;
exports.default = void 0;

var _extends2 = _interopRequireDefault(require("@babel/runtime/helpers/extends"));

var _objectWithoutPropertiesLoose2 = _interopRequireDefault(require("@babel/runtime/helpers/objectWithoutPropertiesLoose"));

var _react = _interopRequireWildcard(require("react"));

var _classnames = _interopRequireDefault(require("classnames"));

var _Overlay = _interopRequireDefault(require("react-overlays/Overlay"));

var _safeFindDOMNode = _interopRequireDefault(require("react-overlays/safeFindDOMNode"));

var _usePopperMarginModifiers = _interopRequireDefault(require("./usePopperMarginModifiers"));

var _Fade = _interopRequireDefault(require("./Fade"));

var defaultProps = {
  transition: _Fade.default,
  rootClose: false,
  show: false,
  placement: 'top'
};

function wrapRefs(props, arrowProps) {
  var ref = props.ref;
  var aRef = arrowProps.ref;

  props.ref = ref.__wrapped || (ref.__wrapped = function (r) {
    return ref((0, _safeFindDOMNode.default)(r));
  });

  arrowProps.ref = aRef.__wrapped || (aRef.__wrapped = function (r) {
    return aRef((0, _safeFindDOMNode.default)(r));
  });
}

function Overlay(_ref) {
  var overlay = _ref.children,
      transition = _ref.transition,
      _ref$popperConfig = _ref.popperConfig,
      popperConfig = _ref$popperConfig === void 0 ? {} : _ref$popperConfig,
      outerProps = (0, _objectWithoutPropertiesLoose2.default)(_ref, ["children", "transition", "popperConfig"]);
  var popperRef = (0, _react.useRef)({});

  var _usePopperMarginModif = (0, _usePopperMarginModifiers.default)(),
      ref = _usePopperMarginModif[0],
      marginModifiers = _usePopperMarginModif[1];

  transition = transition === true ? _Fade.default : transition || null;
  return /*#__PURE__*/_react.default.createElement(_Overlay.default, (0, _extends2.default)({}, outerProps, {
    ref: ref,
    popperConfig: (0, _extends2.default)({}, popperConfig, {
      modifiers: marginModifiers.concat(popperConfig.modifiers || [])
    }),
    transition: transition
  }), function (_ref2) {
    var overlayProps = _ref2.props,
        arrowProps = _ref2.arrowProps,
        show = _ref2.show,
        state = _ref2.state,
        scheduleUpdate = _ref2.scheduleUpdate,
        placement = _ref2.placement,
        outOfBoundaries = _ref2.outOfBoundaries,
        props = (0, _objectWithoutPropertiesLoose2.default)(_ref2, ["props", "arrowProps", "show", "state", "scheduleUpdate", "placement", "outOfBoundaries"]);
    wrapRefs(overlayProps, arrowProps);
    var popper = Object.assign(popperRef.current, {
      state: state,
      scheduleUpdate: scheduleUpdate,
      placement: placement,
      outOfBoundaries: outOfBoundaries
    });
    if (typeof overlay === 'function') return overlay((0, _extends2.default)({}, props, {}, overlayProps, {
      placement: placement,
      show: show,
      popper: popper,
      arrowProps: arrowProps
    }));
    return _react.default.cloneElement(overlay, (0, _extends2.default)({}, props, {}, overlayProps, {
      placement: placement,
      arrowProps: arrowProps,
      popper: popper,
      className: (0, _classnames.default)(overlay.props.className, !transition && show && 'show'),
      style: (0, _extends2.default)({}, overlay.props.style, {}, overlayProps.style)
    }));
  });
}

Overlay.defaultProps = defaultProps;
var _default = Overlay;
exports.default = _default;
module.exports = exports["default"];