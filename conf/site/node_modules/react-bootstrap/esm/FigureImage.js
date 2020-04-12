import _extends from "@babel/runtime/helpers/esm/extends";
import _objectWithoutPropertiesLoose from "@babel/runtime/helpers/esm/objectWithoutPropertiesLoose";
import classNames from 'classnames';
import PropTypes from 'prop-types';
import React from 'react';
import Image from './Image';
var propTypes = {
  /**
   * @default 'img'
   */
  bsPrefix: PropTypes.string,

  /**
   * Sets image as fluid image.
   */
  fluid: PropTypes.bool,

  /**
   * Sets image shape as rounded.
   */
  rounded: PropTypes.bool,

  /**
   * Sets image shape as circle.
   */
  roundedCircle: PropTypes.bool,

  /**
   * Sets image shape as thumbnail.
   */
  thumbnail: PropTypes.bool
};
var defaultProps = {
  fluid: true
};
var FigureImage = React.forwardRef(function (_ref, ref) {
  var className = _ref.className,
      props = _objectWithoutPropertiesLoose(_ref, ["className"]);

  return React.createElement(Image, _extends({
    ref: ref
  }, props, {
    className: classNames(className, 'figure-img')
  }));
});
FigureImage.displayName = 'FigureImage';
FigureImage.propTypes = propTypes;
FigureImage.defaultProps = defaultProps;
export default FigureImage;