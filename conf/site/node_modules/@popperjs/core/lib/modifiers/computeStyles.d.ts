import { PositioningStrategy, Offsets, Modifier, Rect } from "../types";
import { BasePlacement } from "../enums";
declare type Options = {
    gpuAcceleration: boolean;
    adaptive: boolean;
};
export declare function mapToStyles({ popper, popperRect, placement, offsets, position, gpuAcceleration, adaptive }: {
    popper: HTMLElement;
    popperRect: Rect;
    placement: BasePlacement;
    offsets: Offsets;
    position: PositioningStrategy;
    gpuAcceleration: boolean;
    adaptive: boolean;
}): {
    transform: string;
    top: string;
    right: string;
    bottom: string;
    left: string;
    position: PositioningStrategy;
};
declare const _default: Modifier<Options>;
export default _default;
