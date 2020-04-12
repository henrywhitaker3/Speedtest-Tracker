import _extends from "@babel/runtime/helpers/esm/extends";
import _objectWithoutPropertiesLoose from "@babel/runtime/helpers/esm/objectWithoutPropertiesLoose";
import _inheritsLoose from "@babel/runtime/helpers/esm/inheritsLoose";
import contains from 'dom-helpers/contains';
import React, { cloneElement, useCallback, useRef, useState } from 'react';
import useTimeout from '@restart/hooks/useTimeout';
import ReactDOM from 'react-dom';
import warning from 'warning';
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


function handleMouseOverOut(handler, e, relatedNative) {
  var target = e.currentTarget;
  var related = e.relatedTarget || e.nativeEvent[relatedNative];

  if ((!related || related !== target) && !contains(target, related)) {
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
      props = _objectWithoutPropertiesLoose(_ref, ["trigger", "overlay", "children", "popperConfig", "defaultShow", "delay"]);

  var triggerNodeRef = useRef(null);
  var timeout = useTimeout();
  var hoverStateRef = useRef();

  var _useState = useState(!!defaultShow),
      show = _useState[0],
      setShow = _useState[1];

  var delay = normalizeDelay(propsDelay);
  var child = React.Children.only(children);
  var _child$props = child.props,
      onFocus = _child$props.onFocus,
      onBlur = _child$props.onBlur,
      onClick = _child$props.onClick;
  var getTarget = useCallback(function () {
    return ReactDOM.findDOMNode(triggerNodeRef.current);
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
  }, [delay.show, timeout]);
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
  }, [delay.hide, timeout]);
  var handleFocus = useCallback(function (e) {
    handleShow(e);
    if (onFocus) onFocus(e);
  }, [handleShow, onFocus]);
  var handleBlur = useCallback(function (e) {
    handleHide(e);
    if (onBlur) onBlur(e);
  }, [handleHide, onBlur]);
  var handleClick = useCallback(function (e) {
    setShow(function (prevShow) {
      return !prevShow;
    });
    if (onClick) onClick(e);
  }, [onClick]);
  var handleMouseOver = useCallback(function (e) {
    handleMouseOverOut(handleShow, e, 'fromElement');
  }, [handleShow]);
  var handleMouseOut = useCallback(function (e) {
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
    process.env.NODE_ENV !== "production" ? warning(triggers.length > 1, '[react-bootstrap] Specifying only the `"hover"` trigger limits the visibility of the overlay to just mouse users. Consider also including the `"focus"` trigger so that touch and keyboard only users can see the overlay as well.') : void 0;
    triggerProps.onMouseOver = handleMouseOver;
    triggerProps.onMouseOut = handleMouseOut;
  }

  return React.createElement(React.Fragment, null, React.createElement(RefHolder, {
    ref: triggerNodeRef
  }, cloneElement(child, triggerProps)), React.createElement(Overlay, _extends({}, props, {
    popperConfig: _extends({}, popperConfig, {
      modifiers: [ariaModifier].concat(popperConfig.modifiers || [])
    }),
    show: show,
    onHide: handleHide,
    target: getTarget
  }), overlay));
}

OverlayTrigger.defaultProps = defaultProps;
export default OverlayTrigger;