import { Placement, Boundary, RootBoundary } from "../enums";
import { Modifier, Padding } from "../types";
declare type Options = {
    fallbackPlacements: Array<Placement>;
    padding: Padding;
    boundary: Boundary;
    rootBoundary: RootBoundary;
    altBoundary: boolean;
    flipVariations: boolean;
};
declare const _default: Modifier<Options>;
export default _default;
