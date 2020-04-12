import getCompositeRect from "./getCompositeRect.js";
import getWindow from "./getWindow.js";
import getDocumentElement from "./getDocumentElement.js";
import getWindowScroll from "./getWindowScroll.js";
export default function getDocumentRect(element) {
  var win = getWindow(element);
  var winScroll = getWindowScroll(element);
  var documentRect = getCompositeRect(getDocumentElement(element), win);
  documentRect.height = Math.max(documentRect.height, win.innerHeight);
  documentRect.width = Math.max(documentRect.width, win.innerWidth);
  documentRect.x = -winScroll.scrollLeft;
  documentRect.y = -winScroll.scrollTop;
  return documentRect;
}