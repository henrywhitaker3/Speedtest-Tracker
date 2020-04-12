import _extends from "@babel/runtime/helpers/esm/extends";
import _objectWithoutPropertiesLoose from "@babel/runtime/helpers/esm/objectWithoutPropertiesLoose";
import React from 'react';
import PropTypes from 'prop-types';
import Button from './Button';
import ButtonGroup from './ButtonGroup';
import Dropdown from './Dropdown';
var propTypes = {
  /**
   * An html id attribute for the Toggle button, necessary for assistive technologies, such as screen readers.
   * @type {string|number}
   * @required
   */
  id: PropTypes.any,

  /**
   * Accessible label for the toggle; the value of `title` if not specified.
   */
  toggleLabel: PropTypes.string,

  /** An `href` passed to the non-toggle Button */
  href: PropTypes.string,

  /** An anchor `target` passed to the non-toggle Button */
  target: PropTypes.string,

  /** An `onClick` handler passed to the non-toggle Button */
  onClick: PropTypes.func,

  /** The content of the non-toggle Button.  */
  title: PropTypes.node.isRequired,

  /** A `type` passed to the non-toggle Button */
  type: PropTypes.string,

  /** Disables both Buttons  */
  disabled: PropTypes.bool,

  /** An ARIA accessible role applied to the Menu component. When set to 'menu', The dropdown */
  menuRole: PropTypes.string,

  /**
   *  Which event when fired outside the component will cause it to be closed.
   *
   * _see [DropdownMenu](#menu-props) for more details_
   */
  rootCloseEvent: PropTypes.string,

  /** @ignore */
  bsPrefix: PropTypes.string,

  /** @ignore */
  variant: PropTypes.string,

  /** @ignore */
  size: PropTypes.string
};
var defaultProps = {
  toggleLabel: 'Toggle dropdown',
  type: 'button'
};
var SplitButton = React.forwardRef(function (_ref, ref) {
  var id = _ref.id,
      bsPrefix = _ref.bsPrefix,
      size = _ref.size,
      variant = _ref.variant,
      title = _ref.title,
      type = _ref.type,
      toggleLabel = _ref.toggleLabel,
      children = _ref.children,
      onClick = _ref.onClick,
      href = _ref.href,
      target = _ref.target,
      menuRole = _ref.menuRole,
      rootCloseEvent = _ref.rootCloseEvent,
      props = _objectWithoutPropertiesLoose(_ref, ["id", "bsPrefix", "size", "variant", "title", "type", "toggleLabel", "children", "onClick", "href", "target", "menuRole", "rootCloseEvent"]);

  return React.createElement(Dropdown, _extends({
    ref: ref
  }, props, {
    as: ButtonGroup
  }), React.createElement(Button, {
    size: size,
    variant: variant,
    disabled: props.disabled,
    bsPrefix: bsPrefix,
    href: href,
    target: target,
    onClick: onClick,
    type: type
  }, title), React.createElement(Dropdown.Toggle, {
    split: true,
    id: id,
    size: size,
    variant: variant,
    disabled: props.disabled,
    childBsPrefix: bsPrefix
  }, React.createElement("span", {
    className: "sr-only"
  }, toggleLabel)), React.createElement(Dropdown.Menu, {
    role: menuRole,
    rootCloseEvent: rootCloseEvent
  }, children));
});
SplitButton.propTypes = propTypes;
SplitButton.defaultProps = defaultProps;
SplitButton.displayName = 'SplitButton';
export default SplitButton;