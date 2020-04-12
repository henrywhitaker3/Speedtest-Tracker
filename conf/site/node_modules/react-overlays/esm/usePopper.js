import _extends from "@babel/runtime/helpers/esm/extends";
import _objectWithoutPropertiesLoose from "@babel/runtime/helpers/esm/objectWithoutPropertiesLoose";
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import useSafeState from '@restart/hooks/useSafeState';
import { createPopper } from './popper';
var initialPopperStyles = {
  position: 'absolute',
  top: '0',
  left: '0',
  opacity: '0',
  pointerEvents: 'none'
};
var initialArrowStyles = {};
export function toModifierMap(modifiers) {
  var result = {};

  if (!Array.isArray(modifiers)) {
    return modifiers || result;
  } // eslint-disable-next-line no-unused-expressions


  modifiers == null ? void 0 : modifiers.forEach(function (m) {
    result[m.name] = m;
  });
  return result;
}
export function toModifierArray(map) {
  if (Array.isArray(map)) return map;
  return Object.keys(map || {}).map(function (k) {
    map[k].name = k;
    return map[k];
  });
}
/**
 * Position an element relative some reference element using Popper.js
 *
 * @param {HTMLElement} referenceElement The element
 * @param {HTMLElement} popperElement
 * @param {Object}      options
 * @param {Object}      options.modifiers Popper.js modifiers
 * @param {Boolean}     options.enabled toggle the popper functionality on/off
 * @param {String}      options.placement The popper element placement relative to the reference element
 * @param {Boolean}     options.positionFixed use fixed positioning
 * @param {Boolean}     options.eventsEnabled have Popper listen on window resize events to reposition the element
 * @param {Function}    options.onCreate called when the popper is created
 * @param {Function}    options.onUpdate called when the popper is updated
 */

export default function usePopper(referenceElement, popperElement, _ref) {
  if (_ref === void 0) {
    _ref = {};
  }

  var _ref2 = _ref,
      _ref2$enabled = _ref2.enabled,
      enabled = _ref2$enabled === void 0 ? true : _ref2$enabled,
      _ref2$placement = _ref2.placement,
      placement = _ref2$placement === void 0 ? 'bottom' : _ref2$placement,
      _ref2$strategy = _ref2.strategy,
      strategy = _ref2$strategy === void 0 ? 'absolute' : _ref2$strategy,
      _ref2$eventsEnabled = _ref2.eventsEnabled,
      eventsEnabled = _ref2$eventsEnabled === void 0 ? true : _ref2$eventsEnabled,
      userModifiers = _ref2.modifiers,
      popperOptions = _objectWithoutPropertiesLoose(_ref2, ["enabled", "placement", "strategy", "eventsEnabled", "modifiers"]);

  var popperInstanceRef = useRef();
  var scheduleUpdate = useCallback(function () {
    if (popperInstanceRef.current) {
      popperInstanceRef.current.update();
    }
  }, []);

  var _useSafeState = useSafeState(useState({
    placement: placement,
    scheduleUpdate: scheduleUpdate,
    outOfBoundaries: false,
    styles: initialPopperStyles,
    arrowStyles: initialArrowStyles
  })),
      state = _useSafeState[0],
      setState = _useSafeState[1];

  var updateModifier = useMemo(function () {
    return {
      name: 'updateStateModifier',
      enabled: true,
      phase: 'afterWrite',
      requires: ['computeStyles'],
      fn: function fn(data) {
        var _data$state$modifiers, _data$state$styles, _data$state$styles2;

        setState({
          scheduleUpdate: scheduleUpdate,
          outOfBoundaries: (_data$state$modifiers = data.state.modifiersData.hide) == null ? void 0 : _data$state$modifiers.isReferenceHidden,
          placement: data.state.placement,
          styles: _extends({}, (_data$state$styles = data.state.styles) == null ? void 0 : _data$state$styles.popper),
          arrowStyles: _extends({}, (_data$state$styles2 = data.state.styles) == null ? void 0 : _data$state$styles2.arrow),
          state: data.state
        });
      }
    };
  }, [scheduleUpdate, setState]);
  var modifiers = toModifierArray(userModifiers);
  var eventsModifier = modifiers.find(function (m) {
    return m.name === 'eventListeners';
  });

  if (!eventsModifier && eventsEnabled) {
    eventsModifier = {
      name: 'eventListeners',
      enabled: true
    };
    modifiers = [].concat(modifiers, [eventsModifier]);
  } // A placement difference in state means popper determined a new placement
  // apart from the props value. By the time the popper element is rendered with
  // the new position Popper has already measured it, if the place change triggers
  // a size change it will result in a misaligned popper. So we schedule an update to be sure.


  useEffect(function () {
    scheduleUpdate();
  }, [state.placement, scheduleUpdate]);
  useEffect(function () {
    if (!popperInstanceRef.current || !enabled) return;
    popperInstanceRef.current.setOptions({
      placement: placement,
      strategy: strategy,
      modifiers: [].concat(modifiers, [updateModifier])
    }); // intentionally NOT re-running on new modifiers
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [strategy, placement, eventsModifier.enabled, updateModifier, enabled]);
  useEffect(function () {
    if (!enabled || referenceElement == null || popperElement == null) {
      return undefined;
    }

    popperInstanceRef.current = createPopper(referenceElement, popperElement, _extends({}, popperOptions, {
      placement: placement,
      strategy: strategy,
      modifiers: [].concat(modifiers, [updateModifier])
    }));
    return function () {
      if (popperInstanceRef.current !== null) {
        popperInstanceRef.current.destroy();
        popperInstanceRef.current = null;
        setState(function (s) {
          return _extends({}, s, {
            styles: initialPopperStyles,
            arrowStyles: initialArrowStyles
          });
        });
      }
    }; // This is only run once to _create_ the popper
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [enabled, referenceElement, popperElement]);
  return state;
}