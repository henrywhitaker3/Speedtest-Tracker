import _extends from "@babel/runtime/helpers/esm/extends";
import _objectWithoutPropertiesLoose from "@babel/runtime/helpers/esm/objectWithoutPropertiesLoose";
import classNames from 'classnames';
import React from 'react';
import { useBootstrapPrefix } from './ThemeProvider';
var defaultProps = {
  fluid: false,
  rounded: false,
  roundedCircle: false,
  thumbnail: false
};
var Image = React.forwardRef(function (_ref, ref) {
  var bsPrefix = _ref.bsPrefix,
      className = _ref.className,
      fluid = _ref.fluid,
      rounded = _ref.rounded,
      roundedCircle = _ref.roundedCircle,
      thumbnail = _ref.thumbnail,
      props = _objectWithoutPropertiesLoose(_ref, ["bsPrefix", "className", "fluid", "rounded", "roundedCircle", "thumbnail"]);

  bsPrefix = useBootstrapPrefix(bsPrefix, 'img');
  var classes = classNames(fluid && bsPrefix + "-fluid", rounded && "rounded", roundedCircle && "rounded-circle", thumbnail && bsPrefix + "-thumbnail");
  return /*#__PURE__*/React.createElement("img", _extends({
    // eslint-disable-line jsx-a11y/alt-text
    ref: ref
  }, props, {
    className: classNames(className, classes)
  }));
});
Image.displayName = 'Image';
Image.defaultProps = defaultProps;
export default Image;