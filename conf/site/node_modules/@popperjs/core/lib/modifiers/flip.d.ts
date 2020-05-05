import { Placement, Boundary, RootBoundary } from "../enums";
import { Modifier, Padding } from "../types";
declare type Options = {
    fallbackPlacements: Array<Placement>;
    padding: Padding;
    boundary: Boundary;
    rootBoundary: RootBoundary;
    altBoundary: boolean;
    flipVariations: boolean;
    allowedAutoPlacements: Array<Placement>;
};
export declare type FlipModifier = Modifier<"flip", Options>;
declare const _default: Modifier<"flip", Options>;
export default _default;
