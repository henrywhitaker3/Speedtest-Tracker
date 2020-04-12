import getNodeName from "./getNodeName.js";
export default function getParentNode(element) {
  if (getNodeName(element) === 'html') {
    return element;
  }

  return element.parentNode || // DOM Element detected
  // $FlowFixMe: need a better way to handle this...
  element.host || // ShadowRoot detected
  document.ownerDocument || // Fallback to ownerDocument if available
  document.documentElement // Or to documentElement if everything else fails
  ;
}