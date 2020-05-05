"use strict";

var _interopRequireDefault = require("@babel/runtime/helpers/interopRequireDefault");

exports.__esModule = true;
exports.default = void 0;

var _objectWithoutPropertiesLoose2 = _interopRequireDefault(require("@babel/runtime/helpers/objectWithoutPropertiesLoose"));

var _extends2 = _interopRequireDefault(require("@babel/runtime/helpers/extends"));

var _inheritsLoose2 = _interopRequireDefault(require("@babel/runtime/helpers/inheritsLoose"));

var _classnames = _interopRequireDefault(require("classnames"));

var _addEventListener = _interopRequireDefault(require("dom-helpers/addEventListener"));

var _canUseDOM = _interopRequireDefault(require("dom-helpers/canUseDOM"));

var _ownerDocument = _interopRequireDefault(require("dom-helpers/ownerDocument"));

var _removeEventListener = _interopRequireDefault(require("dom-helpers/removeEventListener"));

var _scrollbarSize = _interopRequireDefault(require("dom-helpers/scrollbarSize"));

var _react = _interopRequireDefault(require("react"));

var _Modal = _interopRequireDefault(require("react-overlays/Modal"));

var _BootstrapModalManager = _interopRequireDefault(require("./BootstrapModalManager"));

var _Fade = _interopRequireDefault(require("./Fade"));

var _ModalBody = _interopRequireDefault(require("./ModalBody"));

var _ModalContext = _interopRequireDefault(require("./ModalContext"));

var _ModalDialog = _interopRequireDefault(require("./ModalDialog"));

var _ModalFooter = _interopRequireDefault(require("./ModalFooter"));

var _ModalHeader = _interopRequireDefault(require("./ModalHeader"));

var _ModalTitle = _interopRequireDefault(require("./ModalTitle"));

var _ThemeProvider = require("./ThemeProvider");

var manager;
var defaultProps = {
  show: false,
  backdrop: true,
  keyboard: true,
  autoFocus: true,
  enforceFocus: true,
  restoreFocus: true,
  animation: true,
  dialogAs: _ModalDialog.default
};
/* eslint-disable no-use-before-define, react/no-multi-comp */

function DialogTransition(props) {
  return /*#__PURE__*/_react.default.createElement(_Fade.default, props);
}

function BackdropTransition(props) {
  return /*#__PURE__*/_react.default.createElement(_Fade.default, props);
}
/* eslint-enable no-use-before-define */


