"use strict";

var _interopRequireWildcard = require("@babel/runtime/helpers/interopRequireWildcard");

var _interopRequireDefault = require("@babel/runtime/helpers/interopRequireDefault");

exports.__esModule = true;
exports.default = void 0;

var _extends2 = _interopRequireDefault(require("@babel/runtime/helpers/extends"));

var _objectWithoutPropertiesLoose2 = _interopRequireDefault(require("@babel/runtime/helpers/objectWithoutPropertiesLoose"));

var _useEventCallback = _interopRequireDefault(require("@restart/hooks/useEventCallback"));

var _useUpdateEffect = _interopRequireDefault(require("@restart/hooks/useUpdateEffect"));

var _useTimeout = _interopRequireDefault(require("@restart/hooks/useTimeout"));

var _classnames = _interopRequireDefault(require("classnames"));

var _transitionEnd = _interopRequireDefault(require("dom-helpers/transitionEnd"));

var _Transition = _interopRequireDefault(require("react-transition-group/Transition"));

var _react = _interopRequireWildcard(require("react"));

var _uncontrollable = require("uncontrollable");

var _CarouselCaption = _interopRequireDefault(require("./CarouselCaption"));

var _CarouselItem = _interopRequireDefault(require("./CarouselItem"));

var _ElementChildren = require("./ElementChildren");

var _SafeAnchor = _interopRequireDefault(require("./SafeAnchor"));

var _ThemeProvider = require("./ThemeProvider");

