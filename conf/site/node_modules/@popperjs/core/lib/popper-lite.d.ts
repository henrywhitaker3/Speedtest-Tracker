import { popperGenerator } from "./index";
import detectOverflow from "./utils/detectOverflow";
export * from "./types";
declare const defaultModifiers: (import("./types").Modifier<"popperOffsets", {}> | import("./types").Modifier<"eventListeners", {
    scroll: boolean;
    resize: boolean;
}> | import("./types").Modifier<"computeStyles", {
    gpuAcceleration: boolean;
    adaptive: boolean;
}> | import("./types").Modifier<"applyStyles", {}>)[];
declare const createPopper: <TModifier extends Partial<import("./types").Modifier<any, any>>>(reference: Element | import("./types").VirtualElement, popper: HTMLElement, options?: Partial<import("./types").OptionsGeneric<TModifier>>) => import("./types").Instance;
export { createPopper, popperGenerator, defaultModifiers, detectOverflow };