var Modal = /*#__PURE__*/function (_React$Component) {
  (0, _inheritsLoose2.default)(Modal, _React$Component);

  function Modal() {
    var _this;

    for (var _len = arguments.length, _args = new Array(_len), _key = 0; _key < _len; _key++) {
      _args[_key] = arguments[_key];
    }

    _this = _React$Component.call.apply(_React$Component, [this].concat(_args)) || this;
    _this.state = {
      style: {}
    };
    _this.modalContext = {
      onHide: function onHide() {
        return _this.props.onHide();
      }
    };

    _this.setModalRef = function (ref) {
      _this._modal = ref;
    };

    _this.handleDialogMouseDown = function () {
      _this._waitingForMouseUp = true;
    };

    _this.handleMouseUp = function (e) {
      if (_this._waitingForMouseUp && e.target === _this._modal.dialog) {
        _this._ignoreBackdropClick = true;
      }

      _this._waitingForMouseUp = false;
    };

    _this.handleClick = function (e) {
      if (_this._ignoreBackdropClick || e.target !== e.currentTarget) {
        _this._ignoreBackdropClick = false;
        return;
      }

      _this.props.onHide();
    };

    _this.handleEnter = function (node) {
      var _this$props;

      if (node) {
        node.style.display = 'block';

        _this.updateDialogStyle(node);
      }

      for (var _len2 = arguments.length, args = new Array(_len2 > 1 ? _len2 - 1 : 0), _key2 = 1; _key2 < _len2; _key2++) {
        args[_key2 - 1] = arguments[_key2];
      }

      if (_this.props.onEnter) (_this$props = _this.props).onEnter.apply(_this$props, [node].concat(args));
    };

    _this.handleEntering = function (node) {
      var _this$props2;

      for (var _len3 = arguments.length, args = new Array(_len3 > 1 ? _len3 - 1 : 0), _key3 = 1; _key3 < _len3; _key3++) {
        args[_key3 - 1] = arguments[_key3];
      }

      if (_this.props.onEntering) (_this$props2 = _this.props).onEntering.apply(_this$props2, [node].concat(args)); // FIXME: This should work even when animation is disabled.

      (0, _addEventListener.default)(window, 'resize', _this.handleWindowResize);
    };

    _this.handleExited = function (node) {
      var _this$props3;

      if (node) node.style.display = ''; // RHL removes it sometimes

      for (var _len4 = arguments.length, args = new Array(_len4 > 1 ? _len4 - 1 : 0), _key4 = 1; _key4 < _len4; _key4++) {
        args[_key4 - 1] = arguments[_key4];
      }

      if (_this.props.onExited) (_this$props3 = _this.props).onExited.apply(_this$props3, args); // FIXME: This should work even when animation is disabled.

      (0, _removeEventListener.default)(window, 'resize', _this.handleWindowResize);
    };

    _this.handleWindowResize = function () {
      _this.updateDialogStyle(_this._modal.dialog);
    };

    _this.getModalManager = function () {
      if (_this.props.manager) {
        return _this.props.manager;
      }

      if (!manager) {
        manager = new _BootstrapModalManager.default();
      }

      return manager;
    };

    _this.renderBackdrop = function (props) {
      var _this$props4 = _this.props,
          bsPrefix = _this$props4.bsPrefix,
          backdropClassName = _this$props4.backdropClassName,
          animation = _this$props4.animation;
      return /*#__PURE__*/_react.default.createElement("div", (0, _extends2.default)({}, props, {
        className: (0, _classnames.default)(bsPrefix + "-backdrop", backdropClassName, !animation && 'show')
      }));
    };

    return _this;
  }

  var _proto = Modal.prototype;

  _proto.componentWillUnmount = function componentWillUnmount() {
    // Clean up the listener if we need to.
    (0, _removeEventListener.default)(window, 'resize', this.handleWindowResize);
  };

  _proto.updateDialogStyle = function updateDialogStyle(node) {
    if (!_canUseDOM.default) return;
    var containerIsOverflowing = this.getModalManager().isContainerOverflowing(this._modal);
    var modalIsOverflowing = node.scrollHeight > (0, _ownerDocument.default)(node).documentElement.clientHeight;
    this.setState({
      style: {
        paddingRight: containerIsOverflowing && !modalIsOverflowing ? (0, _scrollbarSize.default)() : undefined,
        paddingLeft: !containerIsOverflowing && modalIsOverflowing ? (0, _scrollbarSize.default)() : undefined
      }
    });
  };

  _proto.render = function render() {
    var _this$props5 = this.props,
        bsPrefix = _this$props5.bsPrefix,
        className = _this$props5.className,
        style = _this$props5.style,
        dialogClassName = _this$props5.dialogClassName,
        children = _this$props5.children,
        Dialog = _this$props5.dialogAs,
        ariaLabelledby = _this$props5['aria-labelledby'],
        show = _this$props5.show,
        animation = _this$props5.animation,
        backdrop = _this$props5.backdrop,
        keyboard = _this$props5.keyboard,
        onEscapeKeyDown = _this$props5.onEscapeKeyDown,
        onShow = _this$props5.onShow,
        onHide = _this$props5.onHide,
        container = _this$props5.container,
        autoFocus = _this$props5.autoFocus,
        enforceFocus = _this$props5.enforceFocus,
        restoreFocus = _this$props5.restoreFocus,
        restoreFocusOptions = _this$props5.restoreFocusOptions,
        onEntered = _this$props5.onEntered,
        onExit = _this$props5.onExit,
        onExiting = _this$props5.onExiting,
        _ = _this$props5.onExited,
        _1 = _this$props5.onEntering,
        _6 = _this$props5.onEnter,
        _4 = _this$props5.onEntering,
        _2 = _this$props5.backdropClassName,
        props = (0, _objectWithoutPropertiesLoose2.default)(_this$props5, ["bsPrefix", "className", "style", "dialogClassName", "children", "dialogAs", "aria-labelledby", "show", "animation", "backdrop", "keyboard", "onEscapeKeyDown", "onShow", "onHide", "container", "autoFocus", "enforceFocus", "restoreFocus", "restoreFocusOptions", "onEntered", "onExit", "onExiting", "onExited", "onEntering", "onEnter", "onEntering", "backdropClassName"]);
    var clickHandler = backdrop === true ? this.handleClick : null;
    var baseModalStyle = (0, _extends2.default)({}, style, {}, this.state.style); // Sets `display` always block when `animation` is false

    if (!animation) baseModalStyle.display = 'block';
    return /*#__PURE__*/_react.default.createElement(_ModalContext.default.Provider, {
      value: this.modalContext
    }, /*#__PURE__*/_react.default.createElement(_Modal.default, {
      show: show,
      backdrop: backdrop,
      container: container,
      keyboard: keyboard,
      autoFocus: autoFocus,
      enforceFocus: enforceFocus,
      restoreFocus: restoreFocus,
      restoreFocusOptions: restoreFocusOptions,
      onEscapeKeyDown: onEscapeKeyDown,
      onShow: onShow,
      onHide: onHide,
      onEntered: onEntered,
      onExit: onExit,
      onExiting: onExiting,
      manager: this.getModalManager(),
      ref: this.setModalRef,
      style: baseModalStyle,
      className: (0, _classnames.default)(className, bsPrefix),
      containerClassName: bsPrefix + "-open",
      transition: animation ? DialogTransition : undefined,
      backdropTransition: animation ? BackdropTransition : undefined,
      renderBackdrop: this.renderBackdrop,
      onClick: clickHandler,
      onMouseUp: this.handleMouseUp,
      onEnter: this.handleEnter,
      onEntering: this.handleEntering,
      onExited: this.handleExited,
      'aria-labelledby': ariaLabelledby
    }, /*#__PURE__*/_react.default.createElement(Dialog, (0, _extends2.default)({}, props, {
      onMouseDown: this.handleDialogMouseDown,
      className: dialogClassName
    }), children)));
  };

  return Modal;
}(_react.default.Component);

Modal.defaultProps = defaultProps;
var DecoratedModal = (0, _ThemeProvider.createBootstrapComponent)(Modal, 'modal');
DecoratedModal.Body = _ModalBody.default;
DecoratedModal.Header = _ModalHeader.default;
DecoratedModal.Title = _ModalTitle.default;
DecoratedModal.Footer = _ModalFooter.default;
DecoratedModal.Dialog = _ModalDialog.default;
DecoratedModal.TRANSITION_DURATION = 300;
DecoratedModal.BACKDROP_TRANSITION_DURATION = 150;
var _default = DecoratedModal;
exports.default = _default;
module.exports = exports["default"];