var _triggerBrowserReflow = _interopRequireDefault(require("./triggerBrowserReflow"));

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
  prevIcon: _react.default.createElement("span", {
    "aria-hidden": "true",
    className: "carousel-control-prev-icon"
  }),
  prevLabel: 'Previous',
  nextIcon: _react.default.createElement("span", {
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

var Carousel = _react.default.forwardRef(function (uncontrolledProps, ref) {
  var _useUncontrolled = (0, _uncontrollable.useUncontrolled)(uncontrolledProps, {
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
      props = (0, _objectWithoutPropertiesLoose2.default)(_useUncontrolled, ["as", "bsPrefix", "slide", "fade", "controls", "indicators", "activeIndex", "onSelect", "onSlide", "onSlid", "interval", "keyboard", "onKeyDown", "pause", "onMouseOver", "onMouseOut", "wrap", "touch", "onTouchStart", "onTouchMove", "onTouchEnd", "prevIcon", "prevLabel", "nextIcon", "nextLabel", "className", "children"]);

  var prefix = (0, _ThemeProvider.useBootstrapPrefix)(bsPrefix, 'carousel');
  var nextDirectionRef = (0, _react.useRef)(null);

  var _useState = (0, _react.useState)('next'),
      direction = _useState[0],
      setDirection = _useState[1];

  var _useState2 = (0, _react.useState)(false),
      isSliding = _useState2[0],
      setIsSliding = _useState2[1];

  var _useState3 = (0, _react.useState)(activeIndex),
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

  var numChildren = _react.default.Children.toArray(children).filter(_react.default.isValidElement).length;

  var prev = (0, _react.useCallback)(function (event) {
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

  var next = (0, _useEventCallback.default)(function (event) {
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
  var elementRef = (0, _react.useRef)();
  (0, _react.useImperativeHandle)(ref, function () {
    return {
      element: elementRef.current,
      prev: prev,
      next: next
    };
  }); // This is used in the setInterval, so it should not invalidate.

  var nextWhenVisible = (0, _useEventCallback.default)(function () {
    if (!document.hidden && isVisible(elementRef.current)) {
      next();
    }
  });
  var slideDirection = direction === 'next' ? 'left' : 'right';
  (0, _useUpdateEffect.default)(function () {
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
  var handleEnter = (0, _react.useCallback)(function (node) {
    (0, _triggerBrowserReflow.default)(node);

    if (onSlide) {
      onSlide(renderedActiveIndex, slideDirection);
    }
  }, [onSlide, renderedActiveIndex, slideDirection]);
  var handleEntered = (0, _react.useCallback)(function () {
    setIsSliding(false);

    if (onSlid) {
      onSlid(renderedActiveIndex, slideDirection);
    }
  }, [onSlid, renderedActiveIndex, slideDirection]);
  var handleKeyDown = (0, _react.useCallback)(function (event) {
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

  var _useState4 = (0, _react.useState)(false),
      pausedOnHover = _useState4[0],
      setPausedOnHover = _useState4[1];

  var handleMouseOver = (0, _react.useCallback)(function (event) {
    if (pause === 'hover') {
      setPausedOnHover(true);
    }

    if (onMouseOver) {
      onMouseOver(event);
    }
  }, [pause, onMouseOver]);
  var handleMouseOut = (0, _react.useCallback)(function (event) {
    setPausedOnHover(false);

    if (onMouseOut) {
      onMouseOut(event);
    }
  }, [onMouseOut]);
  var touchStartXRef = (0, _react.useRef)(0);
  var touchDeltaXRef = (0, _react.useRef)(0);

  var _useState5 = (0, _react.useState)(false),
      pausedOnTouch = _useState5[0],
      setPausedOnTouch = _useState5[1];

  var touchUnpauseTimeout = (0, _useTimeout.default)();
  var handleTouchStart = (0, _react.useCallback)(function (event) {
    touchStartXRef.current = event.touches[0].clientX;
    touchDeltaXRef.current = 0;

    if (touch) {
      setPausedOnTouch(true);
    }

    if (onTouchStart) {
      onTouchStart(event);
    }
  }, [touch, onTouchStart]);
  var handleTouchMove = (0, _react.useCallback)(function (event) {
    if (event.touches && event.touches.length > 1) {
      touchDeltaXRef.current = 0;
    } else {
      touchDeltaXRef.current = event.touches[0].clientX - touchStartXRef.current;
    }

    if (onTouchMove) {
      onTouchMove(event);
    }
  }, [onTouchMove]);
  var handleTouchEnd = (0, _react.useCallback)(function (event) {
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
  var intervalHandleRef = (0, _react.useRef)();
  (0, _react.useEffect)(function () {
    if (!shouldPlay) {
      return undefined;
    }

    intervalHandleRef.current = setInterval(document.visibilityState ? nextWhenVisible : next, interval);
    return function () {
      clearInterval(intervalHandleRef.current);
    };
  }, [shouldPlay, next, interval, nextWhenVisible]);
  var indicatorOnClicks = (0, _react.useMemo)(function () {
    return indicators && Array.from({
      length: numChildren
    }, function (_, index) {
      return function (event) {
        onSelect(index, event);
      };
    });
  }, [indicators, numChildren, onSelect]);
  return _react.default.createElement(Component, (0, _extends2.default)({
    ref: elementRef
  }, props, {
    onKeyDown: handleKeyDown,
    onMouseOver: handleMouseOver,
    onMouseOut: handleMouseOut,
    onTouchStart: handleTouchStart,
    onTouchMove: handleTouchMove,
    onTouchEnd: handleTouchEnd,
    className: (0, _classnames.default)(className, prefix, slide && 'slide', fade && prefix + "-fade")
  }), indicators && _react.default.createElement("ol", {
    className: prefix + "-indicators"
  }, (0, _ElementChildren.map)(children, function (child, index) {
    return _react.default.createElement("li", {
      key: index,
      className: index === renderedActiveIndex ? 'active' : null,
      onClick: indicatorOnClicks[index]
    });
  })), _react.default.createElement("div", {
    className: prefix + "-inner"
  }, (0, _ElementChildren.map)(children, function (child, index) {
    var isActive = index === renderedActiveIndex;
    return slide ? _react.default.createElement(_Transition.default, {
      in: isActive,
      onEnter: isActive ? handleEnter : null,
      onEntered: isActive ? handleEntered : null,
      addEndListener: _transitionEnd.default
    }, function (status) {
      return _react.default.cloneElement(child, {
        className: (0, _classnames.default)(child.props.className, isActive && status !== 'entered' && orderClassName, (status === 'entered' || status === 'exiting') && 'active', (status === 'entering' || status === 'exiting') && directionalClassName)
      });
    }) : _react.default.cloneElement(child, {
      className: (0, _classnames.default)(child.props.className, isActive && 'active')
    });
  })), controls && _react.default.createElement(_react.default.Fragment, null, (wrap || activeIndex !== 0) && _react.default.createElement(_SafeAnchor.default, {
    className: prefix + "-control-prev",
    onClick: prev
  }, prevIcon, prevLabel && _react.default.createElement("span", {
    className: "sr-only"
  }, prevLabel)), (wrap || activeIndex !== numChildren - 1) && _react.default.createElement(_SafeAnchor.default, {
    className: prefix + "-control-next",
    onClick: next
  }, nextIcon, nextLabel && _react.default.createElement("span", {
    className: "sr-only"
  }, nextLabel))));
});

Carousel.displayName = 'Carousel';
Carousel.defaultProps = defaultProps;
Carousel.Caption = _CarouselCaption.default;
Carousel.Item = _CarouselItem.default;
var _default = Carousel;
exports.default = _default;
module.exports = exports["default"];