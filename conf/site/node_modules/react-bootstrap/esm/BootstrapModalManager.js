import _inheritsLoose from "@babel/runtime/helpers/esm/inheritsLoose";
import css from 'dom-helpers/css';
import qsa from 'dom-helpers/querySelectorAll';
import getScrollbarSize from 'dom-helpers/scrollbarSize';
import ModalManager from 'react-overlays/ModalManager';
var Selector = {
  FIXED_CONTENT: '.fixed-top, .fixed-bottom, .is-fixed, .sticky-top',
  STICKY_CONTENT: '.sticky-top',
  NAVBAR_TOGGLER: '.navbar-toggler'
};

var BootstrapModalManager = /*#__PURE__*/function (_ModalManager) {
  _inheritsLoose(BootstrapModalManager, _ModalManager);

  function BootstrapModalManager() {
    var _this;

    for (var _len = arguments.length, args = new Array(_len), _key = 0; _key < _len; _key++) {
      args[_key] = arguments[_key];
    }

    _this = _ModalManager.call.apply(_ModalManager, [this].concat(args)) || this;

    _this.adjustAndStore = function (prop, element, adjust) {
      var _css;

      var actual = element.style[prop];
      element.dataset[prop] = actual;
      css(element, (_css = {}, _css[prop] = parseFloat(css(element, prop)) + adjust + "px", _css));
    };

    _this.restore = function (prop, element) {
      var value = element.dataset[prop];

      if (value !== undefined) {
        var _css2;

        delete element.dataset[prop];
        css(element, (_css2 = {}, _css2[prop] = value, _css2));
      }
    };

    return _this;
  }

  var _proto = BootstrapModalManager.prototype;

  _proto.setContainerStyle = function setContainerStyle(containerState, container) {
    var _this2 = this;

    _ModalManager.prototype.setContainerStyle.call(this, containerState, container);

    if (!containerState.overflowing) return;
    var size = getScrollbarSize();
    qsa(container, Selector.FIXED_CONTENT).forEach(function (el) {
      return _this2.adjustAndStore('paddingRight', el, size);
    });
    qsa(container, Selector.STICKY_CONTENT).forEach(function (el) {
      return _this2.adjustAndStore('margingRight', el, -size);
    });
    qsa(container, Selector.NAVBAR_TOGGLER).forEach(function (el) {
      return _this2.adjustAndStore('margingRight', el, size);
    });
  };

  _proto.removeContainerStyle = function removeContainerStyle(containerState, container) {
    var _this3 = this;

    _ModalManager.prototype.removeContainerStyle.call(this, containerState, container);

    qsa(container, Selector.FIXED_CONTENT).forEach(function (el) {
      return _this3.restore('paddingRight', el);
    });
    qsa(container, Selector.STICKY_CONTENT).forEach(function (el) {
      return _this3.restore('margingRight', el);
    });
    qsa(container, Selector.NAVBAR_TOGGLER).forEach(function (el) {
      return _this3.restore('margingRight', el);
    });
  };

  return BootstrapModalManager;
}(ModalManager);

export { BootstrapModalManager as default };