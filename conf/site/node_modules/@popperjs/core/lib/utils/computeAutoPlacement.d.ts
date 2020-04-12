import { State, Padding } from "../types";
import { Placement, Boundary, RootBoundary, ComputedPlacement } from "../enums";
declare type Options = {
    placement: Placement;
    padding: Padding;
    boundary: Boundary;
    rootBoundary: RootBoundary;
    flipVariations: boolean;
};
export default function computeAutoPlacement(state: Partial<State>, options?: Options): Array<ComputedPlacement>;
export {};
