import _extends from "@babel/runtime/helpers/esm/extends";
import _objectWithoutPropertiesLoose from "@babel/runtime/helpers/esm/objectWithoutPropertiesLoose";
import _inheritsLoose from "@babel/runtime/helpers/esm/inheritsLoose";
import contains from 'dom-helpers/contains';
import React, { cloneElement, useCallback, useRef } from 'react';
import useTimeout from '@restart/hooks/useTimeout';
import safeFindDOMNode from 'react-overlays/safeFindDOMNode';
import warning from 'warning';
import { useUncontrolledProp } from 'uncontrollable';
import Overlay from './Overlay';

var RefHolder = /*#__PURE__*/function (_React$Component) {
  _inheritsLoose(RefHolder, _React$Component);

  function RefHolder() {
    return _React$Component.apply(this, arguments) || this;
  }

  var _proto = RefHolder.prototype;

  _proto.render = function render() {
    return this.props.children;
  };

  return RefHolder;
}(React.Component);

function normalizeDelay(delay) {
  return delay && typeof delay === 'object' ? delay : {
    show: delay,
    hide: delay
  };
} // Simple implementation of mouseEnter and mouseLeave.
// React's built version is broken: https://github.com/facebook/react/issues/4251
// for cases when the trigger is disabled and mouseOut/Over can cause flicker
// moving from one child element to another.


function handleMouseOverOut(handler, args, relatedNative) {
  var e = args[0];
  var target = e.currentTarget;
  var related = e.relatedTarget || e.nativeEvent[relatedNative];

  if ((!related || related !== target) && !contains(target, related)) {
    handler.apply(void 0, args);
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
      propsShow = _ref.show,
      _ref$defaultShow = _ref.defaultShow,
      defaultShow = _ref$defaultShow === void 0 ? false : _ref$defaultShow,
      onToggle = _ref.onToggle,
      propsDelay = _ref.delay,
      placement = _ref.placement,
      _ref$flip = _ref.flip,
      flip = _ref$flip === void 0 ? placement && placement.indexOf('auto') !== -1 : _ref$flip,
      props = _objectWithoutPropertiesLoose(_ref, ["trigger", "overlay", "children", "popperConfig", "show", "defaultShow", "onToggle", "delay", "placement", "flip"]);

  var triggerNodeRef = useRef(null);
  var timeout = useTimeout();
  var hoverStateRef = useRef('');

  var _useUncontrolledProp = useUncontrolledProp(propsShow, defaultShow, onToggle),
      show = _useUncontrolledProp[0],
      setShow = _useUncontrolledProp[1];

  var delay = normalizeDelay(propsDelay);
  var child = React.Children.only(children);
  var _child$props = child.props,
      onFocus = _child$props.onFocus,
      onBlur = _child$props.onBlur,
      onClick = _child$props.onClick;
  var getTarget = useCallback(function () {
    return safeFindDOMNode(triggerNodeRef.current);
  }, []);
  var handleShow = useCallback(function () {
    timeout.clear();
    hoverStateRef.current = 'show';

    if (!delay.show) {
      setShow(true);
      return;
    }

    timeout.set(function () {
      if (hoverStateRef.current === 'show') setShow(true);
    }, delay.show);
  }, [delay.show, setShow, timeout]);
  var handleHide = useCallback(function () {
    timeout.clear();
    hoverStateRef.current = 'hide';

    if (!delay.hide) {
      setShow(false);
      return;
    }

    timeout.set(function () {
      if (hoverStateRef.current === 'hide') setShow(false);
    }, delay.hide);
  }, [delay.hide, setShow, timeout]);
  var handleFocus = useCallback(function () {
    handleShow();
    if (onFocus) onFocus.apply(void 0, arguments);
  }, [handleShow, onFocus]);
  var handleBlur = useCallback(function () {
    handleHide();
    if (onBlur) onBlur.apply(void 0, arguments);
  }, [handleHide, onBlur]);
  var handleClick = useCallback(function () {
    setShow(!show);
    if (onClick) onClick.apply(void 0, arguments);
  }, [onClick, setShow, show]);
  var handleMouseOver = useCallback(function () {
    for (var _len = arguments.length, args = new Array(_len), _key = 0; _key < _len; _key++) {
      args[_key] = arguments[_key];
    }

    handleMouseOverOut(handleShow, args, 'fromElement');
  }, [handleShow]);
  var handleMouseOut = useCallback(function () {
    for (var _len2 = arguments.length, args = new Array(_len2), _key2 = 0; _key2 < _len2; _key2++) {
      args[_key2] = arguments[_key2];
    }

    handleMouseOverOut(handleHide, args, 'toElement');
  }, [handleHide]); // We add aria-describedby in the case where the overlay is a role="tooltip"
  // for other cases describedby isn't appropriate (e.g. a popover with inputs) so we don't add it.

  var ariaModifier = {
    name: 'ariaDescribedBy',
    enabled: true,
    phase: 'afterWrite',
    effect: function effect(_ref2) {
      var state = _ref2.state;
      return function () {
        if ('removeAttribute' in state.elements.reference) state.elements.reference.removeAttribute('aria-describedby');
      };
    },
    fn: function fn(_ref3) {
      var state = _ref3.state;
      var _state$elements = state.elements,
          popper = _state$elements.popper,
          reference = _state$elements.reference;
      if (!show || !reference) return;
      var role = popper.getAttribute('role') || '';

      if (popper.id && role.toLowerCase() === 'tooltip' && 'setAttribute' in reference) {
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
    process.env.NODE_ENV !== "production" ? warning(triggers.length > 1, '[react-bootstrap] Specifying only the `"hover"` trigger limits the visibility of the overlay to just mouse users. Consider also including the `"focus"` trigger so that touch and keyboard only users can see the overlay as well.') : void 0;
    triggerProps.onMouseOver = handleMouseOver;
    triggerProps.onMouseOut = handleMouseOut;
  } // TODO: fix typing
  // @ts-ignore


  var modifiers = [ariaModifier].concat(popperConfig.modifiers || []);
  return /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement(RefHolder, {
    ref: triggerNodeRef
  }, cloneElement(child, triggerProps)), /*#__PURE__*/React.createElement(Overlay, _extends({}, props, {
    popperConfig: _extends({}, popperConfig, {
      modifiers: modifiers
    }),
    show: show,
    onHide: handleHide,
    target: getTarget,
    placement: placement,
    flip: flip
  }), overlay));
}

OverlayTrigger.defaultProps = defaultProps;
export default OverlayTrigger;