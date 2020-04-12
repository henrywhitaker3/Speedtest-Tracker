import _extends from "@babel/runtime/helpers/esm/extends";
import _objectWithoutPropertiesLoose from "@babel/runtime/helpers/esm/objectWithoutPropertiesLoose";
import useEventCallback from '@restart/hooks/useEventCallback';
import useUpdateEffect from '@restart/hooks/useUpdateEffect';
import useTimeout from '@restart/hooks/useTimeout';
import classNames from 'classnames';
import transitionEnd from 'dom-helpers/transitionEnd';
import Transition from 'react-transition-group/Transition';
import React, { useCallback, useEffect, useMemo, useRef, useState, useImperativeHandle } from 'react';
import { useUncontrolled } from 'uncontrollable';
import CarouselCaption from './CarouselCaption';
import CarouselItem from './CarouselItem';
import { map } from './ElementChildren';
import SafeAnchor from './SafeAnchor';
import { useBootstrapPrefix } from './ThemeProvider';
import triggerBrowserReflow from './triggerBrowserReflow';
var SWIPE_THRESHOLD = 40;
var defaultProps = {
  slide: true,
  fade: false,
  controls: true,
  indicators: true,
  defaultActiveIndex: 0,
  interval: 5000,
  keyboard: true,
  pause: 'hover',
  wrap: true,
  touch: true,
  prevIcon: React.createElement("span", {
    "aria-hidden": "true",
    className: "carousel-control-prev-icon"
  }),
  prevLabel: 'Previous',
  nextIcon: React.createElement("span", {
    "aria-hidden": "true",
    className: "carousel-control-next-icon"
  }),
  nextLabel: 'Next'
};

function isVisible(element) {
  if (!element || !element.style || !element.parentNode || !element.parentNode.style) {
    return false;
  }

  var elementStyle = getComputedStyle(element);
  return elementStyle.display !== 'none' && elementStyle.visibility !== 'hidden' && getComputedStyle(element.parentNode).display !== 'none';
}

