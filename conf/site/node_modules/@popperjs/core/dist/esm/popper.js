import { popperGenerator } from "./index.js";
import eventListeners from "./modifiers/eventListeners.js";
import popperOffsets from "./modifiers/popperOffsets.js";
import computeStyles from "./modifiers/computeStyles.js";
import applyStyles from "./modifiers/applyStyles.js";
import offset from "./modifiers/offset.js";
import flip from "./modifiers/flip.js";
import preventOverflow from "./modifiers/preventOverflow.js";
import arrow from "./modifiers/arrow.js";
import hide from "./modifiers/hide.js";
import detectOverflow from "./utils/detectOverflow.js";
var defaultModifiers = [eventListeners, popperOffsets, computeStyles, applyStyles, offset, flip, preventOverflow, arrow, hide];
var createPopper = /*#__PURE__*/popperGenerator({
  defaultModifiers: defaultModifiers
}); // eslint-disable-next-line import/no-unused-modules

export { createPopper, popperGenerator, defaultModifiers, detectOverflow };