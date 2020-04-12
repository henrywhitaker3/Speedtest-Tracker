import _extends from "@babel/runtime/helpers/esm/extends";
import React from 'react';
import classNames from 'classnames';
export default (function (className) {
  return React.forwardRef(function (p, ref) {
    return React.createElement("div", _extends({}, p, {
      ref: ref,
      className: classNames(p.className, className)
    }));
  });
});