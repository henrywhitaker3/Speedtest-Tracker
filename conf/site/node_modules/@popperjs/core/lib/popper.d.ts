import { popperGenerator } from "./index";
import detectOverflow from "./utils/detectOverflow";
export * from "./types";
declare const defaultModifiers: import("./types").Modifier<{}>[];
declare const createPopper: (reference: Element | import("./types").VirtualElement, popper: HTMLElement, options?: Partial<import("./types").Options>) => import("./types").Instance;
export { createPopper, popperGenerator, defaultModifiers, detectOverflow };
