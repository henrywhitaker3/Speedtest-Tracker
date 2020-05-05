"use strict";

var _interopRequireDefault = require("@babel/runtime/helpers/interopRequireDefault");

exports.__esModule = true;
exports.default = usePopperMarginModifiers;

var _react = require("react");

var _hasClass = _interopRequireDefault(require("dom-helpers/hasClass"));

function getMargins(element) {
  var styles = getComputedStyle(element);
  var top = parseFloat(styles.marginTop) || 0;
  var right = parseFloat(styles.marginRight) || 0;
  var bottom = parseFloat(styles.marginBottom) || 0;
  var left = parseFloat(styles.marginLeft) || 0;
  return {
    top: top,
    right: right,
    bottom: bottom,
    left: left
  };
}

function usePopperMarginModifiers() {
  var overlayRef = (0, _react.useRef)(null);
  var margins = (0, _react.useRef)(null);
  return [(0, _react.useCallback)(function (overlay) {
    if (!overlay || !((0, _hasClass.default)(overlay, 'popover') || (0, _hasClass.default)(overlay, 'dropdown-menu'))) return;
    margins.current = getMargins(overlay);
    overlay.style.margin = 0;
    overlayRef.current = overlay;
  }, []), [(0, _react.useMemo)(function () {
    return {
      name: 'offset',
      options: {
        offset: function offset(_ref) {
          var placement = _ref.placement;
          if (!margins.current) return [0, 0];
          var _margins$current = margins.current,
              top = _margins$current.top,
              left = _margins$current.left,
              bottom = _margins$current.bottom,
              right = _margins$current.right;

          switch (placement.split('-')[0]) {
            case 'top':
              return [0, bottom];

            case 'left':
              return [0, right];

            case 'bottom':
              return [0, top];

            case 'right':
              return [0, left];

            default:
              return [0, 0];
          }
        }
      }
    };
  }, [margins])]];
}

module.exports = exports["default"];