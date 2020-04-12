import getParentNode from "./getParentNode.js";
import getComputedStyle from "./getComputedStyle.js";
import getNodeName from "./getNodeName.js";
import { isHTMLElement } from "./instanceOf.js";
export default function getScrollParent(node) {
  if (['html', 'body', '#document'].indexOf(getNodeName(node)) >= 0) {
    // $FlowFixMe: assume body is always available
    return node.ownerDocument.body;
  }

  if (isHTMLElement(node)) {
    // Firefox wants us to check `-x` and `-y` variations as well
    var _getComputedStyle = getComputedStyle(node),
        overflow = _getComputedStyle.overflow,
        overflowX = _getComputedStyle.overflowX,
        overflowY = _getComputedStyle.overflowY;

    if (/auto|scroll|overlay|hidden/.test(overflow + overflowY + overflowX)) {
      return node;
    }
  }

  return getScrollParent(getParentNode(node));
}