import { popperGenerator } from "./index";
import detectOverflow from "./utils/detectOverflow";
export * from "./types";
declare const defaultModifiers: (import("./types").Modifier<"popperOffsets", {}> | import("./types").Modifier<"flip", {
    fallbackPlacements: import("./enums").Placement[];
    padding: import("./types").Padding;
    boundary: import("./enums").Boundary;
    rootBoundary: import("./enums").RootBoundary;
    altBoundary: boolean;
    flipVariations: boolean;
    allowedAutoPlacements: import("./enums").Placement[];
}> | import("./types").Modifier<"hide", {}> | import("./types").Modifier<"offset", {
    offset: ((arg0: {
        popper: import("./types").Rect;
        reference: import("./types").Rect;
        placement: import("./enums").Placement;
    }) => [number, number]) | [number, number];
}> | import("./types").Modifier<"eventListeners", {
    scroll: boolean;
    resize: boolean;
}> | import("./types").Modifier<"computeStyles", {
    gpuAcceleration: boolean;
    adaptive: boolean;
}> | import("./types").Modifier<"arrow", {
    element: string | HTMLElement;
    padding: import("./types").Padding;
}> | import("./types").Modifier<"preventOverflow", {
    mainAxis: boolean;
    altAxis: boolean;
    boundary: import("./enums").Boundary;
    rootBoundary: import("./enums").RootBoundary;
    altBoundary: boolean;
    tether: boolean;
    tetherOffset: (arg0: {
        popper: import("./types").Rect;
        reference: import("./types").Rect;
        placement: import("./enums").Placement;
    }) => number;
    padding: import("./types").Padding;
}> | import("./types").Modifier<"applyStyles", {}>)[];
declare const createPopper: <TModifier extends Partial<import("./types").Modifier<any, any>>>(reference: Element | import("./types").VirtualElement, popper: HTMLElement, options?: Partial<import("./types").OptionsGeneric<TModifier>>) => import("./types").Instance;
export { createPopper, popperGenerator, defaultModifiers, detectOverflow };
