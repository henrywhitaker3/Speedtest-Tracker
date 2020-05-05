"use strict";

var _interopRequireWildcard = require("@babel/runtime/helpers/interopRequireWildcard");

var _interopRequireDefault = require("@babel/runtime/helpers/interopRequireDefault");

exports.__esModule = true;
exports.default = void 0;

var _extends2 = _interopRequireDefault(require("@babel/runtime/helpers/extends"));

var _objectWithoutPropertiesLoose2 = _interopRequireDefault(require("@babel/runtime/helpers/objectWithoutPropertiesLoose"));

var _inheritsLoose2 = _interopRequireDefault(require("@babel/runtime/helpers/inheritsLoose"));

var _contains = _interopRequireDefault(require("dom-helpers/contains"));

var _react = _interopRequireWildcard(require("react"));

var _useTimeout = _interopRequireDefault(require("@restart/hooks/useTimeout"));

var _safeFindDOMNode = _interopRequireDefault(require("react-overlays/safeFindDOMNode"));

var _warning = _interopRequireDefault(require("warning"));

var _Overlay = _interopRequireDefault(require("./Overlay"));

var RefHolder = /*#__PURE__*/function (_React$Component) {
  (0, _inheritsLoose2.default)(RefHolder, _React$Component);

  function RefHolder() {
    return _React$Component.apply(this, arguments) || this;
  }

  var _proto = RefHolder.prototype;

  _proto.render = function render() {
    return this.props.children;
  };

  return RefHolder;
}(_react.default.Component);

function normalizeDelay(delay) {
  return delay && typeof delay === 'object' ? delay : {
    show: delay,
    hide: delay
  };
} // Simple implementation of mouseEnter and mouseLeave.
// React's built version is broken: https://github.com/facebook/react/issues/4251
// for cases when the trigger is disabled and mouseOut/Over can cause flicker
// moving from one child element to another.


function handleMouseOverOut(handler, e, relatedNative) {
  var target = e.currentTarget;
  var related = e.relatedTarget || e.nativeEvent[relatedNative];

  if ((!related || related !== target) && !(0, _contains.default)(target, related)) {
    handler(e);
  }
}

var defaultProps = {
  defaultShow: false,
  trigger: ['hover', 'focus']
};

function OverlayTrigger(_ref) {
  var trigger = _ref.trigger,
      overlay = _ref.overlay,
      children = _ref.children,
      _ref$popperConfig = _ref.popperConfig,
      popperConfig = _ref$popperConfig === void 0 ? {} : _ref$popperConfig,
      defaultShow = _ref.defaultShow,
      propsDelay = _ref.delay,
      placement = _ref.placement,
      _ref$flip = _ref.flip,
      flip = _ref$flip === void 0 ? placement && placement.indexOf('auto') !== -1 : _ref$flip,
      props = (0, _objectWithoutPropertiesLoose2.default)(_ref, ["trigger", "overlay", "children", "popperConfig", "defaultShow", "delay", "placement", "flip"]);
  var triggerNodeRef = (0, _react.useRef)(null);
  var timeout = (0, _useTimeout.default)();
  var hoverStateRef = (0, _react.useRef)();

  var _useState = (0, _react.useState)(!!defaultShow),
      show = _useState[0],
      setShow = _useState[1];

  var delay = normalizeDelay(propsDelay);

  var child = _react.default.Children.only(children);

  var _child$props = child.props,
      onFocus = _child$props.onFocus,
      onBlur = _child$props.onBlur,
      onClick = _child$props.onClick;
  var getTarget = (0, _react.useCallback)(function () {
    return (0, _safeFindDOMNode.default)(triggerNodeRef.current);
  }, []);
  var handleShow = (0, _react.useCallback)(function () {
    timeout.clear();
    hoverStateRef.current = 'show';

    if (!delay.show) {
      setShow(true);
      return;
    }

    timeout.set(function () {
      if (hoverStateRef.current === 'show') setShow(true);
    }, delay.show);
  }, [delay.show, timeout]);
  var handleHide = (0, _react.useCallback)(function () {
    timeout.clear();
    hoverStateRef.current = 'hide';

    if (!delay.hide) {
      setShow(false);
      return;
    }

    timeout.set(function () {
      if (hoverStateRef.current === 'hide') setShow(false);
    }, delay.hide);
  }, [delay.hide, timeout]);
  var handleFocus = (0, _react.useCallback)(function (e) {
    handleShow(e);
    if (onFocus) onFocus(e);
  }, [handleShow, onFocus]);
  var handleBlur = (0, _react.useCallback)(function (e) {
    handleHide(e);
    if (onBlur) onBlur(e);
  }, [handleHide, onBlur]);
  var handleClick = (0, _react.useCallback)(function (e) {
    setShow(function (prevShow) {
      return !prevShow;
    });
    if (onClick) onClick(e);
  }, [onClick]);
  var handleMouseOver = (0, _react.useCallback)(function (e) {
    handleMouseOverOut(handleShow, e, 'fromElement');
  }, [handleShow]);
  var handleMouseOut = (0, _react.useCallback)(function (e) {
    handleMouseOverOut(handleHide, e, 'toElement');
  }, [handleHide]); // We add aria-describedby in the case where the overlay is a role="tooltip"
  // for other cases describedby isn't appropriate (e.g. a popover with inputs) so we don't add it.

  var ariaModifier = {
    name: 'ariaDescribedBy',
    enabled: true,
    phase: 'afterWrite',
    effect: function effect(_ref2) {
      var state = _ref2.state;
      return function () {
        state.elements.reference.removeAttribute('aria-describedby');
      };
    },
    fn: function fn(_ref3) {
      var state = _ref3.state;
      var _state$elements = state.elements,
          popper = _state$elements.popper,
          reference = _state$elements.reference;
      if (!show || !reference) return;
      var role = popper.getAttribute('role') || '';

      if (popper.id && role.toLowerCase() === 'tooltip') {
        reference.setAttribute('aria-describedby', popper.id);
      }
    }
  };
  var triggers = trigger == null ? [] : [].concat(trigger);
  var triggerProps = {};

  if (triggers.indexOf('click') !== -1) {
    triggerProps.onClick = handleClick;
  }

  if (triggers.indexOf('focus') !== -1) {
    triggerProps.onFocus = handleFocus;
    triggerProps.onBlur = handleBlur;
  }

  if (triggers.indexOf('hover') !== -1) {
    process.env.NODE_ENV !== "production" ? (0, _warning.default)(triggers.length > 1, '[react-bootstrap] Specifying only the `"hover"` trigger limits the visibility of the overlay to just mouse users. Consider also including the `"focus"` trigger so that touch and keyboard only users can see the overlay as well.') : void 0;
    triggerProps.onMouseOver = handleMouseOver;
    triggerProps.onMouseOut = handleMouseOut;
  }

  return /*#__PURE__*/_react.default.createElement(_react.default.Fragment, null, /*#__PURE__*/_react.default.createElement(RefHolder, {
    ref: triggerNodeRef
  }, (0, _react.cloneElement)(child, triggerProps)), /*#__PURE__*/_react.default.createElement(_Overlay.default, (0, _extends2.default)({}, props, {
    popperConfig: (0, _extends2.default)({}, popperConfig, {
      modifiers: [ariaModifier].concat(popperConfig.modifiers || [])
    }),
    show: show,
    onHide: handleHide,
    target: getTarget,
    placement: placement,
    flip: flip
  }), overlay));
}

OverlayTrigger.defaultProps = defaultProps;
var _default = OverlayTrigger;
exports.default = _default;
module.exports = exports["default"];