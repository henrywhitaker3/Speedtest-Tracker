import _extends from "@babel/runtime/helpers/esm/extends";
import _objectWithoutPropertiesLoose from "@babel/runtime/helpers/esm/objectWithoutPropertiesLoose";
import React, { useRef } from 'react';
import classNames from 'classnames';
import BaseOverlay from 'react-overlays/Overlay';
import safeFindDOMNode from 'react-overlays/safeFindDOMNode';
import usePopperMarginModifiers from './usePopperMarginModifiers';
import Fade from './Fade';
var defaultProps = {
  transition: Fade,
  rootClose: false,
  show: false,
  placement: 'top'
};

function wrapRefs(props, arrowProps) {
  var ref = props.ref;
  var aRef = arrowProps.ref;

  props.ref = ref.__wrapped || (ref.__wrapped = function (r) {
    return ref(safeFindDOMNode(r));
  });

  arrowProps.ref = aRef.__wrapped || (aRef.__wrapped = function (r) {
    return aRef(safeFindDOMNode(r));
  });
}

function Overlay(_ref) {
  var overlay = _ref.children,
      transition = _ref.transition,
      _ref$popperConfig = _ref.popperConfig,
      popperConfig = _ref$popperConfig === void 0 ? {} : _ref$popperConfig,
      outerProps = _objectWithoutPropertiesLoose(_ref, ["children", "transition", "popperConfig"]);

  var popperRef = useRef({});

  var _usePopperMarginModif = usePopperMarginModifiers(),
      ref = _usePopperMarginModif[0],
      marginModifiers = _usePopperMarginModif[1];

  transition = transition === true ? Fade : transition || null;
  return /*#__PURE__*/React.createElement(BaseOverlay, _extends({}, outerProps, {
    ref: ref,
    popperConfig: _extends({}, popperConfig, {
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
        props = _objectWithoutPropertiesLoose(_ref2, ["props", "arrowProps", "show", "state", "scheduleUpdate", "placement", "outOfBoundaries"]);

    wrapRefs(overlayProps, arrowProps);
    var popper = Object.assign(popperRef.current, {
      state: state,
      scheduleUpdate: scheduleUpdate,
      placement: placement,
      outOfBoundaries: outOfBoundaries
    });
    if (typeof overlay === 'function') return overlay(_extends({}, props, {}, overlayProps, {
      placement: placement,
      show: show,
      popper: popper,
      arrowProps: arrowProps
    }));
    return React.cloneElement(overlay, _extends({}, props, {}, overlayProps, {
      placement: placement,
      arrowProps: arrowProps,
      popper: popper,
      className: classNames(overlay.props.className, !transition && show && 'show'),
      style: _extends({}, overlay.props.style, {}, overlayProps.style)
    }));
  });
}

Overlay.defaultProps = defaultProps;
export default Overlay;