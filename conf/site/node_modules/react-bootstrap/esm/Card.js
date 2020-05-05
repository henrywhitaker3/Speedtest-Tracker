import _extends from "@babel/runtime/helpers/esm/extends";
import _objectWithoutPropertiesLoose from "@babel/runtime/helpers/esm/objectWithoutPropertiesLoose";
import classNames from 'classnames';
import React, { useMemo } from 'react';
import { useBootstrapPrefix } from './ThemeProvider';
import createWithBsPrefix from './createWithBsPrefix';
import divWithClassName from './divWithClassName';
import CardContext from './CardContext';
import CardImg from './CardImg';
var DivStyledAsH5 = divWithClassName('h5');
var DivStyledAsH6 = divWithClassName('h6');
var CardBody = createWithBsPrefix('card-body');
var defaultProps = {
  body: false
};
var Card = React.forwardRef(function (_ref, ref) {
  var bsPrefix = _ref.bsPrefix,
      className = _ref.className,
      bg = _ref.bg,
      text = _ref.text,
      border = _ref.border,
      body = _ref.body,
      children = _ref.children,
      _ref$as = _ref.as,
      Component = _ref$as === void 0 ? 'div' : _ref$as,
      props = _objectWithoutPropertiesLoose(_ref, ["bsPrefix", "className", "bg", "text", "border", "body", "children", "as"]);

  var prefix = useBootstrapPrefix(bsPrefix, 'card');
  var cardContext = useMemo(function () {
    return {
      cardHeaderBsPrefix: prefix + "-header"
    };
  }, [prefix]);
  return /*#__PURE__*/React.createElement(CardContext.Provider, {
    value: cardContext
  }, /*#__PURE__*/React.createElement(Component, _extends({
    ref: ref
  }, props, {
    className: classNames(className, prefix, bg && "bg-" + bg, text && "text-" + text, border && "border-" + border)
  }), body ? /*#__PURE__*/React.createElement(CardBody, null, children) : children));
});
Card.displayName = 'Card';
Card.defaultProps = defaultProps;
Card.Img = CardImg;
Card.Title = createWithBsPrefix('card-title', {
  Component: DivStyledAsH5
});
Card.Subtitle = createWithBsPrefix('card-subtitle', {
  Component: DivStyledAsH6
});
Card.Body = CardBody;
Card.Link = createWithBsPrefix('card-link', {
  Component: 'a'
});
Card.Text = createWithBsPrefix('card-text', {
  Component: 'p'
});
Card.Header = createWithBsPrefix('card-header');
Card.Footer = createWithBsPrefix('card-footer');
Card.ImgOverlay = createWithBsPrefix('card-img-overlay');
export default Card;