var Carousel = React.forwardRef(function (uncontrolledProps, ref) {
  var _useUncontrolled = useUncontrolled(uncontrolledProps, {
    activeIndex: 'onSelect'
  }),
      _useUncontrolled$as = _useUncontrolled.as,
      Component = _useUncontrolled$as === void 0 ? 'div' : _useUncontrolled$as,
      bsPrefix = _useUncontrolled.bsPrefix,
      slide = _useUncontrolled.slide,
      fade = _useUncontrolled.fade,
      controls = _useUncontrolled.controls,
      indicators = _useUncontrolled.indicators,
      activeIndex = _useUncontrolled.activeIndex,
      onSelect = _useUncontrolled.onSelect,
      onSlide = _useUncontrolled.onSlide,
      onSlid = _useUncontrolled.onSlid,
      interval = _useUncontrolled.interval,
      keyboard = _useUncontrolled.keyboard,
      onKeyDown = _useUncontrolled.onKeyDown,
      pause = _useUncontrolled.pause,
      onMouseOver = _useUncontrolled.onMouseOver,
      onMouseOut = _useUncontrolled.onMouseOut,
      wrap = _useUncontrolled.wrap,
      touch = _useUncontrolled.touch,
      onTouchStart = _useUncontrolled.onTouchStart,
      onTouchMove = _useUncontrolled.onTouchMove,
      onTouchEnd = _useUncontrolled.onTouchEnd,
      prevIcon = _useUncontrolled.prevIcon,
      prevLabel = _useUncontrolled.prevLabel,
      nextIcon = _useUncontrolled.nextIcon,
      nextLabel = _useUncontrolled.nextLabel,
      className = _useUncontrolled.className,
      children = _useUncontrolled.children,
      props = _objectWithoutPropertiesLoose(_useUncontrolled, ["as", "bsPrefix", "slide", "fade", "controls", "indicators", "activeIndex", "onSelect", "onSlide", "onSlid", "interval", "keyboard", "onKeyDown", "pause", "onMouseOver", "onMouseOut", "wrap", "touch", "onTouchStart", "onTouchMove", "onTouchEnd", "prevIcon", "prevLabel", "nextIcon", "nextLabel", "className", "children"]);

  var prefix = useBootstrapPrefix(bsPrefix, 'carousel');
  var nextDirectionRef = useRef(null);

  var _useState = useState('next'),
      direction = _useState[0],
      setDirection = _useState[1];

  var _useState2 = useState(false),
      isSliding = _useState2[0],
      setIsSliding = _useState2[1];

  var _useState3 = useState(activeIndex),
      renderedActiveIndex = _useState3[0],
      setRenderedActiveIndex = _useState3[1];

  if (!isSliding && activeIndex !== renderedActiveIndex) {
    if (nextDirectionRef.current) {
      setDirection(nextDirectionRef.current);
      nextDirectionRef.current = null;
    } else {
      setDirection(activeIndex > renderedActiveIndex ? 'next' : 'prev');
    }

    if (slide) {
      setIsSliding(true);
    }

    setRenderedActiveIndex(activeIndex);
  }

  var numChildren = React.Children.toArray(children).filter(React.isValidElement).length;
  var prev = useCallback(function (event) {
    if (isSliding) {
      return;
    }

    var nextActiveIndex = renderedActiveIndex - 1;

    if (nextActiveIndex < 0) {
      if (!wrap) {
        return;
      }

      nextActiveIndex = numChildren - 1;
    }

    nextDirectionRef.current = 'prev';
    onSelect(nextActiveIndex, event);
  }, [isSliding, renderedActiveIndex, onSelect, wrap, numChildren]); // This is used in the setInterval, so it should not invalidate.

  var next = useEventCallback(function (event) {
    if (isSliding) {
      return;
    }

    var nextActiveIndex = renderedActiveIndex + 1;

    if (nextActiveIndex >= numChildren) {
      if (!wrap) {
        return;
      }

      nextActiveIndex = 0;
    }

    nextDirectionRef.current = 'next';
    onSelect(nextActiveIndex, event);
  });
  var elementRef = useRef();
  useImperativeHandle(ref, function () {
    return {
      element: elementRef.current,
      prev: prev,
      next: next
    };
  }); // This is used in the setInterval, so it should not invalidate.

  var nextWhenVisible = useEventCallback(function () {
    if (!document.hidden && isVisible(elementRef.current)) {
      next();
    }
  });
  var slideDirection = direction === 'next' ? 'left' : 'right';
  useUpdateEffect(function () {
    if (slide) {
      // These callbacks will be handled by the <Transition> callbacks.
      return;
    }

    if (onSlide) {
      onSlide(renderedActiveIndex, slideDirection);
    }

    if (onSlid) {
      onSlid(renderedActiveIndex, slideDirection);
    }
  }, [renderedActiveIndex]);
  var orderClassName = prefix + "-item-" + direction;
  var directionalClassName = prefix + "-item-" + slideDirection;
  var handleEnter = useCallback(function (node) {
    triggerBrowserReflow(node);

    if (onSlide) {
      onSlide(renderedActiveIndex, slideDirection);
    }
  }, [onSlide, renderedActiveIndex, slideDirection]);
  var handleEntered = useCallback(function () {
    setIsSliding(false);

    if (onSlid) {
      onSlid(renderedActiveIndex, slideDirection);
    }
  }, [onSlid, renderedActiveIndex, slideDirection]);
  var handleKeyDown = useCallback(function (event) {
    if (keyboard && !/input|textarea/i.test(event.target.tagName)) {
      switch (event.key) {
        case 'ArrowLeft':
          event.preventDefault();
          prev(event);
          return;

        case 'ArrowRight':
          event.preventDefault();
          next(event);
          return;

        default:
      }
    }

    if (onKeyDown) {
      onKeyDown(event);
    }
  }, [keyboard, onKeyDown, prev, next]);

  var _useState4 = useState(false),
      pausedOnHover = _useState4[0],
      setPausedOnHover = _useState4[1];

  var handleMouseOver = useCallback(function (event) {
    if (pause === 'hover') {
      setPausedOnHover(true);
    }

    if (onMouseOver) {
      onMouseOver(event);
    }
  }, [pause, onMouseOver]);
  var handleMouseOut = useCallback(function (event) {
    setPausedOnHover(false);

    if (onMouseOut) {
      onMouseOut(event);
    }
  }, [onMouseOut]);
  var touchStartXRef = useRef(0);
  var touchDeltaXRef = useRef(0);

  var _useState5 = useState(false),
      pausedOnTouch = _useState5[0],
      setPausedOnTouch = _useState5[1];

  var touchUnpauseTimeout = useTimeout();
  var handleTouchStart = useCallback(function (event) {
    touchStartXRef.current = event.touches[0].clientX;
    touchDeltaXRef.current = 0;

    if (touch) {
      setPausedOnTouch(true);
    }

    if (onTouchStart) {
      onTouchStart(event);
    }
  }, [touch, onTouchStart]);
  var handleTouchMove = useCallback(function (event) {
    if (event.touches && event.touches.length > 1) {
      touchDeltaXRef.current = 0;
    } else {
      touchDeltaXRef.current = event.touches[0].clientX - touchStartXRef.current;
    }

    if (onTouchMove) {
      onTouchMove(event);
    }
  }, [onTouchMove]);
  var handleTouchEnd = useCallback(function (event) {
    if (touch) {
      var touchDeltaX = touchDeltaXRef.current;

      if (Math.abs(touchDeltaX) <= SWIPE_THRESHOLD) {
        return;
      }

      if (touchDeltaX > 0) {
        prev(event);
      } else {
        next(event);
      }
    }

    touchUnpauseTimeout.set(function () {
      setPausedOnTouch(false);
    }, interval);

    if (onTouchEnd) {
      onTouchEnd(event);
    }
  }, [touch, prev, next, touchUnpauseTimeout, interval, onTouchEnd]);
  var shouldPlay = interval != null && !pausedOnHover && !pausedOnTouch && !isSliding;
  var intervalHandleRef = useRef();
  useEffect(function () {
    if (!shouldPlay) {
      return undefined;
    }

    intervalHandleRef.current = setInterval(document.visibilityState ? nextWhenVisible : next, interval);
    return function () {
      clearInterval(intervalHandleRef.current);
    };
  }, [shouldPlay, next, interval, nextWhenVisible]);
  var indicatorOnClicks = useMemo(function () {
    return indicators && Array.from({
      length: numChildren
    }, function (_, index) {
      return function (event) {
        onSelect(index, event);
      };
    });
  }, [indicators, numChildren, onSelect]);
  return React.createElement(Component, _extends({
    ref: elementRef
  }, props, {
    onKeyDown: handleKeyDown,
    onMouseOver: handleMouseOver,
    onMouseOut: handleMouseOut,
    onTouchStart: handleTouchStart,
    onTouchMove: handleTouchMove,
    onTouchEnd: handleTouchEnd,
    className: classNames(className, prefix, slide && 'slide', fade && prefix + "-fade")
  }), indicators && React.createElement("ol", {
    className: prefix + "-indicators"
  }, map(children, function (child, index) {
    return React.createElement("li", {
      key: index,
      className: index === renderedActiveIndex ? 'active' : null,
      onClick: indicatorOnClicks[index]
    });
  })), React.createElement("div", {
    className: prefix + "-inner"
  }, map(children, function (child, index) {
    var isActive = index === renderedActiveIndex;
    return slide ? React.createElement(Transition, {
      in: isActive,
      onEnter: isActive ? handleEnter : null,
      onEntered: isActive ? handleEntered : null,
      addEndListener: transitionEnd
    }, function (status) {
      return React.cloneElement(child, {
        className: classNames(child.props.className, isActive && status !== 'entered' && orderClassName, (status === 'entered' || status === 'exiting') && 'active', (status === 'entering' || status === 'exiting') && directionalClassName)
      });
    }) : React.cloneElement(child, {
      className: classNames(child.props.className, isActive && 'active')
    });
  })), controls && React.createElement(React.Fragment, null, (wrap || activeIndex !== 0) && React.createElement(SafeAnchor, {
    className: prefix + "-control-prev",
    onClick: prev
  }, prevIcon, prevLabel && React.createElement("span", {
    className: "sr-only"
  }, prevLabel)), (wrap || activeIndex !== numChildren - 1) && React.createElement(SafeAnchor, {
    className: prefix + "-control-next",
    onClick: next
  }, nextIcon, nextLabel && React.createElement("span", {
    className: "sr-only"
  }, nextLabel))));
});
Carousel.displayName = 'Carousel';
Carousel.defaultProps = defaultProps;
Carousel.Caption = CarouselCaption;
Carousel.Item = CarouselItem;
export default Carousel;